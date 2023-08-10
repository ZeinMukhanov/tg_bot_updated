package main

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

const (
	apiURL = "https://api.telegram.org/bot"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "YOUR_PASSWORD"
	dbname   = "YOUR_DBNAME"
)

var connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

func main() {

	err := storeManagers()
	if err != nil {
		log.Println(err)
	}

	botToken := "YOUR_BOTTOKEN"

	updateID := 0

	for {
		updates, err := getUpdates(botToken, updateID)
		if err != nil {
			log.Println("Error getting updates:", err)
			continue
		}

		for _, update := range updates {
			updateID = update.UpdateID + 1

			if update.Message == nil || strings.HasPrefix(update.Message.Text, "/") {
				continue
			}

			chatID := update.Message.Chat.ID
			text := update.Message.Text

			if isValidEmail(text) {
				name, err := getManagerNameByEmail(text)
				if err != nil {
					log.Println("Error getting manager name:", err)
					continue
				}
				if name == "" {
					err = sendMessage(botToken, chatID, "No manager with this email found!")
					if err != nil {
						log.Println("Error sending message:", err)
					}
					continue
				}

				sipuniData, err := sendSipuniData(name)
				if err != nil {
					log.Println("Error getting sipuni data:", err)
					continue
				}
				err = sendMessage(botToken, chatID, sipuniData)
				if err != nil {
					log.Fatalf("Error sending message: %v", err)
				}
			}
		}
	}
}
