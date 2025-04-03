package mapx

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable](n ...int) *Set[T] {
	x := 0
	if len(n) > 0 {
		x = n[0]
	}
	return &Set[T]{m: make(map[T]struct{}, x)}
}

func (s *Set[T]) Add(v T) {
	s.m[v] = struct{}{}
}

func (s *Set[T]) Remove(v T) {
	delete(s.m, v)
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

func (s *Set[T]) Size() int {
	return len(s.m)
}

func (s *Set[T]) Clear() {
	s.m = make(map[T]struct{})
}

func (s *Set[T]) Values() []T {
	res := make([]T, 0, len(s.m))
	for k := range s.m {
		res = append(res, k)
	}
	return res
}
