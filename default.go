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
	// 可手动调整的参数占用的比特数
	ArgsOrder      = []string{"unused", "machine"}
	ArgsBits       = map[string]uint8{"unused": 1, "machine": 10}
	ArgsOffsetBits = map[string]uint8{}
	ArgsMax        = map[string]int64{}
)

// 新建一个 id 管理
// snowflake id 主要依靠管理来生成
func NewDefaultManager() (Manager, error) {
	// 初始化 OvertimeBits
	var i64 uint8 = 0
	i64 += SequenceBits
	for _, k := range ArgsOrder {
		if _, ok := ArgsBits[k]; !ok {
			return nil, errors.New(fmt.Sprintf("ArgsBits does not have this key \"%v\"", k))
		}
		i64 += ArgsBits[k]
		if i64 > 64 || ArgsBits[k] > 64-SequenceBits {
			return nil, errors.New("OvertimeBits does not have enough storage space")
		}
	}
	OvertimeBits -= i64

	// 初始化 OffsetBits、Max 相关内容
	offsetBits := SequenceBits + SequenceOffsetBits
	OvertimeOffsetBits = offsetBits
	offsetBits += OvertimeBits
	OvertimeMax = -1 ^ (-1 << OvertimeBits)
	SequenceMax = -1 ^ (-1 << SequenceBits)
	for i := len(ArgsOrder) - 1; i >= 0; i-- {
		ArgsOffsetBits[ArgsOrder[i]] = offsetBits
		ArgsMax[ArgsOrder[i]] = -1 ^ (-1 << ArgsBits[ArgsOrder[i]])
		offsetBits += ArgsBits[ArgsOrder[i]]
	}
	return &ManagerByDefault{}, nil
}

// id 管理的结构体
type ManagerByDefault struct {
	Mut         sync.Mutex
	LastUseTime int64
	Sequence    int64
}

// snowflake id 的结构体
type IdByDefault struct {
	Args     map[string]int64
	Overtime int64
	Sequence int64
}

// 新建 snowflake id
// args 为 snowflake id 中存储的内容
// 如果 args 的 key 不存在则默认值为 0
// !!! 事实上这部分是不安全的，因为为了性能原因，没有对最大值做检查。
// 有一个方案是将 args 部分移动到 manager 处设置，或加数据库支持
func (m *ManagerByDefault) New(args map[string]int64) ID {
	m.Mut.Lock()

	// 为了防止在同一秒内容出现重复的 ID 在这里强制休眠一下
	var now int64
	for true {
		now = time.Now().UnixNano() / 1e6
		if m.Sequence < SequenceMax || now > m.LastUseTime {
			break
		} else {
			time.Sleep(time.Nanosecond)
		}
	}
	if m.LastUseTime == now {
		m.Sequence += 1
	} else {
		m.LastUseTime = now
		m.Sequence = 0
	}

	tArgs := map[string]int64{}
	for _, k := range ArgsOrder {
		tArgs[k] = args[k]
	}
	m.Mut.Unlock()
	return &IdByDefault{tArgs, m.LastUseTime - Epoch, m.Sequence}
}

// 新建 snowflake id，并转化为 int64
func (m *ManagerByDefault) NewToInt64(args map[string]int64) (int64, error) {
	return m.New(args).ToInt64()
}

func (m *ManagerByDefault) ParseInt64(i64 int64) (ID, error) {
	sequence := (i64 >> SequenceOffsetBits) & SequenceMax
	overtime := (i64 >> OvertimeOffsetBits) & OvertimeMax
	args := map[string]int64{}
	for i := len(ArgsOrder) - 1; i >= 0; i-- {
		args[ArgsOrder[i]] = (i64 >> ArgsOffsetBits[ArgsOrder[i]]) & ArgsMax[ArgsOrder[i]]
	}
	return &IdByDefault{
		Overtime: overtime,
		Sequence: sequence,
		Args:     args}, nil
}

func (m *ManagerByDefault) NewBytes(args map[string]int64) ([]byte, error) {
	return m.New(args).ToBytes()
}

func (m *ManagerByDefault) ParseBytes(b []byte) (ID, error) {
	i64 := int64(binary.BigEndian.Uint64(b))
	return m.ParseInt64(i64)
}

func (m *ManagerByDefault) NewString(args map[string]int64) (string, error) {
	return m.New(args).ToString()
}

func (m *ManagerByDefault) ParseString(s string) (ID, error) {
	i64, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return nil, err
	}
	return m.ParseInt64(i64)
}

func (id *IdByDefault) ToInt64() (int64, error) {
	var res int64 = (id.Overtime << OvertimeOffsetBits) | (id.Sequence << SequenceOffsetBits)
	for _, k := range ArgsOrder {
		res |= (id.Args[k] << ArgsOffsetBits[k])
	}
	return res, nil
}

func (id *IdByDefault) ToBytes() ([]byte, error) {
	buff := new(bytes.Buffer)
	iRes, _ := id.ToInt64()
	err := binary.Write(buff, binary.BigEndian, iRes)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (id *IdByDefault) ToString() (string, error) {
	iRes, _ := id.ToInt64()
	return fmt.Sprintf("%016s", strconv.FormatInt(iRes, 16)), nil
}

func (id *IdByDefault) CreateTime() int64 {
	return id.Overtime + Epoch
}
