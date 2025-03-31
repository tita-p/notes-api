package model

import (
	"context"
	"fmt"
	"log"

	mongoDb "github.com/tita-p/notes-api/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type filter struct {
	key      string
	value    interface{}
	operator string
}

var dbClient *mongo.Client
var dbContext context.Context

func Filter(key string, operator string, value interface{}) *filter {
	var mongoDbOperator string

	switch operator {
	case "in":
		mongoDbOperator = "$in"
	case "=":
	default:
		mongoDbOperator = "$eq"
	}

	return &filter{
		key:      key,
		value:    value,
		operator: mongoDbOperator,
	}
}

func sliceToInterface[T any](slice []*T, fn func(T) interface{}) []interface{} {
	var result []interface{}
	for _, v := range slice {
		result = append(result, fn(*v))
	}
	return result
}

func init() {
	dbClient = mongoDb.DbClient()
	dbContext = mongoDb.DbContext()
}

func resetDatabase() {
	tags := []*Tag{
		{Name: "Person 1"},
		{Name: "Person 2"},
	}

	dbClient.Database("db").Collection("notes").Drop(dbContext)
	dbClient.Database("db").Collection("tags").Drop(dbContext)

	insertedTags := InsertTags(tags)

	firstTag := &insertedTags[0]

	notes := []*Note{
		{Title: "Note #1", Content: "Content of the note #1", TagId: firstTag.Id},
		{Title: "Note #2", Content: "Content of the note #2", TagId: firstTag.Id},
	}

	InsertNotes(notes)
}

func InsertMany[T any](collectionName string, items []interface{}) []T {
	database := dbClient.Database("db")
	collection := database.Collection(collectionName)

	insertedManyResult, err := collection.InsertMany(dbContext, items)

	if err != nil {
		log.Fatalf("inserted many error : %v", err)
	}

	var insertedItems []T

	for _, id := range insertedManyResult.InsertedIDs {
		if objectID, ok := id.(primitive.ObjectID); ok {
			insertedData, notFound := FindById[T](collectionName, objectID)

			if !notFound {
				insertedItems = append(insertedItems, insertedData)
			}
		}
	}

	return insertedItems
}

func FindById[T any](collectionName string, id primitive.ObjectID) (T, bool) {
	database := dbClient.Database("db")
	collection := database.Collection(collectionName)

	var item T

	filter := bson.D{{Key: "_id", Value: id}}

	err := collection.FindOne(dbContext, filter).Decode(&item)

	notFound := false

	if err != nil {
		if err == mongo.ErrNoDocuments {
			notFound = true
		} else {
			log.Fatal("Error finding item:", err)
		}
	} else {
		fmt.Printf("Found item: %+v\n", item)
	}

	return item, notFound
}

func Find[T any](collectionName string, filters ...*filter) []T {
	var searchFilters bson.D

	if len(filters) > 0 {
		for _, filter := range filters {
			searchFilters = append(searchFilters, bson.E{
				Key: filter.key,
				Value: bson.D{
					bson.E{
						Key:   filter.operator,
						Value: filter.value,
					},
				},
			})
		}
	}

	database := dbClient.Database("db")
	collection := database.Collection(collectionName)

	var cursor *mongo.Cursor
	var err error

	if len(searchFilters) > 0 {
		cursor, err = collection.Find(dbContext, searchFilters)
	} else {
		cursor, err = collection.Find(dbContext, options.Find())
	}

	if err != nil {
		log.Fatalf("Find collection err : %v", err)
	}

	var results []T

	for cursor.Next(context.TODO()) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			log.Fatalf("Error decoding document: %v", err)
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}

	return results
}

func Update[T any](collectionName string, id primitive.ObjectID, item interface{}) {
	collection := dbClient.Database("db").Collection(collectionName)

	filter := bson.M{"_id": id}

	updateManyResult, err := collection.UpdateOne(
		dbContext,
		filter,
		item,
	)
	if err != nil {
		log.Fatalf("update error : %v", err)
		return
	}

	fmt.Println("========= updated modified count ===========")
	fmt.Println(updateManyResult.ModifiedCount)
}

func Delete(collectionName string, id primitive.ObjectID) {
	collection := dbClient.Database("db").Collection(collectionName)

	filter := bson.D{{Key: "_id", Value: id}}
	deleteManyResult, err := collection.DeleteOne(dbContext, filter)

	if err != nil {
		log.Fatalf("delete many data error : %v", err)
		return
	}
	fmt.Println("===== delete many data modified count =====")
	fmt.Println(deleteManyResult.DeletedCount)
}

func stringToObjectId(id string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	return objectID
}
