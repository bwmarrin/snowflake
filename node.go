package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// A Node struct holds the basic information needed for a snowflake generator
// node
type Node struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

// NewNodeWithConfig creates a new snowflake node with the given config
func NewNodeWithConfig(node int64, c Config) (*Node, error) {
	if c.NodeBits == 0 {
		return nil, errors.New("invalid config; NodeBits cannot be 0")
	}

	if c.StepBits == 0 {
		return nil, errors.New("invalid config; StepBits cannot be 0")
	}

	if c.NodeBits+c.StepBits != 22 {
		return nil, errors.New("invalid config; NodeBits + StepBits should be 22")
	}

	nodeMax := int64(-1 ^ (-1 << c.NodeBits))
	if node < 0 || node > nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(nodeMax, 10))
	}

	if c.Epoch == 0 {
		c.Epoch = defaultEpoch
	}

	curTime := time.Now()
	return &Node{
		// add time.Duration to curTime to make sure we use the monotonic clock if available
		epoch:     curTime.Add(time.Unix(c.Epoch/1000, (c.Epoch%1000)*1000000).Sub(curTime)),
		time:      -1,
		node:      node,
		nodeMask:  nodeMax << c.StepBits,
		stepMask:  -1 ^ (-1 << c.StepBits),
		timeShift: c.NodeBits + c.StepBits,
		nodeShift: c.StepBits,
	}, nil
}

// NewNode returns a new snowflake node that can be used to generate snowflake
// IDs
func NewNode(node int64) (*Node, error) {
	return NewNodeWithConfig(node, defaultConfig)
}

// Generate creates and returns a unique snowflake ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node ID
func (n *Node) Generate() ID {

	n.mu.Lock()

	now := time.Since(n.epoch).Nanoseconds() / 1000000

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := n.MakeID(n.time, n.step)

	n.mu.Unlock()
	return r
}

// MakeID makes an ID with the specified time and step
func (n *Node) MakeID(time int64, step int64) ID {
	return ID(time<<n.timeShift | n.node<<n.nodeShift | step)
}
