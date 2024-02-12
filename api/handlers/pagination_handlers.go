package handlers

import (
	"context"
	"fmt"
	"github.com/Egor-Golang-TSM-Course/db-service-homework-vechnonetot/api/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PaginationComment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"postID"`
	UserID  int    `json:"userID"`
	Content string `json:"content"`
}

func GetPaginatedPostsHandler(c *gin.Context) {
	var err error
	var page, pageSize, postID int

	if postID, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	if page, err = strconv.Atoi(c.Query("page")); err != nil {
		page = 1
	}
	if pageSize, err = strconv.Atoi(c.Query("pageSize")); err != nil {
		pageSize = 10
	}

	posts, err := database.GetPaginatedPosts(context.Background(), page, pageSize, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get paginated posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func GetPaginatedCommentsHandler(c *gin.Context) {
	var err error
	var postID, page, pageSize int

	if postID, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	if page, err = strconv.Atoi(c.Query("page")); err != nil {
		page = 1
	}
	if pageSize, err = strconv.Atoi(c.Query("pageSize")); err != nil {
		pageSize = 10
	}

	comments, err := GetPaginatedComments(context.Background(), postID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get paginated comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func GetPaginatedComments(ctx context.Context, postID, page, pageSize int) ([]PaginationComment, error) {

	return nil, fmt.Errorf("not implemented")
}
