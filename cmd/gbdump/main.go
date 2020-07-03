package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/anurse/gogb/pkg/gogb"
	"github.com/jessevdk/go-flags"
)

func main() {
	var opts struct {
		Verbose    []bool `short:"v" long:"verbose" description:"Show verbose logging information."`
		Positional struct {
			Files []string `required:"1" positional-arg-name:"ROM"`
		} `positional-args:"yes"`
	}
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.ParseArgs(os.Args[1:])

	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			return
		} else {
			panic(err)
		}
	}

	if len(opts.Positional.Files) == 0 {
		panic("You must provide at least one ROM to dump!")
	}

	for _, file := range opts.Positional.Files {
		dumpRom(file)
	}
}

func dumpRom(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	fmt.Println("Rom file ", file)

	// Slice the header and parse it
	var header gogb.CartridgeHeader
	headerBytes := content[0x0100:0x0150]
	err = gogb.ParseHeader(headerBytes, &header)
	if errors.Is(err, gogb.ErrHeaderChecksumInvalid) {
		fmt.Fprintln(os.Stderr, "  Warning: Header checksum validation failed.")
	} else if err != nil {
		panic(err)
	}

	fmt.Printf("  Size: 0x%X4\n", len(content))
	fmt.Println("  Title:", header.Title)
	fmt.Println("  Manufacturer Code:", header.ManufacturerCode)
	fmt.Println("  Color GameBoy Support:", header.CGBSupport)
	fmt.Println("  New Licensee Code:", header.NewLicenseeCode)
	fmt.Println("  Old Licensee Code:", header.OldLicenseeCode)
	fmt.Println("  Super GameBoy Support:", header.SGBSupport)
	fmt.Println("  Type:", header.Type)
	fmt.Printf("  ROM Size: %dKB\n", header.ROMSize)
	fmt.Printf("  RAM Size: %dKB\n", header.RAMSize)
	fmt.Println("  Japanese?:", header.Japanese)
	fmt.Println("  Version:", header.VersionNumber)

	// Compute global checksum
	var actualChecksum uint16
	for idx, byt := range content {
		if idx != 0x014E && idx != 0x014F {
			actualChecksum += uint16(byt)
		}
	}
	if actualChecksum == header.GlobalChecksum {
		fmt.Println("  Cartridge Checksum VERIFIED")
	} else {
		fmt.Println("  Cartridge Checksum NOT VERIFIED")
		fmt.Printf("    Expected: 0x%X4, Actual 0x%X4\n", header.GlobalChecksum, actualChecksum)
	}
}
