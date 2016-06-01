// Package flake provides a very simple Twitter Snowflake generator and parser.
// You can optionally set a custom epoch for you use.
package flake

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	timeBits = 41
	nodeBits = 10
	stepBits = 12

	timeMask int64 = -1 ^ (-1 << timeBits)
	nodeMask int64 = -1 ^ (-1 << nodeBits)
	stepMask int64 = -1 ^ (-1 << stepBits)

	nodeMax = -1 ^ (-1 << nodeBits)

	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

// Epoch is set to the twitter snowflake epoch of 2006-03-21:20:50:14 GMT
// You may customize this to set a different epoch for your application.
var Epoch int64 = 1288834974657

// A Node struct holds the basic information needed for a flake generator node
type Node struct {
	sync.Mutex
	time int64
	node int64
	step int64
}

// An ID is a custom type used for a snowflake ID.  This is used so we can
// attach methods onto the ID.
type ID int64

// NewNode returns a new Flake node that can be used to generate flake IDs
func NewNode(node int64) (*Node, error) {

	if node < 0 || node > nodeMax {
		return nil, fmt.Errorf("Node number must be between 0 and 1023")
	}

	return &Node{
		time: 0,
		node: node,
		step: 0,
	}, nil
}

// Generate creates and returns a unique snowflake ID
func (n *Node) Generate() (ID, error) {

	n.Lock()
	defer n.Unlock()

	now := time.Now().UnixNano() / 1000000

	if n.time == now {
		n.step = (n.step + 1) & stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	return ID((now-Epoch)<<timeShift |
		(n.node << nodeShift) |
		(n.step),
	), nil
}

// Int64 returns an int64 of the snowflake ID
func (f ID) Int64() int64 {
	return int64(f)
}

// String returns a string of the snowflake ID
func (f ID) String() string {
	return fmt.Sprintf("%d", f)
}

// Base2 returns a string base2 of the snowflake ID
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// Base36 returns a base36 string of the snowflake ID
func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

// Base64 returns a base64 string of the snowflake ID
func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

// Bytes returns a byte array of the snowflake ID
func (f ID) Bytes() []byte {
	return []byte(f.String())
}

// Time returns an int64 unix timestamp of the snowflake ID time
func (f ID) Time() int64 {
	return (int64(f) >> 22) + Epoch
}

// Node returns an int64 of the snowflake ID node number
func (f ID) Node() int64 {
	return int64(f) & 0x00000000003FF000 >> 12
}

// Step returns an int64 of the snowflake step (or sequence) number
func (f ID) Step() int64 {
	return int64(f) & 0x0000000000000FFF
}

// MarshalJSON returns a json byte array string of the snowflake ID.
func (f ID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + f.String() + `"`), nil
}

// UnmarshalJSON converts a json byte array of a snowflake ID into an ID type.
func (f *ID) UnmarshalJSON(b []byte) error {

	s := strings.Replace(string(b), `"`, ``, 2)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*f = ID(i)

	return nil
}
