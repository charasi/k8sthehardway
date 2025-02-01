/*
	Simple BookStore Service that adds book/customer to a database
	and retrieves book by ISBN and customer by ID
	Charles Asiama
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

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

/*
Add a customer to the system.
This endpoint is called to create the newly registered customer in the system.
A unique numeric ID is generated for the new customer, and the customer is added to
the Customer data table on MySql (the numeric ID is the primary key).
*/
func AddCustomerEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to customerAdded
	var customerAdded CustomerAdded
	var requestBody, _ = io.ReadAll(r.Body)
	json.Unmarshal(requestBody, &customerAdded)

	cus, _ := json.Marshal(&customerAdded)
	KafkaProducer(cus)

	// query for customer
	query := "INSERT INTO customers (userId, name, phone, address, address2, city, state, zipcode)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	// arguments for customer
	args := []interface{}{customerAdded.UserId, customerAdded.Name,
		customerAdded.Phone, customerAdded.Address, customerAdded.Address2,
		customerAdded.City, customerAdded.State, customerAdded.Zipcode}

	// get customer record
	result := AddCustomerTable(query, args)
	// respond with duplicate message if customer exist
	if result == -1 {
		var message Duplicate
		message.Message = "This user ID already exists in the system."
		jsonMessage, _ := json.Marshal(&message)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(422)
		w.Write(jsonMessage)
		return
	}

	// respond with success status code
	id := strconv.Itoa(result)
	customerAdded.ID = id
	jsonResponseBody, _ := json.Marshal(&customerAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "http://"+r.Host+"/customers/"+id)
	w.WriteHeader(201)
	w.Write([]byte(jsonResponseBody))

	// send to kafka broker
	//cus, _ := json.Marshal(&customerAdded)
	//KafkaProducer(cus)
}

/*
Add/Update customer record
*/
func AddCustomerTable(query string, args []interface{}) int {
	db, err := sql.Open("mysql", "awsadmin:awspassword@tcp(assn4-dbauroraa-esuvdyingzvs.c42qw4pddowd.us-east-1.rds.amazonaws.com:3306)/bookstore")
	//db, err := sql.Open("mysql", "charasi:Skittles05@10@tcp(127.0.0.1:3306)/bookstore")
	if err != nil {
		panic(err.Error())
	}

	// close database after exiting function
	defer db.Close()

	// get record, return error if customer already exists
	_, err = db.Exec(query, args...)
	if err != nil {
		message := err.Error()
		isDuplicate := strings.Contains(message, "Duplicate entry")
		if isDuplicate {
			return -1
		}
		panic(err.Error())
	}

	// return customer ID
	var id int
	records := db.QueryRow("SELECT id from customers WHERE userId = ?",
		args[0]).Scan(&id)
	if records == sql.ErrNoRows {
		return -2
	} else {
		return id
	}

}

/*
obtain the data for a customer given its numeric ID.
This endpoint will retrieve the customer data on MySql and send the data in the
response in JSON format. Note that ID is the  numeric ID, not the user-ID.

obtain the data for a customer given its user ID,which is the email address.
This endpoint will retrieve the customer data on MySql and send the data in the
response in JSON format. Note that the ‘@’ character should be encoded in the query
string parameter value (ex.: userId=starlord2002%40gmail.com).
*/
func RetrieveCustomerEndpoint(w http.ResponseWriter, r *http.Request) {
	// parse body to customerAdded
	var customerAdded CustomerAdded
	var urlPath string
	var query string
	var status int

	// determine to query customer by id or userId
	if r.URL.RawQuery == "" {
		urlPath = path.Base(r.URL.Path)
		query = "SELECT * FROM customers where id = ?"
		status = verifyRetrieveScenarios("id", urlPath)
	} else {
		tempPath := path.Base(r.URL.RawQuery)
		arrQuery := strings.Split(tempPath, "=")
		urlPath = arrQuery[1]
		urlPath = strings.Replace(urlPath, "%", "@", 1)
		arrQuery = strings.Split(urlPath, "@")
		if len(arrQuery) > 1 {
			reg := regexp.MustCompile(`\d`)
			domain := reg.ReplaceAllString(arrQuery[1], "")
			urlPath = arrQuery[0] + "@" + domain
		}
		query = "SELECT * FROM customers where userId = ?"
		status = verifyRetrieveScenarios("userId", urlPath)
	}

	if status == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		return
	}

	// get customer record, return error if not found
	result := GetCustomerRecord(query, urlPath, &customerAdded)
	if result == 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		return
	}

	// return customer record
	jsonResponseBody, _ := json.Marshal(&customerAdded)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", r.Host+"/customers/"+customerAdded.ID)
	w.WriteHeader(200)
	w.Write([]byte(jsonResponseBody))
}

/*
retrieves customer record
*/
func GetCustomerRecord(query string, key string, customerAdded *CustomerAdded) int {
	db, err := sql.Open("mysql", "awsadmin:awspassword@tcp(assn4-dbauroraa-esuvdyingzvs.c42qw4pddowd.us-east-1.rds.amazonaws.com:3306)/bookstore")
	//db, err := sql.Open("mysql", "charasi:Skittles05@10@tcp(127.0.0.1:3306)/bookstore")

	if err != nil {
		panic(err.Error())
	}

	// close database after exiting function
	defer db.Close()

	// query customer record, return error if not found
	result := db.QueryRow(query, key).Scan(&customerAdded.ID, &customerAdded.UserId,
		&customerAdded.Name, &customerAdded.Phone, &customerAdded.Address,
		&customerAdded.Address2, &customerAdded.City, &customerAdded.State,
		&customerAdded.Zipcode)

	if result == sql.ErrNoRows {
		return 1
	} else {
		return 0
	}
}

/*
Additional verification scenarios
*/
func verifyRetrieveScenarios(scenario string, key string) int {
	// convert to number and check if its negative
	if scenario == "id" {
		num, _ := strconv.Atoi(key)
		if num <= 0 {
			return 1
		}
	} else {
		// userId must be a valid email address
		_, err := mail.ParseAddress(key)
		if err != nil {
			return 1
		}
	}

	return 0
}

/*
*
kafka producer that sends messages to subscribed topics
*/
func KafkaProducer(jsonResponseBody []byte) {
	// topic
	const (
		KafkaTopic = "casiama.customer.evt"
	)

	// server brokers
	brokers := [3]string{"52.72.198.36:9092", "54.224.217.168:9092", "44.208.221.62:9092"}

	// send messages to the 3 server brokers
	for i := 0; i < 3; i++ {
		broker := brokers[i]

		p, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": broker,
		})
		if err != nil {
			panic(err)
		}

		// handler for produced messages
		go func() {
			for e := range p.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
					} else {
						fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
					}
				}
			}
		}()

		// send message to topics
		topic := KafkaTopic

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          jsonResponseBody,
		}, nil)

		if err != nil {
			panic(err)
		}

		// Wait for message delivery and close
		p.Flush(15000)
		p.Close()
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
	mux.HandleFunc("/customers", RetrieveCustomerEndpoint).Methods("GET")
	mux.HandleFunc("/customers", AddCustomerEndpoint).Methods("POST")
	mux.HandleFunc("/customers/{id}", RetrieveCustomerEndpoint).Methods("GET")
	mux.HandleFunc("/status", status).Methods("GET")
	//http.ListenAndServe(":2345", mux)
	http.ListenAndServe(":3000", mux)
}
