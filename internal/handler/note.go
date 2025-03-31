package handler

import (
	"net/http"

	model "github.com/tita-p/notes-api/internal/models"

	"github.com/gin-gonic/gin"
)

func GetNotes(context *gin.Context) {
	notes := model.GetNotes()

	var formattedNotes []*model.FormattedNote

	for _, note := range notes {
		formattedNote, tagNotFound := formatNote(&note)

		if tagNotFound {
			context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tag not found"})
			return
		}

		formattedNotes = append(formattedNotes, formattedNote)
	}

	context.IndentedJSON(http.StatusOK, formattedNotes)
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

	formattedNote, tagNotFound := formatNote(&note)

	if tagNotFound {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tag not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, formattedNote)
}

func AddNote(context *gin.Context) {
	var note model.Note

	if err := context.ShouldBindJSON(&note); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, notFound := model.GetTagById(note.TagId.Hex())

	if notFound {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tag not found"})
		return
	}

	addedNote := []*model.Note{
		{
			Title:   note.Title,
			TagId:   tag.Id,
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

	tag, tagNotFound := model.GetTagById(note.TagId.Hex())

	if tagNotFound {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tag not found"})
		return
	}

	updatedNote := &model.Note{
		Id:      note.Id,
		Title:   note.Title,
		TagId:   tag.Id,
		Content: note.Content,
	}

	model.UpdateNote(updatedNote)

	var formattedNotes []*model.FormattedNote

	for _, note := range model.GetNotes() {
		formattedNote, tagNotFound := formatNote(&note)

		if tagNotFound {
			context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tag not found"})
			return
		}

		formattedNotes = append(formattedNotes, formattedNote)
	}

	context.IndentedJSON(http.StatusOK, formattedNotes)
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

func formatNote(note *model.Note) (*model.FormattedNote, bool) {
	tag, tagNotFound := model.GetTagById(note.TagId.Hex())

	formattedNote := model.FormattedNote{
		Id:      note.Id,
		Title:   note.Title,
		Content: note.Content,
		Tag:     &tag,
	}

	return &formattedNote, tagNotFound
}
