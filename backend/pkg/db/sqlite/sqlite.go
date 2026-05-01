package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"rtf/help"
	"rtf/models"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)


func InitDB() (*sql.DB, error) {
	m, err := migrate.New("file://pkg/db/migrations/sqlite", "sqlite3://db/db.db")
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "./db/db.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Repo struct {
	Db *sql.DB
}

// try to insert the user into the data base, any invalid input will return an error with a specific message
func (r *Repo) InsertUserDB(user help.U) error {
	hashed, err := help.HashPassword(user.Password)
	if err != nil {
		return errors.New("SERVER ERROR")
	}

	userUUID, err := uuid.NewV4()
	if err != nil {
		return errors.New("SERVER ERROR")
	}

	_, err = r.Db.Exec(
		`INSERT INTO users(
			id, 
			nickname, 
			firstname, 
			lastname, 
			email, 
			password, 
			birthday, 
			gender, 
			profile_image, 
			about_me, 
			account_privacy
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userUUID.String(), user.Nickname, user.Firstname, user.Lastname, user.Email, hashed, user.Birthday, user.Gender, user.Avatar, user.About, 1,
	)
	if err != nil {
		return errors.New("SERVER ERROR")
	}

	return nil
}

func (r *Repo) IsUserAlreadyExist(user *help.U) error {
	var exist int
	err := r.Db.QueryRow("SELECT 1 FROM users WHERE firstname=? OR lastname=? OR email=?", user.Firstname, user.Lastname, user.Email).Scan(&exist)
	if err != nil && err != sql.ErrNoRows {
		return errors.New("SERVER ERROR")
	}
	if exist > 0 {
		return errors.New("user already exist")
	}
	return nil
}

// check existance of the user in the DB
func (r *Repo) IsUserExist(email, password string) (string, error) {
	var id string
	var hashedPassword string

	if len(email) > 60 || len(password) > 60 {
		return "", errors.New("user not exist")
	}

	err := r.Db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not exist")
		}
		return "", errors.New("SERVER ERROR")
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return "", errors.New("password invalid")
	}

	return id, nil
}

// set new session in case of user login
func (r *Repo) SetUserSession(w http.ResponseWriter, userID string) (string, time.Time, error) {
	_, err := r.Db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return "", time.Time{}, errors.New("SERVER ERROR")
	}

	sessionUUID, err := uuid.NewV7()
	if err != nil {
		return "", time.Time{}, errors.New("SERVER ERROR")
	}
	sessionID := sessionUUID.String()

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)

	_, err = r.Db.Exec(`
	INSERT INTO sessions (id, user_id, token, expires_at)
	VALUES (?, ?, ?, ?)`,
		sessionID,
		userID,
		sessionID,
		expiredAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return "", time.Time{}, errors.New("SERVER ERROR")
	}
	return sessionID, expiredAt, nil
}

// delete the session from the DB in case of logout
func (r *Repo) DisconnectUser(userID string) error {
	_, er := r.Db.Exec("UPDATE users SET session_id=NULL, session_created_at=NULL, session_expired_at=NULL WHERE id=?", userID)
	if er != nil {
		return errors.New("SERVER ERROR")
	}
	return nil
}

// this check session sended by the browser if it is included in the DB
func (r *Repo) CheckSessionExistance(req *http.Request) (models.User, error) {
	var user models.User
	var expiresAtStr string

	cookie, err := req.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return user, errors.New("no session cookie")
	}
	user.SessionID = cookie.Value
	//
	var userID string
	err = r.Db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", cookie.Value).
		Scan(&userID, &expiresAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("session not found")
		}
		return user, errors.New("server error")
	}

	//
	expiresAt, err := time.Parse("2006-01-02 15:04:05", expiresAtStr)
	if err != nil {
		return user, errors.New("invalid session format")
	}
	if time.Now().After(expiresAt) {
		return user, errors.New("session expired")
	}

	//
	err = r.Db.QueryRow(`SELECT 
						id, 
						nickname, 
						firstname, 
						lastname, 
						email, 
						birthday, 
						profile_image, 
						about_me, 
						account_privacy 
					FROM users WHERE id = ?`, userID).
		Scan(&user.ID, &user.Nickname, &user.Firstname, &user.Lastname, &user.Email, &user.Birthday, &user.ProfileImage, &user.AboutMe, &user.IsPublic)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, errors.New("server error")
	}

	return user, nil
}

func (r *Repo) InsertPostDB(userID string, post models.Post) (models.Post, error) {
	postUUID, err := uuid.NewV4()
	if err != nil {
		return post, errors.New("SERVER ERROR")
	}

	id := postUUID.String()
	now := time.Now()

	groupID := post.GroupID
	if groupID == "" {
		groupID = ""
	}
	_, err = r.Db.Exec(
		`INSERT INTO posts (
			id, user_id, content, image_url, created_at, privacy, allowed_users, group_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		userID,
		post.Content,
		post.ImageURL,
		now.Format("2006-01-02 15:04:05.000000"),
		post.Privacy,
		post.Alloweduserscreate,
		groupID,
	)
	if err != nil {
		return post, errors.New("SERVER ERROR")
	}

	post.ID = id
	post.UserID = userID
	post.CreatedAt = now

	er := r.Db.QueryRow(`SELECT profile_image FROM users WHERE id = ?`, userID).Scan(&post.UserImageProfile)
	if er != nil {
		return post, errors.New("SERVER ERROR")
	}

	return post, nil
}

// this check if the categories exists in DB,  it return the categories id & bool, bool type is false in case of an element not exist, true otherwise
func (r *Repo) IsCategoryCorrect(category string) ([]int, error) {
	if len(category) == 0 {
		return nil, errors.New("post category is required")
	}

	categories := strings.Split(category, ",")
	seen := make(map[string]bool)
	var ids []int

	for _, c := range categories {
		c = strings.TrimSpace(c)

		if seen[c] {
			return nil, errors.New("duplicate category")
		}
		seen[c] = true

		var id int
		err := r.Db.QueryRow(
			"SELECT id FROM categories WHERE category_name = ?",
			c,
		).Scan(&id)
		if err != nil {
			return nil, errors.New("invalid category")
		}

		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, errors.New("invalid category")
	}

	return ids, nil
}

// insert comment into the DB
func (r *Repo) InsertCommentDB(comment models.Comment) (models.Comment, error) {
	commentUUID, err := uuid.NewV4()
	if err != nil {
		return models.Comment{}, errors.New("SERVER ERROR")
	}

	id := commentUUID.String()
	now := time.Now().Format("2006-01-02 15:04:05.000000")

	err = r.Db.QueryRow("SELECT nickname, firstname, lastname, profile_image FROM users WHERE id = ?", comment.UserID).
		Scan(&comment.AutherName, &comment.FirstName, &comment.LastName, &comment.UserImageProfile)

	_, err = r.Db.Exec(`
		INSERT INTO comments (id, post_id, user_id, content, image_url, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		id, comment.PostID, comment.UserID, comment.Content, comment.ImageURL, now)
	if err != nil {
		return models.Comment{}, errors.New("SERVER ERROR")
	}

	comment.ID = id
	comment.CreatedAt = now

	return comment, nil
}

// to insert a comment to the DB need to check if this post already exist in the DB
func (r *Repo) PostExists(postID string) error {
	var id string
	err := r.Db.QueryRow("SELECT id FROM posts WHERE id = ?", postID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("post not exist")
		}
		return errors.New("SERVER ERROR")
	}
	return nil
}

// get user infos from DB
func (r *Repo) GetPeronalInfoFromDB(userID string) (models.User, error) {
	user := models.User{}
	er := r.Db.QueryRow(`
	SELECT nickname, firstname, lastname, email, birthday, gender, profile_image, about_me, account_privacy
	FROM users WHERE id = ?
	`, userID).Scan(&user.Nickname, &user.Firstname, &user.Lastname, &user.Email, &user.Birthday, &user.Gender, &user.ProfileImage, &user.AboutMe, &user.IsPublic)
	if er != nil {
		return models.User{}, er
	}
	user.ID = userID

	return user, nil
}

// switch account privacy between 0 and 1
func (r *Repo) SwitchAccountPrivacyinDB(userID string) error {
	_, er := r.Db.Exec(`
		UPDATE users
		set account_privacy = 1 - account_privacy
		WHERE id=?`, userID)
	if er != nil {
		return er
	}

	return nil
}

// this function get 10 posts from DB with its comments
func (r *Repo) Get10PostsfromDB(Page, Section, ViewerID, ReqUserID, GroupID string, Offset int) ([]models.Post, error) {
	posts := []models.Post{}
	var rows *sql.Rows
	var err error

	switch Page {
	case "profile-me-posts", "profille-other-posts":
		rows, err = r.Db.Query(`
			SELECT id, user_id, content, image_url, created_at, privacy, allowed_users, group_id
			FROM posts
			WHERE user_id = ?
			ORDER BY created_at DESC
			LIMIT 10 OFFSET ?`, ReqUserID, Offset)

	case "profile-me-activity":
		rows, err = r.Db.Query(`
			SELECT p.id, p.user_id, p.content, p.image_url, p.created_at, p.privacy, p.allowed_users, p.group_id
			FROM posts p
			WHERE p.id IN (
				SELECT post_or_comm_id
				FROM reactions
				WHERE user_id = ?
				AND post_or_comm = 'POST'

				UNION

				SELECT post_id
				FROM comments
				WHERE user_id = ?
			)
			ORDER BY p.created_at DESC
			LIMIT 10 OFFSET ?`, ReqUserID, ReqUserID, Offset)

	case "goups", "groups":
		rows, err = r.Db.Query(`
			SELECT id, user_id, content, image_url, created_at, privacy, allowed_users, group_id
			FROM posts
			WHERE group_id = ?
			ORDER BY created_at DESC
			LIMIT 10 OFFSET ?`, GroupID, Offset)

	default:
		rows, err = r.Db.Query(`
			SELECT id, user_id, content, image_url, created_at, privacy, allowed_users, group_id
			FROM posts
			WHERE group_id = ""
			ORDER BY created_at DESC
			LIMIT 10 OFFSET ?`, Offset)
	}

	if err != nil {
		return nil, errors.New("SERVER ERROR")
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		var createdAtStr string
		var allowedUsersStr string

		err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.ImageURL, &createdAtStr, &p.Privacy, &allowedUsersStr, &p.GroupID)
		if err != nil {
			return nil, errors.New("SERVER ERROR")
		}

		p.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAtStr)
		if err != nil {
			return nil, errors.New("SERVER ERROR")
		}

		p.AllowedUsers = help.ParseAllowedUsers(allowedUsersStr)

		if err := r.Db.QueryRow(
			`SELECT nickname, firstname, lastname, profile_image FROM users WHERE id = ?`,
			p.UserID,
		).Scan(&p.Nickname, &p.Firstname, &p.Lastname, &p.UserImageProfile); err != nil {
			return nil, errors.New("SERVER ERROR")
		}

		show := false
		if ViewerID != p.UserID {
			switch p.Privacy {
			case "public":
				show = true
			case "followers", "almost_private":
				show, _ = r.IsFollower(ViewerID, p.UserID)
			case "private":
				show = slices.Contains(p.AllowedUsers, ViewerID)
			}
			if !show {
				continue
			}
		}

		p.NumberOfComments, err = r.GetTotalComments(p.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, rows.Err()
}

func (r *Repo) IstheUserFreind(user *models.User, mainuserID string) error {
	if user.ID == mainuserID {
		user.IsFreind = true
		user.InteractionStatus = "none"
		return nil
	}

	//
	var status string
	err := r.Db.QueryRow(` SELECT status FROM followers WHERE follower_id=? AND following_id=?`, mainuserID, user.ID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			user.InteractionStatus = "none"
			user.IsFreind = false
		} else {
			return err
		}
	} else {
		if status == "accepted" {
			user.InteractionStatus = "following"
			user.IsFreind = true
		} else if status == "pending" {
			user.InteractionStatus = "requested"
			user.IsFreind = false
		} else {
			user.InteractionStatus = "none"
			user.IsFreind = false
		}
	}

	return nil
}

func (r *Repo) GetTotalComments(PostID string) (int, error) {
	total := 0
	er := r.Db.QueryRow(`SELECT COUNT(*) FROM comments WHERE post_id = ?`, PostID).Scan(&total)
	return total, er
}

func (r *Repo) IsFollower(userID, otherID string) (bool, error) {
	var status string
	err := r.Db.QueryRow(`
		SELECT status 
		FROM followers 
		WHERE follower_id = ? AND following_id = ?`, userID, otherID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	//
	return status == "accepted", nil
}

func (r *Repo) GetUserAuthernameByID(userID string) (string, error) {
	var nickname string
	er := r.Db.QueryRow("SELECT nickname FROM users WHERE id=?", userID).Scan(&nickname)
	return nickname, er
}

// this function get the comments post from DB based on an postID
func (r *Repo) Get10PostComments(postID string, offset int) ([]models.Comment, error) {
	comments := []models.Comment{}

	rows, err := r.Db.Query(`SELECT * FROM comments WHERE post_id = ? ORDER BY created_at DESC LIMIT 10 OFFSET ?`, postID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		er := rows.Scan(&c.ID, &c.Content, &c.UserID, &c.PostID, &c.ImageURL, &c.CreatedAt)
		if er != nil {
			return nil, nil
		}

		er = r.Db.QueryRow("SELECT nickname, firstname, lastname, profile_image FROM users WHERE id = ?", c.UserID).
			Scan(&c.AutherName, &c.FirstName, &c.LastName, &c.UserImageProfile)
		comments = append(comments, c)
	}

	return comments, nil
}

// this func get the categories related to the post from DB
func (r *Repo) GetPostCategory(postID string) (string, error) {
	var categoryIDs string

	// get "1,2,3"
	err := r.Db.QueryRow(`SELECT category_ids FROM posts WHERE id = ?`, postID).Scan(&categoryIDs)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("post not exists")
		}
		return "", errors.New("SERVER ERROR")
	}

	ids := strings.Split(categoryIDs, ",")
	categories := []string{}

	// get categories name
	for _, id := range ids {
		id = strings.TrimSpace(id)
		var name string
		err := r.Db.QueryRow(`SELECT category_name FROM categories WHERE id = ?`, id).Scan(&name)
		if err != nil {
			return "", errors.New("invalid category")
		}
		categories = append(categories, name)
	}

	res := strings.Join(categories, ",")

	return res, nil
}

// get all users exists in the DB
func (r *Repo) Get100UsersFor(userID string, startID int) ([]models.User, error) {
	var users []models.User

	rows, err := r.Db.Query(`SELECT id, nickname, birthday, gender, firstname, lastname, email FROM users LIMIT 100 OFFSET ?`, startID)
	if err != nil {
		return nil, errors.New("SERVER ERROR")
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Nickname, &u.Birthday, &u.Gender, &u.Firstname, &u.Lastname, &u.Email)
		if err != nil {
			return nil, errors.New("SERVER ERROR")
		}
		if u.ID != userID {
			users = append(users, u)
		}
	}
	if len(users) == 0 {
		return nil, errors.New("startID reached the max")
	}
	return users, nil
}

// get the users info with the last message for the message list
func (r *Repo) GetUsersInfoFor(userID string, forAll bool) ([]models.UserInfo, error) {
	usersInfo := []models.UserInfo{}

	rows, err := r.Db.Query(`SELECT id, nickname, firstname, lastname FROM users`)
	if err != nil {
		return nil, errors.New("SERVER ERROR")
	}
	defer rows.Close()

	for rows.Next() {
		var u models.UserInfo
		err := rows.Scan(&u.ID, &u.Nickname, &u.Firstname, &u.Lastname)
		if err != nil {
			return nil, errors.New("SERVER ERROR")
		}

		if u.ID == userID && !forAll {
			continue
		}
		// get the last message between them
		var msg models.Message
		err = r.Db.QueryRow(`
            SELECT id, sender_id, receiver_id, content, is_NOT_read, created_at 
            FROM messages 
            WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
            ORDER BY created_at DESC 
            LIMIT 1
        `, userID, u.ID, u.ID, userID).Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.IsNotRead, &msg.CreatedAt)

		if err != nil && err != sql.ErrNoRows {
			return nil, errors.New("SERVER ERROR")
		}
		u.LastMessage = msg

		// calcul the number of messages not read from the other user
		var count int
		err = r.Db.QueryRow(` SELECT COUNT(*) FROM messages WHERE sender_id = ? AND receiver_id = ? AND is_NOT_read = 1`, u.ID, userID).Scan(&count)
		if err != nil {
			return nil, errors.New("SERVER ERROR")
		}
		u.NumberOfUnreadMessages = count

		usersInfo = append(usersInfo, u)

	}

	if len(usersInfo) == 0 {
		return nil, errors.New("no user exists")
	}

	return usersInfo, nil
}

// get the user id based on the front id from DB
func (r *Repo) GetUserByFrontID(frontID string) (string, error) {
	var id string
	er := r.Db.QueryRow("SELECT id FROM users WHERE front_id=?", frontID).Scan(&id)
	return id, er
}

// insert new message to the DB
func (r *Repo) InsertMessage(msg map[string]interface{}) (models.Message, error) {
	content, ok1 := msg["Content"].(string)
	senderID, ok2 := msg["SenderID"].(string)
	receiverID, ok3 := msg["ReceiverID"].(string)
	if !ok1 || !ok2 || !ok3 {
		return models.Message{}, errors.New("invalid message format")
	}

	m := models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05.000000"),
	}

	res, er := r.Db.Exec("INSERT INTO messages(sender_id, receiver_id, content, is_NOT_read, created_at) VALUES (?, ?, ?, ?, ?)", m.SenderID, m.ReceiverID, m.Content, 1, m.CreatedAt)
	if er != nil {
		return models.Message{}, er
	}

	id, _ := res.LastInsertId()
	m.ID = int(id)
	m.IsNotRead = 1

	m.SenderName, er = r.GetUserAuthernameByID(senderID)
	if er != nil {
		return models.Message{}, er
	}
	m.ReceiverName, er = r.GetUserAuthernameByID(receiverID)
	if er != nil {
		return models.Message{}, er
	}

	return m, nil
}

func (r *Repo) SetMessageRead(senderID, receiverID string) error {
	_, er := r.Db.Exec(` UPDATE messages SET is_NOT_read = 0 WHERE sender_id = ? AND receiver_id = ? AND is_NOT_read = 1`, senderID, receiverID)
	if er != nil {
		return er
	}
	return nil
}

// get user infos from DB
func (r *Repo) GetUserInfos(userID string) (help.U, error) {
	user := help.U{}
	er := r.Db.QueryRow(`
	SELECT 
		nickname, 
		firstname, 
		lastname, 
		email, 
		birthday, 
		gender, 
		profile_image, 
		about_me, 
		account_privacy
	FROM users WHERE id = ?
	`, userID).Scan(&user.Nickname, &user.Firstname, &user.Lastname, &user.Email, &user.Birthday, &user.Gender, &user.Avatar, &user.About, &user.AccountPrivacy)
	if er != nil {
		return help.U{}, er
	}
	fmt.Println(user.Avatar)
	user.ID = userID

	return user, nil
}

// this method get the users with its last message in order from DB for chat
func (r *Repo) GetUsersForChatInOrder(userID string, offset int) ([]models.Message, error) {
	rows, er := r.Db.Query(`
		SELECT sender_id, receiver_id, content, is_NOT_read, created_at FROM messages
		WHERE sender_id = ? OR receiver_id = ?
		ORDER BY created_at DESC`, userID, userID)
	if er != nil {
		return nil, er
	}
	defer rows.Close()

	msgs := []models.Message{}
	seen := make(map[string]bool)
	i := 0

	for rows.Next() {
		m := models.Message{}
		if er := rows.Scan(&m.SenderID, &m.ReceiverID, &m.Content, &m.IsNotRead, &m.CreatedAt); er != nil {
			continue
		}

		if m.SenderID != userID {
			// in case the user is the receiver
			if seen[m.SenderID] {
				continue
			}
			seen[m.SenderID] = true

			if i < offset {
				i++
				continue
			}

			if er := r.GetUserNickNameByID(&m.SenderNickname, m.SenderID); er != nil {
				return nil, errors.New("SERVER ERROR")
			}

		} else {
			// in case the user is the sender
			if seen[m.ReceiverID] {
				continue
			}
			seen[m.ReceiverID] = true

			if i < offset {
				i++
				continue
			}

			if er := r.GetUserNickNameByID(&m.ReceiverNickname, m.ReceiverID); er != nil {
				return nil, errors.New("SERVER ERROR")
			}
		}

		msgs = append(msgs, m)
		if len(msgs) >= 100 {
			break
		}
	}

	return msgs, nil
}

// get the user nickname by its ID from DB
func (r *Repo) GetUserNickNameByID(nickname *string, userID string) error {
	return r.Db.QueryRow("SELECT nickname FROM users WHERE id=?", userID).Scan(nickname)
}

// this method get the messages history between two users from DB in order from the offset provided at ascending order
func (r *Repo) GetMessagesHistoryBetweenTwoUsers(senderID, receiverID string, offset int) ([]models.Message, error) {
	rows, er := r.Db.Query(`
        SELECT sender_id, receiver_id, content, created_at 
        FROM messages
        WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at DESC
        LIMIT 10 OFFSET ?`,
		senderID, receiverID, receiverID, senderID, offset)
	if er != nil {
		return nil, er
	}
	defer rows.Close()

	msgs := []models.Message{}
	for rows.Next() {
		m := models.Message{}
		er := rows.Scan(&m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt)
		if er != nil {
			continue
		}
		m.ReceiverName, er = r.GetUserAuthernameByID(m.ReceiverID)
		if er != nil {
			continue
		}
		m.SenderName, er = r.GetUserAuthernameByID(m.SenderID)
		if er != nil {
			continue
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// =====================
// followers
// =====================

// get users that the (userID) does not follow them yet from db
func (r *Repo) GetSuggestionUsersDB(userID string) ([]models.FollowSuggestion, error) {
	suggestions := []models.FollowSuggestion{}

	rows, er := r.Db.Query(`
		SELECT id, firstname, lastname, profile_image, account_privacy 
		FROM users 
		WHERE id != ? 
		AND id NOT IN (
			SELECT following_id FROM followers WHERE follower_id = ?
		)`, userID, userID)
	if er != nil {
		return nil, er
	}
	defer rows.Close()

	for rows.Next() {
		var follow models.FollowSuggestion
		er := rows.Scan(&follow.Id, &follow.Firstname, &follow.Lastname, &follow.ProfileImage, &follow.AccountPrivacy)
		follow.AccountPrivacy = !follow.AccountPrivacy
		if er != nil {
			return nil, er
		}
		suggestions = append(suggestions, follow)
	}

	return suggestions, nil
}

// this function is to follow the user in db and req if his account is privet
func (r *Repo) FollowUserDB(userID string, followedID string) (string, error) {
	// Check if follow already exists
	var status string
	err := r.Db.QueryRow(`SELECT status FROM followers WHERE follower_id = ? AND following_id = ?`, userID, followedID).Scan(&status)
	if err == nil {
		// If it exists, delete it
		_, err = r.Db.Exec(`DELETE FROM followers WHERE follower_id = ? AND following_id = ?`, userID, followedID)
		if err != nil {
			return "", err
		}
		if status == "pending" {
			return "follow request deleted", nil
		}
		return "follow deleted", nil
	}

	// If it doesn't exist, check if the targeted account is private
	var isPrivate bool
	err = r.Db.QueryRow(`SELECT account_privacy FROM users WHERE id = ?`, followedID).Scan(&isPrivate)
	if err != nil {
		return "", err
	}
	isPrivate = !isPrivate
	// Create the follow
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	finalStatus := "accepted"
	message := "follow have been successfully"
	if isPrivate {
		finalStatus = "pending"
		message = "request have been sent"
	}

	_, err = r.Db.Exec(`
	INSERT INTO followers (id, follower_id, following_id, status)
	VALUES (?, ?, ?, ?)
	`, id.String(), userID, followedID, finalStatus)
	if err != nil {
		return "", err
	}

	return message, nil
}

// this function is to answer to the request wether you accept or now if your profile is privet
func (r *Repo) ManageFollowDB(followerID string, followingID string, decision string) error {
	if decision == "accepted" {
		_, err := r.Db.Exec(`UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND following_id = ?`, followerID, followingID)
		return err
	} else if decision == "rejected" {
		_, err := r.Db.Exec(`DELETE FROM followers WHERE follower_id = ? AND following_id = ?`, followerID, followingID)
		return err
	}
	return nil
}

// this function is to get the all the followers
func (r *Repo) GetFollowersDB(targetID string) ([]models.FollowSuggestion, error) {
	followers := []models.FollowSuggestion{}

	rows, er := r.Db.Query(`
		SELECT u.id, u.firstname, u.lastname, u.profile_image, u.account_privacy 
		FROM users u
		JOIN followers f ON u.id = f.follower_id
		WHERE f.following_id = ? AND f.status = 'accepted'
	`, targetID)
	if er != nil {
		return nil, er
	}
	defer rows.Close()

	for rows.Next() {
		var f models.FollowSuggestion
		er := rows.Scan(&f.Id, &f.Firstname, &f.Lastname, &f.ProfileImage, &f.AccountPrivacy)
		if er != nil {
			return nil, er
		}
		followers = append(followers, f)
	}

	return followers, nil
}

// this function is to get the all the following
func (r *Repo) GetFollowingDB(targetID string) ([]models.FollowSuggestion, error) {
	following := []models.FollowSuggestion{}

	rows, er := r.Db.Query(`
		SELECT u.id, u.firstname, u.lastname, u.profile_image, u.account_privacy 
		FROM users u
		JOIN followers f ON u.id = f.following_id
		WHERE f.follower_id = ? AND f.status = 'accepted'
	`, targetID)
	if er != nil {
		return nil, er
	}
	defer rows.Close()

	for rows.Next() {
		var f models.FollowSuggestion
		er := rows.Scan(&f.Id, &f.Firstname, &f.Lastname, &f.ProfileImage, &f.AccountPrivacy)
		if er != nil {
			return nil, er
		}
		following = append(following, f)
	}

	return following, nil
}

// this fuction is to get the dommands of follows you have from db
func (r *Repo) GetFollowRequestsDB(userID string) ([]models.FollowSuggestion, error) {
	rows, err := r.Db.Query(`
		SELECT id, firstname, lastname, profile_image
		FROM users
		WHERE id IN (
			SELECT follower_id FROM followers WHERE following_id = ? AND status = 'pending'
		)`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.FollowSuggestion
	for rows.Next() {
		var u models.FollowSuggestion
		if err := rows.Scan(&u.Id, &u.Firstname, &u.Lastname, &u.ProfileImage); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// this function is to get the freinds the people who are following and you follow
func (r *Repo) GetFriendsDB(userID string) ([]models.FollowSuggestion, error) {
	q := `
		SELECT u.id, u.profile_image, u.firstname, u.lastname, u.account_privacy
		FROM followers f
		JOIN users u ON u.id = f.following_id
		WHERE f.follower_id = ?
		AND f.status = 'accepted'
		AND f.following_id IN (
			SELECT follower_id FROM followers
			WHERE following_id = ?
			AND status = 'accepted'
		)
	`

	rows, err := r.Db.Query(q, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.FollowSuggestion
	for rows.Next() {
		var friend models.FollowSuggestion
		err := rows.Scan(
			&friend.Id,
			&friend.ProfileImage,
			&friend.Firstname,
			&friend.Lastname,
			&friend.AccountPrivacy,
		)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

// =====================
// End followers
// =====================
