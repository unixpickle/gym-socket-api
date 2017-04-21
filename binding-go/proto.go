package gym

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
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
