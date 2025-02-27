package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Tag struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

func InsertTags(tags []Tag) {
	items := sliceToInterface(tags, func(tag Tag) interface{} {
		return bson.M{
			"name":       tag.Name,
			"create_at":  time.Now(),
			"updated_at": time.Now(),
		}
	})

	InsertMany("tags", items)
}

func GetTags() []Tag {
	return Find[Tag]("tags")
}
