package zdb

import (
	"strconv"
	"strings"
	"testing"
)

func BenchmarkTreeSelect(b *testing.B) {
	// Taking bottom 10 from 10 millions nodes

	// setup
	const N = 10e6
	tree := NewTree()
	for i := 0; i < N; i++ {
		sb := strings.Builder{}
		sb.WriteString("Tree-")
		sb.WriteString(strconv.Itoa(i))
		tree.Add(sb.String(), float64(i))
	}

	for b.Loop() {
		tree.Select(6661488)
	}
}

func BenchmarkTreeSelectReverse(b *testing.B) {
	// Taking bottom 10 from 10 millions nodes

	// setup
	const N = 10e6
	tree := NewTree()
	for i := 0; i < N; i++ {
		sb := strings.Builder{}
		sb.WriteString("Tree-")
		sb.WriteString(strconv.Itoa(i))
		tree.Add(sb.String(), float64(i))
	}

	for b.Loop() {
		tree.SelectReverse(6661488)
	}
}

func BenchmarkTreeRank(b *testing.B) {
	// Taking bottom 10 from 10 millions nodes

	// setup
	const N = 10e6
	tree := NewTree()
	for i := 0; i < N; i++ {
		sb := strings.Builder{}
		sb.WriteString("Tree-")
		sb.WriteString(strconv.Itoa(i))
		tree.Add(sb.String(), float64(i))
	}

	for b.Loop() {
		tree.Rank("Tree-6699988")
	}
}

func BenchmarkTreeRankReverse(b *testing.B) {
	// Taking bottom 10 from 10 millions nodes

	// setup
	const N = 10e6
	tree := NewTree()
	for i := 0; i < N; i++ {
		sb := strings.Builder{}
		sb.WriteString("Tree-")
		sb.WriteString(strconv.Itoa(i))
		tree.Add(sb.String(), float64(i))
	}

	for b.Loop() {
		tree.RankReverse("Tree-6699988")
	}
}

func BenchmarkTreeRangeByIndex(b *testing.B) {
	// Taking bottom 10 from 10 millions nodes

	// setup
	const N = 10e6
	tree := NewTree()
	for i := 0; i < N; i++ {
		sb := strings.Builder{}
		sb.WriteString("Tree-")
		sb.WriteString(strconv.Itoa(i))
		tree.Add(sb.String(), float64(i))
	}

	for b.Loop() {
		tree.RangeByIndex(0, 9)
	}
}

func BenchmarkTreeRangeByIndexReverse(b *testing.B) {
	// Taking top 10 from 10 millions nodes

	// setup
	const N = 10e6
	tree := NewTree()
	for i := 0; i < N; i++ {
		sb := strings.Builder{}
		sb.WriteString("Tree-")
		sb.WriteString(strconv.Itoa(i))
		tree.Add(sb.String(), float64(i))
	}

	for b.Loop() {
		tree.RangeByIndexReverse(0, 9)
	}
}
