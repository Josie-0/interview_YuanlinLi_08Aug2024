package routes

import (
	"GameRoomManagement/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")

	roomHandler := &handlers.RoomHandler{DB: db}
	reservationHandler := &handlers.ReservationHandler{DB: db}

	v1.GET("/rooms", roomHandler.GetRooms)
	v1.POST("/rooms", roomHandler.CreateRoom)
	v1.GET("/rooms/:id", roomHandler.GetRoomByID)
	v1.PUT("/rooms/:id", roomHandler.UpdateRoom)
	v1.DELETE("/rooms/:id", roomHandler.DeleteRoom)

	v1.GET("/reservations", reservationHandler.GetReservations)
	v1.POST("/reservations", reservationHandler.CreateReservation)

	return router
}
