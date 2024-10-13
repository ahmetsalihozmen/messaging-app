package service

import (
	"messaging-app/internal/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockMessageRepository struct{}

func (mdb *MockMessageRepository) GetSentMessages() ([]db.Message, error) {
	return []db.Message{
		{Receiver: "+905551110001", Content: "test1", Status: "sent"},
		{Receiver: "+905551110002", Content: "test2", Status: "sent"},
	}, nil
}

func (mdb *MockMessageRepository) GetUnsentMessages(count int) ([]db.Message, error) {
	return []db.Message{
		{Receiver: "+905551110003", Content: "test3", Status: "unsent"},
		{Receiver: "+905551110004", Content: "test4", Status: "unsent"},
	}, nil
}

func (mdb *MockMessageRepository) UpdateMessageStatus(messages []db.Message, status string) error {
	return nil
}

func TestMessageService_Start(t *testing.T) {
	mockRepo := &MockMessageRepository{}
	mockServer := httptest.NewServer(http.HandlerFunc(mockWebhookSuccess))

	service := NewMessageService(mockRepo, mockServer.URL, 1)

	service.Start()

	if _, ok := service.state.(*StartedState); !ok {
		t.Errorf("Expected StartedState but got %T", service.state)
	}
}

func TestMessageService_Stop(t *testing.T) {
	mockRepo := &MockMessageRepository{}
	mockServer := httptest.NewServer(http.HandlerFunc(mockWebhookSuccess))

	service := NewMessageService(mockRepo, mockServer.URL, 1)

	service.Stop()

	if _, ok := service.state.(*StoppedState); !ok {
		t.Errorf("Expected StoppedState but got %T", service.state)
	}
}

func TestMessageService_SendMessages(t *testing.T) {
	mockRepo := &MockMessageRepository{}
	mockServer := httptest.NewServer(http.HandlerFunc(mockWebhookSuccesEmpty))

	service := NewMessageService(mockRepo, mockServer.URL, 1)
	service.scheduler.StopTimer()

	// Send messages
	err := service.SendMessages()

	if err != nil {
		t.Errorf("Error sending messages: %v", err)
	}
}

func TestMessageService_GetSentMessages(t *testing.T) {
	mockRepo := &MockMessageRepository{}

	service := NewMessageService(mockRepo, "", 1)

	sentMessages, err := service.SentMessages()

	if err != nil {
		t.Errorf("Error fetching sent messages: %v", err)
	}

	if len(sentMessages) != 2 {
		t.Errorf("Expected 2 sent messages, got %d", len(sentMessages))
	}
}
