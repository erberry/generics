package array

type Node[T any] struct {
	Next  *Node[T]
	Prev  *Node[T]
	Val   *T
	Index int
}

// FixSizeArray It is a fixed size array with O(1) time complexity for adding and deleting elements.
type FixSizeArray[T any] struct {
	nodes    []Node[T]
	freeHead *Node[T]
	usedHead *Node[T]
	size     int
	length   int
}

func New[T any](size int) *FixSizeArray[T] {
	nodes := make([]Node[T], size)
	freeHead := &nodes[0]
	head := freeHead
	for i := 1; i < len(nodes); i++ {
		head.Next = &nodes[i]
		head = head.Next
		head.Index = i
	}
	return &FixSizeArray[T]{
		nodes:    nodes,
		freeHead: freeHead,
		usedHead: nil,
		size:     size,
	}
}

func (arr *FixSizeArray[T]) Add(val *T) (index int, ok bool) {
	if arr.freeHead == nil {
		ok = false
		return
	}

	free := arr.freeHead
	arr.freeHead = arr.freeHead.Next
	free.Next = nil

	free.Val = val
	if arr.usedHead == nil {
		arr.usedHead = free
	} else {
		tmp := arr.usedHead.Next
		arr.usedHead.Next = free
		free.Prev = arr.usedHead
		free.Next = tmp
		if tmp != nil {
			tmp.Prev = free
		}
	}
	arr.length++

	return free.Index, true
}

func (arr *FixSizeArray[T]) Del(index int) (exist bool) {
	if index >= arr.size {
		panic("index out of size")
	}
	if arr.nodes[index].Val == nil {
		return false
	}

	used := &arr.nodes[index]
	arr.nodes[index].Val = nil
	if used == arr.usedHead {
		newHead := arr.usedHead.Next
		arr.usedHead.Next = nil
		arr.usedHead = newHead
		if newHead != nil {
			newHead.Prev = nil
		}
	} else {
		prev, next := used.Prev, used.Next
		used.Prev, used.Next = nil, nil
		prev.Next = next
		if next != nil {
			next.Prev = prev
		}
	}
	tmp := arr.freeHead
	arr.freeHead = used
	arr.freeHead.Next = tmp
	arr.length--

	return true
}

func (arr *FixSizeArray[T]) ForRange(f func(val *T)) {
	head := arr.usedHead
	for head != nil {
		f(head.Val)
		head = head.Next
	}
}
