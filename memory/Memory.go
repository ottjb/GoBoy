package memory

import (
	"fmt"
)

type Memory struct {
	cartridge   *Cartridge
	mbc         MBC
	vram        []byte
	externalram []byte
	wram        []byte
	oam         []byte
	io          []byte
	hram        []byte
	ie          byte
}

type MBC struct {
	Type       byte
	ROMBank    int
	RAMBank    int
	RAMEnabled bool
	Mode       bool
}

// NewMemory initializes the memory with a loaded cartridge
func NewMemory(cart *Cartridge) *Memory {
	mem := &Memory{
		cartridge:   cart,
		mbc:         MBC{Type: cart.mbcType, ROMBank: 1}, // Default to ROM bank 1
		vram:        make([]byte, 8*1024),                // 8KB
		externalram: make([]byte, cart.ramSize),
		wram:        make([]byte, 8*1024), // 8KB
		oam:         make([]byte, 160),    // 160 bytes
		io:          make([]byte, 128),    // 128 bytes
		hram:        make([]byte, 127),    // 127 bytes
		ie:          0,
	}

	// Initialize IO registers with default values
	mem.io[0x05] = 0x00 // TIMA
	mem.io[0x06] = 0x00 // TMA
	mem.io[0x07] = 0x00 // TAC
	mem.io[0x10] = 0x80 // NR10
	mem.io[0x11] = 0xBF // NR11
	mem.io[0x12] = 0xF3 // NR12
	mem.io[0x14] = 0xBF // NR14
	mem.io[0x16] = 0x3F // NR21
	mem.io[0x17] = 0x00 // NR22
	mem.io[0x19] = 0xBF // NR24
	mem.io[0x1A] = 0x7F // NR30
	mem.io[0x1B] = 0xFF // NR31
	mem.io[0x1C] = 0x9F // NR32
	mem.io[0x1E] = 0xBF // NR34
	mem.io[0x20] = 0xFF // NR41
	mem.io[0x21] = 0x00 // NR42
	mem.io[0x22] = 0x00 // NR43
	mem.io[0x23] = 0xBF // NR44
	mem.io[0x24] = 0x77 // NR50
	mem.io[0x25] = 0xF3 // NR51
	mem.io[0x26] = 0xF1 // NR52
	mem.io[0x40] = 0x91 // LCDC
	mem.io[0x41] = 0x85 // STAT
	mem.io[0x42] = 0x00 // SCY
	mem.io[0x43] = 0x00 // SCX
	mem.io[0x45] = 0x00 // LYC
	mem.io[0x47] = 0xFC // BGP
	mem.io[0x48] = 0xFF // OBP0
	mem.io[0x49] = 0xFF // OBP1
	mem.io[0x4A] = 0x00 // WY
	mem.io[0x4B] = 0x00 // WX
	mem.io[0xFF] = 0x00 // IE

	return mem
}

// 0x0000 - 0x3FFF: ROM Bank 0
// 0x4000 - 0x7FFF: ROM Bank 01 - NN (switchable)
// 0x8000 - 0x9FFF: Video RAM
// 0xA000 - 0xBFFF: External RAM (cartridge)
// 0xC000 - 0DFFF: Work RAM
// 0xE000 - 0xFDFF: Echo RAM (mirrors 0xC000 - 0xDDFF)
// 0xFE00 - 0xFE9F: Object Attribute Memory
// 0xFEA0 - 0xFEFF: Not usable
// 0xFF00 - 0xFF7F: I/O Registers
// 0xFF80 - 0xFFFE: High RAM
// 0xFFFF - 0xFFFF: Interrupt Enable Register

func (mem *Memory) Read(addr uint16) byte {
	if addr < 0x4000 {
		// ROM Bank 0
		return mem.cartridge.rom[addr]
	} else if addr < 0x8000 {
		// Switchable ROM bank
		offset := uint32(mem.mbc.ROMBank) * 0x4000
		return mem.cartridge.rom[offset+uint32(addr-0x4000)]
	} else if addr < 0xA000 {
		// VRAM
		return mem.vram[addr-0x8000]
	}

	fmt.Println("Unimplemented read")
	return 0
}

func (mem *Memory) Write(addr uint16, value byte) {
	fmt.Println("Unimplemented write")
}
