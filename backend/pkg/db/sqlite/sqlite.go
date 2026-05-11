package sqlite

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"
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

// this function is to check if the user already exists
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
	sessionUUID, err := uuid.NewV7()
	if err != nil {
		return "", time.Time{}, errors.New("SERVER ERROR")
	}
	sessionID := sessionUUID.String()

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)

	// Remove any existing session for this user to enforce single-session policy
	_, err = r.Db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return "", time.Time{}, errors.New("SERVER ERROR")
	}

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
func (r *Repo) DisconnectUser(token string) error {
	_, er := r.Db.Exec("DELETE FROM sessions WHERE token = ?", token)
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

	if groupID != "" {
		var exist int
		err = r.Db.QueryRow("SELECT 1 FROM groups WHERE id = ? ", groupID).Scan(&exist)
		if err != nil || err == sql.ErrNoRows {
			return models.Post{}, errors.New("bad request")
		}
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
			WHERE user_id = ? AND (group_id = "" OR group_id IS NULL)
			ORDER BY created_at DESC
			LIMIT 10 OFFSET ?`, ReqUserID, Offset)

	case "profile-me-activity":
		rows, err = r.Db.Query(`
			SELECT p.id, p.user_id, p.content, p.image_url, p.created_at, p.privacy, p.allowed_users, p.group_id
			FROM posts p
			WHERE (p.group_id = "" OR p.group_id IS NULL) AND p.id IN (
				SELECT post_id
				FROM comments
				WHERE user_id = ?
			)
			ORDER BY p.created_at DESC
			LIMIT 10 OFFSET ?`, ReqUserID, Offset)

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
			WHERE group_id = "" OR group_id IS NULL
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
		er := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ImageURL, &c.CreatedAt)
		if er != nil {
			return nil, nil
		}

		er = r.Db.QueryRow("SELECT nickname, firstname, lastname, profile_image FROM users WHERE id = ?", c.UserID).
			Scan(&c.AutherName, &c.FirstName, &c.LastName, &c.UserImageProfile)
		comments = append(comments, c)
	}

	return comments, nil
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
	user.ID = userID

	return user, nil
}

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

	if finalStatus == "pending" {
		//
	}

	return message, nil
}

// this function is to answer to the request wether you accept or now if your profile is privet
func (r *Repo) ManageFollowDB(followerID string, followingID string, decision string) error {
	if decision == "accepted" {
		_, err := r.Db.Exec(`UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND following_id = ?`, followerID, followingID)
		if err != nil {
			return err
		}
	} else if decision == "rejected" {
		_, err := r.Db.Exec(`DELETE FROM followers WHERE follower_id = ? AND following_id = ?`, followerID, followingID)
		if err != nil {
			return err
		}
	}

	// Dans les deux cas, supprimer la notification follow_request
	_, err := r.Db.Exec(`
		DELETE FROM notifications
		WHERE type = 'follow_request'
		AND sender_id = ?
		AND id IN (
			SELECT notification_id FROM notification_users WHERE user_id = ?
		)
	`, followerID, followingID)
	return err
}

func (r *Repo) GetUsers(userid string, query string) ([]models.FollowSuggestion, error) {
	users := []models.FollowSuggestion{}

	q := `SELECT u.id, u.firstname, u.lastname, u.profile_image, u.account_privacy FROM users u WHERE u.id != ?`
	args := []any{userid}

	if query != "" {
		q += " AND (u.firstname LIKE ? OR u.lastname LIKE ? OR u.firstname || ' ' || u.lastname LIKE ? OR u.nickname LIKE ?)"
		pattern := "%" + query + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	rows, er := r.Db.Query(q, args...)

	if er != nil {
		return nil, er
	}
	for rows.Next() {
		var f models.FollowSuggestion
		er := rows.Scan(&f.Id, &f.Firstname, &f.Lastname, &f.ProfileImage, &f.AccountPrivacy)
		if er != nil {
			return nil, er
		}
		users = append(users, f)
	}

	return users, nil
}

// this function is to get the all the followers
func (r *Repo) GetFollowersDB(targetID string, query string) ([]models.FollowSuggestion, error) {
	followers := []models.FollowSuggestion{}

	q := `
		SELECT u.id, u.firstname, u.lastname, u.profile_image, u.account_privacy 
		FROM users u
		JOIN followers f ON u.id = f.follower_id
		WHERE f.following_id = ? AND f.status = 'accepted'
	`
	args := []any{targetID}

	if query != "" {
		q += " AND (u.firstname LIKE ? OR u.lastname LIKE ? OR u.firstname || ' ' || u.lastname LIKE ? OR u.nickname LIKE ?)"
		p := "%" + query + "%"
		args = append(args, p, p, p, p)
	}

	rows, er := r.Db.Query(q, args...)
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
func (r *Repo) GetFollowingDB(targetID string, query string) ([]models.FollowSuggestion, error) {
	following := []models.FollowSuggestion{}

	q := `
		SELECT u.id, u.firstname, u.lastname, u.profile_image, u.account_privacy 
		FROM users u
		JOIN followers f ON u.id = f.following_id
		WHERE f.follower_id = ? AND f.status = 'accepted'
	`
	args := []any{targetID}

	if query != "" {
		q += " AND (u.firstname LIKE ? OR u.lastname LIKE ? OR u.firstname || ' ' || u.lastname LIKE ? OR u.nickname LIKE ?)"
		pattern := "%" + query + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	rows, er := r.Db.Query(q, args...)
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
func (r *Repo) GetFriendsDB(userID string, query string) ([]models.FollowSuggestion, error) {
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
	args := []any{userID, userID}

	if query != "" {
		q += " AND (u.firstname LIKE ? OR u.lastname LIKE ? OR u.firstname || ' ' || u.lastname LIKE ? OR u.nickname LIKE ?)"
		pattern := "%" + query + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	rows, err := r.Db.Query(q, args...)
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

func (r *Repo) GetUserByIDDB(userID string) (models.User, error) {
	var u models.User
	err := r.Db.QueryRow(`
		SELECT id, nickname, firstname, lastname, email, profile_image
		FROM users WHERE id = ?
	`, userID).Scan(&u.ID, &u.Nickname, &u.Firstname, &u.Lastname, &u.Email, &u.ProfileImage)
	return u, err
}

func (r *Repo) IsFollowExist(sender, receiver string) bool {
	var exists int
	err := r.Db.QueryRow(`SELECT 1 FROM followers WHERE follower_id=? AND following_id=? AND status='accepted' LIMIT 1`, receiver, sender).
		Scan(&exists)

	if err == nil {
		return true
	}

	var Allowed bool
	err = r.Db.QueryRow("SELECT account_privacy from users WHERE id = ?", receiver).Scan(&Allowed)

	return Allowed
}
