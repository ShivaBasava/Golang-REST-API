package handlers

import (
    "encoding/json"
    "fmt"
    "context"
    "log"
    "net/http"
    
    "github.com/ShivaBasava/Golang-REST-API/models"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"    
)

// ConnectDB : This is handler function to connect 'mongoDB'
func ConnectDB() *mongo.Collection {
    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")

    collection := client.Database("go_rest_api").Collection("cartoons")

    return collection
}


// ErrorResponse : This is error model.
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


func getCartoons(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // we created Cartoon array
    var cartoons []models.Cartoon

    //Connection mongoDB 
    collection := ConnectDB()

    // bson.M{},  we passed empty filter. So we want to get all data.
    cur, err := collection.Find(context.TODO(), bson.M{})

    if err != nil {
        GetError(err, w)
        return
    }

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var cartoon models.Cartoon
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

	var cartoon models.Cartoon
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := ConnectDB()

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&cartoon)

	if err != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(cartoon)
}

func createCartoon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cartoon models.Cartoon

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&cartoon)

	// connect db
	collection := ConnectDB()

	// insert our cartoon model.
	result, err := collection.InsertOne(context.TODO(), cartoon)

	if err != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateCartoon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var cartoon models.Cartoon

	collection := ConnectDB()

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&cartoon)

    // prepare update model.
    update := bson.D{
        {"$set", bson.D{
            {"title", cartoon.Title},
            {"genre", cartoon.Genre},
            {"director", cartoon.D{
                {"firstname", cartoon.Director.FirstName},
                {"lastname", cartoon.Director.LastName},
            }},
            {"seasons", cartoon.D{
                {"season_no", cartoon.Seasons.Season_No},
                {"total_episodes", cartoon.Seasons.Total_Episodes},
            }},            
        }},
    }

    err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&cartoon)

    if err != nil {
        GetError(err, w)
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

	collection := ConnectDB()

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func hello(w http.ResponseWriter, r *http.Request) {
    msg := {"Message":"Hello, World!"}
    res, err := Decode(msg)

    if err != nil {
        GetError(err, w)
        return
    }

    json.NewEncoder(w).Encode(res)
}
