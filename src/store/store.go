package store

import (
	"sync"
	"time"

	"github.com/NicolasMartino/data-analysis/src/models"
)

// Type safe typed map with mutex
// from https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c
type SafeMap struct {
	sync.RWMutex
	//TODO use into generics
	internal      map[string]models.CacheUrlInfo
	lastUpdated   map[string]time.Time
	cacheLifespan time.Duration
}

func NewSafeMap(cacheLifespan time.Duration) *SafeMap {
	return &SafeMap{
		internal:      make(map[string]models.CacheUrlInfo),
		lastUpdated:   make(map[string]time.Time),
		cacheLifespan: cacheLifespan,
	}
}

func (m *SafeMap) Load(key string) (value models.CacheUrlInfo, ok bool) {
	m.RLock()
	defer m.RUnlock()

	sinceLastCacheUpdate := time.Since(m.lastUpdated[key])
	if m.cacheLifespan > sinceLastCacheUpdate {
		value, ok = m.internal[key]
	}
	return
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
	m.lastUpdated[key] = time.Now()
}
