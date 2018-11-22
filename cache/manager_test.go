package cache

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

type mockStorage struct {
	sync.Mutex
	items map[string]*bytes.Buffer
}

func (m *mockStorage) Store(key string, r io.Reader) error {
	m.Lock()
	defer m.Unlock()

	b := bytes.NewBuffer(nil)

	io.Copy(b, r)

	m.items[key] = b

	return nil
}

func (m *mockStorage) Load(key string) (r io.Reader, err error) {
	m.Lock()
	defer m.Unlock()

	r = m.items[key]

	return
}

func (m *mockStorage) Remove(key string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.items, key)

	return nil
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		items: make(map[string]*bytes.Buffer),
	}
}

func TestManager_Add_Get(t *testing.T) {
	m := NewManager(newMockStorage(), Config{
		Ttl:             time.Second * 5,
		CleanUpInterval: time.Second * 2,
	})

	const (
		payload = "hello i am string buffer"
		key     = "some_key"
	)

	sbOrig := bytes.NewBufferString(payload)

	err := m.Add(key, sbOrig)
	if err != nil {
		t.Fatal(err)
	}

	r, err := m.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	rdb, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if string(rdb) != payload {
		t.Fatal("content missmatch")
	}
}

func TestManager_Add_Get_Expired_Removed(t *testing.T) {
	m := NewManager(newMockStorage(), Config{
		Ttl:             time.Second * 1,
		CleanUpInterval: time.Second * 1,
	})

	const (
		payload = "hello i am string buffer"
		key     = "some_key"
	)

	sbOrig := bytes.NewBufferString(payload)

	err := m.Add(key, sbOrig)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	_, err = m.Get(key)
	if err != ErrExpired {
		t.Fatal("must be cleanUpInterval")
	}

	m.(*manager).removeExpired()

	_, err = m.Get(key)
	if err != ErrNotFound {
		t.Fatal("must be not found")
	}
}
