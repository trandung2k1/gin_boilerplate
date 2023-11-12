package controllers

import (
	"context"
	"fmt"
	"ginboilerplate/config"
	"ginboilerplate/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = config.GetCollection(config.DB, "users")

func hashPassword(password string) string {
	rounds, _ := strconv.Atoi(os.Getenv("ROUNDS"))
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), rounds)
	return string(hash)
}
func comparePassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func RegisterAccount(c *gin.Context) {
	var newUser models.User
	newUser.Id = primitive.NewObjectID()
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	hashedPassword := hashPassword(newUser.Password)
	newUser.Password = hashedPassword
	indexEmail, err := userCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(indexEmail)
	if indexEmail != "" {
		c.JSON(http.StatusCreated, gin.H{
			"status":  "failure",
			"message": indexEmail + " is created",
		})
		return
	}

	res, err := userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Fatal(err)
		return
	}
	id := res.InsertedID
	fmt.Println(res)
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"_id":    id,
	})
}

func Login(c *gin.Context) {
	var user models.Login
	if err := c.BindJSON(&user); err != nil {
		return
	}
	var u models.User
	filter := bson.D{{Key: "email", Value: user.Email}}
	userCollection.FindOne(context.TODO(), filter).Decode(&u)

	if u.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failure",
			"message": "Invalid email!",
		})
		return
	} else {
		e := comparePassword(u.Password, user.Password)
		if e != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  "failure",
				"message": "Password is not valid!",
			})
			return
		} else {
			userId := string(u.Id.Hex())
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["userId"] = userId

			// claims := &jwt.RegisteredClaims{
			// 	ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			// 	IssuedAt:  jwt.NewNumericDate(time.Now()),
			// 	NotBefore: jwt.NewNumericDate(time.Now()),
			// 	ID:        userId,
			// }

			claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
			mysecret := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
			tokenString, _ := token.SignedString(mysecret)
			var result models.LoginSuccess
			result.Id = u.Id
			result.Username = u.Username
			result.Email = u.Email
			result.Token = tokenString
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   result,
			})
		}

	}
}

func Logout(c *gin.Context) {

}

func DeleteAccount(c *gin.Context) {

}

func UpdateAccount(c *gin.Context) {

}
func ResetPassword(c *gin.Context) {

}

func GetUserById(c *gin.Context) {
	fmt.Println(c.Request.Header.Get("userId"))
	var id = c.Param("id")
	userIdObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": "Id invalid",
		})
		return

	}
	var result models.User
	filter := bson.D{{Key: "_id", Value: userIdObject}}
	err = userCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	if result.Username == "" && result.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failure",
			"message": "User not found!",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"user":   result,
		})
		return
	}
}
