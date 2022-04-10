package store

import (
	"sync"

	"github.com/NicolasMartino/data-analysis/src/models"
)

// Type safe typed map with mutex
// from https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c
type SafeMap struct {
	sync.RWMutex
	internal map[string]models.CacheUrlInfo
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		internal: make(map[string]models.CacheUrlInfo),
	}
}

func (m *SafeMap) Load(key string) (value models.CacheUrlInfo, ok bool) {
	m.RLock()
	defer m.RUnlock()
	result, ok := m.internal[key]
	return result, ok
}

func (m *SafeMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.internal, key)
}

func (m *SafeMap) Store(key string, value models.CacheUrlInfo) {
	m.Lock()
	defer m.Unlock()
	m.internal[key] = value
}
