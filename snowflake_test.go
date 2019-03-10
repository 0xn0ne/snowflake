package snowflake

import (
	"reflect"
	"testing"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name    string
		want    Manager
		wantErr bool
	}{
		{"BaseTest", &ManagerByDefault{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManager()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}
