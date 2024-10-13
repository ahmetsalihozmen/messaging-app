package service

import (
	"fmt"
	"insider-project/internal/cache"
	"insider-project/internal/db"
	"sync"
	"time"
)

type MessageService struct {
	state     ServiceState
	scheduler *MessageScheduler
	sender    *SenderService
	repo      db.MessageRepository
	mutex     sync.Mutex
}

func NewMessageService(repo db.MessageRepository, webhookUrl string) *MessageService {
	service := &MessageService{
		state:     &StartedState{},
		scheduler: NewMessageScheduler(),
		sender:    NewSenderService(webhookUrl),
		repo:      repo,
	}
	service.scheduler.service = service
	service.state.HandleState(service)
	return service
}

func (ms *MessageService) SetState(state ServiceState) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.state = state
}

func (ms *MessageService) Start() {
	ms.SetState(&StartedState{})
	ms.state.HandleState(ms)
}

func (ms *MessageService) Stop() {
	ms.SetState(&StoppedState{})
	ms.state.HandleState(ms)
}

func (ms *MessageService) SentMessages() ([]db.Message, error) {
	sentMessages, err := ms.repo.GetSentMessages()
	if err != nil {
		fmt.Println("Failed to fetch messages ")
		return nil, err
	}

	return sentMessages, nil
}

func (ms *MessageService) SendMessages() error {
	unsentMessages, err := ms.repo.GetUnsentMessages(2)
	if err != nil {
		return fmt.Errorf("failed to fetch messages")
	}

	if len(unsentMessages) != 2 {
		return fmt.Errorf("insufficient number of messages")
	}

	messageIds, err := ms.sender.SendMessagesToWebhook(unsentMessages)

	if cache.RedisClient != nil {
		for _, messageId := range messageIds {
			fmt.Println(messageId.MessageID)
			cache.RedisClient.Set(
				cache.RedisClient.Context(),
				messageId.MessageID,
				messageId.TimeStamp,
				0,
			)
		}
	}

	if err != nil {
		fmt.Println("Failed to send messages")
		return err
	}

	err = ms.repo.UpdateMessageStatus(unsentMessages, "sent")

	if err != nil {
		fmt.Println("Failed to update message status")
		return err
	}

	return nil
}

type MessageScheduler struct {
	service *MessageService
	ticker  *time.Ticker
	done    chan bool
	mutex   sync.Mutex
}

func NewMessageScheduler() *MessageScheduler {
	return &MessageScheduler{
		done: make(chan bool),
	}
}

func (scheduler *MessageScheduler) StartTimer() {
	scheduler.mutex.Lock()
	defer scheduler.mutex.Unlock()

	if scheduler.ticker != nil {
		scheduler.ticker.Stop()
	}

	fmt.Println("Starting timer to send messages every 2 minutes")
	scheduler.ticker = time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-scheduler.ticker.C:
				scheduler.service.SendMessages()
			case <-scheduler.done:
				return
			}
		}
	}()
}

func (scheduler *MessageScheduler) StopTimer() {
	scheduler.mutex.Lock()
	defer scheduler.mutex.Unlock()
	fmt.Println("Stopping the timer and resetting.")
	if scheduler.ticker != nil {
		scheduler.ticker.Stop()
		scheduler.done <- true
		scheduler.ticker = nil
	}
}
