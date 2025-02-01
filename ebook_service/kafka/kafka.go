/**
CRM kafka consumer service that connects to the kafka topic casiama.customer.evt
Upon receiving the message, the CRM service shall parse the content and
send an email to the newly registered customer.
Charles Asiama
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// server email and port
var host string = "smtp.gmail.com"
var port string = "587"

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

// group id, and topic
const (
	KafkaGroupId = "customer-service"
	KafkaTopic   = "casiama.customer.evt"
)

/*
kafka consumer service
loop forever continously wating to for new messages from subscribed topic
*/
func main() {

	// three server brokers
	brokers := [3]string{"52.72.198.36:9092", "54.224.217.168:9092", "44.208.221.62:9092"}

	// loop continously on lsitening  each server broker for new messages on subscribed topic
	for {
		for i := 0; i < 3; i++ {
			broker := brokers[i]
			// subrice to topic
			consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": broker,
				"group.id":          KafkaGroupId,
				"auto.offset.reset": "earliest",
			})

			if err != nil {
				panic(err)
			}

			// topic you are subcribed to
			topic := KafkaTopic
			consumer.SubscribeTopics([]string{topic}, nil)

			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				var customer CustomerAdded
				err := json.Unmarshal(msg.Value, &customer)
				if err != nil {
					fmt.Printf("Error message: %v\n", err)
					continue
				}

				// message to email to customer
				cusName := customer.Name
				cusEmail := customer.UserId
				toList := []string{cusEmail}
				from := "asiamacharles29@gmail.com"
				password := "wgwm lojs gmjt mgsj"
				// This is the message to send in the mail
				msg := "Dear " + cusName + ",\n" + "Welcome to the Book store created by casiama.\n" +
					"Exceptionally this time we won't ask you to click a link to activate your account.\n"

				body := []byte("Subject:Activate your book store account\n" + msg)

				// verify user
				auth := smtp.PlainAuth("", from, password, host)

				// Send mail
				er := smtp.SendMail(host+":"+port, auth, from, toList, body)

				// handle errors
				if er != nil {
					fmt.Println(err)
				}

			} else {
				fmt.Printf("Error: %v\n", err)
			}

			consumer.Close()
		}
	}
}
