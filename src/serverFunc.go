package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var Logout = true
var db *sql.DB

func InitHandlers(database *sql.DB) {
	db = database
}

func Createdposthandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		Logout = true
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	Logout = false
	uuid := c.Value
	row := db.QueryRow("SELECT username FROM sessions WHERE token= ?", uuid)
	var username string
	err = row.Scan(&username)
	if err != nil {
	}
	var createdpost []post
	nrow, _ := db.Query("SELECT id ,username,content,topic,like,dislike,commentcount,create_at FROM posts WHERE username =?", username)
	for nrow.Next() {
		var post post
		var Ctime time.Time
		err = nrow.Scan(&post.Id, &post.Username, &post.Content, &post.Topic, &post.Like, &post.Dislike, &post.Commentcount, &Ctime)
		post.Creation = convertime(time.Now().Unix() - Ctime.Unix())
		if err != nil {
		}
		createdpost = append(createdpost, post)
	}

	Data := data{
		Username: username,
		Posts:    createdpost,
		Logout:   Logout,
	}
	Rendertemplate(w, "filter1", Data)

}

func Likedposthandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	uuid := c.Value
	row := db.QueryRow("SELECT username FROM sessions WHERE token= ?", uuid)
	var username string
	err = row.Scan(&username)
	if err != nil {
	}
	rows, _ := db.Query("SELECT post_id FROM post_reaction WHERE score=? AND username = ?", 1, username)
	var postindex []int
	for rows.Next() {
		var postid int
		err = rows.Scan(&postid)
		if err != nil {
		}
		postindex = append(postindex, postid)
	}
	var likedpost []post
	for _, id := range postindex {
		var post post
		var Ctime time.Time
		nrow := db.QueryRow("SELECT id ,username,content,topic,like,dislike,commentcount,create_at FROM posts WHERE id =?", id)
		err = nrow.Scan(&post.Id, &post.Username, &post.Content, &post.Topic, &post.Like, &post.Dislike, &post.Commentcount, &Ctime)
		post.Creation = convertime(time.Now().Unix() - Ctime.Unix())
		if err != nil {
		}
		likedpost = append(likedpost, post)
	}
	Data := data{
		Username: username,
		Posts:    likedpost,
		Logout:   Logout,
	}

	Rendertemplate(w, "filter1", Data)
}

func Filterhandler(w http.ResponseWriter, r *http.Request) {
	var user string
	topic := r.URL.Query().Get("category")
	c, Err := r.Cookie("session_token")
	if Err != nil {
		Logout = true
	} else {
		uuid := c.Value
		row := db.QueryRow("SELECT username FROM sessions WHERE token =?", uuid)
		row.Scan(&user)
		Logout = false
	}
	rows, err := db.Query("SELECT id ,username,content,topic,like,dislike,commentcount,create_at FROM posts WHERE topic =?", topic)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no results")
		}
	}
	var topicPosts []post
	for rows.Next() {
		var post post
		var Ctime time.Time
		rows.Scan(&post.Id, &post.Username, &post.Content, &post.Topic, &post.Like, &post.Dislike, &post.Commentcount, &Ctime)
		post.Creation = convertime(time.Now().Unix() - Ctime.Unix())
		topicPosts = append(topicPosts, post)
	}

	Data := data{
		Username: user,
		Posts:    topicPosts,
		Logout:   Logout,
	}
	Rendertemplate(w, "filter1", Data)
}

var userComment string

func Commentairehandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c, err := r.Cookie("session_token")
		if err != nil {
			if r.Header.Get("Content-Type") == "application/json" {
				data := Data{Error: true}
				json.NewEncoder(w).Encode(data)
				return
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
		uuid := c.Value
		Id, _ := strconv.Atoi(r.FormValue("postId"))
		rows := db.QueryRow("SELECT username FROM sessions WHERE token =?", uuid)
		err = rows.Scan(&userComment)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No user found with the given token.")
			} else {
				fmt.Println("Error querying the database:", err)
			}
			return
		}
		rtype := r.Header.Get("content-type")
		if rtype == "application/json" {
			commentR := Creation{}
			json.NewDecoder(r.Body).Decode(&commentR)
			if commentId, err := strconv.Atoi(commentR.CommentId); err == nil {
				comment_react, _ := strconv.Atoi(commentR.Rreaction)
				fmt.Println(commentId, comment_react)
				row := db.QueryRow("SELECT score FROM comment_reaction WHERE username =? AND comment_id = ?", userComment, commentId)
				var score int
				err = row.Scan(&score)
				if err != nil {
					if err == sql.ErrNoRows {
						statement, _ := db.Prepare("INSERT INTO comment_reaction (username, comment_id) VALUES (?,?)")
						statement.Exec(userComment, commentId)
					} else {
						panic(err)
					}
				}
				if score == 0 {
					statement, _ := db.Prepare("UPDATE comment_reaction SET score = ? WHERE username = ? AND comment_id = ?")
					statement.Exec(comment_react, userComment, commentId)
					if comment_react > 0 {
						statement, _ := db.Prepare("UPDATE comments SET like = like + ? WHERE  id = ?")
						statement.Exec(comment_react, commentId)
					} else {
						statement, _ := db.Prepare("UPDATE comments SET dislike = dislike - ? WHERE  id = ?")
						statement.Exec(comment_react, commentId)
					}

				} else {
					if score == comment_react {
						if comment_react == 1 {
							statement, _ := db.Prepare("UPDATE comments SET like = like - ? WHERE  id = ?")
							statement.Exec(comment_react, commentId)
							statement, _ = db.Prepare("UPDATE comment_reaction SET score = ? WHERE username = ? AND comment_id = ?")
							statement.Exec(0, userComment, commentId)
						} else {
							statement, _ := db.Prepare("UPDATE comments SET dislike = dislike + ? WHERE  id = ?")
							statement.Exec(comment_react, commentId)
							statement, _ = db.Prepare("UPDATE comment_reaction SET score = ? WHERE username = ? AND comment_id = ?")
							statement.Exec(0, userComment, commentId)
						}
					} else {
						if comment_react == 1 {
							statement, _ := db.Prepare("UPDATE comments SET like = like + ? WHERE  id = ?")
							statement.Exec(comment_react, commentId)
							statement, _ = db.Prepare("UPDATE comments SET dislike = dislike - ? WHERE  id = ?")
							statement.Exec(comment_react, commentId)
							statement, _ = db.Prepare("UPDATE comment_reaction SET score = ? WHERE username = ? AND comment_id = ?")
							statement.Exec(1, userComment, commentId)
						} else {
							statement, _ := db.Prepare("UPDATE comments SET dislike = dislike - ? WHERE  id = ?")
							statement.Exec(comment_react, commentId)
							statement, _ = db.Prepare("UPDATE comments SET like = like + ? WHERE  id = ?")
							statement.Exec(comment_react, commentId)
							statement, _ = db.Prepare("UPDATE comment_reaction SET score = ? WHERE username = ? AND comment_id = ?")
							statement.Exec(-1, userComment, commentId)
						}

					}
				}
				row = db.QueryRow("SELECT `like` ,`dislike` FROM comments WHERE id = ? AND username= ?", commentId, userComment)
				type res struct {
					Like    int
					Dislike int
				}
				Res := res{}
				err = row.Scan(&Res.Like, &Res.Dislike)
				if err != nil {
				}
				type Response struct {
					CommentId int `json:"commentId"`
					Like      int `json:"like"`
					Dislike   int `json:"dislike"`
				}
				response := Response{CommentId: commentId, Like: Res.Like, Dislike: Res.Dislike}
				json.NewEncoder(w).Encode(response)
				return
			}
		} else {
			content := r.FormValue("content")
			if content != "" {
				statement, _ := db.Prepare("INSERT INTO comments (username, content,create_at,post_id) VALUES (?,?,?,?)")
				statement.Exec(userComment, content, time.Now(), Id)
				statement, _ = db.Prepare("UPDATE posts SET commentcount = commentcount + ? WHERE id = ?")
				statement.Exec(1, Id)

			}
			http.Redirect(w, r, fmt.Sprintf("/commentaires?postId=%d", Id), http.StatusSeeOther)
			return
		}
	}
	Id, _ := strconv.Atoi(r.URL.Query().Get("postId"))
	c, _ := r.Cookie("session_token")
	uuid := c.Value
	ro := db.QueryRow("SELECT username FROM sessions WHERE token =?", uuid)
	ro.Scan(&userComment)
	rows, err := db.Query("SELECT id,content,username, like , dislike, create_at FROM comments WHERE post_id = ?", Id)
	if err != nil {
		panic(err)
	}
	var comments = []comment{}
	for rows.Next() {
		var newcomment comment
		var Ctime time.Time
		err := rows.Scan(&newcomment.Id, &newcomment.Content, &newcomment.Username, &newcomment.Like, &newcomment.Dislike, &Ctime)
		newcomment.Creation = convertime(time.Now().Unix() - Ctime.Unix())
		fmt.Println(newcomment.Creation)
		if err == nil {
			newcomment.PostId = Id
			comments = append(comments, newcomment)
		}
	}

	row := db.QueryRow("SELECT  id ,username, topic, content, create_at , like, dislike , commentcount FROM posts WHERE id =? ", Id)
	var Post post
	var Ctime time.Time
	err = row.Scan(&Post.Id, &Post.Username, &Post.Topic, &Post.Content, &Ctime, &Post.Like, &Post.Dislike, &Post.Commentcount)
	Post.Creation = convertime(time.Now().Unix() - Ctime.Unix())
	if err == nil {
	}
	data := Data{
		Username: userComment,
		PostId:   Id,
		Comments: comments,
		Post:     Post,
		Logout:   Logout,
	}
	fmt.Println(userComment)

	Rendertemplate(w, "commentaire", data)
}

func Likehandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		PostID  int  `json:"postId"`
		Like    int  `json:"like"`
		Dislike int  `json:"dislike"`
		Error   bool `json:"error"`
	}
	response := Response{}
	if r.Method == http.MethodPost {
		request := Rreaction{}
		json.NewDecoder(r.Body).Decode(&request)
		c, err := r.Cookie("session_token")

		if err != nil {
			response.Error = true
			json.NewEncoder(w).Encode(response)
			return
		}
		uuid := c.Value
		post_index, _ := strconv.Atoi(request.PostId)
		post_react, _ := strconv.Atoi(request.Reaction)
		row := db.QueryRow("SELECT username FROM sessions WHERE token = ?", uuid)
		var user string
		err = row.Scan(&user)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No user found with the given token.")
			} else {
				fmt.Println("Error querying the database:", err)
			}
			return
		}
		row = db.QueryRow("SELECT score FROM post_reaction WHERE username =? AND post_id = ?", user, post_index)
		var score int
		err = row.Scan(&score)
		if err != nil {
			if err == sql.ErrNoRows {
				statement, _ := db.Prepare("INSERT INTO post_reaction (username, post_id) VALUES (?,?)")
				statement.Exec(user, post_index)
			} else {
				panic(err)
			}
		}
		if score == 0 {
			statement, _ := db.Prepare("UPDATE post_reaction SET score = ? WHERE username = ? AND post_id = ?")
			statement.Exec(post_react, user, post_index)
			if post_react == 1 {
				statement, _ := db.Prepare("UPDATE posts SET like = like + ? WHERE  id = ?")
				statement.Exec(post_react, post_index)
			} else if post_react == -1 {
				statement, _ := db.Prepare("UPDATE posts SET dislike = dislike - ? WHERE  id = ?")
				statement.Exec(post_react, post_index)
			}

		} else {
			if score == post_react {
				if post_react == 1 {
					statement, _ := db.Prepare("UPDATE posts SET like = like - ? WHERE  id = ?")
					statement.Exec(post_react, post_index)
					statement, _ = db.Prepare("UPDATE post_reaction SET score = ? WHERE username = ? AND post_id = ?")
					statement.Exec(0, user, post_index)
				} else {
					statement, _ := db.Prepare("UPDATE posts SET dislike = dislike + ? WHERE  id = ?")
					statement.Exec(post_react, post_index)
					statement, _ = db.Prepare("UPDATE post_reaction SET score = ? WHERE username = ? AND post_id = ?")
					statement.Exec(0, user, post_index)
				}

			} else {
				if post_react == 1 {
					statement, _ := db.Prepare("UPDATE posts SET like = like + ? WHERE  id = ?")
					statement.Exec(post_react, post_index)
					statement, _ = db.Prepare("UPDATE posts SET dislike = dislike - ? WHERE  id = ?")
					statement.Exec(post_react, post_index)
					statement, _ = db.Prepare("UPDATE post_reaction SET score = ? WHERE username = ? AND post_id = ?")
					statement.Exec(1, user, post_index)
				} else {
					statement, _ := db.Prepare("UPDATE posts SET dislike = dislike - ? WHERE  id = ?")
					statement.Exec(post_react, post_index)
					statement, _ = db.Prepare("UPDATE posts SET like = like + ? WHERE  id = ?")
					statement.Exec(post_react, post_index)
					statement, _ = db.Prepare("UPDATE post_reaction SET score = ? WHERE username = ? AND post_id = ?")
					statement.Exec(-1, user, post_index)
				}

			}
		}

		response = Response{
			PostID: post_index,
		}
		row = db.QueryRow("SELECT like, dislike FROM posts WHERE id = ?", post_index)
		row.Scan(&response.Like, &response.Dislike)
		json.NewEncoder(w).Encode(&response)

	}

}

func Postshandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	uuid := c.Value
	row := db.QueryRow("SELECT username FROM sessions WHERE token = ?", uuid)
	var Username string
	err = row.Scan(&Username)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		content := r.FormValue("content")
		topic := r.FormValue("topic")
		statement, _ := db.Prepare("INSERT INTO posts (username,content,topic,create_at) VALUES (?,?,?,?)")
		statement.Exec(Username, content, topic, time.Now())
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
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
		fmt.Println(newPost.Creation)
		posts = append(posts, newPost)
	}

	Rendertemplate(w, "index", posts)

}

func Logouthandler(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	uuid := c.Value
	statement, _ := db.Prepare("DELETE FROM sessions WHERE token = ?")
	statement.Exec(uuid)
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Hashpassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err

}
func Checkpassword(hashed string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func Rendertemplate(w http.ResponseWriter, temp string, data interface{}) {
	tempath := fmt.Sprintf("templates/%s.html", temp)
	t, err := template.ParseFiles(tempath)
	if err != nil {
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)

}

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
		posts = append(posts, newPost)
	}

	Data := data{
		Username: Username,
		Posts:    posts,
		Logout:   Logout,
	}
	Rendertemplate(w, "index", Data)

}
func Loginhandler(w http.ResponseWriter, r *http.Request) {
	Rendertemplate(w, "login", nil)
}
func Signup(w http.ResponseWriter, r *http.Request) {
	Rendertemplate(w, "signup", nil)
}
func Signuphandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	hashed, _ := Hashpassword(password)
	statement, _ := db.Prepare("INSERT INTO users (username,password,email) VALUES(?,?,?)")
	statement.Exec(username, hashed, email)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Connexionhandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Username := r.FormValue("username")
	Password := r.FormValue("password")
	rows := db.QueryRow("SELECT password FROM users WHERE username = ?", Username)
	var password string
	err := rows.Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			Rendertemplate(w, "login", "invalid Credentials")
		} else {
			http.Error(w, "Erreur lors de la recherche de l'utilisateur", http.StatusInternalServerError)
		}
		return
	}
	if Checkpassword(password, Password) {
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
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		Rendertemplate(w, "login", "invalid Credentials")
	}

}
