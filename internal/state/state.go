package state

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
)

type Manager struct {
	filePath string
}

func NewManager(filePath string) *Manager {
	return &Manager{
		filePath: filePath,
	}
}

func (m *Manager) ShouldPostUpdate(newContent string) bool {
	oldContent, err := os.ReadFile(m.filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Warning: Failed to read cache file: %v", err)
		}
		// If file doesn't exist, we assume it's new content (empty vs something).
	}
	return md5.Sum([]byte(newContent)) != md5.Sum(oldContent)
}

func (m *Manager) UpdateCache(content string) error {
	if err := os.WriteFile(m.filePath, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to update cache file: %w", err)
	}
	return nil
}
