package main

import (
	"GoBoy/internal"
	"GoBoy/memory"
	"fmt"
)

func main() {
	fmt.Println("Starting GoBoy Emulator")

	cart, err := memory.LoadCartridge("Tetris.gb")
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		cart.Debug()
	}

	m := memory.NewMemory(cart)

	cpu := internal.NewCPU(m)
	cpu.InitOpcodeTable()
	cpu.InitOpcodeCBTable()

	for {
		err := cpu.Cycle()
		if err {
			break
		}
	}
}
