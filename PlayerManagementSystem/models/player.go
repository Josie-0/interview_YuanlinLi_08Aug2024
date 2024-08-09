package models

type Player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level string `json:"level"`
}

var PlayerMap = map[string]Player{}
