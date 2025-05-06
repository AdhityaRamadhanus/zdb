package zdb

type Node struct {
	key    string
	score  float64
	height int
	count  int
	left   *Node
	right  *Node
}

func NewNode(key string, score float64) *Node {
	return &Node{
		key:    key,
		score:  score,
		count:  1,
		height: 1,
		left:   nil,
		right:  nil,
	}
}

func (n *Node) Key() string {
	if n == nil {
		return ""
	}

	return n.key
}

func (n *Node) Score() (score float64) {
	if n == nil {
		return
	}

	return n.score
}

func (n *Node) Count() int {
	if n == nil {
		return 0
	}

	return n.count
}

func (n *Node) Height() int {
	if n == nil {
		return 0
	}

	return n.height
}

func (n *Node) Balance() int {
	if n == nil {
		return 0
	}

	return n.left.Height() - n.right.Height()
}

func (n *Node) inorderPredecessor() *Node {
	if n == nil || n.left == nil {
		return nil
	}

	curr := n.left
	for curr.right != nil {
		curr = curr.right
	}

	return curr
}

func (n *Node) rankByParent(p *Node, pRank int) (rank int) {
	if p == nil {
		return 1 + n.left.Count()
	}

	if n.score < p.score || (n.score == p.score && n.key < p.key) {
		rank = pRank - n.right.Count() - 1
	} else {
		rank = 1 + pRank + n.left.Count()
	}

	return rank
}

func (n *Node) rankByParentReverse(p *Node, pRank int) (rank int) {
	if p == nil {
		return 1 + n.right.Count()
	}

	if n.score > p.score || (n.score == p.score && n.key > p.key) {
		rank = pRank - n.left.Count() - 1
	} else {
		rank = 1 + pRank + n.right.Count()
	}

	return rank
}

func (n *Node) updateHeightAndCount() {
	n.height = 1 + max(n.left.Height(), n.right.Height())
	n.count = 1 + n.left.Count() + n.right.Count()
}

func compareNode(x, y *Node) int {
	// returns -1 if x < y, 0 if x==y and 1 if x > y
	if x.score < y.score || (x.score == y.score && x.key < y.key) {
		return -1
	} else if x.score > y.score || (x.score == y.score && x.key > y.key) {
		return 1
	} else {
		return 0
	}
}
