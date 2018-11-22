package memory

import (
	"fmt"
	"sort"
)

type block struct {
	start uint64
	stop  uint64
}

func (b *block) size() uint64 {
	return uint64(b.stop - b.start)
}

func (b block) String() string {
	return fmt.Sprintf("block[from: %d, to: %d, size: %d]", b.start, b.stop, b.stop-b.start)
}

func (b *block) reduceSize(s uint64) {
	b.start += s
}

type blocks []block

func (b *blocks) defrag() {
	sblocks := ([]block)(*b)

	sort.Slice(sblocks, func(i, j int) bool {
		return sblocks[i].stop < sblocks[j].start
	})

	for _, b := range sblocks {
		_ = b
	}
}

func (b *blocks) free(bb block) {
	sblocks := ([]block)(*b)

	var current *block

	for idx := range sblocks {
		current = &sblocks[idx]

		if current.stop == bb.start {
			current.stop = bb.stop
			return
		}

		if current.start == bb.stop {
			current.start = bb.start
			return
		}
	}

	sblocks = append(sblocks, bb)
}

func (b *blocks) delete(ids ...int) {
	sblocks := ([]block)(*b)

	for _, id := range ids {
		sblocks = sblocks[:id+copy(sblocks[id:], sblocks[id+1:])]
	}
}

func (b *blocks) get(size uint64) (bb block, ok bool) {
	sblocks := ([]block)(*b)

	var toDelete int
	var current *block

	for idx := range *b {
		current = &sblocks[idx]

		if current.size() >= size {
			bb = block{
				start: sblocks[idx].start,
				stop:  sblocks[idx].start + size,
			}
			ok = true

			current.reduceSize(size)

			if current.size() == 0 {
				toDelete = idx
			}

			break
		}
	}

	b.delete(toDelete)

	return
}
