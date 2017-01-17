package gotyrant

import "fmt"

const (
	CodeExist    = int8(6)
	CodeNotExist = int8(7)
	CodeMisc     = int8(8)
)

type Error struct {
	error
	Code int8
}

func (e Error) Error() string {
	return fmt.Sprintf("Server error code %d", e.Code)
}

func (e Error) IsExist() bool {
	return e.Code == CodeExist
}

func (e Error) IsNotExist() bool {
	return e.Code == CodeNotExist
}
