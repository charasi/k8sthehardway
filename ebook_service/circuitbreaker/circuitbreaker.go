/*
	Circuit breaker for the interaction with the recommendation service
	Expects response within 3 seconds, or sends a 504/503 respond
	Charles Asiama
*/

package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

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

var circBreakerIsOpen bool = false
var errorStatus = 504

/*
access an external recommendation engine service every time the “related books”
endpoint is executed and return recommendations for additional books the customer
may want to purchase.
*/
func RetrieveRelatedBooksEndpoint(w http.ResponseWriter, r *http.Request) {

	// When the circuit is open the service automatically
	// returns an error (503) right away to any requests
	// within 60 seconds of opening the circuit
	if circBreakerIsOpen {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorStatus)
		w.Write([]byte(""))
		return
	}

	// reroutes to external service
	var response = getRelatedBookRequest(r)
	// if request timed out, send error status, set breaker to true
	// and set timer to 60 secondds
	if response == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errorStatus)
		w.Write([]byte(""))
		circBreakerIsOpen = true
		errorStatus = 503
		//
		go setTimer()
		return
	}

	// if book is added, send success response
	var responseBody, _ = io.ReadAll(response.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write([]byte(responseBody))
}

func setTimer() {
	// set time to 60 seconds
	timer := time.NewTimer(60000 * time.Millisecond)
	<-timer.C
	fmt.Println("Timer expired")
	circBreakerIsOpen = false
}

/*
send request to book service
*/
func getRelatedBookRequest(r *http.Request) *http.Response {
	// send request, get response
	var response *http.Response
	baseUrl := "http://52.72.198.36:80/recommended-titles/isbn/"
	path := r.URL.Path
	paths := strings.Split(path, "/")
	isbn := paths[2]
	fmt.Print(paths)
	//url := baseUrl + path
	url := baseUrl + isbn

	client := http.Client{
		Timeout: 3000 * time.Millisecond,
	}

	response, err := client.Get(url)

	if err != nil {
		e, _ := err.(net.Error)
		if e.Timeout() {
			return nil
		} else {
			fmt.Print("Different ERROR! PANIC")
			//panic(err.Error())
		}
	}

	return response
}

/*
Main function to handle routes/path for web server
*/
func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/books/{ISBN}/related-books", RetrieveRelatedBooksEndpoint).Methods("GET")
	http.ListenAndServe(":2345", mux)
}
