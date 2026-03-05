package logger

import (
	"encoding/json"
	"log"

	"github.com/go-telegram/bot/models"
)

func LogUpdate(update *models.Update) {
	data, err := json.MarshalIndent(update, "", "  ")
	if err != nil {
		log.Println("Failed to serialize:", err)
		return
	}

	log.Println(string(data))
}
