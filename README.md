snowflake
====
[![GoDoc](https://godoc.org/github.com/bwmarrin/snowflake?status.svg)](https://godoc.org/github.com/bwmarrin/snowflake) [![Go report](http://goreportcard.com/badge/bwmarrin/snowflake)](http://goreportcard.com/report/bwmarrin/snowflake) [![Build Status](https://travis-ci.org/bwmarrin/snowflake.svg?branch=master)](https://travis-ci.org/bwmarrin/snowflake) [![Discord Gophers](https://img.shields.io/badge/Discord%20Gophers-%23info-blue.svg)](https://discord.gg/0f1SbxBZjYq9jLBk)

snowflake is a [Go](https://golang.org/) package that provides
* A very simple Twitter snowflake generator.
* Methods to parse existing snowflake IDs.
* Methods to convert a snowflake ID into several other data types.
* JSON Marshal/Unmarshal functions to easily use snowflake IDs within a JSON API.

**For help with this package or general Go discussion, please join the [Discord 
Gophers](https://discord.gg/0f1SbxBZjYq9jLBk) chat server.**

## Status @ 2017-02-21
This package should be considered stable and completed.  Any additions in the 
future will strongly avoid API changes to existing functions.  Please see issues
for any remaining TODO items that are planned.
  
### ID Format
By default, the ID format follows the original Twitter snowflake format.
* The ID as a whole is a 63 bit integer stored in an int64
* 41 bits are used to store a timestamp with millisecond precision, using a custom epoch.
* 10 bits are used to store a node id - a range from 0 through 1023.
* 12 bits are used to store a sequence number - a range from 0 through 4095.

### Custom Format
Coming soon.. ish.  I plan to find the best way to modify the existing package to allow altering the number of bits used for Node ID's and Sequence Numbers. The goal is to implement this without breaking the API (1/31/18)

### How it Works.
Each time you generate an ID, it works, like this.
* A timestamp with millisecond precision is stored in the first 41 bits of the ID.
* Then the NodeID is added in subsequent bits.
* Then the Sequence Number is added, starting at 0 and incrementing for each ID generated in the same millisecond. If you generate enough IDs in the same millisecond that the sequence would roll over or overfill then the generate function will pause until the next millisecond.

Using the default settings, this allows for 4096 unique IDs to be generated every millisecond, per Node ID.
## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

```sh
go get github.com/bwmarrin/snowflake
```


### Usage

Import the package into your project then construct a new snowflake Node using a
unique node number from 0 to 1023. With the node object call the Generate() 
method to generate and return a unique snowflake ID. 

Keep in mind that each node you create must have a unique node number, even 
across multiple servers.  If you do not keep node numbers unique the generator 
cannot guarantee unique IDs across all nodes.


**Example Program:**

```go
package main

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
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

	// Print out the ID's timestamp
	fmt.Printf("ID Time  : %d\n", id.Time())

	// Print out the ID's node number
	fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	fmt.Printf("ID Step  : %d\n", id.Step())

  // Generate and print, all in one.
  fmt.Printf("ID       : %d\n", node.Generate().Int64())
}
```

### Performance

This snowflake generator should be sufficiently fast enough on most systems to 
generate 4096 unique ID's per millisecond. This is the maximum that the 
snowflake ID format supports. That is, around 243-244 nanoseconds per operation. 

Since the snowflake generator is single threaded the primary limitation will be
the maximum speed of a single processor on your system.

To benchmark the generator on your system run the following command inside the
snowflake package directory.

```sh
go test -bench=.
```

If your curious, check out this commit that shows benchmarks that compare a few 
different ways of implementing a snowflake generator in Go.
*  https://github.com/bwmarrin/snowflake/tree/9befef8908df13f4102ed21f42b083dd862b5036
