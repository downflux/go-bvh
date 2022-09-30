package stack

type S[T any] []T

func New[T any](data []T) *S[T] { return (*S[T])(&data) }

func (s *S[T]) Len() int { return len(*s) }

func (s *S[T]) Push(p T) { *s = append(*s, p) }

func (s *S[T]) Pop() (T, bool) {
	var d T
	if s.Len() == 0 {
		return d, false
	}

	d, *s = (*s)[s.Len()-1], (*s)[:s.Len()-1]
	return d, true
}
