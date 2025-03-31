package memory

import (
	"encoding/hex"
	"fmt"
	"os"
)

type Cartridge struct {
	rom     []byte
	title   string
	mbcType byte
	romSize int
	ramSize int
}

func (cart *Cartridge) Debug() {
	fmt.Println("Game Title:", cart.title)
	fmt.Println("MBC Type:", hex.EncodeToString([]byte{cart.mbcType}))
	fmt.Println("ROM Size:", cart.romSize, "bytes")
	fmt.Println("RAM Size:", cart.ramSize, "bytes")
}

func LoadCartridge(fileName string) (*Cartridge, error) {
	data, err := os.ReadFile("../roms/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to load ROM: %w", err)
	}

	if len(data) < 0x150 {
		return nil, fmt.Errorf("invalid ROM file, too small")
	}

	return &Cartridge{
		rom:     data,
		title:   parseTitle(data),
		mbcType: data[0x147],
		romSize: getROMSize(data[0x148]),
		ramSize: getRAMSize(data[0x149]),
	}, nil
}

func parseTitle(data []byte) string {
	titleBytes := data[0x134:0x144]
	title := string(titleBytes)
	return title
}

func getROMSize(code byte) int {
	romSizes := []int{32 * 1024, 64 * 1024, 128 * 1024, 256 * 1024, 512 * 1024, 1024 * 1024, 2048 * 1024, 4096 * 1024, 8192 * 1024}
	if int(code) < len(romSizes) {
		return romSizes[code]
	}
	return 0
}

func getRAMSize(code byte) int {
	ramSizes := []int{0, 2048, 8192, 32768, 131072}
	if int(code) < len(ramSizes) {
		return ramSizes[code]
	}
	return 0
}
