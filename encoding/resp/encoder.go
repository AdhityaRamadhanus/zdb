package resp

import (
	"fmt"

	"github.com/AdhityaRamadhanus/zdb"
	"github.com/AdhityaRamadhanus/zdb/miniresp3"
)

func SerializeTree(writer *miniresp3.Writer, tree zdb.OrderStatisticTree, withScores bool) {
	if tree == nil || tree.Root() == nil {
		writer.AppendArrHeader(0)
		return
	}

	it := zdb.NewTreeIterator(tree)
	it.Seek(nil)

	arrLength := tree.Root().Count()
	if withScores {
		arrLength *= 2
	}

	writer.AppendArrHeader(arrLength)
	for {
		next := it.Next()
		if next == nil {
			break
		}

		writer.AppendBulkStr(next.Key())
		if withScores {
			writer.AppendBulkStr(fmt.Sprintf("%.2f", next.Score()))
		}
	}
}

func SerializeNodes(writer *miniresp3.Writer, nodes []zdb.Node, withScores bool) {
	if len(nodes) == 0 {
		writer.AppendArrHeader(0)
		return
	}

	arrLength := len(nodes)
	if withScores {
		arrLength *= 2
	}

	writer.AppendArrHeader(arrLength)
	for _, node := range nodes {
		writer.AppendBulkStr(node.Key())
		if withScores {
			writer.AppendBulkStr(fmt.Sprintf("%.2f", node.Score()))
		}
	}
}
