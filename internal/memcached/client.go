package memcached

import (
	lg "discogo/internal/logger"
	lgm "discogo/internal/logger/messages"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedClient struct {
	client *memcache.Client
}

// NewMemcachedClient initializes a new Memcached client with the given address.
func NewMemcachedClient(address string) *MemcachedClient {
	client := memcache.New(address)
	return &MemcachedClient{client: client}
}

// Set stores a value in Memcached with the specified key and expiration time.
func (m *MemcachedClient) Set(key string, value []byte, expiration time.Duration) error {
	err := m.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(expiration.Seconds()),
	})
	if err != nil {
		lg.Error(lgm.MessageR(lgm.ErrorMemcachedFailedToRetrieveServerAddress, err))
		return err
	}
	return nil
}

// Get retrieves a value from Memcached by its key.
func (m *MemcachedClient) Get(key string) ([]byte, error) {
	item, err := m.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			lg.Info(lgm.MessageR(lgm.InfoMemcachedKeyNotFound, key))
			return nil, nil
		}
		lg.Error(lgm.MessageR(lgm.ErrorMemcachedFailedToRetrieveServerAddress, err))
		return nil, err
	}
	return item.Value, nil
}

// Ping checks the connection to the Memcached server by performing a dummy Set/Get.
func (m *MemcachedClient) Ping() error {
	testKey := "__ping__"
	testValue := []byte("ok")
	err := m.Set(testKey, testValue, 1*time.Second)
	if err != nil {
		lg.Error(lgm.MessageR(lgm.ErrorMemcachedConnectionFailed, err))
		return err
	}
	_, err = m.Get(testKey)
	if err != nil {
		lg.Error(lgm.MessageR(lgm.ErrorMemcachedConnectionFailed, err))
		return err
	}
	return nil
}
