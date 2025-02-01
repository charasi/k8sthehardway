/*
	Simple BookStore Service that adds book/customer to a database
	and retrieves book by ISBN and customer by ID
	Charles Asiama
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// represents a customer record, mapped to JSON structure
type CustomerAdded struct {
	ID       string `json:"id"`
	UserId   string `json:"userId"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zipcode  string `json:"zipcode"`
}

// message structure
type Duplicate struct {
	Message string `json:"message"`
}

// represents a related book record, mapped to JSON structure
type RelatedBooks struct {
	ISBN   string `json:"ISBN"`
	Title  string `json:"title"`
	Author string `json:"Author"`
}

/*
access an external recommendation engine service every time the “related books”
endpoint is executed and return recommendations for additional books the customer
may want to purchase.
*/
func RetrieveRelatedBooksEndpoint(w http.ResponseWriter, r *http.Request) {
	// reroutes to internal load balancer
	//var responseRelatedBook RelatedBooks
	var response = getRelatedBookRequest(r)
	var responseBody, _ = io.ReadAll(response.Body)
	//json.Unmarshal(responseBody, &responseRelatedBook)

	// if book is added, send success response
	//jsonResponseBody, _ := json.Marshal(&responseRelatedBook)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write([]byte(responseBody))
}

/*
send request to book service
*/
func getRelatedBookRequest(r *http.Request) *http.Response {
	// send request, get response
	var response *http.Response
	baseUrl := "http://127.0.0.1:2345"
	path := r.URL.Path
	url := baseUrl + path

	response, _ = http.Get(url)
	//
	return response
}

/*
retrieves a book from mongoDB book table
*/
func GetBookRecord(value string, bookAdded *BookAdded) int {
	// for connection and interacting with mongoDB
	var collection *mongo.Collection
	var ctx = context.TODO()

	// connect to mongoDB
	clientOptions := options.Client().ApplyURI("mongodb+srv://casiama:RIO1yYeZnijK4pJR@assignment4.yvbw3da.mongodb.net/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return 1
	}

	time.Sleep(40 * time.Second)

	// search datbase for specific ISBN
	isbn := bson.D{{Key: "ISBN", Value: value}}
	collection = client.Database("BooksDB").Collection("books_casiama")
	col_err := collection.FindOne(ctx, isbn).Decode(&bookAdded)
	if col_err != nil {
		return 1
	}

	return 0
}

/*
queies the MongoDB collection looking for any book documents
that contain the given keyword anywhere in the document data
*/
func GetBookRecordList(value string) ([]BookAdded, int) {
	var bookAdded []BookAdded
	// for connection and interacting with mongoDB
	var collection *mongo.Collection
	var contex = context.TODO()

	// connect to mongoDB
	clientOptions := options.Client().ApplyURI("mongodb+srv://casiama:RIO1yYeZnijK4pJR@assignment4.yvbw3da.mongodb.net/")
	client, err := mongo.Connect(contex, clientOptions)
	if err != nil {
		return bookAdded, 1
	}

	// create collection
	collection = client.Database("BooksDB").Collection("books_casiama")

	// drop any existing index that was previously created for searching collection
	_, dropIdxErr := collection.Indexes().DropAll(contex)
	if dropIdxErr != nil {
		fmt.Println("ERROR! COULD NOT DELETE INDEXES")
		return bookAdded, 1
	}

	// create index view to be used for searching
	model := []mongo.IndexModel{
		{Keys: bson.D{{Key: "title,", Value: "text"},
			{Key: "Author", Value: "text"},
			{Key: "description", Value: "text"},
			{Key: "genre", Value: "text"}},
		}}

	_, createIdxErr := collection.Indexes().CreateMany(contex, model)
	if createIdxErr != nil {
		fmt.Println("ERROR! COULD NOT CREATE INDEXES")
		return bookAdded, 1
	}

	// word to search for
	search := bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: "\"" + value + "\""}}}}

	// find documents in database
	cursor, col_err := collection.Find(contex, search)
	if col_err != nil {
		return bookAdded, 1
	}

	// parse to array struct book added
	cursor_err := cursor.All(contex, &bookAdded)
	if cursor_err != nil {
		return bookAdded, 1
	}

	// if no document was found on MongoDB
	if len(bookAdded) == 0 {
		return bookAdded, 2
	}

	// document was found
	return bookAdded, 0
}

/*
return a book given its ISBN. Both endpoints shall produce the same response.
*/
func RetrieveBookEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to bookAdded
	var bookAdded BookAdded
	var isbn = path.Base(r.URL.Path)

	// get book record, return error if not found
	result := GetBookRecord(isbn, &bookAdded)
	if result == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		return
	}

	// return book if found
	jsonResponseBody, _ := json.Marshal(&bookAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/books/"+bookAdded.ISBN)
	w.WriteHeader(200)
	w.Write([]byte(jsonResponseBody))
}

/*
return a book given its ISBN. Both endpoints shall produce the same response.
*/
func RetrieveBooksKeywordEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to bookAdded
	var bookAdded []BookAdded

	var keyword = strings.Split(path.Base(r.URL.RawQuery), "=")[1]

	is_alphabet := regexp.MustCompile(`^[a-zA-Z]*$`).MatchString(keyword)

	if !is_alphabet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	// get book record, return error if not found
	bookAdded, result := GetBookRecordList(keyword)
	if result == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		return
	}

	// empty content
	if result == 2 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)
		return
	}

	// return book if found
	jsonResponseBody, _ := json.Marshal(bookAdded)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(jsonResponseBody))
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
	mux := mux.NewRouter()
	mux.HandleFunc("/books/isbn/{ISBN}", RetrieveBookEndpoint).Methods("GET")
	mux.HandleFunc("/books/{ISBN}", RetrieveBookEndpoint).Methods("GET")
	mux.HandleFunc("/books/{ISBN}/related-books", RetrieveRelatedBooksEndpoint).Methods("GET")
	mux.HandleFunc("/books", RetrieveBooksKeywordEndpoint).Queries(
		"keyword", "{keyword}").Methods("GET")
	mux.HandleFunc("/status", status).Methods("GET")
	//http.ListenAndServe(":2345", mux)
	http.ListenAndServe(":3000", mux)
}
