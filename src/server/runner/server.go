package main

import (
	"log"
	"strings"
	"sync"

	nats "github.com/nats-io/nats.go"
	"github.com/teris-io/shortid"
	"github.com/urfave/cli"
)

// User contains the data of each connected user.
type User struct {
	ID       string
	Username string
}

type Workspace struct {
	ID        string
	SessionID string
	Users     map[string]*User
	Code      string
	Schema    string
}

// GeneralMessage contains the data of each messaged shared in the session.
type GeneralMessage struct {
	User    *User
	Content string
	Type    string
	Code    string
	Schema  string
}

var workspaceIDGenerator, err = shortid.New(1, shortid.DefaultABC, 2342)

func NewWorkspace(SessionID string) (*Workspace, error) {
	ID, err := workspaceIDGenerator.Generate()
	if err != nil {
		return nil, err
	}
	w := Workspace{
		ID:        ID,
		SessionID: SessionID,
		Users:     make(map[string]*User),
		Code:      "",
		Schema:    "",
	}
	return &w, nil
}

var sessions = make(map[string]*Workspace)

var encodedNatsConnection *nats.EncodedConn

const (
	userLeaveChannel    = "session.*.workspace.user.*.leave"
	userEnterChannel    = "session.*.workspace.user.*.new"
	workspaceInChannel  = "session.*.workspace.in"
	workspaceOutChannel = "session.*.workspace.out"
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
	encodedNatsConnection.Subscribe(userEnterChannel, handleNewUser)

	// Simple Async Subscriber
	encodedNatsConnection.Subscribe(userLeaveChannel, handleUserLeaving)

	// Simple Async Subscriber
	encodedNatsConnection.Subscribe(workspaceInChannel, handleNewMessage)

	// Wait for a message to come in
	wg.Wait()
	return nil
}

func handleNewUser(subj, reply string, m *GeneralMessage) {
	log.Printf("[1] Received a message from %s\n", string(subj))

	sessionID := strings.Split(subj, ".")[1]
	workspace, sessionExists := sessions[sessionID]
	if !sessionExists {
		workspace, err = NewWorkspace(sessionID)
		if err != nil {
			log.Printf("Can't create workspace for session [%s].", sessionID)
			return
		}
		sessions[sessionID] = workspace
		log.Printf("The workspace [%s] was created for session [%s].", workspace.ID, workspace.SessionID)
	}

	_, userExists := workspace.Users[m.User.ID]
	if !userExists {
		workspace.Users[m.User.ID] = &User{
			ID:       m.User.ID,
			Username: m.User.Username,
		}

		newMessage := &GeneralMessage{
			Content: "Workspace [" + workspace.ID + "] created for session [" + workspace.SessionID + "].",
			Type:    "system",
		}
		outChannel := strings.Replace(workspaceOutChannel, "*", sessionID, 1)
		log.Printf("Sending system message to session [%s] using channel %s.", sessionID, outChannel)
		encodedNatsConnection.Publish(outChannel, newMessage)
	}
}

func handleUserLeaving(subj, reply string, m *GeneralMessage) {
	log.Printf("[2] Received a message from %s\n", string(subj))

	sessionID := strings.Split(subj, ".")[1]
	workspace, sessionExists := sessions[sessionID]
	if !sessionExists {
		log.Printf("The workspace for session [%s] doesn't exists.", sessionID)
		return
	}

	user, userExists := workspace.Users[m.User.ID]
	if !userExists {
		log.Printf("User [%s] is not registered in the workspace for session [%s].", m.User.ID, sessionID)
		return
	}

	delete(workspace.Users, user.ID)

	log.Printf("User [%s] has leave the workspace [%s].", user.ID, workspace.ID)

	newMessage := &GeneralMessage{
		Content: "User " + m.User.Username + " has leave the workspace [" + workspace.ID + "].",
		Type:    "system",
	}
	outChannel := strings.Replace(workspaceOutChannel, "*", sessionID, 1)
	log.Printf("Sending system message to session [%s] using channel [%s].", sessionID, outChannel)
	encodedNatsConnection.Publish(outChannel, newMessage)
}

func handleNewMessage(subj, reply string, m *GeneralMessage) {
	log.Printf("[3] Received a message from %s\n", string(subj))

	sessionID := strings.Split(subj, ".")[1]
	session, sessionExists := sessions[sessionID]
	if !sessionExists {
		log.Printf("There is not a workspace registered for session [%s], can't handle update.", sessionID)
		return
	}

	user, userExists := session.Users[m.User.ID]
	if !userExists {
		log.Printf("The user [%s] is not registered in the workspace for session [%s].", m.User.ID, sessionID)
		return
	}

	newMessage := &GeneralMessage{
		User:   user,
		Code:   "",
		Schema: "",
		Type:   "workspace",
	}
	outChannel := strings.Replace(workspaceOutChannel, "*", sessionID, 1)
	encodedNatsConnection.Publish(outChannel, newMessage)
}
