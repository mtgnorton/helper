package consistenthash

import (
	"testing"
)

func Test_CHash(t *testing.T) {
	cHash := NewCHash(3)
	cHash.AddNode("a", "b", "c")
	t.Log(cHash.Get("1"))
}
