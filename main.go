package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	server "zone/src"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./databases/data.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Error enabling foreign keys:", err)
	}
	sqlStatements, err := os.ReadFile("./databases/mine.sql")
	if err != nil {
		log.Fatal("Error reading SQL file:", err)
	}
	_, err = db.Exec(string(sqlStatements))
	if err != nil {
		log.Fatal("Error executing SQL statements:", err)
	}
	server.InitHandlers(db)
}

func main() {
	defer db.Close()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", server.Homehandelr)
	http.HandleFunc("/signup", server.Signup)
	http.HandleFunc("/login", server.Login)
	
	http.HandleFunc("/createdpost", server.Createdposthandler)
	http.HandleFunc("/likedpost", server.Likedposthandler)
	http.HandleFunc("/filter", server.Filterhandler)
	http.HandleFunc("/logout", server.Logouthandler)
	http.HandleFunc("/posts", server.Postshandler)
	http.HandleFunc("/react", server.Likehandler)
	http.HandleFunc("/commentaires", server.Commentairehandler)

	fmt.Println("http://localhost:9090")
	http.ListenAndServe(":9090", nil)
}
