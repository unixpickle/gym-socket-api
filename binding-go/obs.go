package gym

import "encoding/json"

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

// jsonObs is an observation which was encoded as JSON.
type jsonObs []byte

func (j jsonObs) Unmarshal(dst interface{}) error {
	return json.Unmarshal(j, dst)
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
