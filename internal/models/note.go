package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Note struct {
	Id      string `json:"id" bson:"_id"`
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
	Tag     Tag    `json:"tag" bson:"tag"`
}

func InsertNotes(notes []Note) {
	items := sliceToInterface(notes, func(note Note) interface{} {
		return bson.M{
			"title":      note.Title,
			"content":    note.Content,
			"create_at":  time.Now(),
			"updated_at": time.Now(),
		}
	})

	InsertMany("notes", items)
}

func GetNotes() []Note {
	return Find[Note]("notes")
}
