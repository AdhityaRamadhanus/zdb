//go:build unit

package zdb

import (
	"testing"
)

func TestTreeIterator(t *testing.T) {
	// setup

	tests := []struct {
		Name string
		Tree []Node
		Want []Node
		Seek *Node
		Skip bool
	}{
		{
			Name: "Should return all nodes",
			Tree: []Node{
				*NewNode("A", 50),
				*NewNode("B", 15),
				*NewNode("C", 70),
				*NewNode("D", 10),
				*NewNode("E", 11),
			},
			Want: []Node{
				*NewNode("D", 10),
				*NewNode("E", 11),
				*NewNode("B", 15),
				*NewNode("A", 50),
				*NewNode("C", 70),
			},
			Seek: nil,
		},
		{
			Name: "Should return all nodes",
			Tree: []Node{
				*NewNode("D", 20),
				*NewNode("E", 22),
				*NewNode("B", 30),
				*NewNode("A", 50),
				*NewNode("C", 140),
				*NewNode("F", 12),
			},
			Want: []Node{
				*NewNode("F", 12),
				*NewNode("D", 20),
				*NewNode("E", 22),
				*NewNode("B", 30),
				*NewNode("A", 50),
				*NewNode("C", 140),
			},
			Seek: nil,
		},
		{
			Name: "Should return last nodes",
			Tree: []Node{
				*NewNode("A", 50),
				*NewNode("B", 15),
				*NewNode("C", 70),
				*NewNode("D", 10),
				*NewNode("E", 11),
			},
			Want: []Node{
				*NewNode("B", 15),
				*NewNode("A", 50),
				*NewNode("C", 70),
			},
			Seek: NewNode("B", 15),
		},
		{
			Name: "edge case",
			Tree: []Node{
				*NewNode("A", 50),
				*NewNode("B", 15),
				*NewNode("C", 10),
				*NewNode("D", 70),
				*NewNode("E", 60),
				*NewNode("F", 80),
			},
			Want: []Node{
				*NewNode("B", 15),
				*NewNode("A", 50),
				*NewNode("E", 60),
				*NewNode("D", 70),
				*NewNode("F", 80),
			},
			Seek: NewNode("B", 15),
		},
	}

	for _, test := range tests {
		if test.Skip {
			t.Logf("Skipping test %s\n", test.Name)
			continue
		}

		t.Run(test.Name, func(t *testing.T) {
			tree := NewTree()
			for _, node := range test.Tree {
				tree.Add(node.key, node.score)
			}

			it := NewTreeIterator(tree)
			it.Seek(test.Seek)

			got := []Node{}
			for {
				next := it.Next()
				if next == nil {
					break
				}

				got = append(got, *next)
			}

			checkSliceOfNodes(t, got, test.Want)
		})
	}
}

func checkSliceOfNodes(t *testing.T, got, want []Node) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("got different length of slices got %+v, want %+v\n", got, want)
	}

	for i := 0; i < len(got); i++ {
		checkNilOrWantedNode(t, &got[i], &want[i])
	}
}
