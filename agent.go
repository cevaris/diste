package diste

import (
	"errors"
	"sync"
)

type Enum string

const (
	OFF      Enum = "OFF"
	START    Enum = "START"
	ON       Enum = "ON"
	SHUTDOWN Enum = "SHUTDOWN"
)

type Agent struct {
	State    Enum
	Hostname string
}

func NewAgent(hostname string) *Agent {
	return &Agent{
		State:    OFF,
		Hostname: hostname,
	}
}

type AgentService struct {
	Agent *Agent
	Mutex *sync.RWMutex
}

func (s *AgentService) Ping(
	args ServiceRequest,
	reply *ServiceResponse) error {
	reply.Result = "Pong"
	return nil
}

func (s *AgentService) State(
	args ServiceRequest,
	reply *ServiceResponse) error {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	reply.Result = string(s.Agent.State)
	return nil
}

func (s *AgentService) Start(
	args ServiceRequest,
	reply *ServiceResponse) error {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	s.Agent.State = START
	reply.Result = string(s.Agent.State)
	return nil
}

func (s *AgentService) FakeError(
	args ServiceRequest,
	reply *ServiceResponse) error {

	return errors.New("System Error")
}
