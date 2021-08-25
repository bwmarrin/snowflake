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

	maxOverflowMs int64
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

	if c.MaxOverflowMs < 0 {
		return nil, errors.New("invalid config; max overflow cannot be less than 0")
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
		epoch:         curTime.Add(time.Unix(c.Epoch/1000, (c.Epoch%1000)*1000000).Sub(curTime)),
		time:          -1,
		node:          node,
		nodeMask:      nodeMax << c.StepBits,
		stepMask:      -1 ^ (-1 << c.StepBits),
		timeShift:     c.NodeBits + c.StepBits,
		nodeShift:     c.StepBits,
		maxOverflowMs: c.MaxOverflowMs,
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

	// if MaxOverflow is specified, then any side effect of that should be accounted for when generating ID.
	// Overflow are kept irrelevant when generating individual IDs in order to maintain the current API. The
	// alternative would be to generate an ID if the overflow is within the configured limits. That wouldn't
	// be complete without changing the Generate method to return Overflow as well,
	if n.maxOverflowMs > 0 {
		// Wait for any overflow there is. Note that the overflow mechanism is there to allow faster, but
		// intermittent, bulk generation of IDs. The two use-cases should not be served from the same Node.
		<-n.getOverflow(now).AfterCleared()
	}

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

func (n *Node) getOverflow(now int64) Overflow {
	if n.time <= now {
		return Overflow{}
	}

	return Overflow{
		Time: n.epoch.Add(time.Duration(now) * time.Millisecond).UTC(),
		Ms:   n.time - now,
		Step: n.step,
	}
}

// GenerateN generates up to N number of IDs
// If the time element of the ID goes beyond the overflow limit, then fewer
// IDs are generated.
func (n *Node) GenerateN(N int64) (ids Block, of Overflow) {
	if N < 1 {
		return
	}

	n.mu.Lock()

	now := time.Since(n.epoch).Nanoseconds() / 1000000

	maxN := n.getMaxBlockSize(now)
	if maxN > 0 {
		if N > maxN {
			N = maxN
		}

		ids = Block{
			First:    n.reserveGetFirstID(now, N),
			N:        N,
			NodeMask: n.nodeMask,
			StepMask: n.stepMask,
		}
	}

	of = n.getOverflow(now)

	n.mu.Unlock()

	return
}

func (n *Node) reserveGetFirstID(now int64, N int64) ID {
	// we are here because there room for at least one ID
	if now > n.time {
		n.time = now
		n.step = 0
	} else {
		n.step = (n.step + 1) & n.stepMask
		if n.step == 0 {
			n.time++
		}
	}

	first := n.MakeID(n.time, n.step)

	n.time += (n.step + N - 1) / (n.stepMask + 1)
	n.step = (n.step + N - 1) & n.stepMask

	return first
}

func (n *Node) getMaxBlockSize(now int64) int64 {
	if n.time < now {
		return (n.maxOverflowMs + 1) * (n.stepMask + 1)
	}

	if n.time == now {
		return n.maxOverflowMs*(n.stepMask+1) + n.stepMask - n.step
	}

	ofMs := n.time - now
	if ofMs > n.maxOverflowMs {
		return 0
	}

	return (n.maxOverflowMs-ofMs)*(n.stepMask+1) + n.stepMask - n.step
}
