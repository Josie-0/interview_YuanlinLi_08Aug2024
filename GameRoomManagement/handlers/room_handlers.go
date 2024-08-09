package handlers

import (
	"GameRoomManagement/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type RoomHandler struct {
	DB *gorm.DB
}

func (h *RoomHandler) GetRooms(c *gin.Context) {
	var rooms []models.Room
	if err := h.DB.Preload("Reservations").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": room.ID})
}

func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	id := c.Param("id")
	var room models.Room
	if err := h.DB.Preload("Reservations").First(&room, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room
	if err := h.DB.First(&room, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")

	result := h.DB.Delete(&models.Room{}, "id = ?", id)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}
