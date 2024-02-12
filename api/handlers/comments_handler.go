package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Egor-Golang-TSM-Course/db-service-homework-vechnonetot/api/database"
	"github.com/gin-gonic/gin"
)

type Comment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"postID"`
	UserID  int    `json:"userID"`
	Content string `json:"content"`
}

func CreateCommentHandler(c *gin.Context) {
	var comment database.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commentID, err := database.AddComment(context.Background(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"commentID": commentID})
}

func GetCommentsForHandler(c *gin.Context) {
	var err error
	var postID int

	// Парсим параметр "id" из URL
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

func DeleteCommentHandler(c *gin.Context) {
	var err error
	var commentID int

	if commentID, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	err = database.DeleteComment(context.Background(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
