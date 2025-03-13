package main

import (
	mongoDb "github.com/tita-p/notes-api/internal/database"
	handler "github.com/tita-p/notes-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	defer mongoDb.Disconnect()

	router := gin.Default()
	router.GET("/tags", handler.GetTags)
	router.GET("/notes", handler.GetNotes)
	router.GET("/note/:id", handler.GetNote)

	router.POST("/tag", handler.AddTag)
	router.POST("/note", handler.AddNote)

	router.PATCH("/note", handler.EditNote)
	router.PATCH("/tag", handler.EditTag)

	router.DELETE("/tag/:id", handler.DeleteTag)
	router.DELETE("/note/:id", handler.DeleteNote)

	router.Run("localhost:8080")
}
