package memory

import (
	"errors"
)

// ErrAddressOutOfRange is returned if the address
// provided to a memory operation is out of the bounds of the memory.
var ErrAddressOutOfRange error = errors.New("address out of range")

// An MMU is an object that can provide a read/write interface to memory.
type MMU interface {
	GetByte(addr int) (uint8, error)
	GetWord(addr int) (uint16, error)

	SetByte(addr int, val uint8) error
	SetWord(addr int, val uint16) error
}

// A RAM is a simple byte array that implements the MMU interface
type RAM struct {
	data []uint8
}

// NewRAM creates a new empty RAM of the specified size
func NewRAM(size uint16) RAM {
	return RAM{data: make([]uint8, size)}
}

// GetByte reads a single byte at the specified address.
// Returns ErrAddressOutOfRange if the address is outside the bounds of this RAM
func (r *RAM) GetByte(addr int) (uint8, error) {
	if addr >= len(r.data) {
		return 0, ErrAddressOutOfRange
	}
	return r.data[addr], nil
}

// GetWord reads a 2-byte big-endian word at the specified address.
// Returns ErrAddressOutOfRange if the address is outside the bounds of this RAM
func (r *RAM) GetWord(addr int) (uint16, error) {
	if addr+1 >= len(r.data) {
		return 0, ErrAddressOutOfRange
	}
	return (uint16(r.data[addr]) << 8) | uint16(r.data[addr+1]), nil
}

// SetByte writes a single byte at the specified address.
// Returns ErrAddressOutOfRange if the address is outside the bounds of this RAM
func (r *RAM) SetByte(addr int, val uint8) error {
	if addr >= len(r.data) {
		return ErrAddressOutOfRange
	}
	r.data[addr] = val
	return nil
}

// SetWord writes a 2-byte big-endian word at the specified address.
// Returns ErrAddressOutOfRange if the address is outside the bounds of this RAM
func (r *RAM) SetWord(addr int, val uint16) error {
	if addr+1 >= len(r.data) {
		return ErrAddressOutOfRange
	}
	r.data[addr] = uint8((val & 0xFF00) >> 8)
	r.data[addr+1] = uint8(val & 0x00FF)
	return nil
}
