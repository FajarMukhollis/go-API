package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Articles Tag structure for your database
type Articles struct {
	ID         string `json:"id"`
	Authors_id string `json:"authors_id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Created_at string `json:"created_at"`
}

var db *sql.DB
var err error

// function to open connection to mysql database
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""         // password, leave it like this if there is no password
	dbName := "kumparan" // database name
	dbIP := "127.0.0.1:3306"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbIP+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db

}

//Index func to view all the records
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()

	var articles []Articles

	result, err := db.Query("SELECT * FROM articles")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var post Articles
		err = result.Scan(&post.ID, &post.Authors_id, &post.Title, &post.Body, &post.Created_at)
		if err != nil {
			panic(err.Error())
		}

		articles = append(articles, post)
	}
	json.NewEncoder(w).Encode(articles)
	defer db.Close()
}

func insertArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	vars := mux.Vars(r)
	Title := vars["title"]
	Body := vars["body"]
	Created_at := vars["created_at"]

	// perform a db.Query insert
	stmt, err := db.Prepare("INSERT INTO articles(authors_id ,title, body, created_at) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(Title, Body, Created_at)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New user was created")
	defer db.Close()
}
func getArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)

	// perform a db.Query insert
	stmt, err := db.Query("SELECT * FROM articles WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	var post Articles
	for stmt.Next() {

		err = stmt.Scan(&post.ID, &post.Authors_id, &post.Title, &post.Body, &post.Created_at)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
	defer db.Close()
}
func delArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)

	// perform a db.Query insert
	stmt, err := db.Prepare("DELETE FROM articles WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	fmt.Fprintf(w, "Articles with ID = %s was deleted", params["id"])
	defer db.Close()
}

func updateArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)
	Title := params["title"]
	Body := params["body"]
	Created_at := params["created_at"]

	// perform a db.Query insert
	stmt, err := db.Prepare("Update articles SET authors_id = ?, title = ?, body = ?, created_at = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(Title, Body, Created_at, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Articles with ID = %s was updated", params["id"])
	defer db.Close()
}
func main() {
	log.Println("Server started on: http://localhost:3000")
	router := mux.NewRouter()
	//On postman try http://localhost:3000/all with method GET
	router.HandleFunc("/all", Index).Methods("GET")
	//On postman try http://localhost:3000/add?title=Test&body=LEB&created_at=7777777 with metho POST
	router.HandleFunc("/add", insertArticles).Methods("POST").Queries("authors_id", "{authors_id}", "title", "{title}", "body", "{body}", "created_at", "{created_at}")
	//On postman try http://localhost:3000/get/1 with method GET
	router.HandleFunc("/get/{id}", getArticles).Methods("GET")
	//On postman try http://localhost:3000/update/1 with method PUT
	router.HandleFunc("/update/{id}", updateArticles).Methods("PUT").Queries("authors_id", "{authors_id}", "title", "{title}", "body", "{body}", "created_at", "{created_at}")
	//On postman try http://localhost:3000/del/1 with method DELETE
	router.HandleFunc("/del/{id}", delArticles).Methods("DELETE")

	http.ListenAndServe(":3000", router)

}
