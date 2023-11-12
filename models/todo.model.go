package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	User        primitive.ObjectID `json:"user" bson:"user"`
	IsCompleted bool               `json:"isCompleted" bson:"isCompleted"`
}
