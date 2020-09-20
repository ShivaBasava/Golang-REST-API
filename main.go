package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	handlers "github.com/ShivaBasava/Golang-REST-API/handler"
	"github.com/ShivaBasava/Golang-REST-API/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getCartoons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Cartoon array
	var cartoons []model.Cartoon

	//Connection mongoDB
	collection := handlers.ConnectDB()

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		handlers.GetError(err, w)
		return
	}

	// Closes the cursor once finished
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var cartoon model.Cartoon
		// & character returns the memory address of the following variable.
		err := cur.Decode(&cartoon) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		cartoons = append(cartoons, cartoon)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(cartoons) // encode similar to serialize process.
}

func getCartoon(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var cartoon model.Cartoon
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := handlers.ConnectDB()

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&cartoon)

	if err != nil {
		handlers.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(cartoon)
}

func createCartoon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cartoon model.Cartoon

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&cartoon)

	// connect db
	collection := handlers.ConnectDB()

	// insert our cartoon model.
	result, err := collection.InsertOne(context.TODO(), cartoon)

	if err != nil {
		handlers.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateCartoon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var cartoon model.Cartoon

	collection := handlers.ConnectDB()

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&cartoon)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"title", cartoon.Title},
			{"genre", cartoon.Genre},
			{"director", bson.D{
				{"firstname", cartoon.Director.FirstName},
				{"lastname", cartoon.Director.LastName},
			}},
			{"seasons", bson.D{
				{"season_no", cartoon.Seasons.Season_No},
				{"total_episodes", cartoon.Seasons.Total_Episodes},
			}},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&cartoon)

	if err != nil {
		handlers.GetError(err, w)
		return
	}

	cartoon.ID = id

	json.NewEncoder(w).Encode(cartoon)
}

func deleteCartoon(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	collection := handlers.ConnectDB()

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		handlers.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func main() {

	// The router is now formed by calling the `NewRouter` constructor function
	r := mux.NewRouter()

	
	r.HandleFunc("/api/cartoons", getCartoons).Methods("GET")
	r.HandleFunc("/api/cartoons/{id}", getCartoon).Methods("GET")
	r.HandleFunc("/api/cartoons", createCartoon).Methods("POST")
	r.HandleFunc("/api/cartoons/{id}", updateCartoon).Methods("PUT")
	r.HandleFunc("/api/cartoons/{id}", deleteCartoon).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
