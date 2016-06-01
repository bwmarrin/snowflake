flake
====
[![GoDoc](https://godoc.org/github.com/bwmarrin/flake?status.svg)](https://godoc.org/github.com/bwmarrin/flake) [![Go report](http://goreportcard.com/badge/bwmarrin/flake)](http://goreportcard.com/report/bwmarrin/flake) [![Build Status](https://travis-ci.org/bwmarrin/flake.svg?branch=master)](https://travis-ci.org/bwmarrin/flake) [![Discord Gophers](https://img.shields.io/badge/Discord%20Gophers-general-blue.svg)](https://discord.gg/0f1SbxBZjYoCtNPP)

flake is a [Go](https://golang.org/) package that provides
* A very simple twitter snowflake Generator
* Methods to parse existing snowflake IDs
* Methods to convert a snowflake ID into several other data types.
* JSON Marshal/Unmarshal functions to easily use flake IDs within a JSON API.

## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

```sh
go get github.com/bwmarrin/flake
```

### Usage

Import the package into your project then construct a new flake Node using a
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
        "github.com/bwmarrin/flake"
)

func main() {

        // Create a new Node with a Node number of 1
        node, err := flake.NewNode(1)
        if err != nil {
                fmt.Println(err)
                return
        }

        // Generate a snowflake ID.
        id, err := node.Generate()
        if err != nil {
                fmt.Println(err)
                return
        }

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
}
```

### Performance

The flake generator should be sufficiently fast enough on most systems to 
generate 4096 unique ID's per millisecond. This is the maximum that the 
snowflake ID format supports. That is, around 243-244 nanoseconds per operation. 

Since the flake generator is single threaded the primary limitation will be
the maximum speed of a single processor on your system.

To benchmark the generator on your system run the following command inside the
flake package directory.

```sh
go test -bench=.
```

If your curious, check out this commit that shows benchmarks that compare a few 
different ways of implementing a snowflake generator in Go.
*  https://github.com/bwmarrin/flake/tree/9befef8908df13f4102ed21f42b083dd862b5036
