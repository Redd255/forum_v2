package server

import (
	"database/sql"
	"text/template"
)

var Logout = true
var db *sql.DB

func InitHandlers(database *sql.DB) {
	db = database
}

var tmpl = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/login.html",
	"templates/filter1.html",
	"templates/signup.html",
	"templates/commentaire.html"))

type comment struct {
	Id       int
	Username string
	Content  string
	Creation string
	Like     int
	Dislike  int
	PostId   int
	Score    int
}

type post struct {
	Id           int
	Username     string
	Content      string
	Creation     string
	Topic        string
	Like         int
	Dislike      int
	Commentcount int
	Score        int
}

type data struct {
	Username string
	Posts    []post
	Logout   bool
}

type Data struct {
	Username string
	PostId   int
	Comments []comment
	Post     post
	Logout   bool
	Error    bool
}
type Rreaction struct {
	PostId   string `json:"postId"`
	Reaction string `json:"reaction"`
}

type Creation struct {
	CommentId string `json:"commentId"`
	Rreaction string `json:"reaction"`
}
type Response struct {
	CommentId int `json:"commentId"`
	Like      int `json:"like"`
	Dislike   int `json:"dislike"`
	Score     int `json:"score"`
}

