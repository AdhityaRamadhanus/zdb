package zdb

type TreeIterator struct {
	stack []*Node
	curr  *Node
}

func NewTreeIterator(t OrderStatisticTree) *TreeIterator {
	return &TreeIterator{
		curr: t.Root(),
	}
}

func (it *TreeIterator) moveToStart() {
	for it.curr != nil {
		it.stack = append(it.stack, it.curr)
		it.curr = it.curr.left
	}
}

func (it *TreeIterator) Seek(node *Node) {
	if node == nil {
		// find beginning
		it.moveToStart()
		return
	}

	for it.curr != nil {
		comp := compareNode(it.curr, node)
		if comp > 0 {
			it.stack = append(it.stack, it.curr)
			it.curr = it.curr.left
		} else if comp < 0 {
			it.curr = it.curr.right
		} else { // it.curr == node
			it.stack = append(it.stack, it.curr)
			it.curr = nil
			return
		}
	}

	// not found
	it.curr = nil
	it.stack = []*Node{}
}

func (it *TreeIterator) Next() *Node {
	// inorder traversal with stack
	for it.curr != nil || len(it.stack) > 0 {
		var next *Node
		if it.curr != nil {
			it.stack = append(it.stack, it.curr)
			it.curr = it.curr.left
		} else {
			next = it.stack[len(it.stack)-1]
			it.stack = it.stack[:len(it.stack)-1]
			it.curr = next.right
			return next
		}
	}

	return nil
}
