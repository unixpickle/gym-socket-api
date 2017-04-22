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

func TestFlatten(t *testing.T) {
	ins := []Obs{
		jsonObs("[[1, 2], [3, 4], [5, 6]]"),
		jsonObs("[1, 3, 2]"),
		&uint8Obs{Dims: []int{2, 2}, Values: []uint8{3, 2, 10, 4}},
	}
	outs := [][]float64{
		{1, 2, 3, 4, 5, 6},
		{1, 3, 2},
		{3, 2, 10, 4},
	}
	for i, in := range ins {
		actual, err := Flatten(in)
		if err != nil {
			t.Errorf("case %d: %s", i, err)
			continue
		}
		expected := outs[i]
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("case %d: expected %v but got %v", i, expected, actual)
		}
	}
	failures := []Obs{
		jsonObs("1"),
		jsonObs("[1, 2, [1, 2]]"),
		jsonObs("[[1, 2], 1]"),
	}
	for i, in := range failures {
		if _, err := Flatten(in); err == nil {
			t.Errorf("failure case %d: should fail", i)
		}
	}
}
