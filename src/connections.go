package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Homehandelr(w http.ResponseWriter, r *http.Request) {
	var Username string
	cookie, err := r.Cookie("session_token")
	if err == nil {
		Logout = false
		uuid := cookie.Value
		db.QueryRow("SELECT username FROM  sessions WHERE token =? ", uuid).Scan(&Username)
	} else {
		Logout = true
	}

	rows, err := db.Query("SELECT id, username, content, topic, like , dislike , commentcount, create_at FROM posts ORDER BY create_at DESC")
	if err != nil {
		return
	}

	var posts = []post{}
	for rows.Next() {
		var newPost post
		var Ctime time.Time
		rows.Scan(
			&newPost.Id,
			&newPost.Username,
			&newPost.Content,
			&newPost.Topic,
			&newPost.Like,
			&newPost.Dislike,
			&newPost.Commentcount,
			&Ctime)
		newPost.Creation = convertime(time.Now().Unix() - Ctime.Unix())
		db.QueryRow("SELECT score FROM post_reaction WHERE post_id = ? AND username = ?", newPost.Id, Username).Scan(&newPost.Score)
		posts = append(posts, newPost)
	}

	Data := data{
		Username: Username,
		Posts:    posts,
		Logout:   Logout,
	}
	tmpl.ExecuteTemplate(w, "index.html", Data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	Username := r.FormValue("username")
	Password := r.FormValue("password")
	if Username == ""|| Password == "" {
		errorPage(w, "All fields are required", "login.html")
		return
	}
	var hashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", Username).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		errorPage(w, "Invalid username", "login.html")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Password)); err != nil {
		errorPage(w, "Invalid password", "login.html")
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	uuid, err := GenerateUUID()
	if err != nil {
		log.Printf("UUID generation error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    uuid,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		// SameSite: http.SameSiteLaxMode,
		// MaxAge:   86400,
	}

	http.SetCookie(w, &cookie)

	_, err = db.Exec("INSERT INTO sessions (token,username) VALUES (?,?)", uuid, Username)
	if err != nil {
		log.Printf("Session insertion error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "signup.html", nil)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		errorPage(w, "All fields are required", "signup.html")
		return
	}

	var existingEmail string
	err = db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)
	if err == nil {
		errorPage(w, "Email already in use", "signup.html")
		return
	}

	if err != sql.ErrNoRows {
		log.Println("Database error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var existingUserName string
	err = db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUserName)
	if err == nil {
		errorPage(w, "UserName already in use", "signup.html")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
	if err != nil {
		log.Println("Failed to insert user:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
