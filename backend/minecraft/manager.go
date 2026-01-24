package minecraft

import (
	"sync"
)

type Manager struct {
	instances map[string]*Instance
	mu        sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		instances: make(map[string]*Instance),
	}
}

func (m *Manager) AddInstance(id string, runDir, jarName string) *Instance {
	m.mu.Lock()
	defer m.mu.Unlock()

	inst := NewInstance(id, runDir, jarName)
	m.instances[id] = inst
	return inst
}

func (m *Manager) GetInstance(id string) (*Instance, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inst, exists := m.instances[id]
	return inst, exists
}

func (m *Manager) RemoveInstance(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if inst, exists := m.instances[id]; exists {
		inst.Stop()
		delete(m.instances, id)
	}
}
