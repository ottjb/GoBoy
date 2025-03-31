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
	return uint16(cpu.H)<<8 | uint16(cpu.L)
}

func (cpu *CPU) SetHL(value uint16) {
	cpu.H = uint8(value >> 8)
	cpu.L = uint8(value & 0xFF)
}

func (cpu *CPU) SetCarryFlag(value bool) {
	if value {
		cpu.F |= 0x10
	} else {
		cpu.F &= 0xEF
	}
}

func (cpu *CPU) SetHalfCarryFlag(value bool) {
	if value {
		cpu.F |= 0x20
	} else {
		cpu.F &= 0xDF
	}
}

func (cpu *CPU) SetZeroFlag(value bool) {
	if value {
		cpu.F |= 0x80
	} else {
		cpu.F &= 0x7F
	}
}

func (cpu *CPU) SetSubtractFlag(value bool) {
	if value {
		cpu.F |= 0x40
	} else {
		cpu.F &= 0xBF
	}
}

func (cpu *CPU) GetCarryFlag() bool {
	return cpu.F&0x10 != 0
}

func (cpu *CPU) GetHalfCarryFlag() bool {
	return cpu.F&0x20 != 0
}

func (cpu *CPU) GetZeroFlag() bool {
	return cpu.F&0x80 != 0
}

func (cpu *CPU) GetSubtractFlag() bool {
	return cpu.F&0x40 != 0
}

func (cpu *CPU) PushStack(value uint16) {
	cpu.SP--
	cpu.memory.Write(cpu.SP, byte(value&0xFF))
	cpu.SP--
	cpu.memory.Write(cpu.SP, byte((value>>8)&0xFF))
}

func (cpu *CPU) Cycle() bool {
	opcode := cpu.memory.Read(cpu.PC)
	handler := opcodeTable[opcode]

	if handler != nil {
		fmt.Printf("Executing opcode 0x%02X at PC 0x%04X\n", opcode, cpu.PC)
		handler()
	} else {
		fmt.Printf("Unhandled opcode 0x%02X at PC 0x%04X\n", opcode, cpu.PC)
		return true
	}
	return false
}

// Opcode Handling //
type opcodeFunc func()

var opcodeTable [256]opcodeFunc

func (cpu *CPU) InitOpcodeTable() {
	opcodeTable[0x00] = cpu.NOP
	opcodeTable[0x01] = cpu.LD_BC_u16
	opcodeTable[0x03] = cpu.INC_BC
	opcodeTable[0x07] = cpu.RLCA
	opcodeTable[0x0F] = cpu.RRCA
	opcodeTable[0x11] = cpu.LD_DE_u16
	opcodeTable[0x13] = cpu.INC_DE
	opcodeTable[0x1F] = cpu.RRA
	opcodeTable[0x20] = cpu.JR_NZ_r8
	opcodeTable[0x21] = cpu.LD_HL_u16
	opcodeTable[0x22] = cpu.LD_HLi_A
	opcodeTable[0x25] = cpu.DEC_H
	opcodeTable[0x26] = cpu.LD_H_u8
	opcodeTable[0x2C] = cpu.INC_L
	opcodeTable[0x2D] = cpu.DEC_L
	opcodeTable[0x2F] = cpu.CPL
	opcodeTable[0x30] = cpu.JR_NC_r8
	opcodeTable[0x31] = cpu.LD_SP_u16
	opcodeTable[0x3C] = cpu.INC_A
	opcodeTable[0x3E] = cpu.LD_A_u8
	opcodeTable[0x40] = cpu.LD_B_B
	opcodeTable[0x46] = cpu.LD_B_HL
	opcodeTable[0x47] = cpu.LD_B_A
	opcodeTable[0x4E] = cpu.LD_C_HL
	opcodeTable[0x4F] = cpu.LD_C_A
	opcodeTable[0x55] = cpu.LD_D_L
	opcodeTable[0x56] = cpu.LD_D_HL
	opcodeTable[0x57] = cpu.LD_D_A
	opcodeTable[0x5F] = cpu.LD_E_A
	opcodeTable[0x67] = cpu.LD_H_A
	opcodeTable[0x6F] = cpu.LD_L_A
	opcodeTable[0x70] = cpu.LD_HL_B
	opcodeTable[0x71] = cpu.LD_HL_C
	opcodeTable[0x72] = cpu.LD_HL_D
	opcodeTable[0x78] = cpu.LD_A_B
	opcodeTable[0x79] = cpu.LD_A_C
	opcodeTable[0x7A] = cpu.LD_A_D
	opcodeTable[0x7B] = cpu.LD_A_E
	opcodeTable[0x80] = cpu.ADD_A_B
	opcodeTable[0x81] = cpu.ADD_A_C
	opcodeTable[0x82] = cpu.ADD_A_D
	opcodeTable[0x83] = cpu.ADD_A_E
	opcodeTable[0xAE] = cpu.XOR_A_HL
	opcodeTable[0xB9] = cpu.CP_A_C
	opcodeTable[0xC1] = cpu.POP_BC
	opcodeTable[0xC3] = cpu.JP_u16
	opcodeTable[0xC5] = cpu.PUSH_BC
	opcodeTable[0xC9] = cpu.RET
	opcodeTable[0xCB] = cpu.ExecuteCBOpcode
	opcodeTable[0xCD] = cpu.CALL_u16
	opcodeTable[0xD1] = cpu.POP_DE
	opcodeTable[0xD5] = cpu.PUSH_DE
	opcodeTable[0xE0] = cpu.LD_u8C_A
	opcodeTable[0xE1] = cpu.POP_HL
	opcodeTable[0xE5] = cpu.PUSH_HL
	opcodeTable[0xEA] = cpu.LD_u16_A
	opcodeTable[0xEE] = cpu.XOR_A_u8
	opcodeTable[0xF0] = cpu.LD_A_u8C
	opcodeTable[0xF1] = cpu.POP_AF
	opcodeTable[0xF3] = cpu.DI
	opcodeTable[0xF5] = cpu.PUSH_AF
	opcodeTable[0xF9] = cpu.LD_SP_HL
	opcodeTable[0xFA] = cpu.LD_A_u16
	opcodeTable[0xFF] = cpu.RST_38H
}

var opcodeCBTable [256]opcodeFunc

func (cpu *CPU) InitOpcodeCBTable() {
	opcodeCBTable[0x19] = cpu.RR_C
	opcodeCBTable[0x1A] = cpu.RR_D
	opcodeCBTable[0x38] = cpu.SRL_B
}

func (cpu *CPU) NOP() {
	// 0x00: Do nothing, increment PC
	cpu.PC++
}

func (cpu *CPU) LD_BC_u16() {
	// 0x01: Load next 2 bytes to BC
	cpu.PC++
	low := cpu.memory.Read(cpu.PC)
	cpu.PC++
	high := cpu.memory.Read(cpu.PC)
	value := uint16(high)<<8 | uint16(low)

	cpu.B = uint8(value >> 8)
	cpu.C = uint8(value & 0xFF)
	cpu.PC++
}

func (cpu *CPU) INC_BC() {
	// 0x03: Increment BC
	if cpu.C == 0xFF {
		cpu.C = 0x00
		cpu.B++
	} else {
		cpu.C++
	}
	cpu.PC++
}

func (cpu *CPU) RLCA() {
	// 0x07: Rotate the value of register A left by one bit
	carry := (cpu.A >> 7) & 1
	cpu.A = (cpu.A << 1) | carry
	cpu.SetCarryFlag(carry == 1)
	cpu.PC++
}

func (cpu *CPU) RRCA() {
	// 0x0F: Rotate the value of register A right one bit
	carry := cpu.A & 0x01
	cpu.A = (cpu.A >> 1) | (carry << 7)

	cpu.SetCarryFlag(carry == 1)
	cpu.SetZeroFlag(false)
	cpu.SetHalfCarryFlag(false)
	cpu.SetSubtractFlag(false)

	cpu.PC++
}

func (cpu *CPU) LD_DE_u16() {
	// 0x11: Load next 2 bytes to DE
	cpu.PC++
	low := cpu.memory.Read(cpu.PC)
	cpu.PC++
	high := cpu.memory.Read(cpu.PC)
	value := uint16(high)<<8 | uint16(low)

	cpu.D = uint8(value >> 8)
	cpu.E = uint8(value & 0xFF)
	cpu.PC++
}

func (cpu *CPU) INC_DE() {
	// 0x13: Increment DE
	if cpu.E == 0xFF {
		cpu.E = 0x00
		cpu.D++
	} else {
		cpu.E++
	}
	cpu.PC++
}

func (cpu *CPU) RRA() {
	// 0x1F: Rotate register A right
	carry := cpu.A & 0x01
	cpu.A = (cpu.A >> 1) | (boolToUint8(cpu.GetCarryFlag()) << 7)

	cpu.SetCarryFlag(carry == 1)
	cpu.SetZeroFlag(cpu.A == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag(false)

	cpu.PC++
}

func (cpu *CPU) JR_NZ_r8() {
	// 0x20: If Z flag == 0, jump to PC + r8
	offset := int8(cpu.memory.Read(cpu.PC + 1))

	if !cpu.GetZeroFlag() {
		cpu.PC += uint16(offset) + 2
	} else {
		cpu.PC += 2
	}
}

func (cpu *CPU) LD_HL_u16() {
	// 0x21: Load next 2 bytes to HL
	cpu.PC++
	low := cpu.memory.Read(cpu.PC)
	cpu.PC++
	high := cpu.memory.Read(cpu.PC)
	value := uint16(high)<<8 | uint16(low)

	cpu.H = uint8(value >> 8)
	cpu.L = uint8(value & 0xFF)
	cpu.PC++
}

func (cpu *CPU) LD_HLi_A() {
	// 0x22: Store the value in register A at memory location HL, then increment HL
	cpu.memory.Write(cpu.HL(), cpu.A)
	cpu.SetHL(cpu.HL() + 1)
	cpu.PC++
}

func (cpu *CPU) DEC_H() {
	// 0x25: Decrement register H by 1
	cpu.H--

	cpu.SetSubtractFlag(true)
	cpu.SetZeroFlag(cpu.H == 0)
	cpu.SetHalfCarryFlag((cpu.H & 0x0F) == 0x0F)

	cpu.PC++
}

func (cpu *CPU) LD_H_u8() {
	// 0x26: Load next byte into register H
	cpu.PC++
	cpu.H = cpu.memory.Read(cpu.PC)
	cpu.PC++
}

func (cpu *CPU) INC_L() {
	// 0x2C: Increment register L
	oldL := cpu.L
	cpu.L++

	cpu.SetZeroFlag(cpu.L == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag((oldL&0x0F)+1 > 0x0F)

	cpu.PC++
}

func (cpu *CPU) DEC_L() {
	// 0x2D: Decrement register L by 1
	cpu.L--

	cpu.SetSubtractFlag(true)
	cpu.SetZeroFlag(cpu.L == 0)
	cpu.SetHalfCarryFlag((cpu.L & 0x0F) == 0x0F)

	cpu.PC++
}

func (cpu *CPU) CPL() {
	// 0x2F: Complement A register
	cpu.A = ^cpu.A

	cpu.SetSubtractFlag(true)
	cpu.SetHalfCarryFlag(false)

	cpu.PC++
}

func (cpu *CPU) JR_NC_r8() {
	// 0x30: If C flag == 0, jump to PC + r8
	offset := int8(cpu.memory.Read(cpu.PC + 1))

	if !cpu.GetCarryFlag() {
		cpu.PC += uint16(offset) + 2
	} else {
		cpu.PC += 2
	}
}

func (cpu *CPU) LD_SP_u16() {
	// 0x31: Set SP to next two bytes in program
	cpu.PC++
	low := cpu.memory.Read(cpu.PC)
	cpu.PC++
	high := cpu.memory.Read(cpu.PC)
	address := (uint16(high)<<8 | uint16(low))
	cpu.SP = address
	cpu.PC++
}

func (cpu *CPU) INC_A() {
	// 0x3C: Increment register A
	oldA := cpu.A
	cpu.A++

	cpu.SetZeroFlag(cpu.A == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag((oldA&0x0F)+1 > 0x0F)

	cpu.PC++
}

func (cpu *CPU) LD_A_u8() {
	// 0x3E: Load next byte into register A
	cpu.PC++
	cpu.A = cpu.memory.Read(cpu.PC)
	cpu.PC++
}

func (cpu *CPU) LD_B_B() {
	// 0x40: Load the contents of register B into register B (does nothing)
	cpu.PC++
}

func (cpu *CPU) LD_B_HL() {
	// 0x46: Load the value at memory location HL into register B
	cpu.B = cpu.memory.Read(cpu.HL())

	cpu.PC++
}

func (cpu *CPU) LD_B_A() {
	// 0x47: Load the contents of register A into register B
	cpu.B = cpu.A
	cpu.PC++
}

func (cpu *CPU) LD_C_HL() {
	// 0x4E: Load the value at memory location HL into register C
	cpu.C = cpu.memory.Read(cpu.HL())

	cpu.PC++
}

func (cpu *CPU) LD_C_A() {
	// 0x4F: Load the contents of register A into register C
	cpu.C = cpu.A
	cpu.PC++
}

func (cpu *CPU) LD_D_L() {
	// 0x55: Load the contents of register L into register D
	cpu.D = cpu.L
	cpu.PC++
}

func (cpu *CPU) LD_D_HL() {
	// 0x56: Load the value at memory location HL into register D
	cpu.D = cpu.memory.Read(cpu.HL())

	cpu.PC++
}

func (cpu *CPU) LD_D_A() {
	// 0x57: Load the contents of register A into register D
	cpu.D = cpu.A
	cpu.PC++
}

func (cpu *CPU) LD_E_A() {
	// 0x5F: Load the contents of register A into register E
	cpu.E = cpu.A
	cpu.PC++
}

func (cpu *CPU) LD_H_A() {
	// 0x67: Load the contents of register A into register H
	cpu.H = cpu.A
	cpu.PC++
}

func (cpu *CPU) LD_L_A() {
	// 0x6F: Load the contents of register A into register L
	cpu.L = cpu.A
	cpu.PC++
}

func (cpu *CPU) LD_HL_B() {
	// 0x70: Store the value in register B at memory location HL
	cpu.memory.Write(cpu.HL(), cpu.B)
	cpu.PC++
}

func (cpu *CPU) LD_HL_C() {
	// 0x71: Store the value in register C at memory location HL
	cpu.memory.Write(cpu.HL(), cpu.C)
	cpu.PC++
}

func (cpu *CPU) LD_HL_D() {
	// 0x72: Store the value in register D at memory location HL
	cpu.memory.Write(cpu.HL(), cpu.D)
	cpu.PC++
}

func (cpu *CPU) LD_A_B() {
	// 0x78: Load the contents of register B into register A
	cpu.A = cpu.B
	cpu.PC++
}

func (cpu *CPU) LD_A_C() {
	// 0x79: Load the contents of register C into register A
	cpu.A = cpu.C
	cpu.PC++
}

func (cpu *CPU) LD_A_D() {
	// 0x7A: Load the contents of register D into register A
	cpu.A = cpu.D
	cpu.PC++
}

func (cpu *CPU) LD_A_E() {
	// 0x7B: Load the contents of register E into register A
	cpu.A = cpu.E
	cpu.PC++
}

func (cpu *CPU) ADD_A_B() {
	// 0x80: A = A + B
	result := uint16(cpu.A) + uint16(cpu.B)

	cpu.SetZeroFlag(uint8(result) == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag((cpu.A&0x0F)+(cpu.B&0x0F) > 0x0F)
	cpu.SetCarryFlag(result > 0xFF)

	cpu.A = uint8(result)
	cpu.PC++
}

func (cpu *CPU) ADD_A_C() {
	// 0x81: A = A + C
	result := uint16(cpu.A) + uint16(cpu.C)

	cpu.SetZeroFlag(uint8(result) == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag((cpu.A&0x0F)+(cpu.C&0x0F) > 0x0F)
	cpu.SetCarryFlag(result > 0xFF)

	cpu.A = uint8(result)
	cpu.PC++
}

func (cpu *CPU) ADD_A_D() {
	// 0x82: A = A + D
	result := uint16(cpu.A) + uint16(cpu.D)

	cpu.SetZeroFlag(uint8(result) == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag((cpu.A&0x0F)+(cpu.D&0x0F) > 0x0F)
	cpu.SetCarryFlag(result > 0xFF)

	cpu.A = uint8(result)
	cpu.PC++
}

func (cpu *CPU) ADD_A_E() {
	// 0x83: A = A + E
	result := uint16(cpu.A) + uint16(cpu.E)

	cpu.SetZeroFlag(uint8(result) == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag((cpu.A&0x0F)+(cpu.E&0x0F) > 0x0F)
	cpu.SetCarryFlag(result > 0xFF)

	cpu.A = uint8(result)
	cpu.PC++
}

func (cpu *CPU) XOR_A_HL() {
	// 0xAE: Perform XOR operation with value from memory location stored in HL on register A
	cpu.A ^= cpu.memory.Read(cpu.HL())

	cpu.SetZeroFlag(cpu.A == 0)
	cpu.SetCarryFlag(false)
	cpu.SetHalfCarryFlag(false)
	cpu.SetSubtractFlag(false)

	cpu.PC++
}

func (cpu *CPU) CP_A_C() {
	// 0xB9: Compare A and C, does not modify registers, only flags
	result := int16(cpu.A) - int16(cpu.C)

	cpu.SetZeroFlag(uint8(result) == 0)
	cpu.SetSubtractFlag(true)
	cpu.SetHalfCarryFlag((cpu.A & 0x0F) < (cpu.C & 0x0F))
	cpu.SetCarryFlag(cpu.A < cpu.C)

	cpu.PC++
}

func (cpu *CPU) POP_BC() {
	// 0xC1: Pop two bytes from the stack into register BC
	low := cpu.memory.Read(cpu.SP)
	cpu.SP++
	high := cpu.memory.Read(cpu.SP)
	cpu.SP++

	cpu.B = high
	cpu.C = low
	cpu.PC++
}

func (cpu *CPU) JP_u16() {
	// 0xC3: Set PC to next two bytes in program
	cpu.PC++
	low := cpu.memory.Read(cpu.PC)
	cpu.PC++
	high := cpu.memory.Read(cpu.PC)
	address := (uint16(high)<<8 | uint16(low))
	cpu.PC = address
}

func (cpu *CPU) PUSH_BC() {
	// 0xC5: Push BC to stack
	cpu.PushStack(cpu.BC())
	cpu.PC++
}

func (cpu *CPU) RET() {
	// 0xC9: Return from subroutine
	low := cpu.memory.Read(cpu.SP)
	cpu.SP++
	high := cpu.memory.Read(cpu.SP)
	cpu.SP++

	cpu.PC = uint16(high)<<8 | uint16(low)
}

func (cpu *CPU) ExecuteCBOpcode() {
	// 0xCB: Prefixed opcodes
	cpu.PC++
	opcode := cpu.memory.Read(cpu.PC)
	handler := opcodeCBTable[opcode]

	if handler != nil {
		// Execute the handler corresponding to the second byte
		fmt.Printf("Executing CB-prefixed opcode: 0x%02X at PC 0x%04X\n", opcode, cpu.PC)
		handler()
	} else {
		// Handle unrecognized opcode
		fmt.Printf("Unhandled CB-prefixed opcode: 0x%02X at PC 0x%04X\n", opcode, cpu.PC)
		fmt.Printf("\n\n\n")
	}

	cpu.PC++
}

func (cpu *CPU) CALL_u16() {
	// 0xCD: Push PC to stack, load next two bytes into PC
	cpu.PC++
	low := cpu.memory.Read(cpu.PC)
	cpu.PC++
	high := cpu.memory.Read(cpu.PC)
	address := uint16(high)<<8 | uint16(low)

	cpu.PC++
	cpu.PushStack(cpu.PC)

	cpu.PC = address
}

func (cpu *CPU) POP_DE() {
	// 0xD1: Pop two bytes from the stack into register DE
	low := cpu.memory.Read(cpu.SP)
	cpu.SP++
	high := cpu.memory.Read(cpu.SP)
	cpu.SP++

	cpu.D = high
	cpu.E = low
	cpu.PC++
}

func (cpu *CPU) PUSH_DE() {
	// 0xD5: Push DE to stack
	cpu.PushStack(cpu.DE())
	cpu.PC++
}

func (cpu *CPU) LD_u8C_A() {
	// 0xE0: Load value of register A into memory location of 0xFF00 + register C
	address := 0xFF00 + uint16(cpu.C)
	cpu.memory.Write(address, cpu.A)
	cpu.PC++
}

func (cpu *CPU) POP_HL() {
	// 0xE1: Pop two bytes from the stack into register HL
	low := cpu.memory.Read(cpu.SP)
	cpu.SP++
	high := cpu.memory.Read(cpu.SP)
	cpu.SP++

	cpu.H = high
	cpu.L = low
	cpu.PC++
}

func (cpu *CPU) PUSH_HL() {
	// 0xE5: Push HL to stack
	cpu.PushStack(cpu.HL())
	cpu.PC++
}

func (cpu *CPU) LD_u16_A() {
	// 0xEA: Load register A to memory location of next two bytes
	cpu.PC++
	low := cpu.memory.Read(cpu.PC)
	cpu.PC++
	high := cpu.memory.Read(cpu.PC)
	location := (uint16(high)<<8 | uint16(low))
	cpu.memory.Write(location, cpu.A)
	cpu.PC++
}

func (cpu *CPU) XOR_A_u8() {
	// 0xEE: Perform XOR operation with value of the next byte on register A
	cpu.PC++
	cpu.A ^= cpu.memory.Read(cpu.PC)

	cpu.SetZeroFlag(cpu.A == 0)
	cpu.SetCarryFlag(false)
	cpu.SetHalfCarryFlag(false)
	cpu.SetSubtractFlag(false)

	cpu.PC++
}

func (cpu *CPU) LD_A_u8C() {
	// 0xF0: Load value at 0xFF00 + C into register A
	address := 0xFF00 + uint16(cpu.C)
	cpu.A = cpu.memory.Read(address)
	cpu.PC++
}

func (cpu *CPU) POP_AF() {
	// 0xF1: Pop two bytes from the stack into register AF
	low := cpu.memory.Read(cpu.SP)
	cpu.SP++
	high := cpu.memory.Read(cpu.SP)
	cpu.SP++

	cpu.A = high
	cpu.F = low
	cpu.PC++
}

func (cpu *CPU) DI() {
	// 0xF3: Disable interrupts
	cpu.IME = false
	cpu.PC++
}

func (cpu *CPU) PUSH_AF() {
	// 0xF5: Push AF to stack
	cpu.PushStack(cpu.AF())
	cpu.PC++
}

func (cpu *CPU) LD_SP_HL() {
	// 0xF9: Set SP to the value in HL
	cpu.SP = cpu.HL()
	cpu.PC++
}

func (cpu *CPU) LD_A_u16() {
	// 0xFA: Loads the value from memory address of next two bytes to register A
	cpu.PC++
	low := cpu.memory.Read(cpu.PC)
	cpu.PC++
	high := cpu.memory.Read(cpu.PC)
	address := uint16(high)<<8 | uint16(low)

	cpu.A = cpu.memory.Read(address)
	cpu.PC++
}

func (cpu *CPU) RST_38H() {
	// 0xFF: Save PC to stack and jump to address 0x0038
	cpu.PushStack(cpu.PC)
	cpu.PC = 0x0038
}

// CB-Prefixed Opcodes //
func (cpu *CPU) RR_C() {
	// 0x19: Rotate Right operation on register C
	carry := cpu.C & 0x01
	cpu.C = (cpu.C >> 1) | (boolToUint8(cpu.GetCarryFlag()) << 7)

	cpu.SetCarryFlag(carry == 1)
	cpu.SetZeroFlag(cpu.C == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag(false)
}

func (cpu *CPU) RR_D() {
	// 0x1A: Rotate Right operation on register D
	carry := cpu.D & 0x01
	cpu.D = (cpu.D >> 1) | (boolToUint8(cpu.GetCarryFlag()) << 7)

	cpu.SetCarryFlag(carry == 1)
	cpu.SetZeroFlag(cpu.D == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag(false)
}

func (cpu *CPU) SRL_B() {
	// 0x38: Shift bits in register B one bit to the right
	carry := cpu.B & 0x01
	cpu.B >>= 1

	cpu.SetCarryFlag(carry == 1)
	cpu.SetZeroFlag(cpu.B == 0)
	cpu.SetSubtractFlag(false)
	cpu.SetHalfCarryFlag(false)
}

// Helper functions //
func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
