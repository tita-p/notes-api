package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetIdParam(context *gin.Context) (string, error) {
	id, isFound := context.Params.Get("id")

	if !isFound {
		return "", fmt.Errorf("ID is required")
	}

	return id, nil
}
