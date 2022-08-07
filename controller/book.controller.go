package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kapbyte/book-app-tutorial/database"
	"github.com/kapbyte/book-app-tutorial/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bookCollection *mongo.Collection = database.OpenCollection(database.Client, "book")
var validate = validator.New()

func CreateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var book model.Book

		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(book)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		book.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		book.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		book.ID = primitive.NewObjectID()

		result, insertErr := bookCollection.InsertOne(ctx, book)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item was not created."})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func GetBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		bookId := c.Param("book_id")
		var book model.Book

		objectId, _ := primitive.ObjectIDFromHex(bookId)

		err := bookCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching book."})
			return
		}

		c.JSON(http.StatusOK, book)
	}
}

func UpdateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		bookId := c.Param("book_id")
		var book model.Book

		if err := c.BindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		objectId, _ := primitive.ObjectIDFromHex(bookId)
		filter := bson.M{"_id": objectId}

		var updateObj primitive.D

		if book.Author != nil {
			updateObj = append(updateObj, bson.E{Key: "author", Value: book.Author})
		}

		if book.Title != nil {
			updateObj = append(updateObj, bson.E{Key: "title", Value: book.Title})
		}

		if book.Description != nil {
			updateObj = append(updateObj, bson.E{Key: "description", Value: book.Description})
		}

		book.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: book.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err := bookCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Book item update failed."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Book item updated successfully."})
	}
}

func DeleteBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		bookId := c.Param("book_id")

		objectId, _ := primitive.ObjectIDFromHex(bookId)

		_, err := bookCollection.DeleteOne(ctx, bson.M{"_id": objectId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while deleting book item."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Book item deleted successfully."})
	}
}

func GetAllBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := bookCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching book list."})
			return
		}

		var allBooks []bson.M
		if err := result.All(ctx, &allBooks); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allBooks)
	}
}
