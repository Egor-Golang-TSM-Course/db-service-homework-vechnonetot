package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsersHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "Hello, users!"})
}
