package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/teris-io/shortid"
	"github.com/urfave/cli"
)

// The Session contains the data of the current session.
type Session struct {
	ID         string
	Users      map[string]*Client
	startedAt  int32
	finishedAt int32
	Messages   []*ChatMessage
	Queue      chan *ChatMessage
}

// User contains the data of each connected user.
type User struct {
	ID       string
	Username string
	Status   string
	Token    string
}

// Client contains the data of the connection associated with each user.
type Client struct {
	User     *User
	Status   string
	conn     *websocket.Conn
	Session  *Session
	Token    string
	LastPing time.Time
	Queue    chan *ChatMessage
}

type SyncSessionMessage struct {
	SessionID string
	Users     []*User
	Chat      []*ChatMessage
	Type      string
}

type NewSessionMessage struct {
	UserID    string
	SessionID string
	Status    string
	Username  string
}

var sessions map[string]*Session = map[string]*Session{}

var users map[string]*Client = map[string]*Client{}

var sessionIdGenerator, err = shortid.New(1, shortid.DefaultABC, 2342)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var natConnection *nats.Conn
var encodedNatsConnection *nats.EncodedConn

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func StartListener(c *cli.Context) error {

	natConnection, _ = nats.Connect(nats.DefaultURL)
	encodedNatsConnection, _ = nats.NewEncodedConn(natConnection, nats.JSON_ENCODER)

	listeningPort := c.GlobalString("listening-port")

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/ready", readyCheck)
	http.HandleFunc("/session/new", handleWebSocketRequest)
	http.HandleFunc("/ws", handleMessage)

	log.Printf("Server starting on port %v... \n", listeningPort)
	log.Println("Liveness Endpoint: http://localhost:" + listeningPort + "/health")
	log.Println("Readiness Endpoint: http://localhost:" + listeningPort + "/ready")

	// Simple Async Subscriber
	encodedNatsConnection.Subscribe("session.*.chat.out", handleChatMessage)

	var err = http.ListenAndServe(":"+listeningPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	defer encodedNatsConnection.Close()
	return nil
}

func handleChatMessage(subj, reply string, m *ChatMessage) {
	sessionID := strings.Split(subj, ".")[1]
	session, exists := sessions[sessionID]
	if !exists {
		log.Printf("Session [%s] doesn't exists, can't route message.", sessionID)

	}
	session.Messages = append(session.Messages, m)
	log.Printf("Broadcasting message to session [%s].", sessionID)

	for _, c := range session.Users {
		if c.conn != nil {
			c.conn.WriteJSON(m)
		}
	}
}

// Healthcheck endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Up")
}

// Readiness endpoint
func readyCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Ready")
}

func handleWebSocketRequest(w http.ResponseWriter, r *http.Request) {
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var input struct {
			Username  string `json:"Username"`
			SessionID string `json:"SessionID"`
		}
		err := decoder.Decode(&input)
		if err != nil {
			log.Println(err)
		}

		// check if the session exists
		session, exists := sessions[input.SessionID]
		if !exists {
			// create new session
			sessionID, err := sessionIdGenerator.Generate()
			if err != nil {
				log.Println(err)
				return
			}
			session = &Session{
				ID:         sessionID,
				startedAt:  1,
				finishedAt: 0,
				Users:      map[string]*Client{},
				Messages:   []*ChatMessage{},
				Queue:      make(chan *ChatMessage, 200),
			}
			log.Printf("The session [%s] was created.", sessionID)

			// notify that a new session was created.
			encodedNatsConnection.Publish("session.new", &SessionMessage{
				ID: sessionID,
			})
		}

		userID, err := sessionIdGenerator.Generate()
		if err != nil {
			log.Println(err)
			return
		}

		var newClient = Client{
			User: &User{
				ID:       userID,
				Username: input.Username,
			},
			Session: session,
			Token:   userID,
		}

		users[newClient.User.ID] = &newClient
		session.Users[newClient.User.ID] = &newClient
		sessions[session.ID] = session
		newClient.Session = session

		log.Printf("The user [%s] was added to the session [%s]", input.Username, session.ID)

		go handleSession(session)
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	m, _ := url.ParseQuery(r.URL.RawQuery)
	userID := m["token"][0]

	user, exists := users[userID]
	if !exists {
		log.Println("User " + userID + " not found.")
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	user.conn = c
	user.Status = "connected"
	user.Queue = make(chan *ChatMessage, 200)

	go writews(user)
	go readws(user)
}

func writews(user *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("DEFER")
		ticker.Stop()
		user.conn.Close()
	}()
	for {
		select {
		case message, ok := <-user.Queue:
			log.Println("Message pending detected, sending ...")
			user.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Println("The hub closed the channel.")
				user.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			user.conn.WriteJSON(&message)

		case <-ticker.C:
			user.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := user.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func readws(client *Client) {
	defer func() {
		log.Println("DEFER")
		//c.hub.unregister <- c
		client.conn.Close()
	}()
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	var input = &ClientMessage{}

	for {
		err := client.conn.ReadJSON(input)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		handleClientMessage(client, input)
	}
}

func handleClientMessage(client *Client, input *ClientMessage) {
	switch input.Type {
	case "letswork":
		// notify that a user want to start using the workspace
		encodedNatsConnection.Publish("session."+client.Session.ID+".workspace.user."+client.User.ID+".new", &GeneralMessage{
			User: client.User,
		})
		break
	case "letsfinish":
		// notify that a user want to close his workspace
		encodedNatsConnection.Publish("session."+client.Session.ID+".workspace.user."+client.User.ID+".new", &GeneralMessage{
			User: client.User,
		})
		break
	case "goodbye":
		// notify that a user has leaved the chat
		encodedNatsConnection.Publish("session."+client.Session.ID+".chat.user."+client.User.ID+".leave", &GeneralMessage{
			User: client.User,
		})
		break
	case "message":
		encodedNatsConnection.Publish("session."+client.Session.ID+".chat.in", &GeneralMessage{
			User: client.User,,
			Content: input.Content,
		})
		break
	}
}

func handleSession(session *Session) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-session.Queue:
			log.Println("Broadcast Message pending detected, sending ...")
			if !ok {
				// The hub closed the channel.
				log.Println("The hub closed the channel.")
				for _, u := range session.Users {
					u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				}
				return
			}
			for _, u := range session.Users {
				if u.conn != nil {
					u.conn.SetWriteDeadline(time.Now().Add(writeWait))
					u.conn.WriteJSON(&message)
				}
			}

		case <-ticker.C:
			continue
		}
	}
}

func syncUser(client *Client) {
	var users []*User

	log.Println("User ", client.User.Username, " sync data")

	for _, c := range client.Session.Users {
		users = append(users, &User{
			ID:       c.User.ID,
			Username: c.User.Username,
		})
	}

	var newMessage = &SyncSessionMessage{
		SessionID: client.Session.ID,
		Type:      "sync",
		Users:     users,
		Chat:      client.Session.Messages,
	}

	client.conn.WriteJSON(newMessage)
}

func syncUsers(session *Session) {
	var users []*User

	for _, c := range session.Users {
		users = append(users, &User{
			ID:       c.User.ID,
			Username: c.User.Username,
		})
	}

	var newMessage = &SyncSessionMessage{
		Type:      "sync",
		SessionID: session.ID,
		Users:     users,
	}

	// broadcast the messagefo the users connected to the same session
	for _, u := range session.Users {
		u.conn.WriteJSON(newMessage)
	}
}
