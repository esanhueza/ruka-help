package main

// ChatMessage contains the data of eachessaged shared in the session.
type ChatMessage struct {
	ID      int64
	User    *User
	Content string
	Type    string
}

// SessionMessage2 contains the data of each messaged shared in the session.
type SessionMessage2 struct {
	ID      int64
	User    *User
	Content string
	Type    string
}

// SessionMessage contains the data of each messaged shared in the session.
type SessionMessage struct {
	ID string
}

// GeneralMessage contains the data of each messaged shared in the session.
type GeneralMessage struct {
	User    *User
	Content string
	Type    string
}

type ClientMessage struct {
	SessionID string `json:"SessionID"`
	UserID    string `json:"UserID"`
	Content   string `json:"Content"`
	Type      string `json:"Type"`
	Token     string `json:"Token"`
}
