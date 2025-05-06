package zdb

import "time"

type Shard struct {
	Keys OrderStatisticTree
	DB   []map[string]OrderStatisticTree
	mask uint64
	hash fnv64a
}

func NewShards(shards uint) *Shard {
	shards = max(shards, 1)
	shard := &Shard{
		Keys: NewTree(),
		DB:   []map[string]OrderStatisticTree{},
		mask: Mask64(shards-1) - 1,
		hash: fnv64a{},
	}

	for range shards {
		shard.DB = append(shard.DB, make(map[string]OrderStatisticTree))
	}

	return shard
}

func (s *Shard) GetDBFromKey(key string) OrderStatisticTree {
	shardIdx := int(s.hash.Sum64(key) & s.mask)
	db := s.DB[shardIdx]
	return db[key]
}

func (s *Shard) AddDB(key string) OrderStatisticTree {
	shardIdx := int(s.hash.Sum64(key) & s.mask)
	s.DB[shardIdx][key] = NewTree()
	s.Keys.Add(key, float64(time.Now().Unix()))
	return s.DB[shardIdx][key]
}

func (s *Shard) RemoveDB(key string) {
	shardIdx := int(s.hash.Sum64(key) & s.mask)
	delete(s.DB[shardIdx], key)
	s.Keys.Remove(key)
}
