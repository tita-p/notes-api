package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	Id   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

var tagCollectionName = "tags"

func InsertTags(tags []*Tag) []Tag {
	items := sliceToInterface(tags, func(tag Tag) interface{} {
		return bson.M{
			"name":       tag.Name,
			"create_at":  time.Now(),
			"updated_at": time.Now(),
		}
	})

	return InsertMany[Tag](tagCollectionName, items)
}

func GetTagById(id string) (Tag, bool) {
	objectID := stringToObjectId(id)

	return FindById[Tag](tagCollectionName, objectID)
}

func GetTags(filters ...*filter) []Tag {
	return Find[Tag](tagCollectionName, filters...)
}

func UpdateTag(tag *Tag) {
	updateSet := bson.M{
		"$set": bson.M{
			"name":       tag.Name,
			"updated_at": time.Now(),
		},
	}

	Update[Tag](tagCollectionName, stringToObjectId(tag.Id.Hex()), updateSet)
}

func DeleteTag(id string) {
	Delete(tagCollectionName, stringToObjectId(id))
}
