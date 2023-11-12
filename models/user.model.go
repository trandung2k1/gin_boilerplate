package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID   `json:"id" bson:"_id"`
	Username  string               `json:"username"`
	Email     string               `json:"email"`
	Password  string               `json:"password"`
	Avatar    string               `json:"avatar"`
	IsAdmin   bool                 `json:"isAdmin" bson:"isAdmin"`
	IsBlocked bool                 `json:"isBlocked" bson:"isBlocked"`
	Todos     []primitive.ObjectID `json:"todos"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginSuccess struct {
	Id       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Token    string             `json:"token"`
}
