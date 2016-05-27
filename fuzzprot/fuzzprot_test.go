package fuzzprot

import (
	"reflect"
	"testing"
)

// TestUnpack tests the unpacking routine
func TestUnpack(t *testing.T) {
	var tests = []struct {
		input []byte
		users []user
	}{
		{[]byte("\x01\x0346\x01\x03ADM\x02\x04Bill\x00"), []user{{Age: 46, Type: "ADM", Name: "Bill"}}},
	}

	for i, tt := range tests {
		if u, _ := Unpack(tt.input); !reflect.DeepEqual(tt.users, u) {
			t.Errorf("%d: Unpack(...)=%v want %v\n", i, u, tt.users)
		}
	}
}
