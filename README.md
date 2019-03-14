# snowflake

[![Build Status](https://travis-ci.org/0xNone/snowflake.svg?branch=master)](https://travis-ci.org/0xNone/snowflake) [![Coverage Status](https://coveralls.io/repos/github/0xNone/snowflake/badge.svg?branch=master)](https://coveralls.io/github/0xNone/snowflake?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/0xNone/snowflake)](https://goreportcard.com/report/github.com/0xNone/snowflake)

一个使用 golang 编写并根据 Twitter snowflake id 的原理做了的 snowflake ID 生成工具

关于 snowflake ID 的介绍请[点击这里](https://developer.twitter.com/en/docs/basics/twitter-ids.html)

结构是这样子的：

```
1 bits                         41 bits                           10 bits         12 bits
│  0  │ 0 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 │ 00 0000 0000 │ 0000 0000 0000 │
unused                  time in milliseconds                   machine id      sequence id
```

为了减少计算上带来性能消耗，以及最后生成形式的各种不确定性，生成 snowflake ID 采用了直接赋值的操作形式。测了一下，不错。

## 快速开始

### 安装

```bash
go get -u -v github.com/0xNone/snowflake
```

### 用法

默认的 snowflake id 中默认的属性有 Overtime（创建后经过的时间）、Sequence（id 的序列号）、Args（存储的其他内容）。

支持修改的全局变量有 Epoch、SequenceBits、ArgsOrder、ArgsBits，极大满足高自定义

！！！注意：Snowflake 的总 bits 数为 64，如果修改 SequenceBits、ArgsBits，占用比特过大，留给时间戳的空间就不够了，推荐预留 40 比特的长度给时间戳，即如果 Epoch 设置为当前时间，该 snowflake 可以保持 34 年内没有冲突。

```go
package main

import (
	"fmt"
	"github.com/0xNone/snowflake"
	"log"
)

func main() {
	// 自定义 Snowflake Epoch，这步有没有无所谓，推荐 snowflake.Epoch 设置为项目上线的时间
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

#### 修改 Args

默认情况下 Args 中存储在 unused、machine 两个数据，如果觉得默认的 Args 不好用，可以手动修改成你想要的结构，如：

```go
package main

import (
	"fmt"
	"github.com/0xNone/snowflake"
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
