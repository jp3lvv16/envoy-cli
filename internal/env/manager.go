package env

import (
	"fmt"

	"github.com/user/envoy-cli/internal/config"
)

// Manager provides high-level operations on env sets backed by config.
type Manager struct {
	cfg *config.Config
}

// NewManager creates a Manager wrapping the given Config.
func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// CreateSet adds a new empty set to the config.
func (m *Manager) CreateSet(name string) error {
	if _, err := m.cfg.GetSet(name); err == nil {
		return fmt.Errorf("set %q already exists", name)
	}
	s, err := NewSet(name)
	if err != nil {
		return err
	}
	m.cfg.AddOrUpdateSet(config.EnvSet{Name: s.Name, Vars: s.Vars})
	return nil
}

// SetVar sets a variable inside a named set.
func (m *Manager) SetVar(setName, key, value string) error {
	es, err := m.cfg.GetSet(setName)
	if err != nil {
		return err
	}
	s := &Set{Name: es.Name, Vars: es.Vars}
	if err := s.Put(key, value); err != nil {
		return err
	}
	m.cfg.AddOrUpdateSet(config.EnvSet{Name: s.Name, Vars: s.Vars})
	return nil
}

// GetVar retrieves a variable from a named set.
func (m *Manager) GetVar(setName, key string) (string, error) {
	es, err := m.cfg.GetSet(setName)
	if err != nil {
		return "", err
	}
	s := &Set{Name: es.Name, Vars: es.Vars}
	return s.Get(key)
}

// DeleteVar removes a variable from a named set.
func (m *Manager) DeleteVar(setName, key string) error {
	es, err := m.cfg.GetSet(setName)
	if err != nil {
		return err
	}
	s := &Set{Name: es.Name, Vars: es.Vars}
	if err := s.Delete(key); err != nil {
		return err
	}
	m.cfg.AddOrUpdateSet(config.EnvSet{Name: s.Name, Vars: s.Vars})
	return nil
}
