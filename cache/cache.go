package cache

import (
	"io"
)

type Item interface {
	io.Reader
}

type Manager interface {
	Stop() error
	Start() error
	Add(string, io.Reader) error
	Get(string) (Item, error)
}

// Storage stores data in somewhere, eg: disk, ram
type Storage interface {
	Store(key string, r io.Reader) error
	Load(key string) (io.Reader, error)
	Remove(key string) error
}
