/*
	Simple BookStore Service that adds book to a database
	Charles Asiama
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// BookAdded represents a book record, mapped to JSON structure
type BookAdded struct {
	ISBN        string  `json:"ISBN"`
	Title       string  `json:"title"`
	Author      string  `json:"Author"`
	Description string  `json:"description"`
	Genre       string  `json:"genre"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// Duplicate message structure
type Duplicate struct {
	Message string `json:"message"`
}

var USER = os.Getenv("DB_USER")
var PASS = os.Getenv("DB_PASS")
var HOST = os.Getenv("DB_HOST")
var PORT = os.Getenv("DB_PORT")
var BOOKSTORE = os.Getenv("DB_BOOKSTORE")

/*
AddBookEndpoint
Adds a book to the system. The ISBN will be the unique identifier for the book.
The book is added to the Book data table on MySql (the ISBN is the primary key).
*/
func AddBookEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to bookAdded
	var bookAdded BookAdded
	var requestBody, _ = io.ReadAll(r.Body)
	err := json.Unmarshal(requestBody, &bookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		return
	}

	fmt.Println("BOOKCMD: ", bookAdded)

	// add new book to database
	query := "INSERT INTO book (isbn, title, author, description, genre, price, quantity)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	// arguments for query
	args := []interface{}{bookAdded.ISBN, bookAdded.Title, bookAdded.Author,
		bookAdded.Description, bookAdded.Genre, bookAdded.Price, bookAdded.Quantity}

	// get record
	result := InsertUpdateBookTable(query, args)
	// if book is duplicate return 422 status
	if result == 1 {
		var message Duplicate
		message.Message = "This ISBN already exists in the system."
		jsonMessage, _ := json.Marshal(&message)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(422)
		_, err = w.Write(jsonMessage)
		if err != nil {
			log.Fatalf("Error writing response from server: %v", err)
		}
		return
	}

	// if book is added, send success response
	jsonResponseBody, _ := json.Marshal(&bookAdded)
	fmt.Println("RESPONSE: ", string(jsonResponseBody))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/books/"+bookAdded.ISBN)
	w.WriteHeader(201)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
UpdateBookEndpoint
Update a book’s information in the system.
The ISBN will be the unique identifier for the book.
*/
func UpdateBookEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to bookAdded
	var bookAdded BookAdded
	var requestBody, _ = io.ReadAll(r.Body)
	var isbn string = path.Base(r.URL.Path)
	err := json.Unmarshal(requestBody, &bookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		return
	}

	// update book
	query := "UPDATE book SET title = ?, author = ?, description = ?, " +
		"genre = ?, price = ?, quantity = ? WHERE isbn = ?"

	bookAdded.ISBN = isbn
	// arguments for query
	args := []interface{}{bookAdded.Title, bookAdded.Author, bookAdded.Description,
		bookAdded.Genre, bookAdded.Price, bookAdded.Quantity, bookAdded.ISBN}

	// get record
	result := InsertUpdateBookTable(query, args)
	// if ISBN does not exist
	if result == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		return
	}

	time.Sleep(40 * time.Second)

	// if book is updated, send success response
	jsonResponseBody, _ := json.Marshal(&bookAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/books/"+bookAdded.ISBN)
	w.WriteHeader(200)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
InsertUpdateBookTable
Insert, or Update a table
*/
func InsertUpdateBookTable(query string, args []interface{}) int {

	fmt.Println("About to insert : " + query)

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		USER, PASS, HOST, PORT, BOOKSTORE)
	//db, err := sql.Open("mysql", USER+":"+PASS+"@tcp("+HOST+":"+PORT+")/"+BOOKSTORE)
	//db, err := sql.Open("mysql", "remedy:skincream@tcp(34.145.15.102:3306)/bookstore")
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		panic(err.Error())
	}

	// close database after exiting function
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}(db)

	// execute query
	success, err := db.Exec(query, args...)

	// if error and message is duplicate return error message, else panic and exit
	if err != nil {
		message := err.Error()
		fmt.Println(message)
		isDuplicate := strings.Contains(message, "Duplicate entry")

		if isDuplicate {
			return 1
		}
	}

	// if no error but table was not affected return 1, else 0
	rows, _ := success.RowsAffected()
	if rows > 0 {
		return 0
	} else {
		return 1
	}

}

// In your K8S deployment file, specify a liveness probe
func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Fatalf("Error sending response %v", err)
	}
}

/*
Main function to handle routes/path for web server
*/
func main() {
	//router := http.NewServeMux()
	router := mux.NewRouter()
	router.HandleFunc("/cmd/books", AddBookEndpoint).Methods("POST")
	router.HandleFunc("/cmd/books/{ISBN}", UpdateBookEndpoint).Methods("POST")
	router.HandleFunc("/status", status).Methods("GET")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
