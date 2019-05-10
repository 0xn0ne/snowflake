package snowflake

var (
	// Epoch defines a start time and recommended to be set as the project on-line time.
	Epoch        int64 = 1288834974657

	// OvertimeBits defines the number of bits occupied by Overtime. Automatic initialization.
	OvertimeBits uint8 = 64

	// SequenceBits defines the number of bits occupied by Sequence.
	SequenceBits uint8 = 12

	// OvertimeOffsetBits defines the offset of Overtime. Automatic initialization.
	OvertimeOffsetBits uint8 = SequenceBits

	// SequenceOffsetBits defines the offset of Sequence. Automatic initialization.
	SequenceOffsetBits uint8

	// OvertimeMax defines the maximum value of Overtime. Automatic initialization.
	OvertimeMax int64 = -1 ^ (-1 << OvertimeMax)

	// SequenceMax defines the maximum value of Sequence. Automatic initialization.
	SequenceMax int64 = -1 ^ (-1 << SequenceBits)
)

// ID is the raw interface of snowflake id.
// The following functions represent actions that can be used.
type ID interface {
	// convert id to int64 type
	ToInt64() (int64, error)

	// convert id to []byte type
	ToBytes() ([]byte, error)

	// convert id to string type
	ToString() (string, error)

	// Calculates id create time
	CreateTime() int64
}

// Manager is the raw interface used to create snowflake id.
// The following functions represent the operations that can be used.
// To be honest, I don't know why I don't call it generators.
type Manager interface {
	// New return a Snowflake ID interface.
	// This ID records key data.
	New(map[string]int64) ID

	// NewToInt64 return int64 data that can be used to represent an ID, and return possible errors.
	NewToInt64(map[string]int64) (int64, error)

	// ParseInt64 return parsed int64 data, and return possible errors.
	ParseInt64(int64) (ID, error)

	// NewBytes return []byte data that can be used to represent an ID, and return possible errors.
	NewBytes(map[string]int64) ([]byte, error)

	// ParseInt64 return parsed []byte data, and return possible errors.
	ParseBytes([]byte) (ID, error)

	// NewString return string data that can be used to represent an ID, and return possible errors.
	NewString(map[string]int64) (string, error)

	// ParseInt64 return parsed string data, and return possible errors.
	ParseString(string) (ID, error)
}

// NewManager retrun NewDefaultManager object.
// snowflagid is mainly generated/parsed by manager.
func NewManager() (Manager, error) {
	return NewDefaultManager()
}
