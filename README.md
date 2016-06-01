flake
====
[![GoDoc](https://godoc.org/github.com/bwmarrin/flake?status.svg)](https://godoc.org/github.com/bwmarrin/flake) [![Go report](http://goreportcard.com/badge/bwmarrin/flake)](http://goreportcard.com/report/bwmarrin/flake) [![Build Status](https://travis-ci.org/bwmarrin/flake.svg?branch=master)](https://travis-ci.org/bwmarrin/flake) 
[![Discord Gophers](https://img.shields.io/badge/Discord%20Gophers-%23flake.svg)](https://discord.gg/0f1SbxBZjYoCtNPP)

flake is a [Go](https://golang.org/) package that provides a very simple twitter
snowflake ID generator along with several functions to convert an ID into 
different formats.

## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

```sh
go get github.com/bwmarrin/flake
```

### Usage

Import the package into your project.

```go
import "github.com/bwmarrin/flake"
```

Construct a new flake Node that can be used to generate snowflake IDs then call
the Generate method to get a unique ID. The only argument to the NewNode() 
method is a Node number.  Each node you create must have it's own unique
Node number. A node number can be any number from 0 to 1023.

```go
node, err := flake.NewNode(1)
id := node.Generate()
fmt.Printf("ID: %d, %s\n", id, id.String())
```
