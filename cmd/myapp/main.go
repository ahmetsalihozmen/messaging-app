package main

import (
	"fmt"
	"insider-project/internal/cache"
	"insider-project/internal/db"
	"insider-project/internal/service"
	"insider-project/pkg/handlers"
	"log"
	"net/http"

	_ "insider-project/api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Message Service API
// @version 1.0
// @description This is a sample server for the message service.
// @host localhost:8080
// @BasePath /
func main() {
	// Connect to the database and initialize the Redis client
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}
	cache.InitializeRedisClient()

	messageRepository := db.NewMessageRepository(database)

	messageService := service.NewMessageService(messageRepository, "https://webhook.site/ba6a41ba-a547-425f-ad7c-a833f8837b12")

	// Set up HTTP handlers
	http.HandleFunc("/startstop", func(w http.ResponseWriter, r *http.Request) {
		handlers.StartStopHandler(w, r, messageService)
	})

	http.HandleFunc("/get-sent-messages", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetSentMessagesHandler(w, r, messageService)

	})

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("HTTP server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
