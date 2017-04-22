package gym

import (
	"reflect"
	"testing"
)

func TestUint8Obs(t *testing.T) {
	obs := &uint8Obs{
		Dims: []int{3, 2, 4},
		Values: []uint8{
			117, 13, 176, 41,
			87, 92, 149, 189,

			36, 207, 227, 253,
			34, 13, 37, 48,

			98, 7, 225, 225,
			111, 5, 131, 133,
		},
	}
	var actual [][][]float64
	if err := obs.Unmarshal(&actual); err != nil {
		t.Fatal(err)
	}
	expected := [][][]float64{
		{
			{117, 13, 176, 41},
			{87, 92, 149, 189},
		},
		{
			{36, 207, 227, 253},
			{34, 13, 37, 48},
		},
		{
			{98, 7, 225, 225},
			{111, 5, 131, 133},
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}
