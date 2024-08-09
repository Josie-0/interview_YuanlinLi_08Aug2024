package main

import (
	"PlayerManagementSystem/routes"
)

func main() {
	router := routes.SetupRouter()
	router.Run(":8080")
}
