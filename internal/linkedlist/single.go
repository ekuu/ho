package linkedlist

type Single[T any] interface {
	Len() int
	Shift() (T, bool)
	Append(...T)
}

type node[T any] struct {
	val  T
	next *node[T]
}

type single[T any] struct {
	head *node[T]
	tail *node[T]
	len  int
}

func NewSingle[T any]() *single[T] {
	return new(single[T])
}

func (s *single[T]) Len() int {
	return s.len
}

func (s *single[T]) Shift() (t T, ok bool) {
	if s.head == nil {
		return t, false
	}
	t = s.head.val
	s.head = s.head.next
	if s.head == nil {
		s.tail = nil
	}
	s.len--
	return t, true
}

func (s *single[T]) Append(ts ...T) {
	for _, t := range ts {
		s.append(t)
	}
}

func (s *single[T]) append(t T) {
	n := &node[T]{val: t}
	if s.head == nil {
		s.head = n
	} else {
		s.tail.next = n
	}
	s.tail = n
	s.len++
}
