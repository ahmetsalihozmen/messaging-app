package service

import "fmt"

type ServiceState interface {
	HandleState(service *MessageService)
}

type StartedState struct{}

func (s *StartedState) HandleState(service *MessageService) {
	fmt.Println("Service is running: sending messages...")
	service.scheduler.StartTimer()
}

type StoppedState struct{}

func (s *StoppedState) HandleState(service *MessageService) {
	fmt.Println("Service is stopped: messages not being sent.")
	service.scheduler.StopTimer()
}
