package session

import (
	"context"
	"sync"
	"time"
)

type UserSession struct {
	SessionId    int64
	CancelFunc   context.CancelFunc
	continueChan chan struct{}
	TimeFrom     time.Time
	TimeTo       time.Time
	mu           sync.Mutex
}

func NewUserSession(sessionId int64) *UserSession {
	return &UserSession{
		SessionId:    sessionId,
		continueChan: make(chan struct{}, 1),
		TimeFrom:     time.Now(),
		TimeTo:       time.Now(),
		mu:           sync.Mutex{},
	}
}

func (s *UserSession) Continue() {
	select {
	case s.continueChan <- struct{}{}:
	default:
	}
}

func (s *UserSession) WaitContinue() <-chan struct{} {
	return s.continueChan
}
