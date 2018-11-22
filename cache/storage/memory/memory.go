package memory

import (
	"hash/fnv"
	"io"
)

const (
	SHARDS = 16
)

func shardFnImp(key string, total int) (int, uint64) {
	h := fnv.New64()
	h.Write([]byte(key))

	u := h.Sum64()

	return int(u) % total, u
}

type shardFn func(key string, total int) (int, uint64)

type memory struct {
	shards  [SHARDS]*shard
	shardFn shardFn
}

func (m *memory) Store(key string, r io.Reader) error {
	shardID, hashedKEY := m.shardFn(key, SHARDS)

	shard := m.shards[shardID]

	err := shard.set(hashedKEY, r)
	if err != nil {
		return err
	}

	return nil
}

func (m *memory) Load(key string) (io.Reader, error) {
}

func (m *memory) Remove(key string) error {
}

func New(shardSize int) *memory {
	m := memory{
		shardFn: shardFnImp,
		shards:  [16]*shard{},
	}

	for idx := range m.shards {
		m.shards[idx] = newShard(uint64(shardSize))
	}

	return &m
}
