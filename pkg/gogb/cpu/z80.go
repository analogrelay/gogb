package cpu

import "github.com/anurse/gogb/pkg/gogb/memory"

// A Z80Flags represents a value that can be stored in the Z80's flags register
type Z80Flags uint8

// Set returns a new Z80Flags with the specified flag set. If it is already set, there is no effect.
func (f *Z80Flags) Set(flag Z80Flags) { *f |= flag }

// Clear returns a new Z80Flags with the specified flag clear. If it is already clear, there is no efect.
func (f *Z80Flags) Clear(flag Z80Flags) { *f &= ^flag }

// SetIf returns a new Z80Flags with the specified flag set based on the condition provided.
// If the condition is false, the flag is cleared.
func (f *Z80Flags) SetIf(condition bool, flag Z80Flags) {
	if condition {
		f.Set(flag)
	} else {
		f.Clear(flag)
	}
}

// IsSet returns a boolean indicating if the specified flag is set.
func (f Z80Flags) IsSet(flag Z80Flags) bool { return f&flag != 0 }

// IsClear returns a boolean indicating if the specified flag is set.
func (f Z80Flags) IsClear(flag Z80Flags) bool { return f&flag == 0 }

// Values for Z80Flags
const (
	FlagEmpty     Z80Flags = 0
	FlagCarry     Z80Flags = 1 << 4
	FlagHalfCarry Z80Flags = 1 << 5
	FlagAddSub    Z80Flags = 1 << 6
	FlagZero      Z80Flags = 1 << 7
)

// A State describes the current state of the CPU registers and clock.
type State struct {
	A       uint16
	B       uint16
	C       uint16
	D       uint16
	E       uint16
	H       uint16
	L       uint16
	F       Z80Flags
	PC      uint16
	SP      uint16
	TStates int
}

// A Z80 represents a Zilog 80 processor (configured for the GBA).
type Z80 struct {
	State  State
	Memory memory.MMU
}

// NewZ80 returns a new Z80 with default state and the specified memory unit.
func NewZ80(mem memory.MMU) Z80 {
	return Z80{
		State:  State{},
		Memory: mem,
	}
}
