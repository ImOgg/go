package migrations

import (
	"sort"
	"sync"
)

var (
	registry = make(map[string]Migration)
	mu       sync.RWMutex
)

// Register 註冊 migration
func Register(m Migration) {
	mu.Lock()
	defer mu.Unlock()
	registry[m.Version()] = m
}

// All 獲取所有已註冊的 migrations（按版本排序）
func All() []Migration {
	mu.RLock()
	defer mu.RUnlock()
	
	var migrations []Migration
	for _, m := range registry {
		migrations = append(migrations, m)
	}
	
	// 按版本號排序
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version() < migrations[j].Version()
	})
	
	return migrations
}

// Get 根據版本號獲取 migration
func Get(version string) (Migration, bool) {
	mu.RLock()
	defer mu.RUnlock()
	m, ok := registry[version]
	return m, ok
}
