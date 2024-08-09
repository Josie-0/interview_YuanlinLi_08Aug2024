package controllers

import (
	"PlayerManagementSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetLevels(c *gin.Context) {
	c.JSON(http.StatusOK, models.LevelMap)
}

func CreateLevel(c *gin.Context) {
	var newLevel models.Level
	if err := c.BindJSON(&newLevel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	if _, exists := models.LevelMap[newLevel.Name]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Level already exists"})

	} else {
		newLevel.ID = uuid.New().String()

		models.LevelMap[newLevel.Name] = newLevel
		c.JSON(http.StatusCreated, newLevel)
	}
}
