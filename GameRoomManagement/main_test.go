package main

import (
	"GameRoomManagement/handlers"
	"GameRoomManagement/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	roomHandler := &handlers.RoomHandler{DB: db}
	router.POST("/v1/rooms", roomHandler.CreateRoom)
	router.GET("/v1/rooms/:id", roomHandler.GetRoomByID)
	return router
}

func TestCreateRoom(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&models.Room{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	router := setupRouter(db)

	// Create a room
	payload := `{"name": "Test Room", "description": "Test Room Description"}`
	req, _ := http.NewRequest("POST", "/v1/rooms", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Retrieve the room
	var room models.Room
	result := db.First(&room, "name = ?", "Test Room")
	if result.Error != nil {
		t.Fatalf("failed to find room: %v", result.Error)
	}

	assert.Equal(t, "Test Room", room.Name)
	assert.Equal(t, "Test Room Description", room.Description)
}
