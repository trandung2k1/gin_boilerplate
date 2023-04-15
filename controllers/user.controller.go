package controllers

import (
	"context"
	"ginboilerplate/config"
	"ginboilerplate/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = config.GetCollection(config.DB, "users")

func GetUserById(c *gin.Context) {
	var id = c.Param("id")
	userIdObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": "Id invalid",
		})

	}
	var result models.User
	err = userCollection.FindOne(context.TODO(), bson.D{{"_id", userIdObject}}).Decode(&result)
	if result.Username == "" && result.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failure",
			"message": "User not found!",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"user":   result,
		})
	}
}
