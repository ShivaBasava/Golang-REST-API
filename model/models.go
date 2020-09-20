package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//Creating a Structure for 'Cartoon'
type Cartoon struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title,omitempty"`
	Genre    string             `json:"genre" bson:"genre,omitempty"`
	Director *Director          `json:"director"  bson:"director,omitempty"`
	Seasons  *Seasons           `json:"seasons" bson:"seasons,omitempty"`
}

//Creating a Structure for 'Director'
//Currently Handling only one set of Name
type Director struct {
	FirstName string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

//Creating a Structure for 'Seasons'
type Seasons struct {
	Season_No      int `json:"season_No" bson:"season_no,omitempty"`
	Total_Episodes int `json:"total_episodes" bson:"total_episodes,omitempty"`
}
