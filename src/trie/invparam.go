package trie

import (
	"errors"
	"strconv"
)

// InvParam represent type or inventory parameter
type InvParam uint

const (
	// BytesSize size of the values in bytes
	BytesSize InvParam = iota
	// KeysCount number of keys
	KeysCount
)

// String implements stringer interface to display InvParam as a string
func (p InvParam) String() string {
	switch p {
	case BytesSize:
		return "BytesSize"
	case KeysCount:
		return "KeysCount"
	}

	panic("Unknown InvParam: " + strconv.Itoa(int(p)))
}

// MarshalText renders InvParam as a string when marshalling
func (p InvParam) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

// UnmarshalText decodes InvParam value from string
func (p *InvParam) UnmarshalText(text []byte) error {
	var tmp InvParam
	switch string(text) {
	case "BytesSize":
		tmp = BytesSize
	case "KeysCount":
		tmp = KeysCount
	default:
		return errors.New("Unknown InvParam " + string(text))
	}

	*p = tmp

	return nil
}
