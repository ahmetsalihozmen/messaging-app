package service

import (
	"messaging-app/internal/db"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testMessages = []db.Message{
	{Content: "Test Message 1", Receiver: "receiver1", Status: "unsent"},
	{Content: "Test Message 2", Receiver: "receiver2", Status: "unsent"},
}

func mockWebhookSuccess(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"messageId": "12345"}`))
}

func mockWebhookFailure(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func mockWebhookSuccesEmpty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{}`))
}

func TestSendMessagesToWebhook_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockWebhookSuccess))
	defer mockServer.Close()

	ss := NewSenderService(mockServer.URL)
	result, err := ss.SendMessagesToWebhook(testMessages)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != len(testMessages) {
		t.Errorf("Expected %d messages, got %d", len(testMessages), len(result))
	}

	for _, msg := range result {
		if !strings.HasPrefix(msg.MessageID, "12345") {
			t.Errorf("Expected messageId to be '12345', got %s", msg.MessageID)
		}
	}
}

func TestSendMessagesToWebhook_Failure(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockWebhookFailure))
	defer mockServer.Close()

	ss := NewSenderService(mockServer.URL)

	_, err := ss.SendMessagesToWebhook(testMessages)

	if err == nil {
		t.Error("Expected an error, but got none")
	}

	expectedError := "received an unexpected response"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message to contain '%s', got %v", expectedError, err)
	}
}
