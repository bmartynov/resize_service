package memory

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestShardSetGet(t *testing.T) {
	s := newShard(1024 * 1024 * 100)

	tds := make([]struct {
		key     uint64
		payload []byte
	}, 1000)

	for idx := range tds {
		tds[idx].key = rand.Uint64()
		tds[idx].payload = []byte(uuid.New().String())
	}

	for _, td := range tds {
		s.set(td.key, td.payload)
	}

	for _, td := range tds {
		bts, ok := s.get(td.key)
		if !ok {
			t.Fatal("must be ok: get failed")
		}
		if !reflect.DeepEqual(bts, td.payload) {
			t.Fatal("must be ok: equal failed")
		}
	}
}
