package routes

import (
	"github.com/Egor-Golang-TSM-Course/db-service-homework-vechnonetot/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/users", handlers.GetUsersHandler)

	return router
}
