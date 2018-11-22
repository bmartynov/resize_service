package memory

type stats struct {
	used   uint64
	unused uint64
}

func (s *stats) incUsed(size uint64) {
	s.used += size
	s.unused -= size
}
