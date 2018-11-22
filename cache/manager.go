package cache

import (
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
	ErrExpired  = errors.New("cleanUpInterval")
)

type Config struct {
	Ttl             time.Duration `envconfig:"TTL" required:"true"`
	CleanUpInterval time.Duration `envconfig:"CLEAN_UP_INTERVAL" required:"true"`
}

func IsNoCacheError(e error) bool {
	if e == ErrExpired || e == ErrNotFound {
		return true
	}
	return false
}

type cacheItem struct {
	storageID string
	created   time.Time
}

func (i *cacheItem) isExpired(ttl time.Duration) bool {
	return i.created.Add(ttl).
		Before(time.Now())
}

func newCacheItem(sID string) *cacheItem {
	return &cacheItem{
		storageID: sID,
		created:   time.Now(),
	}
}

type manager struct {
	sync.RWMutex
	config  Config
	storage Storage
	items   map[string]*cacheItem
	close   chan struct{}
}

func (m *manager) removeExpiredLoop() {
	ticker := time.NewTicker(m.config.CleanUpInterval)

	for {
		select {
		case <-m.close:
			ticker.Stop()

		case <-ticker.C:
			m.removeExpired()
		}
	}
}

func (m *manager) removeExpired() {
	m.Lock()

	itemsToDelete := make([]string, 0)

	for key, i := range m.items {
		if i.isExpired(m.config.Ttl) {
			itemsToDelete = append(itemsToDelete, key)
		}
	}

	for _, toDelete := range itemsToDelete {
		delete(m.items, toDelete)
	}

	m.Unlock()
}

func (m *manager) Add(key string, r io.Reader) error {
	m.Lock()
	storageID := uuid.New().String()

	m.items[key] = newCacheItem(storageID)

	err := m.storage.Store(storageID, r)
	if err != nil {
		m.Unlock()
		return err
	}

	m.Unlock()
	return nil
}

func (m *manager) Get(key string) (Item, error) {
	m.RLock()
	i, ok := m.items[key]
	if !ok {
		m.RUnlock()
		return nil, ErrNotFound
	}

	if i.isExpired(m.config.Ttl) {
		m.RUnlock()
		return nil, ErrExpired
	}

	r, err := m.storage.Load(i.storageID)
	if err != nil {
		m.RUnlock()
		return nil, err
	}

	m.RUnlock()

	return r, nil
}

func (m *manager) Start() error {
	go m.removeExpiredLoop()

	return nil
}

func (m *manager) Stop() error {
	close(m.close)

	return nil
}

func NewManager(s Storage, config Config) Manager {
	m := manager{
		config:  config,
		storage: s,
		close:   make(chan struct{}),
		items:   make(map[string]*cacheItem),
	}

	return &m
}
