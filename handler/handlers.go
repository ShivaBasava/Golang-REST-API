package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB : This is handler function to connect 'mongoDB'
func ConnectDB() *mongo.Collection {
	// Setting-up client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")

	// Connecting to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
    //For every request, It prompts a response to Console (Just for checking the connectivity to MondoDB
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("go_rest_api").Collection("cartoons")

	return collection
}

// ErrorResponse : This is error model/Structure.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError : This function is to prepare error model.
func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
