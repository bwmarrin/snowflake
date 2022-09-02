snowflake
====
[![GoDoc](https://godoc.org/github.com/redmatter/snowflake?status.svg)](https://godoc.org/github.com/redmatter/snowflake)
[![Go report](http://goreportcard.com/badge/redmatter/snowflake)](http://goreportcard.com/report/redmatter/snowflake)
[![Coverage](http://gocover.io/_badge/github.com/redmatter/snowflake)](https://gocover.io/github.com/redmatter/snowflake)

### Features
* A very simple Twitter snowflake generator.
* Methods to parse existing snowflake IDs.
* Methods to convert a snowflake ID into several other data types and back.
* JSON Marshal/Unmarshal functions to easily use snowflake IDs within a JSON API.
* Monotonic Clock calculations protect from clock drift.
* Advanced use-case to generate a block of IDs.
  
### ID Format
By default, the ID format follows the original Twitter snowflake format.
* The ID is a 63 bit integer stored in an `int64`
* 41 bits are used to store a millisecond timestamp, using a custom epoch.
* 10 bits are used to store a node id, range from 0 through 1023.
* 12 bits are used to store a sequence number, range from 0 through 4095.

### Custom Format
You can alter the number of bits used for the node id and step number (sequence)
by specifying `Config.NodeBits` and `Config.StepBits` when initialising node
using `NewNodeWithConfig()`. Remember that there is a maximum of 22 bits available
that can be shared between the two. You do not have to use all 22 bits.

### Custom Epoch
By default the Twitter Epoch of 1288834974657 or Nov 04 2010 01:42:54 is used.
You can specify your own epoch value in milliseconds in `Config.Epoch` when
initialising node using `NewNodeWithConfig()`.

### How it Works.
Each time you generate an ID, it works, like this.
* A timestamp with millisecond precision is stored using 41 bits of the ID.
* Then the NodeID is added in subsequent bits.
* Then the Sequence Number is added, starting at 0 and incrementing for each ID
  generated in the same millisecond. If you generate enough IDs in the same
  millisecond, so that the sequence would roll over or overfill, then the generate 
  function will pause until the next millisecond.

The default Twitter format shown below.
```
+--------------------------------------------------------------------------+
| 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
+--------------------------------------------------------------------------+
```

Using the default settings, this allows for 4096 unique IDs to be generated every
millisecond, per Node ID.

## Getting Started

```sh
go get github.com/redmatter/snowflake
```

### Usage

```go
package main

import (
	"fmt"

	"github.com/redmatter/snowflake"
)

func main() {

	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())
}
```

### Performance

With default settings, this snowflake generator should be sufficiently fast 
enough on most systems to generate 4096 unique ID's per millisecond. This is 
the maximum that the snowflake ID format supports. That is, around 243-244 
nanoseconds per operation. 

Since the snowflake generator is single threaded the primary limitation will be
the maximum speed of a single processor on your system.

To benchmark the generator on your system run the following command inside the
snowflake package directory.

```sh
go test -run=^$ -bench=.
```

If your curious, check out this commit that shows benchmarks that compare a few 
different ways of implementing a snowflake generator in Go.
*  https://github.com/redmatter/snowflake/tree/9befef8908df13f4102ed21f42b083dd862b5036
