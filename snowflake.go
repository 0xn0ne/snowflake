package snowflake

var (
	Epoch        int64 = 1288834974657
	OvertimeBits uint8 = 64
	SequenceBits uint8 = 12

	OvertimeOffsetBits uint8 = SequenceBits
	SequenceOffsetBits uint8 = 0

	OvertimeMax int64 = -1 ^ (-1 << SequenceBits)
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
	// Create a new snowflag id.
	// Returns an ID.
	New(map[string]int64) ID

	// Create a new snowflag id and convert it to int64 type.
	// Returns int64 data that can be used to represent an ID.
	NewToInt64(map[string]int64) (int64, error)

	// Parsing snowflake id from int64 type data.
	// Returns the ID represented by the int64 data and possible errors.
	ParseInt64(int64) (ID, error)

	// Create a new snowflag id and convert it to []byte type.
	// Returns []byte data that can be used to represent an ID.
	NewBytes(map[string]int64) ([]byte, error)

	// Parsing snowflake id from []byte type data.
	// Returns the ID represented by the []byte data and possible errors.
	ParseBytes([]byte) (ID, error)

	// Create a new snowflag id and convert it to string type.
	// Returns string data that can be used to represent an ID.
	NewString(map[string]int64) (string, error)

	// Parsing snowflake id from string type data.
	// Returns the ID represented by the string data and possible errors.
	ParseString(string) (ID, error)
}

// create a new default id manager
// snowflagid is mainly generated/parsed by manager.
func NewManager() (Manager, error) {
	return NewDefaultManager()
}
