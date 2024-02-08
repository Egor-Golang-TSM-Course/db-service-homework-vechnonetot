package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Egor-Golang-TSM-Course/db-service-homework-vechnonetot/api/database"
	"github.com/gin-gonic/gin"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func AddTagHandler(c *gin.Context) {
	var tag database.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "postID", postID)

	tagID, err := database.AddTag(ctx, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag added successfully", "tagID": tagID})
}

func GetTagsHandler(c *gin.Context) {

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "postID", postID)

	tags, err := database.GetTags(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}
