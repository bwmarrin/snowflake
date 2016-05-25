// simple generator for..
// Twitter Snowflake with custom definable epoch

package flake

import "encoding/base64"
import "fmt"
import "strconv"
import "sync"
import "time"

const (
	TimeBits = 41
	NodeBits = 10
	StepBits = 12

	TimeMask int64 = -1 ^ (-1 << TimeBits)
	NodeMask int64 = -1 ^ (-1 << NodeBits)
	StepMask int64 = -1 ^ (-1 << StepBits)

	NodeMax = -1 ^ (-1 << NodeBits)

	TimeShift uint8 = NodeBits + StepBits
	NodeShift uint8 = StepBits
)

type Node struct {
	sync.Mutex // TODO: find a way to avoid locks?

	// configurable values
	epoch int64
	node  int64

	// runtime tracking values
	lastTime int64
	step     int64
}

// Start a new Flake factory / server node using the given node number
// sets with default settings, use helper functions to change
// node, epoch, etc.
func NewFlakeNode(node int64) (*Node, error) {

	if node < 0 || node > NodeMax {
		return nil, fmt.Errorf("Invalid node number.")
	}

	return &Node{
		epoch:    time.Date(2016, 1, 0, 0, 0, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond),
		node:     node,
		lastTime: 0,
		step:     0,
	}, nil
}

// high performance generator
// well, that w as the idea...
func (n *Node) Generator(c chan Flake) {

	ticker := time.NewTicker(time.Millisecond)
	now := int64(time.Now().UnixNano() / 1000000)
	for {

		n.step = 0

		select {
		case c <- Flake((now-n.epoch)<<TimeShift | (n.node << NodeShift) | (n.step)):

			n.step = (n.step + 1) & StepMask

			if n.step == 0 {
				// wait for ticker..
			}
		case <-ticker.C:
			now++
			// continue
		}
	}
}

// Return a freshly generated Flake ID
func (n *Node) LockedGenerate() (Flake, error) {

	n.Lock()
	defer n.Unlock()

	now := time.Now().UnixNano() / 1000000

	if n.lastTime > now {
		return 0, fmt.Errorf("Invalid system time.")
	}

	if n.lastTime == now {
		n.step = (n.step + 1) & StepMask

		if n.step == 0 {
			for now <= n.lastTime {
				time.Sleep(100 * time.Microsecond)
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.lastTime = now

	return Flake((now-n.epoch)<<TimeShift |
		(n.node << NodeShift) |
		(n.step),
	), nil
}

// Return a freshly generated Flake ID
func (n *Node) Generate() (Flake, error) {

	now := time.Now().UnixNano() / 1000000

	if n.lastTime > now {
		return 0, fmt.Errorf("Invalid system time.")
	}

	if n.lastTime == now {
		n.step = (n.step + 1) & StepMask

		if n.step == 0 {
			for now <= n.lastTime {
				time.Sleep(100 * time.Microsecond)
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.lastTime = now

	return Flake((now-n.epoch)<<TimeShift |
		(n.node << NodeShift) |
		(n.step),
	), nil
}

type Flake int64
type Flakes []*Flake

func (f Flake) String() string {
	return fmt.Sprintf("%d", f)
}

func (f Flake) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}
func (f Flake) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

func (f Flake) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Byte())
}

func (f Flake) Byte() []byte {
	return []byte(f.String())
}

func (f Flake) Time() int64 {
	// ugh.. TODO
	// epoch is supposed to be configurable.....
	Epoch := time.Date(2016, 1, 0, 0, 0, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond)
	return (int64(f) >> 22) + Epoch
}

func (f Flake) Node() int64 {
	return int64(f) & 0x00000000003FF000 >> 12
}

func (f Flake) Sequence() int64 {
	return int64(f) & 0x0000000000000FFF
}
