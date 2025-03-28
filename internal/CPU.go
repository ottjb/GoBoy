package internal

import (
	"GoBoy/memory"
	"fmt"
)

type CPU struct {
	A, B, C, D, E, H, L, F uint8
	SP, PC                 uint16
	IME                    bool
	memory                 *memory.Memory
}

func NewCPU(m *memory.Memory) *CPU {
	fmt.Println("Initializing CPU")
	return &CPU{PC: 0x0100, memory: m}
}

func (cpu *CPU) AF() uint16 {
	return uint16(cpu.A)<<8 | uint16(cpu.F)
}

func (cpu *CPU) SetAF(value uint16) {
	cpu.A = uint8(value >> 8)
	cpu.F = uint8(value & 0xFF)
}

func (cpu *CPU) BC() uint16 {
	return uint16(cpu.B)<<8 | uint16(cpu.C)
}

func (cpu *CPU) SetBC(value uint16) {
	cpu.B = uint8(value >> 8)
	cpu.C = uint8(value & 0xFF)
}

func (cpu *CPU) DE() uint16 {
	return uint16(cpu.D)<<8 | uint16(cpu.E)
}

func (cpu *CPU) SetDE(value uint16) {
	cpu.D = uint8(value >> 8)
	cpu.E = uint8(value & 0xFF)
}

func (cpu *CPU) HL() uint16 {
	return uint16(cpu.B)<<8 | uint16(cpu.C)
}

func (cpu *CPU) SetHL(value uint16) {
	cpu.H = uint8(value >> 8)
	cpu.L = uint8(value & 0xFF)
}

func (cpu *CPU) SetCarryFlag(value byte) {
	if value == 1 {
		cpu.F |= 0x10
	} else {
		cpu.F &= 0xEF
	}
}

func (cpu *CPU) PushStack(value uint16) {
	cpu.SP--
	cpu.memory.StoreByte(cpu.SP, byte(value&0xFF))
	cpu.SP--
	cpu.memory.StoreByte(cpu.SP, byte((value>>8)*0xFF))
}

func (cpu *CPU) Cycle() bool {
	opcode := cpu.memory.GetByte(cpu.PC)
	handler := opcodeTable[opcode]

	if handler != nil {
		fmt.Printf("Executing opcode: 0x%02X\n", opcode)
		handler()
	} else {
		fmt.Printf("Unhandled opcode: 0x%02X\n", opcode)
		return true
	}
	return false
}

// Opcode Handling //
type opcodeFunc func()

var opcodeTable [256]opcodeFunc

func (cpu *CPU) InitOpcodeTable() {
	opcodeTable[0x00] = cpu.NOP
	opcodeTable[0x07] = cpu.RLCA
	opcodeTable[0x0F] = cpu.RRCA
	opcodeTable[0x31] = cpu.LD_SP_u16
	opcodeTable[0x3E] = cpu.LD_A_u8
	opcodeTable[0x55] = cpu.LD_D_L
	opcodeTable[0xC3] = cpu.JP_u16
	opcodeTable[0xE0] = cpu.LD_C_A
	opcodeTable[0xEA] = cpu.LD_u16_A
	opcodeTable[0xF3] = cpu.DI
	opcodeTable[0xFF] = cpu.RST_38H
}

func (cpu *CPU) NOP() {
	// 0x00: Do nothing, increment PC
	cpu.PC++
}

func (cpu *CPU) RLCA() {
	// 0x07: Rotate the value of register A left by one bit
	carry := (cpu.A >> 7) & 1
	cpu.A = (cpu.A << 1) | carry
	cpu.SetCarryFlag(carry)
	cpu.PC++
}

func (cpu *CPU) RRCA() {
	// 0x0F: Rotate the value of register A right one bit
	carry := cpu.F & 0x10
	carry = carry >> 4

	cpu.F &= 0xEF
	cpu.A = (cpu.A >> 1) | (carry << 7)
	cpu.SetCarryFlag(cpu.A & 0x01)
	cpu.PC++
}

func (cpu *CPU) LD_SP_u16() {
	// 0x31: Set SP to next two bytes in program
	cpu.PC++
	low := cpu.memory.GetByte(cpu.PC)
	cpu.PC++
	high := cpu.memory.GetByte(cpu.PC)
	address := (uint16(high)<<8 | uint16(low))
	cpu.SP = address
	cpu.PC++
}

func (cpu *CPU) LD_A_u8() {
	// 0x3E: Load next byte into register A
	cpu.PC++
	cpu.A = cpu.memory.GetByte(cpu.PC)
	cpu.PC++
}

func (cpu *CPU) LD_D_L() {
	// 0x55: Load the contents of register L into register D
	cpu.D = cpu.L
	cpu.PC++
}

func (cpu *CPU) JP_u16() {
	// 0xC3: Set PC to next two bytes in program
	cpu.PC++
	low := cpu.memory.GetByte(cpu.PC)
	cpu.PC++
	high := cpu.memory.GetByte(cpu.PC)
	address := (uint16(high)<<8 | uint16(low))
	cpu.PC = address
}

func (cpu *CPU) LD_C_A() {
	// 0xE0: Load value of register A into memory location of 0xFF00 + register C
	address := 0xFF00 + uint16(cpu.C)
	cpu.memory.StoreByte(address, cpu.A)
	cpu.PC++
}

func (cpu *CPU) LD_u16_A() {
	// 0xEA: Load register A to memory location of next two bytes
	cpu.PC++
	low := cpu.memory.GetByte(cpu.PC)
	cpu.PC++
	high := cpu.memory.GetByte(cpu.PC)
	location := (uint16(high)<<8 | uint16(low))
	cpu.memory.StoreByte(location, cpu.A)
	cpu.PC++
}

func (cpu *CPU) DI() {
	// 0xF3: Disable interrupts
	cpu.IME = false
	cpu.PC++
}

func (cpu *CPU) RST_38H() {
	// 0xFF: Save PC to stack and jump to address 0x0038
	cpu.PushStack(cpu.PC)
	cpu.PC = 0x0038
}
