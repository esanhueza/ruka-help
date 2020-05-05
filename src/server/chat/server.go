package main

import (
	"log"
	"strings"
	"sync"

	nats "github.com/nats-io/nats.go"
	"github.com/urfave/cli"
)

// User contains the data of each connected user.
type User struct {
	ID       string
	Username string
}

type ChatMessage struct {
	ID      int
	User    *User
	Content string
	Type    string
}

type Session struct {
	ID       string
	Messages []*ChatMessage
	Users    map[string]*User
}

// GeneralMessage contains the data of each messaged shared in the session.
type GeneralMessage struct {
	User    *User
	Content string
	Type    string
}

func NewSession(ID string) *Session {
	s := Session{
		ID:       ID,
		Messages: []*ChatMessage{},
		Users:    make(map[string]*User),
	}
	return &s
}

var sessions = make(map[string]*Session)

var encodedNatsConnection *nats.EncodedConn

const (
	userLeaveChannel  = "session.*.chat.user.*.leave"
	newUserChannel    = "session.*.chat.user.*.new"
	inGeneralChannel  = "session.*.chat.in"
	outGeneralChannel = "session.*.chat.out"
	outUserChannel    = "session.*.{UserID}.out"
)

// StartListener start
func StartListener(c *cli.Context) error {
	log.Printf("Connecting to server : %s", nats.DefaultURL)

	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}

	encodedNatsConnection, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()
	defer encodedNatsConnection.Close()

	// Use a WaitGroup to wait for a message to arrive
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Simple Async Subscriber
	encodedNatsConnection.Subscribe(newUserChannel, handleNewUser)

	// Simple Async Subscriber
	encodedNatsConnection.Subscribe(userLeaveChannel, handleUserLeaving)

	// Simple Async Subscriber
	encodedNatsConnection.Subscribe(inGeneralChannel, handleNewMessage)

	// Wait for a message to come in
	wg.Wait()
	return nil
}

func handleNewUser(subj, reply string, m *GeneralMessage) {
	log.Printf("[1] Received a message from %s\n", string(subj))

	sessionID := strings.Split(subj, ".")[1]
	session, sessionExists := sessions[sessionID]
	if !sessionExists {
		session = NewSession(sessionID)
		sessions[sessionID] = session
		log.Printf("The session [%s] was created.", sessionID)
	}

	_, userExists := session.Users[m.User.ID]
	if !userExists {
		session.Users[m.User.ID] = &User{
			ID:       m.User.ID,
			Username: m.User.Username,
		}

		newMessage := &ChatMessage{
			Content: "User " + m.User.Username + " has entered the workspace.",
			Type:    "system",
		}
		session.Messages = append(session.Messages, newMessage)
		outChannel := strings.Replace(outGeneralChannel, "*", sessionID, 1)
		log.Printf("Sending system message to session [%s] using channel %s.", sessionID, outChannel)
		encodedNatsConnection.Publish(outChannel, newMessage)
	}
}

func handleUserLeaving(subj, reply string, m *GeneralMessage) {
	log.Printf("[2] Received a message from %s\n", string(subj))

	sessionID := strings.Split(subj, ".")[1]
	session, sessionExists := sessions[sessionID]
	if !sessionExists {
		log.Printf("The session [%s] doesn't exists.", sessionID)
		return
	}

	log.Println(len(session.Users))
	user, userExists := session.Users[m.User.ID]
	if !userExists {
		log.Printf("User [%s] is not registered in session [%s].", m.User.ID, sessionID)
		return
	}

	delete(session.Users, user.ID)

	log.Printf("User [%s] has leave the session [%s].", user.ID, sessionID)

	newMessage := &ChatMessage{
		Content: "User " + m.User.Username + " has leave the workspace.",
		Type:    "system",
	}
	session.Messages = append(session.Messages, newMessage)
	outChannel := strings.Replace(outGeneralChannel, "*", sessionID, 1)
	log.Printf("Sending system message to session [%s] using channel [%s].", sessionID, outChannel)
	encodedNatsConnection.Publish(outChannel, newMessage)
}

func handleNewMessage(subj, reply string, m *ChatMessage) {
	log.Printf("[3] Received a message from %s\n", string(subj))

	sessionID := strings.Split(subj, ".")[1]
	session, sessionExists := sessions[sessionID]
	if !sessionExists {
		log.Printf("The session [%s] is not registered, can't handle message.", sessionID)
		return
	}

	user, userExists := session.Users[m.User.ID]
	if !userExists {
		log.Printf("The user [%s] is not registered in the session [%s].", m.User.ID, sessionID)
		return
	}

	newMessage := &ChatMessage{
		ID:      len(session.Messages),
		User:    user,
		Content: m.Content,
		Type:    "message",
	}
	session.Messages = append(session.Messages, newMessage)
	outChannel := strings.Replace(outGeneralChannel, "*", sessionID, 1)
	encodedNatsConnection.Publish(outChannel, newMessage)
}
