package zdb

import "errors"

var (
	errNotFound = errors.New("err not found")
)

type AggFunc func(score1, score2 float64) float64

var (
	SumAggFunc = func(wt1, wt2 float64) AggFunc {
		return func(score1, score2 float64) float64 {
			return score1*wt1 + score2*wt2
		}
	}

	MinAggFunc = func(wt1, wt2 float64) AggFunc {
		return func(score1, score2 float64) float64 {
			return min(score1*wt1, score2*wt2)
		}
	}

	MaxAggFunc = func(wt1, wt2 float64) AggFunc {
		return func(score1, score2 float64) float64 {
			return max(score1*wt1, score2*wt2)
		}
	}
)

type OrderStatisticTree interface {
	// CRUD
	Root() *Node
	GetScore(key string) (float64, error)
	Add(key string, score float64)
	Remove(key string)

	// ordering
	Select(idx int) *Node
	SelectReverse(idx int) *Node
	Rank(key string) int
	RankReverse(key string) int

	// filter operation
	RangeByIndex(start, stop int) []Node
	RangeByIndexReverse(start, stop int) []Node
	RangeByScore(min, max float64) []Node
	RangeByScoreReverse(max, min float64) []Node
	RangeByLex(min, max string) []Node
	RangeByLexReverse(min, max string) []Node
	CountByScore(min, max float64) int

	// set operation
	Diff(other OrderStatisticTree) OrderStatisticTree
	Inter(other OrderStatisticTree, aggFunc AggFunc) OrderStatisticTree
	Union(other OrderStatisticTree, aggFunc AggFunc) OrderStatisticTree
}
