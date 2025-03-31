package memory

import (
	"fmt"
	"time"
)

type Memory struct {
	cartridge  *Cartridge
	rtc        *RTC
	vram       [0x2000]byte
	wram       [0x2000]byte
	oam        [0xA0]byte // Sprite attribute table
	hram       [0x7F]byte
	io         [0x80]byte
	romBank    int
	ramBank    int
	ramEnabled bool
}

func NewMemory(cart *Cartridge, r *RTC) *Memory {
	fmt.Println("Initializing Memory")
	return &Memory{
		cartridge: cart,
		rtc:       r,
		romBank:   1,
		ramBank:   0,
	}
}

func (m *Memory) Read(addr uint16) byte {
	switch {
	case addr < 0x4000: // Fixed ROM bank 0
		return m.cartridge.rom[addr]

	case addr >= 0x4000 && addr < 0x8000: // Switchable ROM bank
		offset := int(addr-0x4000) + (m.romBank * 0x4000)
		return m.cartridge.rom[offset]

	case addr >= 0x8000 && addr < 0xA000: // Video RAM
		return m.vram[addr-0x8000]

	case addr >= 0xA000 && addr < 0xC000: // Cartridge RAM or RTC Registers
		if !m.ramEnabled {
			return 0xFF
		}

		if m.cartridge.mbcType == 3 && m.ramBank >= 0x08 && m.ramBank <= 0x0C {
			if !m.rtc.latched {
				m.UpdateRTC()
			}

			switch m.ramBank {
			case 0x08: // Seconds
				return m.rtc.seconds
			case 0x09: // Minutes
				return m.rtc.minutes
			case 0x0A: // Hours
				return m.rtc.hours
			case 0x0B: // Days Low
				return m.rtc.daysLow
			case 0x0C: // Days High
				return m.rtc.daysHigh
			}
			return 0xFF
		} else if m.cartridge.rom != nil && int(addr-0xA000) < len(m.cartridge.rom) {
			offset := int(addr-0xA000) + (m.ramBank * 0x2000)
			if offset < len(m.cartridge.rom) {
				return m.cartridge.rom[offset]
			}
		}
		return 0xFF

	case addr >= 0xC000 && addr < 0xE000: // Work RAM
		return m.wram[addr-0xC000]

	case addr >= 0xE000 && addr < 0xFE00: // Echo RAM (mirrors VRAM)
		return m.wram[addr-0xE000]

	case addr >= 0xFE00 && addr < 0xFEA0: // OAM
		return m.oam[addr-0xFE00]

	case addr >= 0xFF00 && addr < 0xFF80: // I/O Registers
		return m.io[addr-0xFF00]

	case addr >= 0xFF80 && addr < 0xFFFF: // High RAM
		return m.hram[addr-0xFF80]

	default:
		fmt.Printf("Unhandled read at 0x%04X\n", addr)
		return 0xFF
	}
}

func (m *Memory) Write(addr uint16, value byte) {
	switch {
	case addr == 0x6000: // Latch RTC
		if value == 0x00 {
			m.rtc.latched = false
		} else if value == 0x01 {
			m.UpdateRTC()
			m.rtc.latched = true
		}
	case addr < 0x2000: // Enable/Disable RAM
		m.ramEnabled = (value & 0x0F) == 0x0A

	case addr >= 0x2000 && addr < 0x4000: // ROM Bank Switch
		m.romBank = m.selectROMBank(value)

	case addr >= 0x4000 && addr < 0x6000: // RAM Bank Switch / Mode select
		m.ramBank = int(value & 0x03)

	case addr >= 0x6000 && addr < 0x8000: // Latch clock data
		if m.cartridge.mbcType == 3 {
			static := value & 0x01
			if static == 0x01 && !m.rtc.latched {
				m.UpdateRTC()
				m.rtc.latched = true
			} else if static == 0x00 {
				m.rtc.latched = false
			}
		}

	case addr >= 0x8000 && addr < 0xA000: // VRAM
		m.vram[addr-0x8000] = value

	case addr >= 0xA000 && addr < 0xC000: // Cartridge RAM or RTC Registers (MBC3)
		if !m.ramEnabled {
			return
		}

		if m.cartridge.mbcType == 3 && m.ramBank >= 0x08 && m.ramBank <= 0x0C {
			switch m.ramBank {
			case 0x08: // Seconds
				m.rtc.seconds = value
			case 0x09: // Minutes
				m.rtc.minutes = value
			case 0x0A: // Hours
				m.rtc.hours = value
			case 0x0B: // Days Low
				m.rtc.daysLow = value
			case 0x0C: // Writing to Days High register
				m.rtc.daysHigh = value
				m.rtc.latched = (value & 0x40) != 0
			}
		} else if m.cartridge.rom != nil {
			// Write to Cartridge RAM
			offset := int(addr-0xA000) + (m.ramBank * 0x2000)
			if offset < len(m.cartridge.rom) {
				m.cartridge.rom[offset] = value
			}
		}

	case addr >= 0xC000 && addr < 0xE000: // Work RAM
		m.wram[addr-0xC000] = value

	case addr >= 0xE000 && addr < 0xFE00: // Echo RAM
		m.wram[addr-0xFE00] = value

	case addr >= 0xFE00 && addr < 0xFEA0: // OAM
		m.oam[addr-0xFE00] = value

	case addr >= 0xFF00 && addr < 0xFF80: // I/O Registers
		m.io[addr-0xFF00] = value

	case addr >= 0xFF80 && addr < 0xFFFF: // High RAM
		m.hram[addr-0xFF80] = value

	default:
		fmt.Printf("Unhandled write at 0x%04X: %02X\n", addr, value)
	}
}

func (m *Memory) selectROMBank(value byte) int {
	switch m.cartridge.mbcType {
	case 1: // MBC1
		bank := int(value & 0x1F)
		if bank == 0 {
			bank = 1
		}
		return bank

	case 3: // MBC3
		return int(value & 0x7F)

	case 5: // MBC5
		return int(value)

	default:
		return 1
	}
}

func (m *Memory) UpdateRTC() {
	if m.rtc.daysHigh&0x40 != 0 {
		return // No update if latched
	}

	now := time.Now()
	elapsed := now.Sub(m.rtc.lastSync)

	seconds := int(m.rtc.seconds) + int(elapsed.Seconds())
	m.rtc.seconds = byte(seconds % 60)

	minutes := int(m.rtc.minutes) + (seconds / 60)
	m.rtc.minutes = byte(minutes % 60)

	hours := int(m.rtc.hours) + (minutes / 60)
	m.rtc.hours = byte(hours % 24)

	days := int(m.rtc.daysLow) + (hours / 24)
	m.rtc.daysLow = byte(days & 0xFF)
	if days > 0xFF {
		m.rtc.daysHigh = (m.rtc.daysHigh & 0xFE) | 0x01
		if (days >> 8) > 1 {
			m.rtc.daysHigh |= 0x80
		}
	}

	m.rtc.lastSync = now
}
