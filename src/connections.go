package server

import (
	"database/sql"
	"log"
	"net/http"
	"slices"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Homehandelr(w http.ResponseWriter, r *http.Request) {
	var Username string
	a, err := r.Cookie("session_token")
	if err == nil {
		Logout = false
		uuid := a.Value
		rows := db.QueryRow("SELECT username FROM  sessions WHERE token =? ", uuid)
		err := rows.Scan(&Username)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	} else {
		Logout = true
	}
	rows, err := db.Query("SELECT id ,username,content,topic,like,dislike,commentcount,create_at FROM posts")
	if err != nil {
		panic(err)
	}
	var posts = []post{}
	for rows.Next() {
		var newPost post
		var Ctime time.Time
		rows.Scan(&newPost.Id, &newPost.Username, &newPost.Content, &newPost.Topic, &newPost.Like, &newPost.Dislike, &newPost.Commentcount, &Ctime)
		newPost.Creation = convertime(time.Now().Unix() - Ctime.Unix())
		r := db.QueryRow("SELECT score FROM post_reaction WHERE post_id = ? AND username = ?", newPost.Id, Username)
		r.Scan(&newPost.Score)
		posts = append(posts, newPost)
	}
	slices.Reverse(posts)
	Data := data{
		Username: Username,
		Posts:    posts,
		Logout:   Logout,
	}
	tmpl.ExecuteTemplate(w, "index.html", Data)

}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	Username := r.FormValue("username")
	Password := r.FormValue("password")

	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", Username).Scan(&hashedPassword)
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

	uuid, _ := GenerateUUID()
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    uuid,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
	statement, err := db.Prepare("INSERT INTO sessions (token,username) VALUES (?,?)")
	if err != nil {
		log.Println("Erreur lors de la préparation de la requête :", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	statement.Exec(uuid, Username)
	defer statement.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Signup(w http.ResponseWriter, r *http.Request) {
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
	err := db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)
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
