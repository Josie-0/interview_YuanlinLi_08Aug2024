package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"GameRoomManagement/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupDatabase initializes an in-memory SQLite database for testing.
func setupDatabase() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate the models
	db.AutoMigrate(&models.Room{}, &models.Reservation{})

	// Return a teardown function
	return db, func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

// setupRouter initializes the router with handlers
func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	roomHandler := &RoomHandler{DB: db}
	router.GET("/v1/rooms", roomHandler.GetRooms)
	router.POST("/v1/rooms", roomHandler.CreateRoom)
	router.GET("/v1/rooms/:id", roomHandler.GetRoomByID)
	router.PUT("/v1/rooms/:id", roomHandler.UpdateRoom)
	router.DELETE("/v1/rooms/:id", roomHandler.DeleteRoom)
	return router
}

func TestCreateRoom(t *testing.T) {
	db, teardown := setupDatabase()
	defer teardown()

	router := setupRouter(db)

	room := models.Room{Name: "Test Room", Description: "Test Room Description"}
	payload, _ := json.Marshal(room)
	req, _ := http.NewRequest("POST", "/v1/rooms", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["id"])

	var createdRoom models.Room
	db.First(&createdRoom, "id = ?", response["id"])
	assert.Equal(t, room.Name, createdRoom.Name)
	assert.Equal(t, room.Description, createdRoom.Description)
}

func TestGetRooms(t *testing.T) {
	db, teardown := setupDatabase()
	defer teardown()

	router := setupRouter(db)

	room := models.Room{Name: "Test Room", Description: "Test Room Description"}
	db.Create(&room)

	req, _ := http.NewRequest("GET", "/v1/rooms", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rooms []models.Room
	err := json.Unmarshal(w.Body.Bytes(), &rooms)
	assert.NoError(t, err)
	assert.Len(t, rooms, 1)
	assert.Equal(t, room.Name, rooms[0].Name)
}

func TestGetRoomByID(t *testing.T) {
	db, teardown := setupDatabase()
	defer teardown()

	router := setupRouter(db)

	room := models.Room{Name: "Test Room", Description: "Test Room Description"}
	db.Create(&room)

	req, _ := http.NewRequest("GET", "/v1/rooms/"+room.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Room
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, room.Name, response.Name)
}

func TestUpdateRoom(t *testing.T) {
	db, teardown := setupDatabase()
	defer teardown()

	router := setupRouter(db)

	room := models.Room{Name: "Test Room", Description: "Test Room Description"}
	db.Create(&room)

	updatedRoom := models.Room{Name: "Updated Room", Description: "Updated Room Description"}
	payload, _ := json.Marshal(updatedRoom)
	req, _ := http.NewRequest("PUT", "/v1/rooms/"+room.ID, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Room
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updatedRoom.Name, response.Name)
	assert.Equal(t, updatedRoom.Description, response.Description)
}

func TestDeleteRoom(t *testing.T) {
	db, teardown := setupDatabase()
	defer teardown()

	router := setupRouter(db)

	room := models.Room{Name: "Test Room", Description: "Test Room Description"}
	db.Create(&room)

	req, _ := http.NewRequest("DELETE", "/v1/rooms/"+room.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	var deletedRoom models.Room
	result := db.First(&deletedRoom, "id = ?", room.ID)
	assert.Error(t, result.Error)
}
