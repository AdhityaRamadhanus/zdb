package zdb

import (
	"github.com/AdhityaRamadhanus/zdb/commands"
)

type ZDB struct {
	// TODO: Abstract out shards
	shards Shard
}

func NewZDB(shards uint) *ZDB {
	return &ZDB{
		shards: *NewShards(shards),
	}
}

func (zdb *ZDB) ShardStats() []int {
	lengths := []int{}
	//TODO: encapsulate this better
	for _, shard := range zdb.shards.DB {
		lengths = append(lengths, len(shard))
	}

	return lengths
}

func (zdb *ZDB) Scan(cmd *commands.ScanCmd) (keys []string, nextCursor string, err error) {
	it := NewTreeIterator(zdb.shards.Keys)
	var cursor *Node
	// TODO: Create a type for 0 cursor
	if cmd.Cursor != "0" {
		score, err := zdb.shards.Keys.GetScore(cmd.Cursor)
		if err != nil {
			return keys, nextCursor, err
		}

		cursor = NewNode(cmd.Cursor, score)
	}
	it.Seek(cursor)

	// TODO: Match glob pattern
	count := 0
	for {
		if count == cmd.Count {
			break
		}

		next := it.Next()
		if next == nil {
			break
		}
		keys = append(keys, next.key)
		count += 1
	}

	if len(keys) > 1 {
		nextCursor = keys[len(keys)-1]
	} else {
		nextCursor = "0"
	}

	return keys, nextCursor, nil
}

func (zdb *ZDB) ZAdd(cmd *commands.ZADDCmd) int {
	tree := zdb.shards.GetDBFromKey(cmd.Key)
	if tree == nil {
		tree = zdb.shards.AddDB(cmd.Key)
	}

	success := 0
	for _, z := range cmd.Members {
		tree.Add(z.Key, z.Score)
		success += 1
	}

	return success
}

func (zdb *ZDB) ZCard(cmd *commands.ZCardCmd) int {
	tree := zdb.shards.GetDBFromKey(cmd.Key)
	if tree == nil {
		return 0
	}

	return tree.Root().Count()
}

func (zdb *ZDB) ZCount(cmd *commands.ZCountCmd) int {
	tree := zdb.shards.GetDBFromKey(cmd.Key)
	if tree == nil {
		return 0
	}

	return tree.CountByScore(cmd.Min, cmd.Max)
}

func (zdb *ZDB) ZRank(cmd *commands.ZRankCmd) int {
	tree := zdb.shards.GetDBFromKey(cmd.Key)
	if tree == nil {
		return -1
	}

	return tree.Rank(cmd.Member)
}

func (zdb *ZDB) ZRange(cmd *commands.ZRangeCmd) []Node {
	// TODO: Pagination
	tree := zdb.shards.GetDBFromKey(cmd.Key)
	if tree == nil {
		return []Node{}
	}

	if cmd.ByIndex {
		if cmd.StopIndex == -1 {
			cmd.StopIndex = tree.Root().Count()
		}

		if cmd.Reverse {
			return tree.RangeByIndexReverse(cmd.StartIndex, cmd.StopIndex)
		}

		return tree.RangeByIndex(cmd.StartIndex, cmd.StopIndex)
	}

	if cmd.ByScore {
		if cmd.Reverse {
			return tree.RangeByScoreReverse(cmd.MinScore, cmd.MaxScore)
		}

		return tree.RangeByScore(cmd.MinScore, cmd.MaxScore)
	}

	if cmd.ByLex {
		if cmd.Reverse {
			return tree.RangeByLexReverse(cmd.MinKey, cmd.MaxKey)
		}

		return tree.RangeByLex(cmd.MinKey, cmd.MaxKey)
	}

	return []Node{}
}

func (zdb *ZDB) ZRem(cmd *commands.ZRemCmd) int {
	tree := zdb.shards.GetDBFromKey(cmd.Key)
	if tree == nil {
		return 0
	}

	success := 0
	for _, key := range cmd.Members {
		tree.Remove(key)
		success += 1
	}

	if tree.Root().Count() == 0 {
		zdb.shards.RemoveDB(cmd.Key)
	}

	return success
}

func (zdb *ZDB) ZScan(cmd *commands.ZScanCmd) (keys []string, nextCursor string) {
	tree := zdb.shards.GetDBFromKey(cmd.Key)
	if tree == nil {
		return keys, "0"
	}

	it := NewTreeIterator(tree)
	var cursor *Node
	if cmd.ScanCmd.Cursor != "0" {
		score, err := tree.GetScore(cmd.ScanCmd.Cursor)
		if err != nil {
			return
		}

		cursor = NewNode(cmd.ScanCmd.Cursor, score)
	}
	it.Seek(cursor)

	// TODO: Match glob pattern
	count := 0
	for {
		if count == cmd.ScanCmd.Count {
			break
		}

		next := it.Next()
		if next == nil {
			break
		}
		keys = append(keys, next.key)
		count += 1
	}

	if len(keys) > 1 {
		nextCursor = keys[len(keys)-1]
	} else {
		nextCursor = "0"
	}

	return keys, nextCursor
}

func (zdb *ZDB) ZScore(cmd *commands.ZScoreCmd) (float64, error) {
	tree := zdb.shards.GetDBFromKey(cmd.Key)
	if tree == nil {
		return 0, errNotFound
	}

	return tree.GetScore(cmd.Member)
}
