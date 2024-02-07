package main

import (
	"fmt"
	"github.com/Egor-Golang-TSM-Course/db-service-homework-vechnonetot/api/routes"
)

func main() {

	router := routes.SetupRouter()

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}

	fmt.Println("Server is running on :8080")
}
