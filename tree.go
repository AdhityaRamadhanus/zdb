package zdb

type Tree struct {
	root    *Node
	HashMap map[uint64]float64
	hasher  fnv64a
}

func NewTree() OrderStatisticTree {
	return &Tree{
		HashMap: make(map[uint64]float64),
		hasher:  fnv64a{},
	}
}

func (t *Tree) IsEmpty() bool {
	return t == nil || t.root == nil
}

func (t *Tree) Root() *Node {
	if t == nil {
		return nil
	}

	return t.root
}

func (t *Tree) GetScore(key string) (float64, error) {
	hashedKey := t.hasher.Sum64(key)
	score, exists := t.HashMap[hashedKey]
	if !exists {
		return 0, errNotFound
	}

	return score, nil
}

func (t *Tree) Add(key string, score float64) {
	hashedKey := t.hasher.Sum64(key)
	if oldScore, exists := t.HashMap[hashedKey]; exists {
		// delete and insert new
		t.root = t.deleteRec(t.root, NewNode(key, oldScore))
	}
	t.root = t.insertRec(t.root, NewNode(key, score))
	t.HashMap[hashedKey] = score
}

func (t *Tree) Remove(key string) {
	hashedKey := t.hasher.Sum64(key)
	score, exists := t.HashMap[hashedKey]
	if !exists {
		return
	}
	t.root = t.deleteRec(t.root, NewNode(key, score))
	delete(t.HashMap, hashedKey)
}

func (t *Tree) Rank(key string) int {
	if t.root == nil {
		return -1
	}

	hashedKey := t.hasher.Sum64(key)
	score, exists := t.HashMap[hashedKey]
	if !exists {
		return -1
	}

	seek := NewNode(key, score)
	curr := t.root
	var parent *Node
	parentRank := -1

	for curr != nil {
		currRank := curr.rankByParent(parent, parentRank)
		currComparison := compareNode(curr, seek)
		if currComparison > 0 {
			parent = curr
			parentRank = currRank
			curr = curr.left
		} else if currComparison < 0 {
			parent = curr
			parentRank = currRank
			curr = curr.right
		} else {
			return currRank
		}
	}

	return -1
}

func (t *Tree) RankReverse(key string) int {
	if t.root == nil {
		return -1
	}

	hashedKey := t.hasher.Sum64(key)
	score, exists := t.HashMap[hashedKey]
	if !exists {
		return -1
	}

	seek := NewNode(key, score)
	curr := t.root
	var parent *Node
	parentRank := -1

	for curr != nil {
		currRank := curr.rankByParentReverse(parent, parentRank)
		currComparison := compareNode(curr, seek)

		if currComparison > 0 {
			parent = curr
			parentRank = currRank
			curr = curr.left
		} else if currComparison < 0 {
			parent = curr
			parentRank = currRank
			curr = curr.right
		} else {
			return currRank
		}
	}

	return -1
}

func (t *Tree) Select(idx int) *Node {
	if t.root == nil || idx > t.root.Count() {
		return nil
	}

	curr := t.root
	var parent *Node
	parentRank := -1

	for curr != nil {
		currRank := curr.rankByParent(parent, parentRank)

		if idx < currRank {
			parent = curr
			parentRank = currRank
			curr = curr.left
		} else if idx > currRank {
			parent = curr
			parentRank = currRank
			curr = curr.right
		} else {
			return curr
		}
	}

	return nil
}

func (t *Tree) SelectReverse(idx int) *Node {
	if t.root == nil || idx > t.root.Count() {
		return nil
	}

	curr := t.root
	var parent *Node
	parentRank := -1

	for curr != nil {
		currRank := curr.rankByParentReverse(parent, parentRank)

		if idx < currRank {
			parent = curr
			parentRank = currRank
			curr = curr.right
		} else if idx > currRank {
			parent = curr
			parentRank = currRank
			curr = curr.left
		} else {
			return curr
		}
	}

	return nil
}

func (t *Tree) RangeByIndex(start, stop int) (nodes []Node) {
	start += 1
	stop += 1

	if t.root == nil {
		return
	}

	type stackElmt struct {
		Curr     *Node
		CurrRank int
	}
	stack := []stackElmt{}

	curr := t.root
	var parent *Node
	parentRank := -1

	for curr != nil || len(stack) != 0 {
		if curr == nil && len(stack) != 0 {
			se := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			curr = se.Curr
			currRank := se.CurrRank

			nodes = append(nodes, *curr)
			parent = curr
			parentRank = currRank
			curr = curr.right
			continue
		}

		currRank := curr.rankByParent(parent, parentRank)
		parent = curr
		parentRank = currRank

		if currRank < start {
			curr = curr.left
		} else if currRank >= start && currRank <= stop {
			stack = append(stack, stackElmt{Curr: curr, CurrRank: currRank})
			curr = curr.left
		} else { // currRank > stop
			curr = curr.left
		}
	}

	return nodes
}

func (t *Tree) RangeByIndexReverse(start, stop int) (nodes []Node) {
	start += 1
	stop += 1

	if t.root == nil {
		return
	}

	type stackElmt struct {
		Curr     *Node
		CurrRank int
	}
	stack := []stackElmt{}

	curr := t.root
	var parent *Node
	parentRank := -1

	for curr != nil || len(stack) != 0 {
		if curr == nil && len(stack) != 0 {
			se := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			curr = se.Curr
			currRank := se.CurrRank

			nodes = append(nodes, *curr)
			parent = curr
			parentRank = currRank
			curr = curr.left
			continue
		}

		currRank := curr.rankByParentReverse(parent, parentRank)
		parent = curr
		parentRank = currRank

		if currRank < start {
			curr = curr.left
		} else if currRank >= start && currRank <= stop {
			stack = append(stack, stackElmt{Curr: curr, CurrRank: currRank})
			curr = curr.right
		} else { // currRank > stop
			curr = curr.right
		}
	}

	return nodes
}

func (t *Tree) rankByScoreLowerBound(score float64) int {
	if t.root == nil {
		return 0
	}

	// find min
	curr := t.root
	var parent *Node
	parentRank := -1

	for curr != nil {
		currRank := curr.rankByParent(parent, parentRank)
		parent = curr
		parentRank = currRank
		if score <= curr.score {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}

	if score <= parent.score {
		return parentRank
	} else {
		return parentRank + 1
	}
}

func (t *Tree) rankByScoreUpperBound(score float64) int {
	if t.root == nil {
		return 0
	}

	// find max
	curr := t.root
	var parent *Node
	parentRank := -1

	for curr != nil {
		currRank := curr.rankByParent(parent, parentRank)
		parent = curr
		parentRank = currRank
		if score < curr.score {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}

	if score < parent.score {
		return parentRank - 1
	} else {
		return parentRank + 1
	}
}

func (t *Tree) CountByScore(min, max float64) int {
	if t.root == nil {
		return 0
	}

	minRank := t.rankByScoreLowerBound(min)
	maxRank := t.rankByScoreUpperBound(max)

	return maxRank - minRank
}

func (t *Tree) RangeByScore(min, max float64) (nodes []Node) {
	if t.root == nil {
		return
	}

	stack := []*Node{}

	curr := t.root
	for curr != nil || len(stack) != 0 {
		if curr == nil && len(stack) != 0 {
			curr = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			nodes = append(nodes, *curr)

			curr = curr.right
			continue
		}

		if curr.score < min {
			curr = curr.right
		} else if curr.score >= min && curr.score <= max {
			stack = append(stack, curr)
			curr = curr.left
		} else { // curr.score > max
			curr = curr.left
		}
	}

	return nodes
}

func (t *Tree) RangeByScoreReverse(min, max float64) (nodes []Node) {
	if t.root == nil {
		return
	}

	stack := []*Node{}

	curr := t.root
	for curr != nil || len(stack) != 0 {
		if curr == nil && len(stack) != 0 {
			curr = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			nodes = append(nodes, *curr)

			curr = curr.left
			continue
		}

		if curr.score < min {
			curr = curr.right
		} else if curr.score >= min && curr.score <= max {
			stack = append(stack, curr)
			curr = curr.right
		} else { // curr.score > max
			curr = curr.left
		}
	}

	return nodes
}

func (t *Tree) RangeByLex(minKey, maxKey string) (nodes []Node) {
	if t.root == nil {
		return
	}

	stack := []*Node{}

	curr := t.root
	for curr != nil || len(stack) != 0 {
		if curr == nil && len(stack) != 0 {
			curr = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			nodes = append(nodes, *curr)

			curr = curr.right
			continue
		}

		if curr.key < minKey {
			curr = curr.right
		} else if curr.key >= minKey && curr.key <= maxKey {
			stack = append(stack, curr)
			curr = curr.left
		} else { // curr.key > max
			curr = curr.left
		}
	}

	return nodes
}

func (t *Tree) RangeByLexReverse(minKey, maxKey string) (nodes []Node) {
	if t.root == nil {
		return
	}

	stack := []*Node{}

	curr := t.root
	for curr != nil || len(stack) != 0 {
		if curr == nil && len(stack) != 0 {
			curr = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			nodes = append(nodes, *curr)

			curr = curr.left
			continue
		}

		if curr.key < minKey {
			curr = curr.right
		} else if curr.key >= minKey && curr.key <= maxKey {
			stack = append(stack, curr)
			curr = curr.right
		} else { // curr.key > max
			curr = curr.left
		}
	}

	return nodes
}

func (t *Tree) Diff(other OrderStatisticTree) (diff OrderStatisticTree) {
	it := NewTreeIterator(t)
	it.Seek(nil)

	diff = NewTree()
	for {
		next := it.Next()
		if next == nil {
			break
		}

		_, err := other.GetScore(next.key)
		// either key doesn't exist in other or exists with different score
		if err == errNotFound {
			diff.Add(next.key, next.score)
		}
	}

	return diff
}

func (t *Tree) Inter(other OrderStatisticTree, aggFunc AggFunc) (inter OrderStatisticTree) {
	it := NewTreeIterator(t)
	it.Seek(nil)

	inter = NewTree()
	for {
		next := it.Next()
		if next == nil {
			break
		}

		score, err := other.GetScore(next.key)
		// either key doesn't exist in other or exists with different score
		if err != errNotFound {
			inter.Add(next.key, aggFunc(next.score, score))
		}
	}

	return inter
}

func (t *Tree) Union(other OrderStatisticTree, aggFunc AggFunc) (union OrderStatisticTree) {
	it := NewTreeIterator(t)
	it.Seek(nil)

	it2 := NewTreeIterator(other)
	it2.Seek(nil)

	union = NewTree()
	for {
		next := it.Next()
		if next == nil {
			break
		}

		score, err := other.GetScore(next.key)
		inserted := NewNode(next.key, next.score)
		// either key doesn't exist in other or exists with different score
		if err != errNotFound {
			inserted.score = aggFunc(next.score, score)
		}

		union.Add(inserted.key, inserted.score)
	}

	for {
		next := it2.Next()
		if next == nil {
			break
		}

		_, err := union.GetScore(next.key)
		if err == errNotFound {
			union.Add(next.key, next.score)
		}
	}

	return union
}

func (t *Tree) insertRec(root *Node, inserted *Node) *Node {
	if root == nil {
		return inserted
	}

	rootComparison := compareNode(root, inserted)

	if rootComparison > 0 { // insert left
		root.left = t.insertRec(root.left, inserted)
	} else if rootComparison < 0 { // insert right
		root.right = t.insertRec(root.right, inserted)
	} else {
		return root
	}

	root.updateHeightAndCount()
	return t.rebalance(root)
}

func (t *Tree) deleteRec(root, deleted *Node) *Node {
	if root == nil {
		return nil
	}

	rootComparison := compareNode(root, deleted)
	if rootComparison > 0 {
		root.left = t.deleteRec(root.left, deleted)
	} else if rootComparison < 0 {
		root.right = t.deleteRec(root.right, deleted)
	} else {
		if root.left == nil && root.right == nil {
			return nil
		}

		if root.left == nil {
			return root.right
		}

		if root.right == nil {
			return root.left
		}

		predecessor := root.inorderPredecessor()
		root.key, predecessor.key = predecessor.key, root.key
		root.score, predecessor.score = predecessor.score, root.score

		return t.deleteRec(root.left, deleted)
	}

	root.updateHeightAndCount()
	return t.rebalance(root)
}

func (t *Tree) rebalance(root *Node) *Node {
	if root == nil {
		return nil
	}

	balanceFactor := root.Balance()
	if balanceFactor < -1 {
		if root.right.Balance() > 0 {
			root.right = t.rightRotate(root.right)
		}

		return t.leftRotate(root)
	}

	if balanceFactor > 1 {
		if root.left.Balance() < 0 {
			root.left = t.leftRotate(root.left)
		}

		return t.rightRotate(root)
	}

	return root
}

func (t *Tree) leftRotate(root *Node) *Node {
	newRoot := root.right
	root.right = newRoot.left
	newRoot.left = root

	root.updateHeightAndCount()
	newRoot.updateHeightAndCount()

	return newRoot
}

func (ln *Tree) rightRotate(root *Node) *Node {
	newRoot := root.left
	root.left = newRoot.right
	newRoot.right = root

	root.updateHeightAndCount()
	newRoot.updateHeightAndCount()

	return newRoot
}
