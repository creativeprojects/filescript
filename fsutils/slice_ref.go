package fsutils

type sliceRef struct {
	ref []string
}

func newSliceRef(size int) *sliceRef {
	return &sliceRef{
		ref: make([]string, 0, size),
	}
}

func (s *sliceRef) append(value string) {
	s.ref = append(s.ref, value)
}
