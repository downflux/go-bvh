package stack

type S[T any] struct {
	data []T
}

func New[T any](data []T) *S[T] {
	return &S[T]{
		data: data,
	}
}

func (s *S[T]) Push(p T) {
	s.data = append(s.data, p)
}

func (s *S[T]) Pop() (T, bool) {
	var d T
	if len(s.data) == 0 {
		return d, false
	}
	d, s.data = s.data[len(s.data)-1], s.data[:len(s.data)-1]
	return d, true
}

func (s *S[T]) Data() []T {
	data := make([]T, len(s.data))
	copy(data, s.data)
	return data
}
