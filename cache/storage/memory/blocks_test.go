package memory

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestBlocksGetFree(t *testing.T) {
	b := blocks{block{0, 1024}}

	bb, ok := b.get(512)

	t.Log(bb, ok)

	t.Log("after get", b)

	b.free(bb)

	t.Log("after free", b)
}

func TestBlocksDefrag(t *testing.T) {
	for _, tc := range []struct {
		in       blocks
		expected blocks
	}{
		{
			in: blocks{
				{start: 100, stop: 200},
				{start: 200, stop: 300},
				{start: 300, stop: 400},
			},
			expected: blocks{{start: 100, stop: 400}},
		},
	} {
		tc.in.defrag()
		if !reflect.DeepEqual(tc.expected, tc.in) {
			t.Fatal("missmatch")
		}
	}

}

func BenchmarkBlockDelete(b *testing.B) {
	blocks := blocks{}

	blocksNum := 10000

	for i := 0; i < blocksNum; i++ {
		blocks = append(blocks, block{})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {

		b.StopTimer()
		idx := rand.Intn(blocksNum)
		b.StartTimer()

		blocks.delete(idx)
	}
}
