package controllers

import (
	"context"
	"fmt"
	"ginboilerplate/config"
	"ginboilerplate/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var todoCollection = config.GetCollection(config.DB, "todos")

func GetAllTodo(c *gin.Context) {
	fmt.Println(c.Request.Header.Get("userId"))
	cur, err := todoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	var results []models.Todo
	if err = cur.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   results,
	})
}

func CreateTodo(c *gin.Context) {
	var userId = c.Request.Header.Get("userId")
	var newTodo models.Todo
	newTodo.Id = primitive.NewObjectID()
	userIdObject, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": "Id invalid",
		})
		return
	}
	filter := bson.D{{Key: "_id", Value: userIdObject}}
	var user models.User
	userCollection.FindOne(context.TODO(), filter).Decode(&user)

	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}
	newTodo.User = userIdObject
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err,
		})
		return
	}
	res, err := todoCollection.InsertOne(context.Background(), newTodo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err,
		})
		return
	}

	id := res.InsertedID

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"_id":    id,
	})
}

func GetTodoById(c *gin.Context) {
	var id = c.Param("id")
	todoIdObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": "Id invalid",
		})
		return

	}
	filter := bson.D{{Key: "_id", Value: todoIdObject}}
	var result models.Todo
	err = todoCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	if result.Title == "" {
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"status": "failure",
		"data":   result,
	})
}
