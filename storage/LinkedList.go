package storage

type List[T comparable] interface {
	Append(value T)
	Get(key T) *LinkedNode[T]
	GetHead() *LinkedNode[T]
	RemoveHead()
	RemoveNext(*LinkedNode[T])
}

type LinkedNode[T comparable] struct {
	Value T
	Next  *LinkedNode[T]
}

type LinkedList[T comparable] struct {
	Head *LinkedNode[T]
	Tail *LinkedNode[T]
}

func (l *LinkedList[T]) Append(value T) {
	node := &LinkedNode[T]{Value: value}
	if l.Head == nil {
		l.Head = node
		l.Tail = node
	} else {
		l.Tail.Next = node
		l.Tail = node
	}
}

func (l *LinkedList[T]) Get(key T) *LinkedNode[T] {
	node := l.Head

	for node != nil && node.Value != key {
		node = node.Next
	}

	return node
}

func (l *LinkedList[T]) GetHead() *LinkedNode[T] {
	return l.Head
}

func (l *LinkedList[T]) RemoveHead() {
	if l.Head == nil {
		return
	}

	if l.Head == l.Tail {
		l.Tail = l.Tail.Next
	}
	l.Head = l.Head.Next
}

func (l *LinkedList[T]) RemoveNext(node *LinkedNode[T]) {
	if node.Next == l.Tail {
		l.Tail = node
	}
	node.Next = node.Next.Next
}
