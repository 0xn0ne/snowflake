# snowflake

**ENGLISH** | [简体中文](./README-zh.md)

[![Build Status](https://travis-ci.org/0xNone/snowflake.svg?branch=master)](https://travis-ci.org/0xNone/snowflake) [![Coverage Status](https://coveralls.io/repos/github/0xNone/snowflake/badge.svg?branch=master)](https://coveralls.io/github/0xNone/snowflake?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/0xn0ne/snowflake)](https://goreportcard.com/report/github.com/0xn0ne/snowflake)

A snowflag id generation tool written by golang and based on the principle of twitter snowflag id.

For more introduction to snowflake ID, please [click me!](https://developer.twitter.com/en/docs/basics/twitter-ids.html)

The structure is like this:

```
1 bits                         41 bits                           10 bits         12 bits
│  0  │ 0 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 │ 00 0000 0000 │ 0000 0000 0000 │
unused                  time in milliseconds                   machine id      sequence id
```

In order to reduce the performance consumption caused by calculation and various uncertainties of the final generation form, the operation form of direct assignment is adopted to generate snowflake ID.

It was tested. It was good.

## Quick start

### Install

```bash
go get -u -v github.com/0xn0ne/snowflake
```

### Usage

The default attributes in the default snowflake id are `Overtime` (elapsed time after creation), `Sequence` (ID sequence number generated in the same microsecond), `Args` (other data stored).

Global variables that support modification include `Epoch`, `SequenceBits`, `ArgsOrder`, `ArgsBits`, high customization.

! ! ! Note: the total number of bits of Snowflake id is 64. if the modified `SequenceBits` and `ArgsBits` are occupy too many bits, leaving insufficient bits for timestamps. 

recommended to reserve a length of 40 bits for timestamps, i.e. if Epoch is set to the current time, Snowflake id can maintain no conflicts for 34 years.

```go
package main

import (
	"fmt"
	"github.com/0xn0ne/snowflake"
	"log"
)

func main() {
	// Custom Snowflake Epoch, it doesn't matter whether this step is taken or not
	// but it is recommended that Snowflake Epoch is set as the time when the project goes online.
	snowflake.Epoch = 1547963303708

	m, err := snowflake.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	sid := m.New(map[string]int64{}).(*snowflake.IdByDefault)

	fmt.Println("Raw:", sid)
	fmt.Println()

	fmt.Println("Base info:")
	fmt.Println("Over time      :", sid.Overtime)
	fmt.Println("Args           :", sid.Args)
	fmt.Println("Sequence       :", sid.Sequence)
	fmt.Println("Create time    :", sid.CreateTime())
	fmt.Println()

	fmt.Println("Conversion:")
	i64, err := sid.ToInt64()
	fmt.Println("int64          :", i64)
	b, err := sid.ToBytes()
	fmt.Println("[]byte         :", b)
	s, err := sid.ToString()
	fmt.Println("hex string     :", s)
	fmt.Println()

	sFromI64, _ := m.ParseInt64(i64)
	sFromByt, _ := m.ParseBytes(b)
	sFromStr, _ := m.ParseString(s)
	fmt.Println("Restore:")
	fmt.Println("from int64     :", sFromI64)
	fmt.Println("from hex string:", sFromByt)
	fmt.Println("from []byte    :", sFromStr)
}
```

#### Modify Args

By default, Args stores unused and machine data. if you think the default Args is not useful, you can manually modify it to the structure you want, such as:

```go
package main

import (
	"fmt"
	"github.com/0xn0ne/snowflake"
	"log"
)

func main() {
	snowflake.ArgsOrder = []string{"thread", "machine"}
	snowflake.ArgsBits = map[string]uint8{"thread": 4, "machine": 8}

	m, err := snowflake.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	sid := m.New(map[string]int64{"thread": 6, "machine": 4}).(*snowflake.IdByDefault)
	fmt.Println("Snowfalke id:", sid)
	fmt.Println("machine     :", sid.Args["machine"])
	fmt.Println("thread      :", sid.Args["thread"])

	i64, _ := sid.ToInt64()
	fmt.Println("to int64    :", i64)
	id, _ := m.ParseInt64(i64)
	fmt.Println("from int64  :", id)
}
```
