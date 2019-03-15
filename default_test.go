package snowflake

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewDefaultManager(t *testing.T) {
	tests := []struct {
		name    string
		want    Manager
		wantErr bool
	}{
		{"BaseTest", &ManagerByDefault{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDefaultManager()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultManager() = %v, want %v", got, tt.want)
			}
		})
	}

	srcArgsOrder := ArgsOrder
	srcArgsBits := ArgsBits

	ArgsOrder = []string{"a", "b"}
	_, err := NewDefaultManager()
	if (err != nil) != true {
		t.Errorf("NewDefaultManager() error = %v, wantErr %v", err, true)
		return
	}
	ArgsBits = map[string]uint8{"a": 64, "b": 64}
	_, err = NewDefaultManager()
	if (err != nil) != true {
		t.Errorf("NewDefaultManager() error = %v, wantErr %v", err, true)
		return
	}

	ArgsOrder = srcArgsOrder
	ArgsBits = srcArgsBits
}

func TestManagerByDefault_New(t *testing.T) {
	type fields struct {
		Mut         sync.Mutex
		LastUseTime int64
		Sequence    int64
	}
	type args struct {
		args map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ID
	}{
		{"BaseTest", fields{}, args{}, nil},
		{"BaseTest", fields{}, args{map[string]int64{"unused": 0, "machine": 12}}, nil},
		{"BaseTest", fields{}, args{map[string]int64{"unused": 0, "machine": 12}}, nil},
	}
	for _, tt := range tests {
		tArgs := map[string]int64{}
		for _, k := range ArgsOrder {
			tArgs[k] = tt.args.args[k]
		}
		t.Run(tt.name, func(t *testing.T) {
			m := &ManagerByDefault{
				Mut:         tt.fields.Mut,
				LastUseTime: tt.fields.LastUseTime,
				Sequence:    tt.fields.Sequence,
			}
			// Anti-interference sleep
			time.Sleep(time.Millisecond)
			tt.want = &IdByDefault{Args: tArgs, Overtime: time.Now().UnixNano()/1e6 - Epoch, Sequence: 0}
			if got := m.New(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ManagerByDefault.New() = %v, want %v", got, tt.want)
			}
		})
	}
	m := &ManagerByDefault{}
	sSnow := []*IdByDefault{}
	// Anti-interference sleep, goroutines test
	time.Sleep(time.Millisecond)
	for i := 0; i < 1<<SequenceBits/1024+2; i++ {
		go newSnowflakeId(sSnow, m)
	}

	for i := range sSnow {
		if i == 1<<SequenceBits && sSnow[i].Sequence != 0 {
			t.Errorf("ManagerByDefault.New() Next Time Seq = %v, want %v", sSnow[i].Sequence, 0)
		} else if sSnow[i+1].Sequence-1 != sSnow[i].Sequence {
			t.Errorf("ManagerByDefault.New() Seq = %v, want %v", sSnow[i].Sequence, sSnow[i-1].Sequence+1)
		}
	}
}

func newSnowflakeId(s []*IdByDefault, m *ManagerByDefault) {
	for i := 0; i < 1024; i++ {
		s = append(s, m.New(map[string]int64{}).(*IdByDefault))
	}
}

func TestManagerByDefault_NewToInt64(t *testing.T) {
	type fields struct {
		Mut         sync.Mutex
		LastUseTime int64
		Sequence    int64
	}
	type args struct {
		args map[string]int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{"BaseTest", fields{}, args{}, 0, false},
		{"BaseTest", fields{}, args{map[string]int64{"machine": 10}}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ManagerByDefault{
				Mut:         tt.fields.Mut,
				LastUseTime: tt.fields.LastUseTime,
				Sequence:    tt.fields.Sequence,
			}
			got, err := m.NewToInt64(tt.args.args)
			id, _ := m.ParseInt64(got)
			tt.want, _ = id.ToInt64()
			if (err != nil) != tt.wantErr {
				t.Errorf("ManagerByDefault.NewToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ManagerByDefault.NewToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManagerByDefault_NewBytes(t *testing.T) {
	type fields struct {
		Mut         sync.Mutex
		LastUseTime int64
		Sequence    int64
	}
	type args struct {
		args map[string]int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{"BaseTest", fields{}, args{}, []byte{}, false},
		{"BaseTest", fields{}, args{map[string]int64{"machine": 45}}, []byte{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ManagerByDefault{
				Mut:         tt.fields.Mut,
				LastUseTime: tt.fields.LastUseTime,
				Sequence:    tt.fields.Sequence,
			}
			got, err := m.NewBytes(tt.args.args)
			id, _ := m.ParseBytes(got)
			tt.want, _ = id.ToBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("ManagerByDefault.NewBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ManagerByDefault.NewBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManagerByDefault_NewString(t *testing.T) {
	type fields struct {
		Mut         sync.Mutex
		LastUseTime int64
		Sequence    int64
	}
	type args struct {
		args map[string]int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"BaseTest", fields{}, args{}, "", false},
		{"BaseTest", fields{}, args{map[string]int64{"machine": 998}}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ManagerByDefault{
				Mut:         tt.fields.Mut,
				LastUseTime: tt.fields.LastUseTime,
				Sequence:    tt.fields.Sequence,
			}
			got, err := m.NewString(tt.args.args)
			id, _ := m.ParseString(got)
			tt.want, _ = id.ToString()
			if (err != nil) != tt.wantErr {
				t.Errorf("ManagerByDefault.NewString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ManagerByDefault.NewString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdByDefault_CreateTime(t *testing.T) {
	type fields struct {
		Mut         sync.Mutex
		LastUseTime int64
		Sequence    int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"BaseTest", fields{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ManagerByDefault{
				Mut:         tt.fields.Mut,
				LastUseTime: tt.fields.LastUseTime,
				Sequence:    tt.fields.Sequence,
			}
			// Anti-interference sleep, goroutines test
			time.Sleep(time.Millisecond)
			tt.want = time.Now().UnixNano() / 1e6
			id := m.New(map[string]int64{})
			if got := id.CreateTime(); got != tt.want {
				t.Errorf("IdByDefault.CreateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
