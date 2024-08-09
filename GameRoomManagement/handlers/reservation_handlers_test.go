package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"GameRoomManagement/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 初始化测试路由
func setupRoomRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	reservationHandler := &ReservationHandler{DB: db}
	router.GET("/v1/reservations", reservationHandler.GetReservations)
	router.POST("/v1/reservations", reservationHandler.CreateReservation)
	return router
}

// 初始化测试数据库
func setupRoomDatabase() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&models.Room{}, &models.Reservation{}); err != nil {
		panic("failed to migrate database")
	}

	// 清理测试数据库
	return db, func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

func TestCreateReservation(t *testing.T) {
	db, teardown := setupRoomDatabase()
	defer teardown()

	router := setupRoomRouter(db)

	// Create a room for reservation
	room := models.Room{Name: "Test Room", Description: "Test Room Description"}
	db.Create(&room)

	// Create a reservation
	reservation := models.Reservation{
		RoomID:         room.ID,
		WholeHourStart: time.Now().Add(1 * time.Hour).Truncate(time.Hour), // 1小时后整点
	}
	payload, _ := json.Marshal(reservation)
	req, _ := http.NewRequest("POST", "/v1/reservations", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify the reservation in the database
	var res models.Reservation
	result := db.First(&res, "room_id = ? AND whole_hour_start = ?", room.ID, reservation.WholeHourStart)
	assert.NoError(t, result.Error)

	// Compare only the time part, ignoring time zone differences
	assert.True(t, res.WholeHourStart.Equal(reservation.WholeHourStart.In(res.WholeHourStart.Location())))
}

func TestGetReservations(t *testing.T) {
	db, teardown := setupRoomDatabase()
	defer teardown()

	router := setupRoomRouter(db)

	// Create a room
	room := models.Room{Name: "Test Room", Description: "Test Room Description"}
	db.Create(&room)

	// Create reservations
	reservation1 := models.Reservation{
		RoomID:         room.ID,
		WholeHourStart: time.Now().Add(1 * time.Hour).Truncate(time.Hour),
	}
	reservation2 := models.Reservation{
		RoomID:         room.ID,
		WholeHourStart: time.Now().Add(2 * time.Hour).Truncate(time.Hour),
	}
	db.Create(&reservation1)
	db.Create(&reservation2)

	// Query reservations
	req, _ := http.NewRequest("GET", "/v1/reservations?room_id="+room.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var reservations []models.Reservation
	err := json.Unmarshal(w.Body.Bytes(), &reservations)
	assert.NoError(t, err)
	assert.Len(t, reservations, 2)
}

func TestGetReservationsInvalidDate(t *testing.T) {
	db, teardown := setupRoomDatabase()
	defer teardown()

	router := setupRoomRouter(db)

	// Create a room
	room := models.Room{Name: "Test Room", Description: "Test Room Description"}
	db.Create(&room)

	// Create a reservation
	reservation := models.Reservation{
		RoomID:         room.ID,
		WholeHourStart: time.Now().Add(1 * time.Hour).Truncate(time.Hour),
	}
	db.Create(&reservation)

	// Query reservations with invalid date
	req, _ := http.NewRequest("GET", "/v1/reservations?room_id="+room.ID+"&date=invalid-date", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateReservationInvalidTime(t *testing.T) {
	db, teardown := setupRoomDatabase()
	defer teardown()

	router := setupRoomRouter(db)

	// Create a room
	room := models.Room{Name: "Test Room", Description: "Test Room Description"}
	db.Create(&room)

	// Create a reservation with invalid time (not on the hour)
	reservation := models.Reservation{
		RoomID:         room.ID,
		WholeHourStart: time.Now().Add(1 * time.Hour).Truncate(time.Hour).Add(30 * time.Minute), // 不整点
	}
	payload, _ := json.Marshal(reservation)
	req, _ := http.NewRequest("POST", "/v1/reservations", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
