package main

import (
	"GoBoy/internal"
	"GoBoy/memory"
	"fmt"
)

func main() {
	fmt.Println("Starting GoBoy Emulator")
	memory := memory.NewMemory()
	cpu := internal.NewCPU(memory)
	cpu.InitOpcodeTable()

	memory.LoadROM("cpu_instrs.gb")

	for {
		err := cpu.Cycle()
		if err {
			break
		}
	}
}
