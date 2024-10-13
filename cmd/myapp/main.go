package main

import (
	"fmt"
	"log"
	"messaging-app/internal/cache"
	"messaging-app/internal/config"
	"messaging-app/internal/db"
	"messaging-app/internal/handlers"
	"messaging-app/internal/service"
	"net/http"

	_ "messaging-app/api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Message Service API
// @version 1.0
// @description This is a sample server for the message service.
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.LoadConfig()

	database, err := db.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}
	cache.InitializeRedisClient(cfg.RedisAddr)

	messageRepository := db.NewMessageRepository(database)

	messageService := service.NewMessageService(messageRepository, cfg.WebhookURL, cfg.MessagingPeriod)

	// Set up HTTP handlers
	http.HandleFunc("/startstop", func(w http.ResponseWriter, r *http.Request) {
		handlers.StartStopHandler(w, r, messageService)
	})

	http.HandleFunc("/get-sent-messages", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetSentMessagesHandler(w, r, messageService)

	})

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("HTTP server listening on :8080")
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
