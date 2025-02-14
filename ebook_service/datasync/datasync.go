/*
	Cron Job that queries Books from
	the RDS table and store them as JSON documents on MongoDB
	Charles Asiama
*/

package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
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

/*
Cron job to synchronize between the RDS and mongoDB
*/
func main() {

	// for connection and interacting with mongoDB
	var collection *mongo.Collection
	var ctx = context.TODO()

	// connect to mongoDB
	clientOptions := options.Client().ApplyURI("mongodb+srv://casiama:RIO1yYeZnijK4pJR@assignment4.yvbw3da.mongodb.net/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	collection = client.Database("BooksDB").Collection("books_casiama")

	db, err := sql.Open("mysql", "awsadmin:awspassword@tcp(assn4-dbauroraa-esuvdyingzvs.c42qw4pddowd.us-east-1.rds.amazonaws.com:3306)/bookcmd")
	//db, err := sql.Open("mysql", "charasi:Skittles05@10@tcp(127.0.0.1:3306)/bookcmd")
	if err != nil {
		panic(err.Error())
	}

	// close database after exiting function
	defer db.Close()
	//var bookAddedList []BookAdded
	var bookAdded BookAdded
	query := "SELECT * FROM book"
	// get record, return error if book does not exist
	results, resultErr := db.Query(query)
	fmt.Println(results)

	if resultErr == sql.ErrNoRows {
		return
	}

	// update or insert each record to mongoDB
	for results.Next() {
		results.Scan(&bookAdded.ISBN, &bookAdded.Title, &bookAdded.Author,
			&bookAdded.Description, &bookAdded.Genre, &bookAdded.Price, &bookAdded.Quantity)

		// isbn to update or insert
		isbn := bookAdded.ISBN
		// search in mongoDB for isbn
		filterDoc := bson.D{{Key: "ISBN", Value: isbn}}
		// construct document
		updateDoc := bson.D{{Key: "$set", Value: bson.D{{Key: "ISBN", Value: bookAdded.ISBN},
			{Key: "title", Value: bookAdded.Title}, {Key: "Author", Value: bookAdded.Author},
			{Key: "description", Value: bookAdded.Description},
			{Key: "genre", Value: bookAdded.Genre}, {Key: "price", Value: bookAdded.Price},
			{Key: "quantity", Value: bookAdded.Quantity}}}}

		// set upsert to true
		options := options.Update().SetUpsert(true)
		// update or insert with document
		collection.UpdateOne(context.TODO(), filterDoc, updateDoc, options)
	}
}
