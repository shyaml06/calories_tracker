package routes

import (
	"calorie_tracker_backend/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var entryCollection *mongo.Collection = OpenColeection(Client, "calories")
var validate = validator.New()

func AddEntry(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry

	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return

	}
	//validate := validator.New()

	validationerr := validate.Struct(entry)
	if validationerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationerr.Error()})
		fmt.Println(validationerr)
		return

	}

	entry.ID = primitive.NewObjectID()
	result, inserterror := entryCollection.InsertOne(ctx, entry)
	if inserterror != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": inserterror.Error()})
		fmt.Println(inserterror)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)

}
func GetEntries(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &entries); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return

	}
	defer cancel()
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)

}

func GetEntriesByIngredient(c *gin.Context) {

	ingredient := c.Params.ByName("id")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M

	cursor, err := entryCollection.Find(ctx, bson.M{"ingredient": ingredient})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return

	}
	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return

	}
	defer cancel()
	fmt.Println(entries)

	c.JSON(http.StatusOK, entries)

}
func GetEntryById(c *gin.Context) {
	EntyID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(EntyID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry bson.M
	if err := entryCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(entry)
	c.JSON(http.StatusOK, entry)

}
func UpdateEntry(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry
	if err := entryCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	//validate := validator.New()

	validationerr := validate.Struct(entry)
	if validationerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationerr.Error()})
		fmt.Println(validationerr)
		return

	}
	result, err := entryCollection.ReplaceOne(ctx, bson.M{"_id": docID}, bson.M{
		"dish":        entry.Dish,
		"fat":         entry.Fat,
		"ingredients": entry.Ingredients,
		"calories":    entry.Calories,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return

	}
	defer cancel()
	c.JSON(http.StatusOK, result)

}
func UpdateIngredient(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	type Ingredient struct {
		Ingredients *string `json:"ingredients"`
	}
	var ingredient Ingredient
	if err := c.BindJSON(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return

	}
	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID},
		bson.D{{"$set", bson.D{{"ingredients", ingredient.Ingredients}}}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result.ModifiedCount)

}
func DeleteEntry(c *gin.Context) {

	entryID := c.Params.ByName("id")

	docID, _ := primitive.ObjectIDFromHex(entryID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)

}
func Land(c *gin.Context) {
	c.JSON(http.StatusOK, 200)

}
