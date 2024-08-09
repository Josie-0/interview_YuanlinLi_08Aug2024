package handlers

import (
	"GameRoomManagement/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type ReservationHandler struct {
	DB *gorm.DB
}

func (h *ReservationHandler) GetReservations(c *gin.Context) {
	roomID := c.Query("room_id")
	dateStr := c.Query("date")
	limitStr := c.Query("limit")

	var reservations []models.Reservation
	query := h.DB.Model(&models.Reservation{})

	// 验证房间 ID
	if roomID != "" {
		if err := h.validateRoomID(roomID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		query = query.Where("room_id = ?", roomID)
	}

	if dateStr != "" {
		date, err := h.validateDate(dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query = query.Where("DATE(whole_hour_start) = ?", date.Format("2006-01-02"))
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
			return
		}
		query = query.Limit(limit)
	}

	if err := query.Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reservations)
}

// 用于验证仅包含日期的字符串: "2024-08-19"
func (h *ReservationHandler) validateDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format")
	}
	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		return time.Time{}, fmt.Errorf("date cannot be in the past")
	}
	return date, nil
}

func (h *ReservationHandler) validateRoomID(roomID string) error {
	var room models.Room
	if err := h.DB.First(&room, "id = ?", roomID).Error; err != nil {
		return fmt.Errorf("invalid room ID")
	}
	return nil
}

func (h *ReservationHandler) validateReservationTime(roomID string, wholeHourStart time.Time) error {
	// 为了简化，时间只能是整点开始，且预定时长固定为1h
	if wholeHourStart.Minute() != 0 || wholeHourStart.Second() != 0 {
		return fmt.Errorf("WholeHourStart must be at a whole hour")
	}

	// 检查时间是否在过去
	if !wholeHourStart.After(time.Now()) {
		return fmt.Errorf("datetime cannot be in the past")
	}

	// 检查时间是否与现有预订冲突
	var existingReservation models.Reservation
	result := h.DB.Where("room_id = ? AND whole_hour_start = ?", roomID, wholeHourStart).First(&existingReservation)
	if result.RowsAffected > 0 {
		return fmt.Errorf("a reservation already exists at the given time")
	}

	return nil
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证房间 ID
	if err := h.validateRoomID(reservation.RoomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证预订时间
	if err := h.validateReservationTime(reservation.RoomID, reservation.WholeHourStart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建预订
	if err := h.DB.Create(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": reservation.ID})
}
