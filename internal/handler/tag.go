package handler

import (
	"net/http"

	model "github.com/tita-p/notes-api/internal/models"

	"github.com/gin-gonic/gin"
)

func GetTags(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, model.GetTags())
}

func AddTag(context *gin.Context) {
	var tag model.Tag

	if err := context.ShouldBindJSON(&tag); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addedTag := []*model.Tag{
		{
			Name: tag.Name,
		},
	}

	model.InsertTags(addedTag)

	context.IndentedJSON(http.StatusOK, model.GetTags())
}

func EditTag(context *gin.Context) {
	var tag model.Tag

	if err := context.ShouldBindJSON(&tag); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, notFound := model.GetTagById(tag.Id.Hex())

	if notFound {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Tag not found"})
		return
	}

	updatedTag := &model.Tag{
		Id:   tag.Id,
		Name: tag.Name,
	}

	model.UpdateTag(updatedTag)

	updateResult, _ := model.GetTagById(tag.Id.Hex())

	context.IndentedJSON(http.StatusOK, updateResult)
}

func DeleteTag(context *gin.Context) {
	id, error := GetIdParam(context)

	if error != nil {
		context.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": error.Error()})
		return
	}

	model.DeleteTag(id)

	context.IndentedJSON(http.StatusOK, model.GetTags())
}
