package snowflake

import (
	"testing"
	"time"
)

func TestBlockIterator_Next(t *testing.T) {
	epoch := time.Now().UnixNano() / 1000000
	node, err := NewNodeWithConfig(0, Config{
		Epoch:         epoch,
		NodeBits:      10,
		StepBits:      12,
		MaxOverflowMs: 25,
	})

	if err != nil {
		t.Errorf("failed to create node; error: %v", err)
		return
	}

	b, o := node.GenerateN(4000000)

	t.Logf("b:%v, o:%v", b, o)

	bi := NewBlockIterator(b)
	i := int64(0)
	dup := make(map[ID]bool, b.N)
	for ; ; i++ {
		id, ok := bi.Next()
		if !ok {
			break
		}

		if _, found := dup[id]; found {
			t.Fatalf("duplicate ID found; %d", id)
		}

		dup[id] = true
	}

	if i != b.N {
		t.Errorf("block did not iterate 100 times")
	}
}
