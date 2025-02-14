/*
	Simple Bookstore Service that adds books/customers to a database
	and retrieves books by ISBN and customers by ID.
	Charles Asiama
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"
)

// BookAdded represents the structure of a book record in the system.
type BookAdded struct {
	ISBN        string  `json:"ISBN"`
	Title       string  `json:"title"`
	Author      string  `json:"Author"`
	Description string  `json:"description"`
	Genre       string  `json:"genre"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// CustomerAdded represents the structure of a customer record in the system.
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

// Map of US state codes for validation.
var states = map[string]bool{
	"AL": true, "AK": true, "AZ": true, "AR": true, "CA": true, "CO": true, "CT": true,
	"DE": true, "FL": true, "GA": true, "HI": true, "ID": true, "IL": true, "IN": true,
	"IA": true, "KS": true, "KY": true, "LA": true, "ME": true, "MD": true, "MA": true,
	"MI": true, "MN": true, "MS": true, "MO": true, "MT": true, "NE": true, "NV": true,
	"NH": true, "NJ": true, "NM": true, "NY": true, "NC": true, "ND": true, "OH": true,
	"OK": true, "OR": true, "PA": true, "RI": true, "SC": true, "SD": true, "TN": true,
	"TX": true, "UT": true, "VT": true, "VA": true, "WA": true, "WV": true, "WI": true,
	"WY": true, "DC": true, "PR": true,
}

/*
AddBookEndpoint handles the endpoint to add a new book to the system.
Validates the data and interacts with the internal load balancer to add the book.
*/
func AddBookEndpoint(w http.ResponseWriter, r *http.Request) {
	var bookAdded BookAdded
	var requestBody, _ = io.ReadAll(r.Body)

	// Unmarshal JSON body into BookAdded struct
	err := json.Unmarshal(requestBody, &bookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Validate book data
	status := verifyBookData(bookAdded)
	if status == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	// Verify request authorization
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Send the request to book service
	var responseBookAdded BookAdded
	var response = bookRequest(r, &bookAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &responseBookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Send success response with book information
	jsonResponseBody, _ := json.Marshal(&responseBookAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/books/"+responseBookAdded.ISBN)
	w.WriteHeader(response.StatusCode)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
UpdateBookEndpoint handles updating a book’s information in the system.
This will update the details of a book identified by its ISBN.
*/
func UpdateBookEndpoint(w http.ResponseWriter, r *http.Request) {
	var bookAdded BookAdded
	var requestBody, _ = io.ReadAll(r.Body)
	err := json.Unmarshal(requestBody, &bookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Validate book data
	status := verifyBookData(bookAdded)
	if status == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	// Verify request authorization
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Send the request to internal load balancer
	var responseBookAdded BookAdded
	var response = bookRequest(r, &bookAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &responseBookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	// Send success response with updated book information
	jsonResponseBody, _ := json.Marshal(&responseBookAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/books/"+responseBookAdded.ISBN)
	w.WriteHeader(response.StatusCode)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
verifyBookData validates the book data for required fields and formatting.
Returns 0 for valid data, 1 for invalid data.
*/
func verifyBookData(data BookAdded) int {
	if len(data.ISBN) == 0 || len(data.Title) == 0 || len(data.Author) == 0 ||
		len(data.Description) == 0 || len(data.Genre) == 0 {
		return 1
	}
	if data.Price <= 0 || data.Quantity <= 0 {
		return 1
	}

	// Validate price precision (up to two decimal places)
	d := decimal.NewFromFloat(data.Price)
	exp := d.Exponent()
	if exp != 0 && exp != -1 && exp != -2 {
		return 1
	}
	return 0
}

/*
RetrieveBookEndpoint retrieves a book by its ISBN from the system.
It interacts with the internal load balancer to fetch the book details.
*/
func RetrieveBookEndpoint(w http.ResponseWriter, r *http.Request) {
	var bookAdded BookAdded

	// Verify request authorization
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Send the request to internal load balancer
	var responseBookAdded BookAdded
	var response = bookRequest(r, &bookAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	err := json.Unmarshal(responseBody, &responseBookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Return the found book
	jsonResponseBody, _ := json.Marshal(&responseBookAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/books/"+responseBookAdded.ISBN)
	w.WriteHeader(response.StatusCode)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
RetrieveBooksKeywordEndpoint searches for books using a keyword in the database.
This query looks for the keyword anywhere in the book data.
*/
func RetrieveBooksKeywordEndpoint(w http.ResponseWriter, r *http.Request) {
	var bookAdded BookAdded

	// Verify request authorization
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Send the request to internal load balancer
	var responseBookAdded []BookAdded
	var response = bookRequest(r, &bookAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	err := json.Unmarshal(responseBody, &responseBookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Return books if found
	jsonResponseBody, _ := json.Marshal(&responseBookAdded)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
AddCustomerEndpoint handles adding a new customer to the system.
This endpoint validates customer data and adds them to the database.
*/
func AddCustomerEndpoint(w http.ResponseWriter, r *http.Request) {
	var customerAdded CustomerAdded
	var requestBody, err = io.ReadAll(r.Body)
	err = json.Unmarshal(requestBody, &customerAdded)
	if err != nil {
		panic(err.Error())
	}

	// Validate customer data
	status := verifyCustomerData(customerAdded)
	if status == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	// Verify request authorization
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Send the request to internal load balancer
	var response = customerRequest(r, &customerAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	var responseCustomer CustomerAdded
	err = json.Unmarshal(responseBody, &responseCustomer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Respond with customer ID
	id := responseCustomer.ID
	jsonResponseBody, _ := json.Marshal(&responseCustomer)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "http://"+r.Host+"/customers/"+id)
	w.WriteHeader(response.StatusCode)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
RetrieveRelatedBooksEndpoint fetches related books based on the current book's details.
This interacts with a recommendation engine.
*/
func RetrieveRelatedBooksEndpoint(w http.ResponseWriter, r *http.Request) {
	var bookAdded BookAdded

	// Verify request authorization
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Send request to internal load balancer
	var response = bookRequest(r, &bookAdded)
	var responseBody, _ = io.ReadAll(response.Body)

	// Return related books response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	_, err := w.Write([]byte(responseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
verifyCustomerData validates the customer data for required fields and formats.
Returns 0 for valid data, 1 for invalid data.
*/
func verifyCustomerData(data CustomerAdded) int {
	// Validate userId as a valid email address
	_, err := mail.ParseAddress(data.UserId)
	if err != nil {
		return 1
	}

	// Validate state as a valid US state abbreviation
	_, ok := states[data.State]
	if !ok {
		return 1
	}

	// All fields in the request body are mandatory except address2
	if len(data.Name) == 0 || len(data.Phone) == 0 || len(data.Address) == 0 ||
		len(data.City) == 0 || len(data.Zipcode) == 0 {
		return 1
	}

	return 0
}

/*
RetrieveCustomerEndpoint retrieves customer data by ID or userId.
It interacts with the internal load balancer to fetch customer details.
*/
func RetrieveCustomerEndpoint(w http.ResponseWriter, r *http.Request) {
	var customerAdded CustomerAdded

	// Verify request authorization
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Send request to internal load balancer
	var response = customerRequest(r, &customerAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	var responseCustomer CustomerAdded
	err := json.Unmarshal(responseBody, &responseCustomer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// Return customer data
	jsonResponseBody, _ := json.Marshal(&responseCustomer)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/customers/"+customerAdded.ID)
	w.WriteHeader(response.StatusCode)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
VerifyHeaderAuth verifies that the header contains valid authorization.
Returns 0 for valid, 1 for invalid.
*/
func VerifyHeaderAuth(r *http.Request) int {
	// convert headers to json object
	requestHeaders, _ := json.Marshal(r.Header)
	// for mapping json object
	headers := make(map[string]interface{})
	// map headers
	err := json.Unmarshal([]byte(requestHeaders), &headers)
	if err != nil {
		return 1
	}
	// return 1 if Authorization header is not present
	tokenHeader, status := r.Header["Authorization"]
	if !status {
		return 1
	}
	// get token string and decode
	tokenString := tokenHeader[0][len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(""), nil
	})

	// skip returning error if error is related to invalid token signature
	if err != nil {
		msg := err.Error()
		isInvalid := strings.Contains(msg, "token signature is invalid: signature is invalid")
		fmt.Println(isInvalid)
	}

	// get and map claims, return if claims could not be parsed
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 1
	}

	// return if claim sub does not exist
	sub, subExist := claims["sub"]
	if !subExist {
		return 1
	}

	// return if sub claim does not contain the following
	if sub != "starlord" && sub != "gamora" && sub != "drax" &&
		sub != "rocket" && sub != "groot" {
		return 1
	}

	// return if claim exp does not exist
	exp, expExist := claims["exp"]
	if !expExist {
		return 1
	}

	// get date if not identified as float then exit
	date, notExp := exp.(float64)
	if !notExp {
		return 1
	}

	// claim expiration date
	var claimDate int64 = int64(date)

	// current date
	var currentDate int64 = time.Now().Unix()

	// error if current date is greater than claim date
	if currentDate > claimDate {
		return 1
	}

	// return if claim iss does not exist
	iss, issExist := claims["iss"]
	if !issExist {
		return 1
	}

	// return if iss claim does not contain the following
	if iss != "cmu.edu" {
		return 1
	}

	return 0

}

/*
bookRequest is a helper function to send GET requests to load balancers or internal services.
*/
func bookRequest(r *http.Request, requestBookAdded *BookAdded) *http.Response {
	// send request, get response
	var response *http.Response
	method := r.Method

	switch method {
	case "GET":
		host := os.Getenv("GET_BOOKS")
		baseUrl := "http://" + host + ":3000"
		path := r.RequestURI
		url := baseUrl + path
		response, _ = http.Get(url)
	case "POST", "PUT":
		host := os.Getenv("ADD_BOOKS")
		baseUrl := "http://" + host + ":3000"
		path := r.URL.Path
		url := baseUrl + path
		jsonRequestBody, _ := json.Marshal(requestBookAdded)
		response, _ = http.Post(url, "application/json", bytes.NewBuffer(jsonRequestBody))
	}
	//
	return response
}

/*
customerRequest is a helper function to send GET requests to internal services for customer data.
*/
func customerRequest(r *http.Request, requestCustomerAdded *CustomerAdded) *http.Response {
	// send request, get response
	var response *http.Response
	method := r.Method
	host := os.Getenv("CUSTOMERS")
	baseUrl := "http://" + host + ":3000"
	path := r.RequestURI
	url := baseUrl + path

	switch method {
	case "GET":
		response, _ = http.Get(url)
	case "POST", "PUT":
		jsonRequestBody, _ := json.Marshal(requestCustomerAdded)
		response, _ = http.Post(url, "application/json", bytes.NewBuffer(jsonRequestBody))
	}
	//
	return response
}

// status In your K8S deployment file, specify a liveness probe
func statusEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	//_, err := w.Write([]byte("OKkkkk"))
	w.WriteHeader(200)
	_, err := w.Write([]byte("can you hear me"))

	if err != nil {
		log.Fatalf("Error sending response %v", err)
	}
}

/*
Main function initializes the API routes and starts the server.
*/
func main() {
	// Initialize and register routes
	r := mux.NewRouter()
	r.HandleFunc("/cmd/books", AddBookEndpoint).Methods("POST")
	r.HandleFunc("/cmd/books/{ISBN}", UpdateBookEndpoint).Methods("PUT")
	r.HandleFunc("/retrieveBook/{isbn}", RetrieveBookEndpoint).Methods("GET")
	r.HandleFunc("/retrieveBooks", RetrieveBooksKeywordEndpoint).Methods("GET")
	r.HandleFunc("/customers", AddCustomerEndpoint).Methods("POST")
	r.HandleFunc("/customers/{id}", RetrieveCustomerEndpoint).Methods("GET")
	r.HandleFunc("/customers", RetrieveCustomerEndpoint).Methods("GET")
	r.HandleFunc("/retrieveRelatedBooks/{isbn}", RetrieveRelatedBooksEndpoint).Methods("GET")
	// Add a /status route for liveness probe
	r.HandleFunc("/status", statusEndPoint).Methods("GET")

	// Start the server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
