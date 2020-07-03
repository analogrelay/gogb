package cpu

import (
	"testing"

	"github.com/anurse/gogb/pkg/gogb/memory"
	"github.com/stretchr/testify/assert"
)

func createStack() (uint16, memory.RAM) {
	var sp uint16 = 0x00FF
	mem := memory.NewRAM(0xFF)
	return sp, mem
}

// Tests for push/pop mean we can safely use them in other tests
func TestPush(t *testing.T) {
	sp, mem := createStack()
	assert.NoError(t, push(0xBEEF, &sp, &mem))
	assert.Equal(t, uint16(0xFD), sp)

	val, err := mem.GetWord(int(sp))
	assert.NoError(t, err)
	assert.Equal(t, uint16(0xBEEF), val)
}

func TestPushFailsWhenOutOfStackSpace(t *testing.T) {
	var sp uint16 = 0x00
	mem := memory.NewRAM(0xFF)
	assert.EqualError(t, push(0xBEEF, &sp, &mem), ErrStackOverflow.Error())
}

func TestPop(t *testing.T) {
	var sp uint16 = 0x00FD
	mem := memory.NewRAM(0xFF)
	assert.NoError(t, mem.SetWord(int(sp), 0xBEEF))

	val, err := pop(&sp, &mem)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0xBEEF), val)
}

func TestPopFailsWhenStackUnderflows(t *testing.T) {
	sp, mem := createStack()
	res, err := pop(&sp, &mem)
	assert.EqualError(t, err, ErrStackUnderflow.Error())
	assert.Equal(t, uint16(0), res)
}

// Arithmetic operations
func TestAdd8AddsOperandToAd(t *testing.T) {
	var a uint8 = 40
	f := FlagEmpty
	add8(&a, 2, &f, false)
	assert.Equal(t, uint8(42), a)
	assert.True(t, f.IsClear(FlagAddSub))
}

func TestAdcCarriesIn(t *testing.T) {
	var a uint8 = 0x01
	f := FlagCarry
	add8(&a, 1, &f, true)
	assert.Equal(t, uint8(3), a)
}

func TestAdd8SetsZeroFlagIfResultIsZero(t *testing.T) {
	var a uint8 = 0xFF
	f := FlagEmpty
	add8(&a, 1, &f, false)
	assert.Equal(t, uint8(0), a)
	assert.True(t, f.IsSet(FlagZero))
}

func TestAdd8ClearsZeroFlagIfResultIsNonZero(t *testing.T) {
	var a uint8 = 40
	f := FlagZero
	add8(&a, 2, &f, false)
	assert.Equal(t, uint8(42), a)
	assert.False(t, f.IsSet(FlagZero))
}

func TestAdd8SetsCarryOutIfResultOverflows(t *testing.T) {
	var a uint8 = 0xFF
	f := FlagEmpty
	add8(&a, 40, &f, false)
	assert.Equal(t, uint8(39), a)
	assert.True(t, f.IsSet(FlagCarry))
}

func TestAdd8ClearsCarryOutIfResultDoesNotOverflow(t *testing.T) {
	var a uint8 = 0x01
	f := FlagCarry
	add8(&a, 1, &f, false)
	assert.False(t, f.IsSet(FlagCarry))
}

func TestAdd8SetsHalfCarryIfLowNybbleOverflows(t *testing.T) {
	var a uint8 = 0x0A
	f := FlagEmpty
	add8(&a, 0x0A, &f, false)
	assert.Equal(t, uint8(20), a)
	assert.True(t, f.IsSet(FlagHalfCarry))
}

func TestAdd8ClearsHalfCarryIfLowNybbleDoesNotOverflow(t *testing.T) {
	var a uint8 = 0x01
	f := FlagHalfCarry
	add8(&a, 0x01, &f, false)
	assert.Equal(t, uint8(2), a)
	assert.True(t, f.IsClear(FlagHalfCarry))
}

func TestAdd16AddsValues(t *testing.T) {
	var hl uint16 = 0xAA00
	var sp uint16 = 0x00AA
	f := FlagEmpty
	add16(&hl, sp, &f)
	assert.Equal(t, uint16(0xAAAA), hl)
}

func TestAdd16SetsCarryIfResultOverflows(t *testing.T) {
	var hl uint16 = 0xFFFF
	var sp uint16 = 0x0002
	f := FlagEmpty
	add16(&hl, sp, &f)
	assert.Equal(t, uint16(0x0001), hl)
	assert.True(t, f.IsSet(FlagCarry))
}

func TestAdd16ClearsCarryIfResultDoesNotOverflow(t *testing.T) {
	var hl uint16 = 0x0003
	var sp uint16 = 0x0002
	f := FlagCarry
	add16(&hl, sp, &f)
	assert.Equal(t, uint16(0x0005), hl)
	assert.True(t, f.IsClear(FlagCarry))
}

func TestAdd16SetsHalfCarryIfLowNybbleOfHighByteOverflows(t *testing.T) {
	var hl uint16 = 0x0F00
	var sp uint16 = 0x0100
	f := FlagEmpty
	add16(&hl, sp, &f)
	assert.Equal(t, uint16(0x1000), hl)
	assert.True(t, f.IsSet(FlagHalfCarry))
}

func TestAdd16ClearsHalfCarryIfLowNybbleOfHighByteDoesNotOverflow(t *testing.T) {
	var hl uint16 = 0x00FF
	var sp uint16 = 0x0001
	f := FlagHalfCarry
	add16(&hl, sp, &f)
	assert.Equal(t, uint16(0x0100), hl)
	assert.True(t, f.IsClear(FlagHalfCarry))
}

func TestAnd(t *testing.T) {
	var a uint8 = 0b1010_1010
	f := FlagAddSub | FlagCarry | FlagZero
	and(&a, 0b1111_0000, &f)
	assert.Equal(t, uint8(0b1010_0000), a)
	assert.True(t, f.IsClear(FlagZero))
	assert.True(t, f.IsClear(FlagAddSub))
	assert.True(t, f.IsClear(FlagCarry))
	assert.True(t, f.IsSet(FlagHalfCarry))
}

func TestAndSetsZeroFlag(t *testing.T) {
	var a uint8 = 0b1010_1010
	f := FlagAddSub | FlagCarry | FlagZero
	and(&a, 0b0101_0101, &f)
	assert.Equal(t, uint8(0), a)
	assert.True(t, f.IsSet(FlagZero))
}

func TestBitSetsZeroFlagIfBitIsNotSet(t *testing.T) {
	f := FlagEmpty
	bit(6, 0b1011_1111, &f)
	assert.True(t, f.IsSet(FlagZero))
}

func TestBitClearsZeroFlagIfBitIsSet(t *testing.T) {
	f := FlagZero
	bit(6, 0b0100_0000, &f)
	assert.True(t, f.IsClear(FlagZero))
}

// Call/Jump Operations
func TestCallPushesPCValueToStack(t *testing.T) {
	sp, mem := createStack()
	var pc uint16 = 0xFEED
	assert.NoError(t, call(0xBEEF, &pc, &sp, &mem))
	res, err := pop(&sp, &mem)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0xFEED), res)
}

func TestCallSetsPCToSpecifiedValue(t *testing.T) {
	sp, mem := createStack()
	var pc uint16 = 0xFEED
	assert.NoError(t, call(0xBEEF, &pc, &sp, &mem))
	assert.Equal(t, uint16(0xBEEF), pc)
}
