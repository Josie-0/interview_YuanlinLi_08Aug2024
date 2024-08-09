package routes

import (
	"PlayerManagementSystem/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/players", controllers.GetPlayers)
		v1.POST("/players", controllers.CreatePlayer)
		v1.GET("/players/:id", controllers.GetPlayer)
		v1.PUT("/players/:id", controllers.UpdatePlayer)
		v1.DELETE("/players/:id", controllers.DeletePlayer)

		v1.GET("/levels", controllers.GetLevels)
		v1.POST("/levels", controllers.CreateLevel)
	}

	return router
}
