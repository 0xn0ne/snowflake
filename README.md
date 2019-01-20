# snowflake

一个使用 golang 编写并做了部分魔改的 snowflake ID 生成工具

关于 snowflake ID 的介绍请[点击这里](https://developer.twitter.com/en/docs/basics/twitter-ids.html)

因为对每毫秒内的序列号和节点比较看重，魔改后现在的结构是这样子的：

```
                          42 bits                           10 bits         12 bits
│ 00 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 │ 00 0000 0000 │ 0000 0000 0000 │
                    time in milliseconds                    node id       sequence id
```

为了减少计算上带来性能消耗，以及最后生成形式的各种不确定性，生成 snowflake ID 采用了直接赋值的操作形式。让我值得高兴的是，测了一下，1 毫秒内最大序列号 4095 随随便便跑满。

## 快速开始

### 安装

```bash
go get -u github.com/0xNone/snowflake
```

### 用法

```cgo
package main

import (
	"fmt"
	"github.com/0xNone/snowflake"
)

func main() {
	// 自定义 Snowflake Epoch， 这步有没有无所谓， 推荐 SnowflakeEpoch 设置为项目上线的时间
	snowflake.SnowflakeEpoch = 1547963303708
	
	snowflakeManager, err := snowflake.NewSnowflakeManager(998)
	if err != nil {
		log.Fatal(err)
	}
	snowflakeID := snowflakeManager.NewSnowflakeID()

	fmt.Println("Raw:", *snowflakeID)
	fmt.Println()

	fmt.Println("Base info:")
	fmt.Println("Over time      :", snowflakeID.Overtime)
	fmt.Println("Node ID        :", snowflakeID.Node)
	fmt.Println("Sequence       :", snowflakeID.Sequence)
	fmt.Println("Create time    :", snowflakeID.CreateTime())
	fmt.Println()

	fmt.Println("Conversion:")
	fmt.Println("int64          :", snowflakeID.Int64())
	fmt.Println("hex string     :", snowflakeID.String())
	snowBytes, err := snowflakeID.Bytes()
	fmt.Println("[]byte         :", snowBytes)
	fmt.Println()

	tmpSnowflakeID, _ := snowflakeManager.ParseString(snowflakeID.String())
	fmt.Println("Restore:")
	fmt.Println("from int64     :", *snowflakeManager.ParseInt64(snowflakeID.Int64()))
	fmt.Println("from hex string:", *tmpSnowflakeID)
	fmt.Println("from []byte    :", *snowflakeManager.ParseBytes(snowBytes))
}
```

## TODO

补充测试