package storage

type List[T any] interface {
	Append(key string, value T)
	Get(key string) *T
	GetHead() *LinkedNode[T]
	RemoveHead()
	RemoveNext(*LinkedNode[T])
}

type LinkedNode[T any] struct {
	Key   string
	Value T
	Next  *LinkedNode[T]
}

type LinkedList[T any] struct {
	Head *LinkedNode[T]
	Tail *LinkedNode[T]
}

func (l *LinkedList[T]) Append(key string, value T) {
	node := &LinkedNode[T]{Key: key, Value: value}
	if l.Head == nil {
		l.Head = node
		l.Tail = node
	} else {
		l.Tail.Next = node
		l.Tail = node
	}
}

func (l *LinkedList[T]) Get(key string) *T {
	node := l.Head

	for node != nil && node.Key != key {
		node = node.Next
	}

	return &node.Value
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
