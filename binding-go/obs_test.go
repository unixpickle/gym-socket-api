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

func TestUnpackTuple(t *testing.T) {
	obj := jsonObs("[1, 2, [1, 2, 3]]")
	obses, err := UnpackTuple(obj)
	if err != nil {
		t.Fatal(err)
	} else if len(obses) != 3 {
		t.Fatalf("expected 3 observations but got %d", len(obses))
	}

	var obs1, obs2 int
	var obs3 []int
	ptrs := []interface{}{&obs1, &obs2, &obs3}
	for i, ptr := range ptrs {
		if err := obses[i].Unmarshal(ptr); err != nil {
			t.Fatalf("unmarshal obs %d: %s", i, err)
		}
	}
	expected := []interface{}{1, 2, []int{1, 2, 3}}
	for i, x := range expected {
		a := reflect.ValueOf(ptrs[i]).Elem().Interface()
		if !reflect.DeepEqual(a, x) {
			t.Errorf("obs %d: should be %v but got %v", i, x, a)
		}
	}
}
