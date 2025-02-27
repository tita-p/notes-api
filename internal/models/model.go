package model

import (
	"context"
	mongoDb "firstGoPro/internal/database"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type filter struct {
	key       string
	value     string
	operation string
}

var collectionName string
var dbClient *mongo.Client
var dbContext context.Context

func sliceToInterface[T any](slice []T, fn func(T) interface{}) []interface{} {
	var result []interface{}
	for _, v := range slice {
		result = append(result, fn(v))
	}
	return result
}

func init() {
	tags := []Tag{
		{Id: "1", Name: "Person 1"},
		{Id: "2", Name: "Person 2"},
	}

	notes := []Note{
		{Id: "1", Title: "Note #1", Content: "Content of the note #1", Tag: tags[0]},
		{Id: "2", Title: "Note #2", Content: "Content of the note #2", Tag: tags[1]},
	}

	dbClient = mongoDb.DbClient()
	dbContext = mongoDb.DbContext()

	dbClient.Database("db").Drop(dbContext)
	InsertNotes(notes)
	InsertTags(tags)
}

func InsertMany(collectionName string, items []interface{}) {
	database := dbClient.Database("db")
	collection := database.Collection(collectionName)

	insertedManyResult, err := collection.InsertMany(dbContext, items)

	if err != nil {
		log.Fatalf("inserted many error : %v", err)
		return
	}

	for _, doc := range insertedManyResult.InsertedIDs {
		fmt.Println(doc)
	}
}

func Find[T any](collectionName string, filters ...filter) []T {
	var searchFilters bson.D

	if len(filters) > 0 {
		for _, filter := range filters {
			searchFilters = append(searchFilters, bson.E{
				Key: filter.key,
				Value: bson.D{
					bson.E{
						Key:   filter.operation,
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

func GetFilter(key string, value string, operation string) *filter {
	return &filter{
		key:       key,
		value:     value,
		operation: operation,
	}
}
