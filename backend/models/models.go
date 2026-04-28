package models

// user
type User struct {
	ID                string `json:"id"`
	Nickname          string `json:"nickname"`
	Birthday          string `json:"birthday"`
	Gender            string `json:"gender"`
	Firstname         string `json:"firstname"`
	Lastname          string `json:"lastname"`
	AboutMe           string `json:"aboutme"`
	Email             string `json:"email"`
	Password          string
	ProfileImage      string `json:"profile_image"`
	IsPublic          bool   `json:"is_public"`
	IsFreind          bool   `json:"is_freind"`
	InteractionStatus string `json:"interaction_status"`

	SessionID      string `json:"session_id"`
	SessionExpired string
	SessionCreated string
}

// posts
type Post struct {
	ID             string
	AutherName     string
	UserID         string
	Content        string
	CategoryType   string
	CreatedAt      string
	NbrOfComments  int
	NbrOfLikes     int
	NbrOfDislikes  int
	NbrOfReactions int
	UserReaction   int
	Comments       []Comment
	ImageURL       string
}

// comments
type Comment struct {
	ID             string
	AutherName     string
	Content        string
	UserID         string
	PostID         string
	CreatedAt      string
	NbrOfReactions int
	UserReaction   int
	Offset         int
}

// messages
type Message struct {
	ID           int
	SenderID     string
	SenderName   string
	ReceiverName string
	ReceiverID   string
	Content      string
	IsNotRead    int
	CreatedAt    string

	SenderNickname   string
	ReceiverNickname string
}

// include (psot or comment ID), (PostOrComment: "POST" or "COMMENT"), (type : 0 -> 6) ...etc
type Reaction struct {
	PostorcommentID string
	PostOrComment   string
	Type            int
	UserID          int
	CreatedAt       string
}

// this is for the ws
type UserInfo struct {
	ID                     string
	Nickname               string
	Firstname              string
	Lastname               string
	LastMessage            Message
	NumberOfUnreadMessages int
	IsOnline               bool
}
