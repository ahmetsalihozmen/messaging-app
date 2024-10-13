package handlers

import (
	"encoding/json"
	"fmt"
	"messaging-app/internal/service"
	"net/http"
)

// StartHandler godoc
// @Summary Starts and stops the message service
// @Description Starts or stops the message sending service that sends messages every 2 minutes by sending the action in query
// @Param action query string true "Action to perform: 'start' or 'stop'"
// @Tags Service
// @Success 200 {string} string "Service started"
// @Router /startstop [get]
func StartStopHandler(w http.ResponseWriter, r *http.Request, msgService *service.MessageService) {
	action := r.URL.Query().Get("action")

	if action == "start" {
		msgService.Start()
		fmt.Fprintln(w, "Service started.")
	} else if action == "stop" {
		msgService.Stop()
		fmt.Fprintln(w, "Service stopped.")
	} else {
		http.Error(w, "Invalid action. Please specify 'start' or 'stop'.", http.StatusBadRequest)
	}
}

// GetSentMessagesHandler
// @Summary Getting all messages that are sent
// @Description Getting all messages that are sent
// @Tags Service
// @Success 200 {string} string "Messages "
// @Router /get-sent-messages [get]
func GetSentMessagesHandler(w http.ResponseWriter, r *http.Request, msgService *service.MessageService) {
	w.Header().Set("Content-Type", "application/json")

	sentMessages, err := msgService.SentMessages()
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	if len(sentMessages) == 0 {
		w.Write([]byte("[]"))
		return
	}

	jsonMessages, err := json.Marshal(sentMessages)
	if err != nil {
		http.Error(w, "Failed to convert messages to JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonMessages)
}
