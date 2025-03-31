package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	Id      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
	TagId   primitive.ObjectID `json:"tagId" bson:"tagId"`
}

type FormattedNote struct {
	Id      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
	Tag     *Tag               `json:"tag" bson:"tag"`
}

var noteCollectionName = "notes"

func InsertNotes(notes []*Note) []Note {
	items := sliceToInterface(notes, func(note Note) interface{} {
		return bson.M{
			"title":      note.Title,
			"content":    note.Content,
			"tagId":      note.TagId,
			"create_at":  time.Now(),
			"updated_at": time.Now(),
		}
	})

	return InsertMany[Note](noteCollectionName, items)
}

func GetNotes(filters ...*filter) []Note {
	return Find[Note](noteCollectionName, filters...)
}

func UpdateNote(note *Note) {
	updateSet := bson.M{
		"$set": bson.M{
			"title":      note.Title,
			"content":    note.Content,
			"tagId":      note.TagId,
			"updated_at": time.Now(),
		},
	}

	Update[Note](noteCollectionName, note.Id, updateSet)
}

func DeleteNote(id string) {
	Delete(noteCollectionName, stringToObjectId(id))
}

func GetNoteById(id string) (Note, bool) {
	objectID := stringToObjectId(id)

	return FindById[Note](noteCollectionName, objectID)
}
