package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoDb struct {
	client  *mongo.Client
	context context.Context
	cancel  context.CancelFunc
}

var db *mongoDb

func init() {
	context, cancel := context.WithCancel(context.Background())

	client, err := mongo.Connect(
		context,
		options.Client().ApplyURI("mongodb://admin:secret@localhost:27017/"),
	)

	if err != nil {
		log.Fatalf("connection error :%v", err)
	}

	err = client.Ping(context, readpref.Primary())

	if err != nil {
		log.Fatalf("ping mongodb error :%v", err)
	}

	db = &mongoDb{
		client:  client,
		context: context,
		cancel:  cancel,
	}
}

func Disconnect() {
	if db == nil {
		log.Println("No active database connection to disconnect.")
		return
	}
	defer func() {
		db.cancel()
		if err := db.client.Disconnect(db.context); err != nil {
			log.Fatalf("mongodb disconnect error : %v", err)
		}
		log.Println("MongoDB disconnected successfully")
	}()
}

func DbClient() *mongo.Client {
	return db.client
}

func DbContext() context.Context {
	return db.context
}
