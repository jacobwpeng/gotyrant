package gotyrant

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/jacobwpeng/goutil"
)

type Request struct {
	magic  uint16
	buffer bytes.Buffer
	bin    [][]byte
}

func (req *Request) SetMagic(magic uint16) {
	req.magic = magic
}

func (req *Request) AddUInt32(v uint32) {
	err := binary.Write(&req.buffer, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
}

func (req *Request) AddBinary(data []byte) {
	req.bin = append(req.bin, data)
}

func (req *Request) WriteTo(out io.Writer) (n int64, err error) {
	w := goutil.NewCountWriter(out)
	err = binary.Write(w, binary.BigEndian, req.magic)
	if err != nil {
		return w.Count(), fmt.Errorf("Write magic error: %s", err)
	}
	_, err = io.Copy(w, &req.buffer)
	if err != nil {
		return w.Count(), fmt.Errorf("Write buffer error: %s", err)
	}
	for index, data := range req.bin {
		_, err := w.Write(data)
		if err != nil {
			return w.Count(), fmt.Errorf("Write data#%d error: %s", index, err)
		}
	}
	return w.Count(), nil
}
