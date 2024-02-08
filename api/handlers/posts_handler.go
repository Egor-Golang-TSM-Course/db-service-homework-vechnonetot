package handlers

import (
	"context"
	"github.com/Egor-Golang-TSM-Course/db-service-homework-vechnonetot/api/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreatePostForHandler(c *gin.Context) {
	user, err := database.GetUserFromContext(c.Request.Context())
	if err != nil {
		// Обработка ошибки, например, отправка HTTP-ответа с ошибкой
		return
	}

	newPost := database.Post{
		UserID: user.ID,
		// Другие поля поста
	}

	createdPost, err := database.CreatePost(newPost)
	if err != nil {

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post": createdPost})
}

func GetCommentsForPostHandler(c *gin.Context) {
	var err error
	var postID int

	if postID, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	comments, err := database.GetCommentsForPost(context.Background(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
