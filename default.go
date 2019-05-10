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
	// ArgsOrder defines the order in which data is stored.
	// If the keys of ArgsBits are not in ArgsOrder, these keys will not be used.
	ArgsOrder      = []string{"unused", "machine"}

	// ArgsBits defines number of bits of snowflake id occupied by key.
	// It is recommended to set aside 40 bits for Overtime.
	// Please note that the default Sequence takes up 12 bits.
	ArgsBits       = map[string]uint8{"unused": 1, "machine": 10}

	// ArgsOffsetBits defines offset of each Args, Automatically calculate when new manager is created.
	ArgsOffsetBits = map[string]uint8{}

	// ArgsMax defines max of each Args, Automatically calculate when new manager is created.
	ArgsMax        = map[string]int64{}
)

// ManagerByDefault returns a Manager object and an error caused by some abnormal data
// All the generation and analysis are operated by the manager.
func NewDefaultManager() (Manager, error) {
	// Initialize OvertimeBits
	var i64 uint8 = 0
	i64 += SequenceBits
	for _, k := range ArgsOrder {
		if _, ok := ArgsBits[k]; !ok {
			return nil, fmt.Errorf("ArgsBits does not have this key \"%v\"", k)
		}
		i64 += ArgsBits[k]
		if i64 > 64 || ArgsBits[k] > 64-SequenceBits {
			return nil, errors.New("OvertimeBits does not have enough storage space")
		}
	}
	OvertimeBits -= i64

	// Initialize OffsetBits, Max related content
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

// id manager
// Mut dirty data lock.
// LastUseTime records the last time the manager was used.
// Sequence records the id serial number that was last generated using the manager.
type ManagerByDefault struct {
	Mut         sync.Mutex
	LastUseTime int64
	Sequence    int64
}

// Create a new snowflag id
// Args is snowflake id Args
// !!!In fact, this part is unsafe because the maximum value is not checked for performance reasons.
// One solution is to move the args part to the manager.
func (m *ManagerByDefault) New(args map[string]int64) ID {
	m.Mut.Lock()
	// Forced sleep, in order to prevent duplicate ID in the same microsecond content
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

// Create a new snowflag id and convert it to int64 type
func (m *ManagerByDefault) NewToInt64(args map[string]int64) (int64, error) {
	return m.New(args).ToInt64()
}

// Parsing snowflake id from int64 type data
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

// Create a new snowflag id and convert it to []byte type
func (m *ManagerByDefault) NewBytes(args map[string]int64) ([]byte, error) {
	return m.New(args).ToBytes()
}

// Parsing snowflake id from []byte type data
func (m *ManagerByDefault) ParseBytes(b []byte) (ID, error) {
	i64 := int64(binary.BigEndian.Uint64(b))
	return m.ParseInt64(i64)
}

// Create a new snowflag id and convert it to string type
func (m *ManagerByDefault) NewString(args map[string]int64) (string, error) {
	return m.New(args).ToString()
}

// Parsing snowflake id from string type data
func (m *ManagerByDefault) ParseString(s string) (ID, error) {
	i64, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return nil, err
	}
	return m.ParseInt64(i64)
}

// snowflake id
// Args is the content stored in snowflake id. If the key of args does not exist, the default value is 0.
// Overtime = now - Epoch (in ms)
// Sequence is the serial number in the same microsecond.
type IdByDefault struct {
	Args     map[string]int64
	Overtime int64
	Sequence int64
}

// convert id to int64 type
func (id *IdByDefault) ToInt64() (int64, error) {
	var res int64 = (id.Overtime << OvertimeOffsetBits) | (id.Sequence << SequenceOffsetBits)
	for _, k := range ArgsOrder {
		res |= (id.Args[k] << ArgsOffsetBits[k])
	}
	return res, nil
}

// convert id to []byte type
func (id *IdByDefault) ToBytes() ([]byte, error) {
	buff := new(bytes.Buffer)
	iRes, _ := id.ToInt64()
	err := binary.Write(buff, binary.BigEndian, iRes)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// convert id to string type
func (id *IdByDefault) ToString() (string, error) {
	iRes, _ := id.ToInt64()
	return fmt.Sprintf("%016s", strconv.FormatInt(iRes, 16)), nil
}

// Calculates id create time
// Please do not modify Epoch at will, otherwise CreateTime will be confused.
func (id *IdByDefault) CreateTime() int64 {
	return id.Overtime + Epoch
}

