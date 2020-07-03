package cpu

import (
	"errors"

	"github.com/anurse/gogb/pkg/gogb/memory"
)

// ErrStackOverflow occurs when a push would cause the stack pointer to wrap around.
var ErrStackOverflow error = errors.New("stack overflow")

// ErrStackUnderflow occurs when a push would cause the stack pointer to wrap around.
var ErrStackUnderflow error = errors.New("stack underflow")

func push(val uint16, sp *uint16, mem memory.MMU) error {
	if *sp <= 2 {
		return ErrStackOverflow
	}
	*sp -= 2
	return mem.SetWord(int(*sp), val)
}

func pop(sp *uint16, mem memory.MMU) (uint16, error) {
	res, err := mem.GetWord(int(*sp))
	if errors.Is(err, memory.ErrAddressOutOfRange) {
		return 0, ErrStackUnderflow
	}
	*sp += 2
	return res, err
}

func add8(left *uint8, right uint8, f *Z80Flags, withCarry bool) {
	if withCarry && f.IsSet(FlagCarry) {
		right++
	}

	result := int(*left) + int(right)

	f.SetIf(result > 0xFF, FlagCarry)
	result = result & 0xFF

	f.Clear(FlagAddSub)
	f.SetIf(result == 0, FlagZero)
	f.SetIf((right&0x0F)+(*left&0x0F) > 0x0F, FlagHalfCarry)

	*left = uint8(result & 0xFF)
}

func add16(left *uint16, right uint16, f *Z80Flags) {
	result := int(*left) + int(right)

	f.SetIf(result > 0xFFFF, FlagCarry)
	result = result & 0xFFFF

	f.Clear(FlagAddSub)
	f.SetIf((right&0xFFF)+(*left&0xFFF) > 0xFFF, FlagHalfCarry)

	*left = uint16(result & 0xFFFF)
}

func and(left *uint8, right uint8, f *Z80Flags) {
	*left = *left & right
	f.SetIf(*left == 0, FlagZero)
	f.Clear(FlagAddSub)
	f.Set(FlagHalfCarry)
	f.Clear(FlagCarry)
}

func bit(b uint8, val uint8, f *Z80Flags) {
	f.SetIf(val&(1<<b) == 0, FlagZero)
}

func call(addr uint16, pc *uint16, sp *uint16, mem memory.MMU) error {
	err := push(*pc, sp, mem)
	if err != nil {
		return err
	}
	*pc = addr
	return nil
}
