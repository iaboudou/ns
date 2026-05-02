package models

import "time"

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
	ID               string `json:"id"`
	Nickname         string `json:"nickname"`
	Firstname        string `json:"firstname"`
	Lastname         string `json:"lastname"`
	UserImageProfile string `json:"profile_image"`

	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ImageURL  string    `json:"image_url"`
	Privacy   string    `json:"privacy"`

	NumberOfComments int `json:"number_of_comments"`

	AllowedUsers       []string `json:"allowed_users"`
	Alloweduserscreate string
	GroupID            string    `json:"group_id"`
	Comments           []Comment `json:"comments"`
}

// comments
type Comment struct {
	ID               string `json:"id"`
	AutherName       string `json:"auther_name"`
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	Content          string `json:"content"`
	UserID           string `json:"user_id"`
	PostID           string `json:"post_id"`
	CreatedAt        string `json:"created_at"`
	ImageURL         string `json:"image_url"`
	UserImageProfile string `json:"profile_image"`
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

// group struct
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// followers struct
type FollowSuggestion struct {
	Id             string `json:"id"`
	ProfileImage   string `json:"profile_image"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	AccountPrivacy bool   `json:"account_privacy"`
}

type GroupeInfo struct {
	Title              string `json:"title"`
	Description        string `json:"description"`
	IsCreator          bool   `json:"isCreator"`
	MemberAmount       int    `json:"members"`
	Posts              []Post `json:"posts"`
	UnreadMessageCount int    `json:"unreadMsg"`
}


type Group struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Event struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type EventResponse struct {
	ID      string
	EventID string
	UserID  string
	Status  string
}
