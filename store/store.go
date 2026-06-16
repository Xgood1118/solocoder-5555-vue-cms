package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"cms/config"
)

type Store struct {
	mu       sync.RWMutex
	dataDir  string
	filename string
}

func NewStore(filename string) *Store {
	return &Store{
		dataDir:  config.AppConfig.DataDir,
		filename: filename,
	}
}

func (s *Store) filePath() string {
	return filepath.Join(s.dataDir, s.filename)
}

func (s *Store) ReadAll(v interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filePath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read file error: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	return nil
}

func (s *Store) WriteAll(v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(s.dataDir, 0755); err != nil {
		return fmt.Errorf("mkdir error: %w", err)
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := os.WriteFile(s.filePath(), data, 0644); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}
