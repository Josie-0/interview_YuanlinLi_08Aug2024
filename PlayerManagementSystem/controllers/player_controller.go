package controllers

import (
	"PlayerManagementSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetPlayers(c *gin.Context) {
	c.JSON(http.StatusOK, models.PlayerMap)
}

func GetPlayer(c *gin.Context) {
	id := c.Param("id")
	if player, exists := models.PlayerMap[id]; exists {
		c.JSON(http.StatusOK, player)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Player not found"})
}

func CreatePlayer(c *gin.Context) {
	var newPlayer models.Player
	if err := c.BindJSON(&newPlayer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	if _, ok := models.LevelMap[newPlayer.Level]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid level"})
		return
	}

	newPlayer.ID = uuid.New().String()
	models.PlayerMap[newPlayer.ID] = newPlayer
	c.JSON(http.StatusCreated, newPlayer)
}

func UpdatePlayer(c *gin.Context) {
	id := c.Param("id")

	if player, exists := models.PlayerMap[id]; exists {
		if err := c.BindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		models.PlayerMap[id] = player
		c.JSON(http.StatusOK, player)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Player not found"})
}

func DeletePlayer(c *gin.Context) {
	id := c.Param("id")

	if _, exists := models.PlayerMap[id]; exists {
		delete(models.PlayerMap, id)
		c.JSON(http.StatusOK, gin.H{"message": "Player deleted"})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Player not found"})
}
