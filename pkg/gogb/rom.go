package gogb

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// A CartridgeType represents the hardware present in a GBA cartridge.
type CartridgeType uint8

// Defines the known cartridge types
const (
	ROMOnly                    CartridgeType = 0x00
	Mbc1                                     = 0x01
	Mbc1Ram                                  = 0x02
	Mbc1RamBattery                           = 0x03
	Mbc2                                     = 0x05
	Mbc2Battery                              = 0x06
	ROMRAM                                   = 0x08
	ROMRAMBattery                            = 0x09
	Mmm01                                    = 0x0B
	Mmm01RAM                                 = 0x0C
	Mmm01RAMBattery                          = 0x0D
	Mbc3TimerBattery                         = 0x0F
	Mbc3TimerRAMBattery                      = 0x10
	Mbc3                                     = 0x11
	Mbc3Ram                                  = 0x12
	Mbc3RamBattery                           = 0x13
	Mbc5                                     = 0x19
	Mbc5Ram                                  = 0x1A
	Mbc5RamBattery                           = 0x1B
	Mbc5Rumble                               = 0x1C
	Mbc5RumbleRAM                            = 0x1D
	Mbc5RumbleRAMBattery                     = 0x1E
	Mbc6                                     = 0x20
	Mbc7SensorRumbleRAMBattery               = 0x22
	PocketCamera                             = 0xFC
	BandaiTama5                              = 0xFD
	Huc3                                     = 0xFE
	Huc1RamBattery                           = 0xFF
)

func (v CartridgeType) String() string {
	switch v {
	case ROMOnly:
		return "ROMOnly"
	case Mbc1:
		return "Mbc1"
	case Mbc1Ram:
		return "Mbc1Ram"
	case Mbc1RamBattery:
		return "Mbc1RamBattery"
	case Mbc2:
		return "Mbc2"
	case Mbc2Battery:
		return "Mbc2Battery"
	case ROMRAM:
		return "ROMRAM"
	case ROMRAMBattery:
		return "ROMRAMBattery"
	case Mmm01:
		return "Mmm01"
	case Mmm01RAM:
		return "Mmm01RAM"
	case Mmm01RAMBattery:
		return "Mmm01RAMBattery"
	case Mbc3TimerBattery:
		return "Mbc3TimerBattery"
	case Mbc3TimerRAMBattery:
		return "Mbc3TimerRAMBattery"
	case Mbc3:
		return "Mbc3"
	case Mbc3Ram:
		return "Mbc3Ram"
	case Mbc3RamBattery:
		return "Mbc3RamBattery"
	case Mbc5:
		return "Mbc5"
	case Mbc5Ram:
		return "Mbc5Ram"
	case Mbc5RamBattery:
		return "Mbc5RamBattery"
	case Mbc5Rumble:
		return "Mbc5Rumble"
	case Mbc5RumbleRAM:
		return "Mbc5RumbleRAM"
	case Mbc5RumbleRAMBattery:
		return "Mbc5RumbleRAMBattery"
	case Mbc6:
		return "Mbc6"
	case Mbc7SensorRumbleRAMBattery:
		return "Mbc7SensorRumbleRAMBattery"
	case PocketCamera:
		return "PocketCamera"
	case BandaiTama5:
		return "BandaiTama5"
	case Huc3:
		return "Huc3"
	case Huc1RamBattery:
		return "Huc1RamBattery"
	default:
		return fmt.Sprintf("Unknown(0x%X2)", uint8(v))
	}
}

// A CgbSupport value indicates if a cartridge supports CGB (Color GameBoy) features
type CgbSupport uint8

// Values that indicate if a cartridge supports CGB (Color GameBoy) features
const (
	CgbNotSupported CgbSupport = iota
	CgbSupported
	CgbRequired
)

func (v CgbSupport) String() string {
	return [...]string{"NotSupported", "Supported", "Required"}[v]
}

// A CartridgeHeader represents the header of a GBA ROM.
type CartridgeHeader struct {
	// The title of the ROM.
	Title string

	// The Manufacturer Code of the ROM.
	ManufacturerCode string

	// A boolean indicating if CGB features are required.
	CGBSupport CgbSupport

	// The New Licensee Code of the ROM.
	NewLicenseeCode string

	// A boolean indicating if SGB features are supported.
	SGBSupport bool

	// A value indicating the type of the cartridge.
	Type CartridgeType

	// A value indicating the ROM size (in KB) of the cartridge.
	ROMSize int

	// A value indicating the RAM size (in KB) of the cartridge.
	RAMSize int

	// A boolean indicating if this version of the game is to be sold in Japan.
	Japanese bool

	// The Old Licensee Code of the ROM. A value of 0x33 indicates that the NewLicenseeCode should be used.
	OldLicenseeCode byte

	// The version number of the game.
	VersionNumber byte

	// A global checksum over the cartridge data.
	GlobalChecksum uint16
}

// ErrHeaderLengthInvalid indicates that the provided header data was not the correct size.
var ErrHeaderLengthInvalid error = errors.New("header data is not 0x50 bytes long")

// ErrHeaderChecksumInvalid indicates that the header checksum does not match the actual header data.
var ErrHeaderChecksumInvalid error = errors.New("the header checksum could not be validated")

// ParseHeader parses the provided header and fills in the CartridgeHeader struct provided.
// If ErrHeaderChecksumInvalid is returned, the provided CartridgeHeader struct will **still** be filled in with data!
func ParseHeader(inp []byte, header *CartridgeHeader) error {
	if len(inp) != 0x50 {
		return ErrHeaderLengthInvalid
	}

	// Read the title
	header.Title = strings.TrimRight(string(inp[0x34:0x3F]), "\x00")
	header.ManufacturerCode = strings.TrimRight(string(inp[0x3F:0x42]), "\x00")

	cgbVal := inp[0x43]
	if cgbVal == 0x80 {
		header.CGBSupport = CgbSupported
	} else if cgbVal == 0xC0 {
		header.CGBSupport = CgbRequired
	} else {
		header.CGBSupport = CgbNotSupported
	}

	header.NewLicenseeCode = string(inp[0x44:0x45])

	header.SGBSupport = inp[0x46] == 0x03

	header.Type = CartridgeType(inp[0x47])

	header.ROMSize = getROMSize(inp[0x48])
	header.RAMSize = getRAMSize(inp[0x49])

	header.Japanese = inp[0x4A] == 0x00

	header.OldLicenseeCode = inp[0x4B]
	header.VersionNumber = inp[0x4C]

	// Read global checksum
	header.GlobalChecksum = binary.BigEndian.Uint16(inp[0x4E:0x50])

	// Compute checksum. We still fill the header structure even if the checksum fails, but we want to return an error so the user knows
	headerChecksum := inp[0x4D]
	var actualChecksum byte

	for x := 0x34; x <= 0x4C; x++ {
		actualChecksum = actualChecksum - inp[x] - 1
	}

	if headerChecksum != actualChecksum {
		return ErrHeaderChecksumInvalid
	}

	return nil
}

func getROMSize(size byte) int {
	switch size {
	case 0x00:
		return 2 * 16
	case 0x01:
		return 4 * 16
	case 0x02:
		return 8 * 16
	case 0x03:
		return 16 * 16
	case 0x04:
		return 32 * 16
	case 0x05:
		return 128 * 16
	case 0x06:
		return 256 * 16
	case 0x08:
		return 512 * 16
	case 0x52:
		return 72 * 16
	case 0x53:
		return 80 * 16
	case 0x54:
		return 96 * 16
	default:
		return 0
	}
}

func getRAMSize(size byte) int {
	switch size {
	case 0x01:
		return 2
	case 0x02:
		return 8
	case 0x03:
		return 32
	case 0x04:
		return 128
	case 0x05:
		return 64
	default:
		return 0
	}
}
