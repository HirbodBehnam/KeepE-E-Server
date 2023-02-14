package session

import "sync"

type InMemorySession struct {
	data map[string]uint64
	mu   sync.RWMutex
}

func NewInMemorySession() *InMemorySession {
	return &InMemorySession{
		data: make(map[string]uint64),
	}
}

func (s *InMemorySession) Store(userID uint64) (string, error) {
	token := newSessionID()
	s.mu.Lock()
	s.data[token] = userID
	s.mu.Unlock()
	return token, nil
}

func (s *InMemorySession) Authorize(token string) (uint64, error) {
	s.mu.RLock()
	userID, ok := s.data[token]
	s.mu.RUnlock()
	if !ok {
		return 0, UnauthorizedTokenErr
	}
	return userID, nil
}

func (s *InMemorySession) Delete(token string) error {
	s.mu.Lock()
	delete(s.data, token)
	s.mu.Unlock()
	return nil
}
