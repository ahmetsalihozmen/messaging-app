package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"messaging-app/internal/db"
	"net/http"
	"sync"
	"time"
)

type SenderService struct {
	url string
}

func NewSenderService(url string) *SenderService {
	return &SenderService{
		url: url,
	}
}

type ResponseMessage struct {
	MessageID string `json:"messageId"`
}

type MessageData struct {
	MessageID string
	TimeStamp time.Time
}

func (ss *SenderService) SendMessagesToWebhook(messages []db.Message) ([]MessageData, error) {
	var wg sync.WaitGroup
	resultCh := make(chan MessageData, len(messages))
	errorCh := make(chan error, len(messages))

	for _, message := range messages {
		wg.Add(1)

		go func(message db.Message) {
			defer wg.Done()

			body := fmt.Sprintf(`{"message": "%s", "to": "%s"}`, message.Content, message.Receiver)
			reqBody := bytes.NewBuffer([]byte(body))

			timeStamp := time.Now()

			resp, err := http.Post(ss.url, "application/json", reqBody)
			if err != nil {
				errorCh <- fmt.Errorf("failed to send request for message to %s: %v", message.Receiver, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusAccepted {
				errorCh <- fmt.Errorf("received an unexpected response: %d", resp.StatusCode)
				return
			}

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				errorCh <- fmt.Errorf("failed to read response body: %v", err)
				return
			}

			var jsonResponse ResponseMessage
			err = json.Unmarshal(responseBody, &jsonResponse)
			if err != nil {
				errorCh <- fmt.Errorf("failed to unmarshal response body: %v", err)
				return
			}

			resultCh <- MessageData{
				MessageID: jsonResponse.MessageID,
				TimeStamp: timeStamp,
			}

			fmt.Printf("Successfully sent message to %s\n", message.Receiver)
		}(message)
	}

	wg.Wait()

	// Close the channels for not blocking the main goroutine in endless loop in range
	close(resultCh)
	close(errorCh)

	var sentMsgsResp []MessageData
	for result := range resultCh {
		sentMsgsResp = append(sentMsgsResp, result)
	}

	if len(errorCh) > 0 {
		return nil, <-errorCh
	}

	return sentMsgsResp, nil
}
