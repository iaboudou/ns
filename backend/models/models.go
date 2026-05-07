package models

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

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
	ID           string `json:"id"`
	Type         string `json:"type"`
	SenderID     string `json:"sender_id"`
	SenderName   string `json:"sender_name"`
	ReceiverName string `json:"receiver_name"`
	ReceiverID   string `json:"receiver_id"`
	Content      string `json:"content"`
	CreatedAt    int64  `json:"created_at"`

	PortKey      string `json:"portKey"`
	LastReadID   string `json:"last_read_Id"`
	LastReadTime int64  `json:"last_read_time"`
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
	Image              string `json:"image"`
}

type Group struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type Event struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Vote        string `json:"vote"`
}

// represent user in ws connection
type Client struct {
	ID          string
	Ws          *websocket.Conn
	Mu          *sync.Mutex
	Description *User
}

// represent the websocket hub
type Hub struct {
	Connect    chan *Client
	Disconnect chan *Client
	Broadcast  chan Message
	Notify     chan FollowNotif
}

// FollowNotif is sent through Hub.Notify when a follow event occurs
type FollowNotif struct {
	FromUser  User
	ToUserID  string
	NotifType string // "follow_request" or "follow_accepted"
	GroupID   string // for groups notif
}
