package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AdhityaRamadhanus/zdb"
)

func main() {
	const N = 1e6
	tree := zdb.NewTree()

	// randSrc := rand.NewSource(rand.Int63())
	scores := []float64{}
	keys := []string{}
	for i := 0; i < N; i++ {
		sb := strings.Builder{}
		sb.WriteString("Tree-")
		sb.WriteString(strconv.Itoa(i))
		keys = append(keys, sb.String())
		scores = append(scores, float64(i))
		// scores = append(scores, randSrc.Int63())
	}

	start := time.Now()
	for i := 0; i < N; i++ {
		tree.Add(keys[i], scores[i])
	}

	end := time.Since(start)
	fmt.Printf("Inserting %.2f nodes into AVL tree %.2f seconds \n", N, end.Seconds())
}
