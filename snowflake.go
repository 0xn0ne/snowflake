package snowflake

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// SnowflakeEpoch 默认设置为 twitter snowflake Epoch 即 Nov 04 2010 01:42:54 UTC
	// 可以自定义此选项，推荐 SnowflakeEpoch 设置为项目上线的时间（单位：毫秒）
	SnowflakeEpoch int64 = 1288834974657

	// 储存节点号的位数
	SnowflakeNodeBits uint8 = 10
	// 储存序列号的位数
	SnowflakeSequenceBits uint8 = 12
	// 储存经过时间的位数
	SnowflakeOvertimeBits uint8 = 64 - SnowflakeNodeBits - SnowflakeSequenceBits
	// 节点号最大数量
	SnowflakeNodeMax int64 = -1 ^ (-1 << SnowflakeNodeBits)
	// 序列号最大数量
	SnowflakeSequenceMax int64 = -1 ^ (-1 << SnowflakeSequenceBits)
	// 计算经过时间的时候偏移的位数
	SnowflakeOvertimeOffsetBits uint8 = SnowflakeNodeBits + SnowflakeSequenceBits
	// 计算节点号的时候偏移的位数
	SnowflakeNodeOffsetBits uint8 = SnowflakeSequenceBits
	// 计算节点号的时候偏移的位数
	SnowflakeSequenceOffsetBits uint8 = 0
)

// SnowflakeManager snowflake ID 的生成器用于保存所需的基本信息
type SnowflakeManager struct {
	Mut         sync.Mutex
	LastUseTime int64
	Node        int64
	Sequence    int64
}

// SnowflakeID
type SnowflakeID struct {
	Overtime int64
	Node     int64
	Sequence int64
}

// 新建可用于生成的新 Snowflake ID 的生成器
func NewSnowflakeManager(nodeNumber int64) (*SnowflakeManager, error) {
	if nodeNumber < 0 || nodeNumber > SnowflakeNodeMax {
		return nil, errors.New(fmt.Sprintf("SnowflakeManager number Mutst be between 0 and %d", SnowflakeNodeMax))
	}

	return &SnowflakeManager{
		LastUseTime: 0,
		Node:        nodeNumber,
		Sequence:    0,
	}, nil
}

// 创建一个 snowflake ID
func (self *SnowflakeManager) NewSnowflakeID() *SnowflakeID {
	self.Mut.Lock()

	// 为了防止在同一秒内容出现重复的 ID 在这里强制休眠一下
	var now int64
	for true {
		now = time.Now().UnixNano() / 1e6
		if self.Sequence < SnowflakeSequenceMax || now > self.LastUseTime {
			break
		} else {
			time.Sleep(time.Nanosecond)
		}
	}
	//if self.Sequence >= SnowflakeSequenceMax {
	//	for time.Now().UnixNano()/1e6 <= self.LastUseTime {
	//		time.Sleep(time.Nanosecond)
	//	}
	//}
	//now = time.Now().UnixNano() / 1e6
	if self.LastUseTime == now {
		self.Sequence += 1
	} else {
		self.LastUseTime = now
		self.Sequence = 0
	}

	snowflakeID := &SnowflakeID{Overtime: self.LastUseTime - SnowflakeEpoch, Node: self.Node, Sequence: self.Sequence}

	self.Mut.Unlock()
	return snowflakeID
}

// 创建一个 snowflake ID，并转化为 int64 类型
func (self *SnowflakeManager) NewSnowflakeIDInt64() int64 {
	snowflakeID := self.NewSnowflakeID()
	return snowflakeID.Int64()
}

// SnowflakeID 转化为 int64 类型
func (self *SnowflakeID) Int64() int64 {
	return (self.Overtime << SnowflakeOvertimeOffsetBits) | (self.Node << SnowflakeNodeOffsetBits) | (self.Sequence)
}

// 将一个 int64 类型变量解析为 SnowflakeID
func (self *SnowflakeManager) ParseInt64(it int64) *SnowflakeID {
	return &SnowflakeID{
		Overtime: (it >> SnowflakeOvertimeOffsetBits) & (1<<SnowflakeOvertimeBits - 1),
		Node:     (it >> SnowflakeNodeOffsetBits) & (1<<SnowflakeNodeBits - 1),
		Sequence: (it >> SnowflakeSequenceOffsetBits) & (1<<SnowflakeSequenceBits - 1)}
}

// 创建一个 snowflake ID，并转化为 []byte 类型
func (self *SnowflakeManager) NewSnowflakeIDBytes() ([]byte, error) {
	snowflakeID := self.NewSnowflakeID()
	it, err := snowflakeID.Bytes()
	if err != nil {
		return nil, err
	}
	return it, nil
}

// SnowflakeID 转化为 []byte 类型
func (self *SnowflakeID) Bytes() ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, self.Int64())
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// 将一个 []byte 类型变量解析为 SnowflakeID
func (self *SnowflakeManager) ParseBytes(it []byte) *SnowflakeID {
	snowflakeIDInt64 := int64(binary.BigEndian.Uint64(it))
	return self.ParseInt64(snowflakeIDInt64)
}

// 创建一个 snowflake ID，并转化为 string 类型
func (self *SnowflakeManager) NewSnowflakeIDString() string {
	snowflakeID := self.NewSnowflakeID()
	return snowflakeID.String()
}

// SnowflakeID 转化为 string 类型
func (self *SnowflakeID) String() string {
	return strconv.FormatInt(self.Int64(), 16)
}

// 将一个 string 类型变量解析为 SnowflakeID
func (self *SnowflakeManager) ParseString(it string) (*SnowflakeID, error) {
	snowflakeInt64, err := strconv.ParseInt(it, 16, 64)
	if err != nil {
		return nil, err
	}
	return self.ParseInt64(snowflakeInt64), nil
}
