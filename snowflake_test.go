package snowflake

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var (
	SnfkMng     *SnowflakeManager
	NodeNumber  int64 = 0
	OvertimeMap       = make(map[int64]map[string]int64)
)

func init() {
	var err error
	SnfkMng, err = NewSnowflakeManager(NodeNumber)
	if err != nil {
		log.Fatal("Create generater error: ", err)
	}
}

func TestGenerater_NewSnowflakeID(t *testing.T) {
	snfkID := SnfkMng.NewSnowflakeID()
	if snfkID.Node != NodeNumber {
		t.Error("Node id is different")
	}
}

func TestGenerater_NewSnowflakeIDInt64(t *testing.T) {
	// 防干扰休眠
	time.Sleep(time.Millisecond)

	snfkID := SnfkMng.NewSnowflakeID()
	snfkIDInt64 := SnfkMng.NewSnowflakeIDInt64()
	snfkIDFromInt64 := SnfkMng.ParseInt64(snfkIDInt64)
	if snfkIDFromInt64.Overtime-snfkID.Overtime > 1 {
		t.Error("SnowflakeManager parsing error")
	}
	if !((snfkID.Overtime != snfkIDFromInt64.Overtime && snfkID.Sequence == snfkIDFromInt64.Sequence) ||
		(snfkID.Overtime == snfkIDFromInt64.Overtime && snfkID.Sequence != snfkIDFromInt64.Sequence)) {
		t.Error("Error creating SnowflakeID of int64 type")
	}
}

func TestSnowflakeManager_NewSnowflakeIDBytes(t *testing.T) {
	// 防干扰休眠
	time.Sleep(time.Millisecond)

	snfkID := SnfkMng.NewSnowflakeID()
	snfkIDBytes, _ := SnfkMng.NewSnowflakeIDBytes()
	snfkIDFromBytes := SnfkMng.ParseBytes(snfkIDBytes)
	if snfkIDFromBytes.Overtime-snfkID.Overtime > 1 {
		t.Error("SnowflakeManager parsing error")
	}
	if !((snfkID.Overtime != snfkIDFromBytes.Overtime && snfkID.Sequence == snfkIDFromBytes.Sequence) ||
		(snfkID.Overtime == snfkIDFromBytes.Overtime && snfkID.Sequence != snfkIDFromBytes.Sequence)) {
		t.Error("Error creating SnowflakeID of int64 type")
	}
}

func TestSnowflakeManager_NewSnowflakeIDString(t *testing.T) {
	// 防干扰休眠
	time.Sleep(time.Millisecond)

	snfkID := SnfkMng.NewSnowflakeID()
	snfkIDString := SnfkMng.NewSnowflakeIDString()
	snfkIDFromString, _ := SnfkMng.ParseString(snfkIDString)
	if snfkIDFromString.Overtime-snfkID.Overtime > 1 {
		t.Error("SnowflakeManager parsing error")
	}
	if !((snfkID.Overtime != snfkIDFromString.Overtime && snfkID.Sequence == snfkIDFromString.Sequence) ||
		(snfkID.Overtime == snfkIDFromString.Overtime && snfkID.Sequence != snfkIDFromString.Sequence)) {
		t.Error("Error creating SnowflakeID of int64 type")
	}
}

func BenchmarkGenerater_NewSnowflakeID(b *testing.B) {
	var lSnfkID []SnowflakeID
	// 防干扰休眠
	time.Sleep(time.Millisecond)
	for i := 0; i < 8192; i++ {
		lSnfkID = append(lSnfkID, *SnfkMng.NewSnowflakeID())
	}

	for _, snfkID := range lSnfkID {
		if _, ok := OvertimeMap[snfkID.Overtime]; !ok {
			OvertimeMap[snfkID.Overtime] = map[string]int64{"seq_min": SnowflakeSequenceMax, "seq_max": 0}
		}
		if snfkID.Sequence > OvertimeMap[snfkID.Overtime]["seq_max"] {
			OvertimeMap[snfkID.Overtime]["seq_max"] = snfkID.Sequence
		}
		if snfkID.Sequence < OvertimeMap[snfkID.Overtime]["seq_min"] {
			OvertimeMap[snfkID.Overtime]["seq_min"] = snfkID.Sequence
		}
	}

	b.Log("Generate raw benchmark results:", )
	for k, v := range OvertimeMap {
		if v["seq_min"] != 0 {
			b.Error("Wrong starting sequence number")
		}
		if v["seq_max"] > SnowflakeSequenceMax {
			b.Error("The same sequence number occurred in the same millisecond")
		}

		b.Logf(fmt.Sprintf("Sequence: %v\nmax %6v - min %6v\n", k, v["seq_max"], v["seq_min"]))
		delete(OvertimeMap, k)
	}
}

func BenchmarkSnowflakeManager_NewSnowflakeIDInt64(b *testing.B) {
	var lSnfkID []int64
	// 防干扰休眠
	time.Sleep(time.Millisecond)
	for i := 0; i < 8192; i++ {
		lSnfkID = append(lSnfkID, SnfkMng.NewSnowflakeIDInt64())
	}

	for _, snfkID := range lSnfkID {
		tmpSnfkID := SnfkMng.ParseInt64(snfkID)

		if _, ok := OvertimeMap[tmpSnfkID.Overtime]; !ok {
			OvertimeMap[tmpSnfkID.Overtime] = map[string]int64{"seq_min": SnowflakeSequenceMax, "seq_max": 0}
		}
		if tmpSnfkID.Sequence > OvertimeMap[tmpSnfkID.Overtime]["seq_max"] {
			OvertimeMap[tmpSnfkID.Overtime]["seq_max"] = tmpSnfkID.Sequence
		}
		if tmpSnfkID.Sequence < OvertimeMap[tmpSnfkID.Overtime]["seq_min"] {
			OvertimeMap[tmpSnfkID.Overtime]["seq_min"] = tmpSnfkID.Sequence
		}
	}

	b.Log("Generate int64 benchmark results:", )
	for k, v := range OvertimeMap {
		if v["seq_min"] != 0 {
			b.Error("Wrong starting sequence number")
		}
		if v["seq_max"] > SnowflakeSequenceMax {
			b.Error("The same sequence number occurred in the same millisecond")
		}

		b.Logf(fmt.Sprintf("Sequence: %v\nmax %6v - min %6v\n", k, v["seq_max"], v["seq_min"]))
		delete(OvertimeMap, k)
	}
}

func BenchmarkSnowflakeManager_NewSnowflakeIDBytes(b *testing.B) {
	var lSnfkID [][]byte
	// 防干扰休眠
	time.Sleep(time.Millisecond)
	for i := 0; i < 8192; i++ {
		it, _ := SnfkMng.NewSnowflakeIDBytes()
		lSnfkID = append(lSnfkID, it)
	}

	for _, snfkID := range lSnfkID {
		tmpSnfkID := SnfkMng.ParseBytes(snfkID)

		if _, ok := OvertimeMap[tmpSnfkID.Overtime]; !ok {
			OvertimeMap[tmpSnfkID.Overtime] = map[string]int64{"seq_min": SnowflakeSequenceMax, "seq_max": 0}
		}
		if tmpSnfkID.Sequence > OvertimeMap[tmpSnfkID.Overtime]["seq_max"] {
			OvertimeMap[tmpSnfkID.Overtime]["seq_max"] = tmpSnfkID.Sequence
		}
		if tmpSnfkID.Sequence < OvertimeMap[tmpSnfkID.Overtime]["seq_min"] {
			OvertimeMap[tmpSnfkID.Overtime]["seq_min"] = tmpSnfkID.Sequence
		}
	}

	b.Log("Generate []byte benchmark results:", )
	for k, v := range OvertimeMap {
		if v["seq_min"] != 0 {
			b.Error("Wrong starting sequence number")
		}
		if v["seq_max"] > SnowflakeSequenceMax {
			b.Error("The same sequence number occurred in the same millisecond")
		}

		b.Logf(fmt.Sprintf("Sequence: %v\nmax %6v - min %6v\n", k, v["seq_max"], v["seq_min"]))
		delete(OvertimeMap, k)
	}
}

func BenchmarkSnowflakeManager_NewSnowflakeIDString(b *testing.B) {
	var lSnfkID []string
	// 防干扰休眠
	time.Sleep(time.Millisecond)
	for i := 0; i < 8192; i++ {
		lSnfkID = append(lSnfkID, SnfkMng.NewSnowflakeIDString())
	}

	for _, snfkID := range lSnfkID {
		tmpSnfkID, _ := SnfkMng.ParseString(snfkID)

		if _, ok := OvertimeMap[tmpSnfkID.Overtime]; !ok {
			OvertimeMap[tmpSnfkID.Overtime] = map[string]int64{"seq_min": SnowflakeSequenceMax, "seq_max": 0}
		}
		if tmpSnfkID.Sequence > OvertimeMap[tmpSnfkID.Overtime]["seq_max"] {
			OvertimeMap[tmpSnfkID.Overtime]["seq_max"] = tmpSnfkID.Sequence
		}
		if tmpSnfkID.Sequence < OvertimeMap[tmpSnfkID.Overtime]["seq_min"] {
			OvertimeMap[tmpSnfkID.Overtime]["seq_min"] = tmpSnfkID.Sequence
		}
	}

	b.Log("Generate string benchmark results:", )
	for k, v := range OvertimeMap {
		if v["seq_min"] != 0 {
			b.Error("Wrong starting sequence number")
		}
		if v["seq_max"] > SnowflakeSequenceMax {
			b.Error("The same sequence number occurred in the same millisecond")
		}

		b.Logf(fmt.Sprintf("Sequence: %v\nmax %6v - min %6v\n", k, v["seq_max"], v["seq_min"]))
		delete(OvertimeMap, k)
	}
}
