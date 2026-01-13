package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type State struct {
	path string
	mu   sync.RWMutex
	Map  map[string]string
}

func Load(rootPath string) (*State, error) {
	path := filepath.Join(rootPath, ".debtbomb", "jira-map.json")
	s := &State{
		path: path,
		Map:  make(map[string]string),
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return s, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return s, nil
	}

	if err := json.Unmarshal(data, &s.Map); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *State) GetTicket(bombID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Map[bombID]
}

func (s *State) SetTicket(bombID, ticket string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Map[bombID] = ticket
}

func (s *State) RemoveTicket(bombID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Map, bombID)
}

func (s *State) Snapshot() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]string)
	for k, v := range s.Map {
		copy[k] = v
	}
	return copy
}

func (s *State) Save() error {

	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.Map, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}