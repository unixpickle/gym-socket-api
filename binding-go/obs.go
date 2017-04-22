package gym

import (
	"encoding/json"
	"errors"

	"github.com/unixpickle/essentials"
)

// Obs is an observation from an environment.
type Obs interface {
	// Unmarshal unmarshals the observation into the
	// given object.
	//
	// This works the same way as json.Unmarshal.
	Unmarshal(dst interface{}) error
}

// Uint8Obs is an observation which can be converted to a
// flattened slice of raw 8-bit unsigned integers.
// It is useful for things like pixel data.
//
// When available, Uint8Obs() is typically much faster
// than Unmarshal().
//
// The slice returned by Uint8Obs is read-only.
// The caller should not modify it.
type Uint8Obs interface {
	Uint8Obs() []uint8
}

// Flatten turns a tensor observation into a 1-dimensional
// vector.
// This fails if the observation is not a tensor.
func Flatten(o Obs) ([]float64, error) {
	if u8, ok := o.(Uint8Obs); ok {
		nums := u8.Uint8Obs()
		res := make([]float64, len(nums))
		for i, x := range nums {
			res[i] = float64(x)
		}
		return res, nil
	}

	var sliceObs []interface{}
	if err := o.Unmarshal(&sliceObs); err != nil {
		return nil, essentials.AddCtx("flatten", err)
	}

	if res, ok := flatten(sliceObs); ok {
		return res, nil
	} else {
		return nil, errors.New("flatten: bad observation type")
	}
}

func flatten(obj []interface{}) ([]float64, bool) {
	if len(obj) == 0 {
		return nil, true
	}
	switch obj[0].(type) {
	case float64:
		res := make([]float64, len(obj))
		for i, x := range obj {
			f64, ok := x.(float64)
			if !ok {
				return nil, false
			}
			res[i] = f64
		}
		return res, true
	case []interface{}:
		var res []float64
		for _, child := range obj {
			subList, ok := child.([]interface{})
			if !ok {
				return nil, false
			}
			subRes, ok := flatten(subList)
			if !ok {
				return nil, false
			}
			res = append(res, subRes...)
		}
		return res, true
	default:
		return nil, false
	}
}

// jsonObs is an observation which was encoded as JSON.
type jsonObs []byte

func (j jsonObs) Unmarshal(dst interface{}) error {
	return json.Unmarshal(j, dst)
}

func (j jsonObs) String() string {
	return string(j)
}

// uint8Obs is an observation which was encoded as a raw
// array of 8-bit unsigned integers.
type uint8Obs struct {
	Dims   []int
	Values []uint8
}

// Unmarshal produces a JSON-compatible multi-dimensional
// array for the observation.
//
// This should be avoided for high-performance code.
// It is much more efficient to use the []uint8 directly.
func (u *uint8Obs) Unmarshal(dst interface{}) error {
	obj := u.jsonObject()
	data, _ := json.Marshal(obj)
	return json.Unmarshal(data, dst)
}

func (u *uint8Obs) Uint8Obs() []uint8 {
	return u.Values
}

func (u *uint8Obs) jsonObject() interface{} {
	if len(u.Dims) == 1 {
		res := make([]float64, len(u.Values))
		for i, x := range u.Values {
			res[i] = float64(x)
		}
		return res
	}
	chunkSize := len(u.Values) / u.Dims[0]
	var res []interface{}
	for i := 0; i < u.Dims[0]; i++ {
		chunk := &uint8Obs{
			Dims:   u.Dims[1:],
			Values: u.Values[i*chunkSize : (i+1)*chunkSize],
		}
		res = append(res, chunk.jsonObject())
	}
	return res
}
