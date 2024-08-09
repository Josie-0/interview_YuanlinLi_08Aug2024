package models

import "github.com/google/uuid"

type Level struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func InitializeLevelMap() map[string]Level {
	return map[string]Level{
		"Beginner":     {ID: uuid.New().String(), Name: "Beginner"},
		"Intermediate": {ID: uuid.New().String(), Name: "Intermediate"},
		"Advanced":     {ID: uuid.New().String(), Name: "Advanced"},
	}
}

var LevelMap = InitializeLevelMap()
