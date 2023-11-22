package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFun func(data []byte) uint32

type Node string

type CHash struct {
	hashFun     HashFun
	vNodeNumber int
	ring        []int // Sorted
	vAMapping   map[int]Node
}

func NewCHash(vNodeNumber int, fns ...HashFun) *CHash {
	fn := crc32.ChecksumIEEE
	if len(fns) > 0 {
		fn = fns[0]
	}
	m := &CHash{
		vNodeNumber: vNodeNumber,
		hashFun:     fn,
		vAMapping:   make(map[int]Node),
	}
	return m
}

func (c *CHash) AddNode(nodes ...Node) {
	for _, node := range nodes {
		for i := 0; i < c.vNodeNumber; i++ {
			nodeHash := int(c.hashFun([]byte(string(node) + strconv.Itoa(i))))
			c.ring = append(c.ring, nodeHash)
			c.vAMapping[nodeHash] = node
		}
	}
	sort.Ints(c.ring)
}

func (c *CHash) Get(key string) Node {
	if len(c.ring) == 0 {
		return ""
	}
	hashKey := int(c.hashFun([]byte(key)))
	idx := sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= hashKey
	})
	return c.vAMapping[c.ring[idx%len(c.ring)]]
}
