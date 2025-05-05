package zdb

import (
	"slices"
	"testing"
)

func TestAVLAdd(t *testing.T) {
	tests := []struct {
		Name  string
		Nodes []Node
		Want  []string
		Skip  bool
	}{
		{
			Name: "Right Rotate",
			Nodes: []Node{
				*NewNode("A", 105),
				*NewNode("B", 72),
				*NewNode("C", 34),
				*NewNode("D", 17),
				*NewNode("E", 12),
			},
			Want: []string{"B", "D", "E", "C", "A"},
		},
		{
			Name: "Left Rotate",
			Nodes: []Node{
				*NewNode("A", 1),
				*NewNode("B", 2),
				*NewNode("C", 3),
				*NewNode("D", 4),
				*NewNode("E", 5),
			},
			Want: []string{"B", "A", "D", "C", "E"},
		},
		{
			Name: "Left Right Rotate",
			Nodes: []Node{
				*NewNode("A", 15),
				*NewNode("B", 5),
				*NewNode("C", 1),
				*NewNode("D", 7),
				*NewNode("E", 9),
			},
			Want: []string{"B", "C", "E", "D", "A"},
		},
		{
			Name: "Right Left Rotate",
			Nodes: []Node{
				*NewNode("A", 15),
				*NewNode("B", 5),
				*NewNode("C", 1),
				*NewNode("D", 20),
				*NewNode("E", 17),
			},
			Want: []string{"B", "C", "E", "A", "D"},
		},
		{
			Name: "Same score different key",
			Nodes: []Node{
				*NewNode("A", 15),
				*NewNode("B", 15),
				*NewNode("C", 15),
				*NewNode("D", 1),
				*NewNode("E", 0),
			},
			Want: []string{"B", "D", "E", "A", "C"},
		},
	}

	for _, test := range tests {
		if test.Skip {
			continue
		}
		t.Run(test.Name, func(t *testing.T) {
			tree := NewTree()
			for _, node := range test.Nodes {
				tree.Add(node.key, node.score)
			}

			checkPreOrderKeyTree(t, tree, test.Want)
		})
	}
}

func TestAVLRemove(t *testing.T) {
	tests := []struct {
		Name         string
		Nodes        []Node
		DeletedNodes []Node
		Want         []string
		Skip         bool
	}{
		{
			Name: "Right Rotate",
			Nodes: []Node{
				*NewNode("A", 105),
				*NewNode("B", 72),
				*NewNode("C", 34),
				*NewNode("D", 17),
				*NewNode("E", 12),
			},
			DeletedNodes: []Node{*NewNode("A", 105)},
			Want:         []string{"D", "E", "B", "C"},
		},
		{
			Name: "Left Rotate",
			Nodes: []Node{
				*NewNode("A", 105),
				*NewNode("B", 72),
				*NewNode("C", 106),
				*NewNode("D", 17),
			},
			DeletedNodes: []Node{*NewNode("D", 17)},
			Want:         []string{"A", "B", "C"},
		},
		{
			Name: "Left Right Rotate",
			Nodes: []Node{
				*NewNode("A", 105),
				*NewNode("B", 72),
				*NewNode("C", 34),
				*NewNode("D", 17),
				*NewNode("E", 12),
			},
			DeletedNodes: []Node{*NewNode("E", 12), *NewNode("A", 105)},
			Want:         []string{"C", "D", "B"},
		},
		{
			Name: "Right Left Rotate",
			Nodes: []Node{
				*NewNode("A", 105),
				*NewNode("B", 72),
				*NewNode("C", 107),
				*NewNode("D", 17),
				*NewNode("E", 89),
			},
			DeletedNodes: []Node{*NewNode("C", 107), *NewNode("D", 17)},
			Want:         []string{"E", "B", "A"},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Log("Skipping test ", test.Name)
			continue
		}
		t.Run(test.Name, func(t *testing.T) {
			tree := NewTree()
			for _, node := range test.Nodes {
				tree.Add(node.key, node.score)
			}

			for _, deletedNode := range test.DeletedNodes {
				tree.Remove(deletedNode.key)
			}

			checkPreOrderKeyTree(t, tree, test.Want)
		})
	}
}

func TestAVLRank(t *testing.T) {
	// setup

	avl := NewTree()
	avl.Add("A", 50)
	avl.Add("A1", 50)
	avl.Add("B", 15)
	avl.Add("C", 70)
	avl.Add("D", 10)
	avl.Add("E", 11)

	tests := []struct {
		Name      string
		SearchKey string
		Want      int
		Skip      bool
	}{
		{
			Name:      "Rank 50",
			SearchKey: "A",
			Want:      4,
		},
		{
			Name:      "Rank 50, A1",
			SearchKey: "A1",
			Want:      5,
		},
		{
			Name:      "Rank B",
			SearchKey: "B",
			Want:      3,
		},
		{
			Name:      "Rank C",
			SearchKey: "C",
			Want:      6,
		},
		{
			Name:      "Rank E",
			SearchKey: "E",
			Want:      2,
		},
		{
			Name:      "Rank D",
			SearchKey: "D",
			Want:      1,
		},
		{
			Name:      "Search for rank of non-existing element",
			SearchKey: "NOT",
			Want:      -1,
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := avl.Rank(test.SearchKey)

			if got != test.Want {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLRankReverse(t *testing.T) {
	// setup

	avl := NewTree()
	avl.Add("A", 50)
	avl.Add("A1", 50)
	avl.Add("B", 15)
	avl.Add("C", 70)
	avl.Add("D", 10)
	avl.Add("E", 11)

	tests := []struct {
		Name      string
		SearchKey string
		Want      int
		Skip      bool
	}{
		{
			Name:      "Rank 50",
			SearchKey: "A",
			Want:      3,
		},
		{
			Name:      "Rank 50, A1",
			SearchKey: "A1",
			Want:      2,
		},
		{
			Name:      "Rank B",
			SearchKey: "B",
			Want:      4,
		},
		{
			Name:      "Rank C",
			SearchKey: "C",
			Want:      1,
		},
		{
			Name:      "Rank E",
			SearchKey: "E",
			Want:      5,
		},
		{
			Name:      "Rank D",
			SearchKey: "D",
			Want:      6,
		},
		{
			Name:      "Search for rank of non-existing element",
			SearchKey: "NOT",
			Want:      -1,
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := avl.RankReverse(test.SearchKey)

			if got != test.Want {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLSelect(t *testing.T) {
	// setup

	avl := NewTree()
	avl.Add("A", 50)
	avl.Add("B", 15)
	avl.Add("C", 70)
	avl.Add("D", 10)
	avl.Add("E", 11)

	tests := []struct {
		Name      string
		SelectIdx int
		Want      *Node
		Skip      bool
	}{
		{
			Name:      "Search for 2nd ranked",
			SelectIdx: 2,
			Want:      NewNode("E", 11),
		},
		{
			Name:      "Search for last ranked",
			SelectIdx: 4,
			Want:      NewNode("A", 50),
		},
		{
			Name:      "Search for out of bound index",
			SelectIdx: 10,
			Want:      nil,
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := avl.Select(test.SelectIdx)
			checkNilOrWantedNode(t, got, test.Want)
		})
	}
}

func TestAVLSelectReverse(t *testing.T) {
	// setup

	avl := NewTree()
	avl.Add("A", 50)
	avl.Add("B", 15)
	avl.Add("C", 70)
	avl.Add("D", 10)
	avl.Add("E", 11)

	tests := []struct {
		Name      string
		SelectIdx int
		Want      *Node
		Skip      bool
	}{
		{
			Name:      "Search for 2nd ranked",
			SelectIdx: 2,
			Want:      NewNode("A", 50),
		},
		{
			Name:      "Search for last ranked",
			SelectIdx: 5,
			Want:      NewNode("D", 10),
		},
		{
			Name:      "Search for out of bound index",
			SelectIdx: 10,
			Want:      nil,
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := avl.SelectReverse(test.SelectIdx)
			checkNilOrWantedNode(t, got, test.Want)
		})
	}
}

func TestAVLRangeByIndex(t *testing.T) {
	// setup

	avl := NewTree()
	avl.Add("A", 50)
	avl.Add("B", 15)
	avl.Add("C", 70)
	avl.Add("D", 10)
	avl.Add("E", 11)

	tests := []struct {
		Name  string
		Start int
		End   int
		Want  []string
		Skip  bool
	}{
		{
			Name:  "Get 0-2 ranked elements",
			Start: 0,
			End:   2,
			Want:  []string{"D", "E", "B"},
		},
		{
			Name:  "Get 2 ranked elements",
			Start: 1,
			End:   1,
			Want:  []string{"E"},
		},
		{
			Name:  "Get 0-100 ranked elements, out of bound end",
			Start: 0,
			End:   100,
			Want:  []string{"D", "E", "B", "A", "C"},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := Map(
				avl.RangeByIndex(test.Start, test.End),
				func(n Node) string { return n.key },
			)

			if slices.Compare(got, test.Want) != 0 {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLRangeByIndexReverse(t *testing.T) {
	// setup
	avl := NewTree()
	avl.Add("A", 50)
	avl.Add("B", 15)
	avl.Add("C", 70)
	avl.Add("D", 10)
	avl.Add("E", 11)

	tests := []struct {
		Name  string
		Start int
		End   int
		Want  []string
		Skip  bool
	}{
		{
			Name:  "Get 0-2 ranked elements",
			Start: 0,
			End:   2,
			Want:  []string{"C", "A", "B"},
		},
		{
			Name:  "Get 2 ranked elements",
			Start: 1,
			End:   1,
			Want:  []string{"A"},
		},
		{
			Name:  "Get 0-100 ranked elements, out of bound end",
			Start: 0,
			End:   100,
			Want:  []string{"C", "A", "B", "E", "D"},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := Map(
				avl.RangeByIndexReverse(test.Start, test.End),
				func(n Node) string { return n.key },
			)

			if slices.Compare(got, test.Want) != 0 {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLRangeByScore(t *testing.T) {
	// setup

	avl := NewTree()
	avl.Add("A", 50)
	avl.Add("B", 15)
	avl.Add("C", 70)
	avl.Add("D", 10)
	avl.Add("E", 11)

	tests := []struct {
		Name  string
		Start float64
		End   float64
		Want  []string
		Skip  bool
	}{
		{
			Name:  "Get elements with score range within limit",
			Start: 10,
			End:   70,
			Want:  []string{"D", "E", "B", "A", "C"},
		},
		{
			Name:  "Get elements with score range high out of limit",
			Start: 50,
			End:   100,
			Want:  []string{"A", "C"},
		},
		{
			Name:  "Get elements with score range low and high out of limit",
			Start: 100,
			End:   100,
			Want:  []string{},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := Map(
				avl.RangeByScore(test.Start, test.End),
				func(n Node) string { return n.key },
			)

			if slices.Compare(got, test.Want) != 0 {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLRangeByScoreReverse(t *testing.T) {
	// setup

	avl := NewTree()
	avl.Add("A", 50)
	avl.Add("B", 15)
	avl.Add("C", 70)
	avl.Add("D", 10)
	avl.Add("E", 11)

	tests := []struct {
		Name  string
		Start float64
		End   float64
		Want  []string
		Skip  bool
	}{
		{
			Name:  "Get elements with score range within limit",
			Start: 10,
			End:   70,
			Want:  []string{"C", "A", "B", "E", "D"},
		},
		{
			Name:  "Get elements with score range high out of limit",
			Start: 50,
			End:   100,
			Want:  []string{"C", "A"},
		},
		{
			Name:  "Get elements with score range low and high out of limit",
			Start: 100,
			End:   100,
			Want:  []string{},
		},
		{
			Name:  "Get elements with score range low and high out of limit",
			Start: 70,
			End:   90,
			Want:  []string{"C"},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := Map(
				avl.RangeByScoreReverse(test.Start, test.End),
				func(n Node) string { return n.key },
			)

			if slices.Compare(got, test.Want) != 0 {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLRangeByLex(t *testing.T) {
	// setup

	avl := NewTree()
	avl.Add("A", 0)
	avl.Add("B", 0)
	avl.Add("C", 0)
	avl.Add("D", 0)
	avl.Add("E", 0)

	tests := []struct {
		Name  string
		Start string
		End   string
		Want  []string
		Skip  bool
	}{
		{
			Name:  "Get elements with score range within limit",
			Start: "A",
			End:   "Z",
			Want:  []string{"A", "B", "C", "D", "E"},
		},
		{
			Name:  "Get elements with score range within limit",
			Start: "C",
			End:   "D",
			Want:  []string{"C", "D"},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := Map(
				avl.RangeByLex(test.Start, test.End),
				func(n Node) string { return n.key },
			)

			if slices.Compare(got, test.Want) != 0 {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLRangeByLexReverse(t *testing.T) {
	// setup

	avl := NewTree()
	avl.Add("A", 0)
	avl.Add("B", 0)
	avl.Add("C", 0)
	avl.Add("D", 0)
	avl.Add("E", 0)

	tests := []struct {
		Name  string
		Start string
		End   string
		Want  []string
		Skip  bool
	}{
		{
			Name:  "Get elements with score range within limit",
			Start: "A",
			End:   "Z",
			Want:  []string{"E", "D", "C", "B", "A"},
		},
		{
			Name:  "Get elements with score range within limit",
			Start: "C",
			End:   "D",
			Want:  []string{"D", "C"},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			got := Map(
				avl.RangeByLexReverse(test.Start, test.End),
				func(n Node) string { return n.key },
			)

			if slices.Compare(got, test.Want) != 0 {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLDiff(t *testing.T) {
	tests := []struct {
		Name  string
		Tree1 []Node
		Tree2 []Node
		Want  []Node
		Skip  bool
	}{
		{
			Name: "Diff key diff score",
			Tree1: []Node{
				*NewNode("A", 50),
				*NewNode("B", 15),
				*NewNode("C", 70),
				*NewNode("D", 10),
				*NewNode("E", 11),
			},
			Tree2: []Node{
				*NewNode("B", 15),
				*NewNode("C", 70),
				*NewNode("D", 10),
				*NewNode("E", 11),
			},
			Want: []Node{
				*NewNode("A", 50),
			},
		},
		{
			Name: "Same key different score",
			Tree1: []Node{
				*NewNode("A", 50),
				*NewNode("B", 15),
				*NewNode("C", 70),
				*NewNode("D", 10),
				*NewNode("E", 11),
			},
			Tree2: []Node{
				*NewNode("B", 10),
				*NewNode("E", 11),
			},
			Want: []Node{
				*NewNode("D", 10),
				*NewNode("A", 50),
				*NewNode("C", 70),
			},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			t1 := NewTree()
			for _, node := range test.Tree1 {
				t1.Add(node.key, node.score)
			}

			t2 := NewTree()
			for _, node := range test.Tree2 {
				t2.Add(node.key, node.score)
			}

			diff := t1.Diff(t2)
			it := NewTreeIterator(diff)
			it.Seek(nil)

			diffNodes := []*Node{}
			for {
				next := it.Next()
				if next == nil {
					break
				}
				diffNodes = append(diffNodes, next)
			}

			for i, node := range diffNodes {
				checkNilOrWantedNode(t, node, &test.Want[i])
			}
		})
	}
}

func TestAVLInter(t *testing.T) {
	// setup

	sumAggFunc := func(score1, score2 float64) float64 {
		return score1 + score2
	}

	tests := []struct {
		Name    string
		Tree1   []Node
		Tree2   []Node
		Want    []Node
		AggFunc func(score1, score2 float64) float64
		Skip    bool
	}{
		{
			Name: "Diff key diff score",
			Tree1: []Node{
				*NewNode("A", 50),
				*NewNode("B", 15),
				*NewNode("C", 70),
				*NewNode("D", 10),
				*NewNode("E", 11),
			},
			Tree2: []Node{
				*NewNode("B", 15),
				*NewNode("C", 70),
				*NewNode("D", 10),
				*NewNode("E", 11),
			},
			AggFunc: sumAggFunc,
			Want: []Node{
				*NewNode("D", 20),
				*NewNode("E", 22),
				*NewNode("B", 30),
				*NewNode("C", 140),
			},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			t1 := NewTree()
			for _, node := range test.Tree1 {
				t1.Add(node.key, node.score)
			}

			t2 := NewTree()
			for _, node := range test.Tree2 {
				t2.Add(node.key, node.score)
			}

			diff := t1.Inter(t2, test.AggFunc)
			it := NewTreeIterator(diff)
			it.Seek(nil)

			diffNodes := []*Node{}
			for {
				next := it.Next()
				if next == nil {
					break
				}
				diffNodes = append(diffNodes, next)
			}

			for i, node := range diffNodes {
				checkNilOrWantedNode(t, node, &test.Want[i])
			}
		})
	}
}

func TestAVLUnion(t *testing.T) {
	// setup

	sumAggFunc := func(score1, score2 float64) float64 {
		return score1 + score2
	}

	tests := []struct {
		Name    string
		Tree1   []Node
		Tree2   []Node
		Want    []Node
		AggFunc func(score1, score2 float64) float64
		Skip    bool
	}{
		{
			Name: "Union 2 Trees",
			Tree1: []Node{
				*NewNode("A", 50),
				*NewNode("B", 15),
				*NewNode("C", 70),
				*NewNode("D", 10),
				*NewNode("E", 11),
			},
			Tree2: []Node{
				*NewNode("B", 15),
				*NewNode("C", 70),
				*NewNode("D", 10),
				*NewNode("E", 11),
				*NewNode("F", 12),
			},
			AggFunc: sumAggFunc,
			Want: []Node{
				*NewNode("F", 12),
				*NewNode("D", 20),
				*NewNode("E", 22),
				*NewNode("B", 30),
				*NewNode("A", 50),
				*NewNode("C", 140),
			},
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			t1 := NewTree()
			for _, node := range test.Tree1 {
				t1.Add(node.key, node.score)
			}

			t2 := NewTree()
			for _, node := range test.Tree2 {
				t2.Add(node.key, node.score)
			}

			union := t1.Union(t2, test.AggFunc)
			it := NewTreeIterator(union)
			it.Seek(nil)

			idx := 0
			for {
				next := it.Next()
				if next == nil {
					break
				}

				checkNilOrWantedNode(t, next, &test.Want[idx])
				idx += 1
			}
		})
	}
}

func TestAVLrankByScoreLowerBound(t *testing.T) {
	// setup

	tree := NewTree().(*Tree)
	tree.Add("A", 50)
	tree.Add("B", 15)
	tree.Add("C", 70)
	tree.Add("D", 10)
	tree.Add("E", 11)
	tree.Add("E1", 11)

	tests := []struct {
		Name  string
		Score float64
		Want  int
	}{
		{
			Name:  "Lower than first rank",
			Score: 9,
			Want:  1,
		},
		{
			Name:  "Lower than first rank",
			Score: 10,
			Want:  1,
		},
		{
			Name:  "Lower than first rank",
			Score: 45,
			Want:  5,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := tree.rankByScoreLowerBound(test.Score)
			if got != test.Want {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLrankByScoreUpperBound(t *testing.T) {
	// setup

	tree := NewTree().(*Tree)
	tree.Add("A", 50)
	tree.Add("B", 15)
	tree.Add("C", 70)
	tree.Add("D", 10)
	tree.Add("E", 11)
	tree.Add("E1", 11)

	tests := []struct {
		Name  string
		Score float64
		Want  int
	}{
		{
			Name:  "Lower than first rank",
			Score: 9,
			Want:  0,
		},
		{
			Name:  "Lower than first rank",
			Score: 11,
			Want:  4,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := tree.rankByScoreUpperBound(test.Score)
			if got != test.Want {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func TestAVLCountByScore(t *testing.T) {
	// setup

	tree := NewTree()
	tree.Add("A", 50)
	tree.Add("B", 15)
	tree.Add("C", 70)
	tree.Add("D", 10)
	tree.Add("E", 11)
	tree.Add("E1", 11)

	tests := []struct {
		Name string
		Min  float64
		Max  float64
		Want int
	}{
		{
			Name: "Lower than first rank",
			Min:  9,
			Max:  11,
			Want: 3,
		},
		{
			Name: "Lower than first rank",
			Min:  10,
			Max:  11,
			Want: 3,
		},
		{
			Name: "Lower than first rank",
			Min:  10,
			Max:  95,
			Want: 6,
		},
		{
			Name: "Lower than first rank",
			Min:  0,
			Max:  95,
			Want: 6,
		},
		{
			Name: "Lower than first rank",
			Min:  11,
			Max:  95,
			Want: 5,
		},
		{
			Name: "Lower than first rank",
			Min:  11,
			Max:  70,
			Want: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := tree.CountByScore(test.Min, test.Max)
			if got != test.Want {
				t.Errorf("got %v, want %v\n", got, test.Want)
			}
		})
	}
}

func checkNilOrWantedNode(t *testing.T, got, want *Node) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Errorf("got %v, want nil\n", got)
		}
	} else {
		if got.key != want.key || got.score != want.score {
			t.Errorf("got %v, want %v\n", got, want)
		}
	}
}

func checkPreOrderKeyTree(t *testing.T, tree OrderStatisticTree, want []string) {
	t.Helper()

	traversalFunc := func() []string {
		keys := []string{}
		stack := []*Node{}
		curr := tree.Root()
		for curr != nil || len(stack) > 0 {
			if curr == nil {
				curr = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				curr = curr.right
				continue
			}

			keys = append(keys, curr.key)
			stack = append(stack, curr)
			curr = curr.left
		}
		return keys
	}

	got := traversalFunc()
	if slices.Compare(got, want) != 0 {
		t.Errorf("got %v, want %v", got, want)
	}
}
