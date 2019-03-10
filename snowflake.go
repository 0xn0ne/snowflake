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

// ID 是 snowflake id 的原始接口
// 以下函数表示可以使用的操作
type ID interface {
	// 转换为 int64 类型，并返回可能出现的错误
	ToInt64() (int64, error)

	// 转换为 []byte 类型，并返回可能出现的错误
	ToBytes() ([]byte, error)

	// 转换为 string 类型，并返回可能出现的错误
	ToString() (string, error)

	// 返回当前 snowflake 的创建时间
	CreateTime() int64
}

// Manager 是用于创建 snowflake id 的原始接口
// 以下函数表示可以使用的操作，说实在我也不知道为什么不叫他生成器
type Manager interface {
	// 新建 snowflake id
	New(map[string]int64) ID

	// 新建 snowflake id，并转化为 int64
	NewToInt64(map[string]int64) (int64, error)

	// 从 int64 中解析 snowflake id
	ParseInt64(int64) (ID, error)

	// 新建 snowflake id，并转化为 []byte
	NewBytes(map[string]int64) ([]byte, error)

	// 从 []byte 中解析 snowflake id
	ParseBytes([]byte) (ID, error)

	// 新建 snowflake id，并转化为 string
	NewString(map[string]int64) (string, error)

	// 从 string 中解析 snowflake id
	ParseString(string) (ID, error)
}

// 新建一个 id 管理
// snowflake id 主要依靠管理来生成
func NewManager() (Manager, error) {
	return NewDefaultManager()
}
