package scheduler

import (
	"log/slog"
	"sync"
)

type Key string

const (
	SnapshotKey Key = "snapshot"
	DateKey     Key = "date"
	TimeFromKey Key = "time_from"
	TimeToKey   Key = "time_to"
)

type Storage struct {
	scheduleCourts map[string][]Court
	sessionCache   map[int64]map[Key]string
	mu             sync.Mutex
}

func (s *Storage) GetCourtsByDate(date string) []Court {
	s.mu.Lock()
	defer s.mu.Unlock()

	courtsFromStorage := s.scheduleCourts[date]

	if courtsFromStorage == nil {
		return []Court{}
	}

	return courtsFromStorage
}

func (s *Storage) GetSessionCacheByKey(sessionId int64, key Key) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionCache, ok := s.sessionCache[sessionId][key]
	if ok {
		return sessionCache
	}

	return ""
}

// CacheByKey Get cache value from storage
func (s *Storage) CacheByKey(sessionId int64, key Key) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.sessionCache[sessionId][key]

	if !ok {
		slog.Warn("no cache for key")
		value = ""
	}

	return value
}

// AddCache Method helps to add key to cache storage
func (s *Storage) AddCache(sessionId int64, key Key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessionCache[sessionId]
	if !ok {
		session = make(map[Key]string)
	}
	session[key] = value

	s.sessionCache[sessionId] = session
}

func NewStorage() *Storage {
	return &Storage{
		scheduleCourts: make(map[string][]Court),
		sessionCache:   make(map[int64]map[Key]string),
	}
}

func (s *Storage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scheduleCourts = make(map[string][]Court)
}

func (s *Storage) AddCourt(date string, court Court) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.scheduleCourts[date] = append(s.scheduleCourts[date], court)
}
