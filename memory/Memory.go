package memory

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Memory struct {
	Data []uint8
}

func NewMemory() *Memory {
	fmt.Println("Initializing Memory")
	return &Memory{Data: make([]uint8, 0xFFFF)}
}

func (m Memory) LoadROM(rom string) {
	fmt.Println("Loading ROM: " + rom)
	fileName := "../roms/" + rom
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening ROM file: ")
		log.Fatal(err)
		return
	}
	romData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading ROM file: ")
		log.Fatal(err)
	}
	defer file.Close()

	copy(m.Data[0x0000:], romData)
	fmt.Println("ROM Loaded")
}

func (m *Memory) GetByte(PC uint16) uint8 {
	return m.Data[PC]
}

func (m *Memory) StoreByte(location uint16, value uint8) {
	m.Data[location] = value
}
