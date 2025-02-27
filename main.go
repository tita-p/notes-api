package main

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	mongoDb "github.com/tita-p/notes-api/internal/database"
	model "github.com/tita-p/notes-api/internal/models"

	"github.com/gin-gonic/gin"
)

type ItemWithId interface {
	GetIte() int
}

var tags = []model.Tag{
	{Id: "1", Name: "Person 1"},
	{Id: "2", Name: "Person 2"},
}

var notes = []model.Note{
	{Id: "1", Title: "Note #1", Content: "Content of the note #1", Tag: tags[0]},
	{Id: "2", Title: "Note #2", Content: "Content of the note #2", Tag: tags[1]},
}

func main() {
	defer mongoDb.Disconnect()

	router := gin.Default()
	router.GET("/tags", getTags)
	router.GET("/notes", getNotes)
	router.GET("/note/:id", getNote)

	router.POST("/tag", addTag)
	router.POST("/note", addNote)

	router.PATCH("/note", editNote)
	router.PATCH("/tag", editTag)

	router.DELETE("/tag/:id", deleteTag)
	router.DELETE("/note/:id", deleteNote)

	router.Run("localhost:8080")
}

func getTags(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, model.GetTags())
}

func getNotes(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, model.GetNotes())
}

func getNote(context *gin.Context) {
	id, isFound := context.Params.Get("id")

	if !isFound {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "ID is required"})
		return
	}

	index := slices.IndexFunc(notes, func(note model.Note) bool { return note.Id == id })

	if index < 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	note := notes[index]

	context.IndentedJSON(http.StatusOK, note)
}

func addTag(context *gin.Context) {
	var tag model.Tag

	if err := context.ShouldBindJSON(&tag); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addedTag := model.Tag{
		Id:   fmt.Sprintf("%d", getHighestTagId()+1),
		Name: tag.Name,
	}

	tags = append(tags, addedTag)

	context.IndentedJSON(http.StatusOK, tags)
}

func addNote(context *gin.Context) {
	var note model.Note

	if err := context.ShouldBindJSON(&note); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagIndex := slices.IndexFunc(tags, func(tag model.Tag) bool { return tag.Id == note.Tag.Id })

	if tagIndex < 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "tag not found"})
		return
	}

	addedNote := []model.Note{
		{
			Title:   note.Title,
			Tag:     tags[tagIndex],
			Content: note.Content,
		},
	}

	model.InsertNotes(addedNote)

	context.IndentedJSON(http.StatusOK, model.GetNotes())
}

func editNote(context *gin.Context) {
	var note model.Note

	if err := context.ShouldBindJSON(&note); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := note.Id

	index := slices.IndexFunc(notes, func(note model.Note) bool { return note.Id == id })

	if index < 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	tagIndex := slices.IndexFunc(tags, func(tag model.Tag) bool { return tag.Id == note.Tag.Id })

	if tagIndex < 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "tag not found"})
		return
	}

	updatedNote := model.Note{
		Id:      note.Id,
		Title:   note.Title,
		Tag:     tags[tagIndex],
		Content: note.Content,
	}

	notes[index] = updatedNote

	context.IndentedJSON(http.StatusOK, notes[index])
}

func editTag(context *gin.Context) {
	var tag model.Tag

	if err := context.ShouldBindJSON(&tag); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := tag.Id

	index := slices.IndexFunc(tags, func(tag model.Tag) bool { return tag.Id == id })

	if index < 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	updatedTag := model.Tag{
		Id:   tag.Id,
		Name: tag.Name,
	}

	tags[index] = updatedTag

	context.IndentedJSON(http.StatusOK, tags[index])
}

func deleteTag(context *gin.Context) {
	id, error := getIdParam(context)

	if error != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": error.Error()})
		return
	}

	index := slices.IndexFunc(tags, func(tag model.Tag) bool { return tag.Id == id })

	if index < 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	tags = removeIndex(tags, index)
	context.IndentedJSON(http.StatusOK, tags)
}

func deleteNote(context *gin.Context) {
	id, error := getIdParam(context)

	if error != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": error.Error()})
		return
	}

	index := slices.IndexFunc(notes, func(note model.Note) bool { return note.Id == id })

	if index < 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	notes = removeIndex(notes, index)
	context.IndentedJSON(http.StatusOK, notes)
}

func getIdParam(context *gin.Context) (string, error) {
	id, isFound := context.Params.Get("id")

	if !isFound {
		return "", fmt.Errorf("ID is required")
	}

	return id, nil
}

func removeIndex[T comparable](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func getHighestTagId() int {
	highestId := 0

	for _, item := range tags {
		id, error := strconv.Atoi(item.Id)

		if error != nil {
			panic(error)
		}

		if id > highestId {
			highestId = id
		}
	}

	return highestId
}
