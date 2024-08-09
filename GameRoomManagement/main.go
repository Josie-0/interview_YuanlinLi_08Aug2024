package main

import (
	"fmt"
	"log"
	"os"

	"GameRoomManagement/models"
	"GameRoomManagement/routes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "file:game_room_management.db?cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Room{}, &models.Reservation{}); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	createDefaultRooms(db)

	router := routes.SetupRouter(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func createDefaultRooms(db *gorm.DB) {
	defaultRooms := []models.Room{
		{Name: "Room A", Description: "Room A."},
		{Name: "Room B", Description: "Room B."},
	}

	for _, room := range defaultRooms {
		var existingRoom models.Room
		result := db.Where("name = ?", room.Name).First(&existingRoom)
		if result.RowsAffected == 0 {
			if err := db.Create(&room).Error; err != nil {
				log.Printf("failed to create default room %s: %v", room.Name, err)
			} else {
				log.Printf("created default room %s", room.Name)
			}
		}
	}
}
