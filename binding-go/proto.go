package gym

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var byteOrder = binary.LittleEndian

const (
	packetReset        = 0
	packetStep         = 1
	packetGetSpace     = 2
	packetSampleAction = 3
	packetMonitor      = 4
	packetRender       = 5
)

const (
	observationJSON     = 0
	observationByteList = 1
)

const (
	actionJSON = 0
)

const (
	actionSpace      = 0
	observationSpace = 1
)

func handshake(rw *bufio.ReadWriter, envName string) error {
	if err := rw.WriteByte(0); err != nil {
		return err
	}
	if err := writeByteField(rw, []byte(envName)); err != nil {
		return err
	}
	if err := rw.Flush(); err != nil {
		return err
	}

	if errBytes, err := readByteField(rw); err != nil {
		return err
	} else if len(errBytes) > 0 {
		return errors.New(string(errBytes))
	}

	return nil
}

func writeByteField(w io.Writer, b []byte) error {
	if err := binary.Write(w, byteOrder, uint32(len(b))); err != nil {
		return err
	}
	_, err := w.Write(b)
	return err
}

func readByteField(r io.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(r, byteOrder, &length); err != nil {
		return nil, err
	}
	if length == 0 {
		return nil, nil
	}

	res := make([]byte, int(length))
	if _, err := io.ReadFull(r, res); err != nil {
		return nil, err
	}
	return res, nil
}

func writePacketType(w io.Writer, typeID int) error {
	_, err := w.Write([]byte{byte(typeID)})
	return err
}

func writeSpaceType(w io.Writer, typeID int) error {
	return writePacketType(w, typeID)
}

func readObservation(r io.Reader) (Obs, error) {
	var typeID uint8
	if err := binary.Read(r, byteOrder, &typeID); err != nil {
		return nil, err
	}
	obsData, err := readByteField(r)
	if err != nil {
		return nil, err
	}
	switch typeID {
	case observationJSON:
		return jsonObs(obsData), nil
	case observationByteList:
		return decodeUint8Obs(obsData)
	default:
		return nil, fmt.Errorf("unknown observation type: %d", typeID)
	}
}

func decodeUint8Obs(data []byte) (Obs, error) {
	r := bytes.NewReader(data)
	var numDims uint32
	if err := binary.Read(r, byteOrder, &numDims); err != nil {
		return nil, err
	}
	if numDims == 0 {
		return nil, errors.New("byte list has 0 dimensions")
	}
	dims := make([]int, int(numDims))
	product := 1
	for i := range dims {
		var dim uint32
		if err := binary.Read(r, byteOrder, &dim); err != nil {
			return nil, err
		}
		dims[i] = int(dim)
		product *= dims[i]
	}
	if int(product) != r.Len() {
		return nil, errors.New("incorrect byte list size")
	}
	return &uint8Obs{
		Dims:   dims,
		Values: data[len(data)-product:],
	}, nil
}

func readAction(r io.Reader, dst interface{}) error {
	var typeID uint8
	if err := binary.Read(r, byteOrder, &typeID); err != nil {
		return err
	}
	if typeID != 0 {
		return fmt.Errorf("unsupported action type: %d", typeID)
	}
	jsonData, err := readByteField(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, dst)
}

func writeAction(w io.Writer, act interface{}) error {
	jsonData, err := json.Marshal(act)
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte{actionJSON}); err != nil {
		return err
	}
	return writeByteField(w, jsonData)
}

func readReward(r io.Reader) (float64, error) {
	var res float64
	if err := binary.Read(r, byteOrder, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func readBool(r io.Reader) (bool, error) {
	var b uint8
	if err := binary.Read(r, byteOrder, &b); err != nil {
		return false, err
	}
	if b != 0 && b != 1 {
		return false, fmt.Errorf("invalid bool: %d", b)
	}
	return b != 0, nil
}

func writeBool(w io.Writer, b bool) error {
	var err error
	if b {
		_, err = w.Write([]byte{1})
	} else {
		_, err = w.Write([]byte{0})
	}
	return err
}
