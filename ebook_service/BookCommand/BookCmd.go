/*
	Simple BookStore Service that adds book to a database
	Charles Asiama
*/

package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// represents a book record, mapped to JSON structure
type BookAdded struct {
	ISBN        string  `json:"ISBN"`
	Title       string  `json:"title"`
	Author      string  `json:"Author"`
	Description string  `json:"description"`
	Genre       string  `json:"genre"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// message structure
type Duplicate struct {
	Message string `json:"message"`
}

/*
Adds a book to the system. The ISBN will be the unique identifier for the book.
The book is added to the Book data table on MySql (the ISBN is the primary key).
*/
func AddBookEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to bookAdded
	var bookAdded BookAdded
	var requestBody, _ = io.ReadAll(r.Body)
	json.Unmarshal(requestBody, &bookAdded)

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
		w.Write(jsonMessage)
		return
	}

	// if book is added, send success response
	jsonResponseBody, _ := json.Marshal(&bookAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/books/"+bookAdded.ISBN)
	w.WriteHeader(201)
	w.Write([]byte(jsonResponseBody))
}

/*
Update a book’s information in the system.
The ISBN will be the unique identifier for the book.
*/
func UpdateBookEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to bookAdded
	var bookAdded BookAdded
	var requestBody, _ = io.ReadAll(r.Body)
	var isbn string = path.Base(r.URL.Path)
	json.Unmarshal(requestBody, &bookAdded)

	// update book
	query := "UPDATE book SET title = ?, author = ?, description = ?, " +
		"genre = ?, price = ?, quantity = ? WHERE isbn = ?"

	bookAdded.ISBN = isbn
	// arguments for query
	args := []interface{}{bookAdded.Title, bookAdded.Author, bookAdded.Description,
		bookAdded.Genre, bookAdded.Price, bookAdded.Quantity, bookAdded.ISBN}

	// get record
	result := InsertUpdateBookTable(query, args)
	// if ISBN does not exists
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
	w.Write([]byte(jsonResponseBody))
}

/*
Insert, or Update a table
*/
func InsertUpdateBookTable(query string, args []interface{}) int {
	db, err := sql.Open("mysql", "awsadmin:awspassword@tcp(assn4-dbauroraa-esuvdyingzvs.c42qw4pddowd.us-east-1.rds.amazonaws.com:3306)/bookstore")
	//db, err := sql.Open("mysql", "charasi:Skittles05@10@tcp(127.0.0.1:3306)/bookstore")
	if err != nil {
		panic(err.Error())
	}

	// close database after exiting function
	defer db.Close()

	// execute query
	success, err := db.Exec(query, args...)

	// if error and message is duplicate return error message, else panic and exit
	if err != nil {
		message := err.Error()
		isDuplicate := strings.Contains(message, "Duplicate entry")

		if isDuplicate {
			return 1
		}
		//panic(err.Error())
	}

	// if no error but table was not affected return 1, else 0
	rows, _ := success.RowsAffected()
	if rows > 0 {
		return 0
	} else {
		return 1
	}

}

//	monitor the health of the REST service within EKS.
//
// In your K8S deployment file, specify a liveness probe
func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK"))
}

/*
Main function to handle routes/path for web server
*/
func main() {
	//mux := http.NewServeMux()
	mux := mux.NewRouter()
	mux.HandleFunc("/cmd/books", AddBookEndpoint).Methods("POST")
	mux.HandleFunc("/cmd/books/{ISBN}", UpdateBookEndpoint).Methods("POST")
	mux.HandleFunc("/status", status).Methods("GET")
	//http.ListenAndServe(":2345", mux)
	http.ListenAndServe(":3000", mux)
}
