package handler

import (
	"net/http"

	model "github.com/tita-p/notes-api/internal/models"

	"github.com/gin-gonic/gin"
)

func GetNotes(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, model.GetNotes())
}

func GetNote(context *gin.Context) {
	id, isFound := context.Params.Get("id")

	if !isFound {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "ID is required"})
		return
	}

	note, notFound := model.GetNoteById(id)

	if notFound {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Note not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, note)
}

func AddNote(context *gin.Context) {
	var note model.Note

	if err := context.ShouldBindJSON(&note); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, notFound := model.GetTagById(note.Tag.Id.Hex())

	if notFound {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tag not found"})
		return
	}

	addedNote := []*model.Note{
		{
			Title:   note.Title,
			Tag:     &tag,
			Content: note.Content,
		},
	}

	model.InsertNotes(addedNote)

	context.IndentedJSON(http.StatusOK, model.GetNotes())
}

func EditNote(context *gin.Context) {
	var note model.Note

	if err := context.ShouldBindJSON(&note); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, notFound := model.GetNoteById(note.Id.Hex())

	if notFound {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Note not found"})
		return
	}

	tag, tagNotFound := model.GetTagById(note.Tag.Id.Hex())

	if tagNotFound {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tag not found"})
		return
	}

	updatedNote := &model.Note{
		Id:      note.Id,
		Title:   note.Title,
		Tag:     &tag,
		Content: note.Content,
	}

	model.UpdateNote(updatedNote)

	context.IndentedJSON(http.StatusOK, model.GetNotes())
}

func DeleteNote(context *gin.Context) {
	id, error := GetIdParam(context)

	if error != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": error.Error()})
		return
	}

	model.DeleteNote(id)

	context.IndentedJSON(http.StatusOK, model.GetNotes())
}
