package memory

import (
	"io"
	"sync"

	"github.com/pkg/errors"
)

var (
	errNoEnoughMem = errors.New("not enough mem")
	errAllocFailed = errors.New("alloc failed")
)

type shard struct {
	sync.RWMutex
	stats  stats  // memory stats
	blocks blocks // available mem blocks
	data   []byte // contains all data
	buff   []byte // for read from src
	keys   map[uint64]block
}

func (s *shard) read(start, stop int) io.Reader {

}

func (s *shard) allKeys() (keys []uint64) {
	s.RLock()

	keys = make([]uint64, 0, len(s.keys))
	for k, _ := range s.keys {
		keys = append(keys, k)
	}
	s.RUnlock()

	return keys
}

func (s *shard) set(key uint64, r io.Reader) error {
	s.Lock()

	size, err := r.Read(s.buff)
	if err != nil {
		s.Unlock()
		return err
	}

	if s.stats.unused < uint64(size) {
		s.Unlock()
		return errNoEnoughMem
	}

	b, ok := s.blocks.get(uint64(size))
	if !ok {
		// todo: handle defrag
		return errAllocFailed
	}

	copy(s.data[b.start:b.stop], s.buff[0:size])

	s.keys[key] = b
	s.stats.incUsed(uint64(size))

	s.Unlock()

	return nil
}

func (s *shard) get(key uint64) (bts []byte, ok bool) {
	s.RLock()

	b, ok := s.keys[key]
	if !ok {
		s.RUnlock()
		return nil, false
	}

	bts, ok = s.data[b.start:b.stop], true
	s.RUnlock()

	return
}

func (s *shard) delete(key uint64) {
	s.Lock()
	b, ok := s.keys[key]
	if !ok {
		s.Unlock()
		return
	}

	delete(s.keys, key)
	s.blocks.free(b)

	s.Unlock()
}

func newShard(size, maxSize uint64) *shard {
	return &shard{
		stats: stats{
			used:   0,
			unused: uint64(size),
		},
		data:   make([]byte, size),
		buff:   make([]byte, maxSize),
		keys:   make(map[uint64]block),
		blocks: blocks{{0, size}},
	}
}
