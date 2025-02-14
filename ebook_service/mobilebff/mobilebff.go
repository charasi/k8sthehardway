/*
	Simple BookStore Service that adds book/customer to a database
	and retrieves book by ISBN and customer by ID
	Charles Asiama
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
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

// CustomerAdded represents a customer record, mapped to JSON structure
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

// map of US state codes
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
		w.WriteHeader(401)
		return
	}

	// book data validation
	status := verifyBookData(bookAdded)
	if status == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	// if header is invalid/missing return 401 status
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// reroutes to internal load balancer
	var responseBookAdded BookAdded
	var response = getBookRequest(r, &bookAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &responseBookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// if book is added, send success response
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
UpdateBookEndpoint
Update a book’s information in the system.
The ISBN will be the unique identifier for the book.
*/
func UpdateBookEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to bookAdded
	var bookAdded BookAdded
	var requestBody, _ = io.ReadAll(r.Body)
	err := json.Unmarshal(requestBody, &bookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// book data validation
	status := verifyBookData(bookAdded)
	if status == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	// if header is invalid/missing return 401 status
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// reroutes to internal load balancer
	var responseBookAdded BookAdded
	var response = getBookRequest(r, &bookAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	err = json.Unmarshal(responseBody, &responseBookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// if book is updated, send success response
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
request body data validation
*/
func verifyBookData(data BookAdded) int {
	// all fields in the request body are mandatory
	if len(data.ISBN) == 0 || len(data.Title) == 0 || len(data.Author) == 0 ||
		len(data.Description) == 0 || len(data.Genre) == 0 {
		return 1
	}

	// all fields in the request body are mandatory
	if data.Price <= 0 || data.Quantity <= 0 {
		return 1
	}

	d := decimal.NewFromFloat(data.Price)
	exp := d.Exponent()
	if exp != 0 && exp != -1 && exp != -2 {
		return 1
	}
	return 0

}

/*
RetrieveBookEndpoint
return a book given its ISBN. Both endpoints shall produce the same response.
*/
func RetrieveBookEndpoint(w http.ResponseWriter, r *http.Request) {
	// // updated response body for book
	type UpdatedBookAdded struct {
		ISBN        string  `json:"ISBN"`
		Title       string  `json:"title"`
		Author      string  `json:"Author"`
		Description string  `json:"description"`
		Genre       int64   `json:"genre"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
	}

	// parse body to bookAdded
	var bookAdded BookAdded
	// to convert genre to int and set value to 3
	var updatedBookAdded UpdatedBookAdded
	// if header is invalid return 401 status
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// reroutes to internal load balancer
	var responseBookAdded BookAdded
	var response = getBookRequest(r, &bookAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	err := json.Unmarshal(responseBody, &responseBookAdded)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// replace in the response body the word “non-fiction” with the numeric value 3
	updatedBookAdded.Author = responseBookAdded.Author
	updatedBookAdded.Description = responseBookAdded.Description
	updatedBookAdded.Genre = 3
	updatedBookAdded.ISBN = responseBookAdded.ISBN
	updatedBookAdded.Price = responseBookAdded.Price
	updatedBookAdded.Quantity = responseBookAdded.Quantity
	updatedBookAdded.Title = responseBookAdded.Title

	// return book if found
	jsonResponseBody, _ := json.Marshal(&updatedBookAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/books/"+updatedBookAdded.ISBN)
	w.WriteHeader(response.StatusCode)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
AddCustomerEndpoint
Add a customer to the system.
This endpoint is called to create the newly registered customer in the system.
A unique numeric ID is generated for the new customer, and the customer is added to
the Customer data table on MySql (the numeric ID is the primary key).
*/
func AddCustomerEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to customerAdded
	var customerAdded CustomerAdded
	var requestBody, err = io.ReadAll(r.Body)
	err = json.Unmarshal(requestBody, &customerAdded)
	// exit if request body cannot be extracted
	if err != nil {
		panic(err.Error())
	}

	// customer data validation
	status := verifyCustomerData(customerAdded)
	if status == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	// if header is invalid return 401 status
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// reroutes to internal load balancer
	var response = getCustomerRequest(r, &customerAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	var responseCustomer CustomerAdded
	err = json.Unmarshal(responseBody, &responseCustomer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// respond with success status code
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
RetrieveRelatedBooksEndpoint
access an external recommendation engine service every time the “related books”
endpoint is executed and return recommendations for additional books the customer
may want to purchase.
*/
func RetrieveRelatedBooksEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to bookAdded
	var bookAdded BookAdded

	// if header is invalid return 401 status
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// response body for related books to return to customer
	//var responseRelatedBooks RelatedBooks
	// reroutes to internal load balancer
	var response = getBookRequest(r, &bookAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	//json.Unmarshal(responseBody, &responseRelatedBooks)

	// return book if found
	//jsonResponseBody, _ := json.Marshal(&responseRelatedBooks)
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Location", r.Host+"/books/"+responseBookAdded.ISBN)
	w.WriteHeader(response.StatusCode)
	_, err := w.Write([]byte(responseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
verifyCustomerData
request customer data validation
*/
func verifyCustomerData(data CustomerAdded) int {
	// userId must be a valid email address
	_, err := mail.ParseAddress(data.UserId)
	if err != nil {
		return 1
	}

	// state must be a valid 2-letter US state abbreviation
	_, ok := states[data.State]
	if !ok {
		return 1
	}
	// all fields in the request body are mandatory except address2
	if len(data.Name) == 0 || len(data.Phone) == 0 || len(data.Address) == 0 ||
		len(data.City) == 0 || len(data.Zipcode) == 0 {
		return 1
	}

	return 0

}

/*
RetrieveCustomerEndpoint
obtain the data for a customer given its numeric ID.
This endpoint will retrieve the customer data on MySql and send the data in the
response in JSON format. Note that ID is the  numeric ID, not the user-ID.

obtain the data for a customer given its user ID,which is the email address.
This endpoint will retrieve the customer data on MySql and send the data in the
response in JSON format. Note that the ‘@’ character should be encoded in the query
string parameter value (ex.: userId=starlord2002%40gmail.com).
*/
func RetrieveCustomerEndpoint(w http.ResponseWriter, r *http.Request) {
	// updated response body for customer
	type UpdatedCustomerAdded struct {
		ID     string `json:"id"`
		UserId string `json:"userId"`
		Name   string `json:"name"`
		Phone  string `json:"phone"`
	}
	// parse body to customerAdded
	var customerAdded CustomerAdded
	//
	var updatedCustomerAdded UpdatedCustomerAdded
	// if header is invalid return 401 status
	var success = VerifyHeaderAuth(r)
	if success == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// reroutes to internal load balancer
	var response = getCustomerRequest(r, &customerAdded)
	var responseBody, _ = io.ReadAll(response.Body)
	var responseCustomer CustomerAdded
	err := json.Unmarshal(responseBody, &responseCustomer)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		return
	}

	// update customer response
	updatedCustomerAdded.ID = responseCustomer.ID
	updatedCustomerAdded.Name = responseCustomer.Name
	updatedCustomerAdded.Phone = responseCustomer.Phone
	updatedCustomerAdded.UserId = responseCustomer.UserId

	// return customer record
	jsonResponseBody, _ := json.Marshal(&updatedCustomerAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/customers/"+updatedCustomerAdded.ID)
	w.WriteHeader(response.StatusCode)
	_, err = w.Write([]byte(jsonResponseBody))
	if err != nil {
		log.Fatalf("Error writing response from server: %v", err)
	}
}

/*
VerifyHeaderAuth
Verifies if request from client is from a web device and not mobile
returns 0 if client type is web, else 1
*/
func VerifyHeaderAuth(r *http.Request) int {
	// convert headers to json object
	requestHeaders, _ := json.Marshal(r.Header)
	// for mapping json object
	headers := make(map[string]interface{})
	// map headers
	err := json.Unmarshal([]byte(requestHeaders), &headers)
	if err != nil {

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
		fmt.Print(isInvalid)
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

	// return if claim does not contain the following
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
	var claimDate = int64(date)
	//claimDate := time.Unix(convertDate, 0)

	// current date
	var currentDate = time.Now().Unix()
	//currentDate := time.Unix(currentTime, 0)

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
send request to book service
*/
func getCustomerRequest(r *http.Request, requestCustomerAdded *CustomerAdded) *http.Response {
	// send request, get response
	var response *http.Response
	method := r.Method
	//baseUrl := "http://127.0.0.1:4000"
	baseUrl := "http://aa0dec50b455f46ffb742f085498278e-1773798306.us-east-1.elb.amazonaws.com:3000"
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

/*
send request to book service
*/
func getBookRequest(r *http.Request, requestBookAdded *BookAdded) *http.Response {
	// send request, get response
	var response *http.Response
	method := r.Method
	//baseUrl := "http://127.0.0.1:3000"

	switch method {
	case "GET":
		baseUrl := "http://ac8dc363275454af9932233199b0fd78-392305053.us-east-1.elb.amazonaws.com:3000"
		path := r.RequestURI
		url := baseUrl + path
		response, _ = http.Get(url)
	case "POST", "PUT":
		baseUrl := "http://a953e35b792784b378087d1c5f8c6e05-1741789884.us-east-1.elb.amazonaws.com:3000"
		//path := r.URL.Path
		path := r.RequestURI
		url := baseUrl + path
		jsonRequestBody, _ := json.Marshal(requestBookAdded)
		response, _ = http.Post(url, "application/json", bytes.NewBuffer(jsonRequestBody))
	}
	//
	return response
}

//	monitor the health of the REST service within EKS.
//
// In your K8S deployment file, specify a liveness probe
func statusEndpoint(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("/cmd/books/{ISBN}", UpdateBookEndpoint).Methods("PUT")
	router.HandleFunc("/books/isbn/{ISBN}", RetrieveBookEndpoint).Methods("GET")
	router.HandleFunc("/books/{ISBN}", RetrieveBookEndpoint).Methods("GET")
	router.HandleFunc("/books", RetrieveBookEndpoint).Queries(
		"keyword", "{keyword}").Methods("GET")
	router.HandleFunc("/customers", RetrieveCustomerEndpoint).Methods("GET")
	router.HandleFunc("/customers", AddCustomerEndpoint).Methods("POST")
	router.HandleFunc("/customers/{id}", RetrieveCustomerEndpoint).Methods("GET")
	router.HandleFunc("/books/{ISBN}/related-books", RetrieveRelatedBooksEndpoint).Methods("GET")
	router.HandleFunc("/status", statusEndpoint).Methods("GET")
	//http.ListenAndServe(":2345", router)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
