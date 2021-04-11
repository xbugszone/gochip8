package chip8

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

var drawFlag bool
var opcode uint16
var memory [4096]uint8
var V [16]uint8
var I uint16
var pc uint16

var gfx [32][64]uint8

var delay_timer uint8
var sound_timer uint8

var stack [16]uint16
var sp uint16

var key [16]uint8

var chip8_fontset = [80]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}
func Initialize() {

	// Initialize registers and memory once
	pc = 0x200
	opcode = 0
	I = 0
	sp = 0

	// Clear display
	// Clear stack
	// Clear registers V0-VF
	// Clear memory

	//Load fontset
	for i := 1; i < 80; i++ {
		memory[i] = chip8_fontset[i]
	}


	// Reset timers


}

func LoadGame(name string) {

	file, err := os.Open(name)

	if err != nil {

	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {

	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_,err = bufr.Read(bytes)

	for  i := 0; i < len(bytes); i++ {
		memory[i + 512] = bytes[i]
	}
}

func EmulateCycle() {
	opcode = (uint16(memory[pc]) << 8) | uint16(memory[pc+1])

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x000F {
		case 0x0000: // 0x00E0 Clears screen
			for i := 0; i < len(gfx); i++ {
				for j := 0; j < len(gfx[i]); j++ {
					gfx[i][j] = 0x0
				}
			}
			drawFlag = true
			pc = pc + 2
		case 0x000E: // 0x00EE Returns from a subroutine
			sp = sp - 1
			pc = stack[sp]
			pc = pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", opcode)
		}
	case 0x1000: // 0x1NNN Jump to address NNN
		pc = opcode & 0x0FFF
	case 0x2000: // 0x2NNN Calls subroutine at NNN
		stack[sp] = pc // store current program counter
		sp = sp + 1      // increment stack pointer
		pc = opcode & 0x0FFF // jump to address NNN
	case 0x3000: // 0x3XNN Skips the next instruction if VX equals NN
		if uint16(V[(opcode&0x0F00)>>8]) == opcode&0x00FF {
			pc = pc + 4
		} else {
			pc = pc + 2
		}
	case 0x4000: // 0x4XNN Skips the next instruction if VX doesn't equal NN
		if uint16(V[(opcode&0x0F00)>>8]) != opcode&0x00FF {
			pc = pc + 4
		} else {
			pc = pc + 2
		}
	case 0x5000: // 0x5XY0 Skips the next instruction if VX equals VY
		if V[(opcode&0x0F00)>>8] == V[(opcode&0x00F0)>>4] {
			pc = pc + 4
		} else {
			pc = pc + 2
		}
	case 0x6000: // 0x6XNN Sets VX to NN
		V[(opcode&0x0F00)>>8] = uint8(opcode & 0x00FF)
		pc = pc + 2
	case 0x7000: // 0x7XNN Adds NN to VX
		V[(opcode&0x0F00)>>8] = V[(opcode&0x0F00)>>8] + uint8(opcode&0x00FF)
		pc = pc + 2
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000: // 0x8XY0 Sets VX to the value of VY
			V[(opcode&0x0F00)>>8] = V[(opcode&0x00F0)>>4]
			pc = pc + 2
		case 0x0001: // 0x8XY1 Sets VX to VX or VY
			V[(opcode&0x0F00)>>8] = V[(opcode&0x0F00)>>8] | V[(opcode&0x00F0)>>4]
			pc = pc + 2
		case 0x0002: // 0x8XY2 Sets VX to VX and VY
			V[(opcode&0x0F00)>>8] = V[(opcode&0x0F00)>>8] & V[(opcode&0x00F0)>>4]
			pc = pc + 2
		case 0x0003: // 0x8XY3 Sets VX to VX xor VY
			V[(opcode&0x0F00)>>8] = V[(opcode&0x0F00)>>8] ^ V[(opcode&0x00F0)>>4]
			pc = pc + 2
		case 0x0004: // 0x8XY4 Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't
			if V[(opcode&0x00F0)>>4] > 0xFF-V[(opcode&0x0F00)>>8] {
				V[0xF] = 1
			} else {
				V[0xF] = 0
			}
			V[(opcode&0x0F00)>>8] = V[(opcode&0x0F00)>>8] + V[(opcode&0x00F0)>>4]
			pc = pc + 2
		case 0x0005: // 0x8XY5 VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if V[(opcode&0x00F0)>>4] > V[(opcode&0x0F00)>>8] {
				V[0xF] = 0
			} else {
				V[0xF] = 1
			}
			V[(opcode&0x0F00)>>8] = V[(opcode&0x0F00)>>8] - V[(opcode&0x00F0)>>4]
			pc = pc + 2
		case 0x0006: // 0x8XY6 Shifts VY right by one and stores the result to VX (VY remains unchanged). VF is set to the value of the least significant bit of VY before the shift
			V[0xF] = V[(opcode&0x0F00)>>8] & 0x1
			V[(opcode&0x0F00)>>8] = V[(opcode&0x0F00)>>8] >> 1
			pc = pc + 2
		case 0x0007: // 0x8XY7 Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if V[(opcode&0x0F00)>>8] > V[(opcode&0x00F0)>>4] {
				V[0xF] = 0
			} else {
				V[0xF] = 1
			}
			V[(opcode&0x0F00)>>8] = V[(opcode&0x00F0)>>4] - V[(opcode&0x0F00)>>8]
			pc = pc + 2
		case 0x000E: // 0x8XYE Shifts VY left by one and copies the result to VX. VF is set to the value of the most significant bit of VY before the shift
			V[0xF] = V[(opcode&0x0F00)>>8] >> 7
			V[(opcode&0x0F00)>>8] = V[(opcode&0x0F00)>>8] << 1
			pc = pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", opcode)
		}
	case 0x9000: // 9XY0 Skips the next instruction if VX doesn't equal VY
		if V[(opcode&0x0F00)>>8] != V[(opcode&0x00F0)>>4] {
			pc = pc + 4
		} else {
			pc = pc + 2
		}
	case 0xA000: // 0xANNN Sets I to the address NNN
		I = opcode & 0x0FFF
		pc = pc + 2
	case 0xB000: // 0xBNNN Jumps to the address NNN plus V0
		pc = (opcode & 0x0FFF) + uint16(V[0x0])
	case 0xC000: // 0xCXNN Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
		V[(opcode&0x0F00)>>8] = uint8(rand.Intn(256)) & uint8(opcode&0x00FF)
		pc = pc + 2
	case 0xD000: // 0xDXYN Draws a sprite at coordinate (VX, VY)
		x := V[(opcode&0x0F00)>>8]
		y := V[(opcode&0x00F0)>>4]
		h := opcode & 0x000F
		V[0xF] = 0
		var j uint16 = 0
		var i uint16 = 0
		for j = 0; j < h; j++ {
			pixel := memory[I+j]
			for i = 0; i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					if gfx[(y + uint8(j))][x+uint8(i)] == 1 {
						V[0xF] = 1
					}
					gfx[(y + uint8(j))][x+uint8(i)] ^= 1
				}
			}
		}
		drawFlag = true
		pc = pc + 2
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E: // 0xEX9E Skips the next instruction if the key stored in VX is pressed
			if key[V[(opcode&0x0F00)>>8]] == 1 {
				pc = pc + 4
			} else {
				pc = pc + 2
			}
		case 0x00A1: // 0xEXA1 Skips the next instruction if the key stored in VX isn't pressed
			if key[V[(opcode&0x0F00)>>8]] == 0 {
				pc = pc + 4
			} else {
				pc = pc + 2
			}
		default:
			fmt.Printf("Invalid opcode %X\n", opcode)
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007: // 0xFX07 Sets VX to the value of the delay timer
			V[(opcode&0x0F00)>>8] = delay_timer
			pc = pc + 2
		case 0x000A: // 0xFX0A A key press is awaited, and then stored in VX
			pressed := false
			for i := 0; i < len(key); i++ {
				if key[i] != 0 {
					V[(opcode&0x0F00)>>8] = uint8(i)
					pressed = true
				}
			}
			if !pressed {
				return
			}
			pc = pc + 2
		case 0x0015: // 0xFX15 Sets the delay timer to VX
			delay_timer = V[(opcode&0x0F00)>>8]
			pc = pc + 2
		case 0x0018: // 0xFX18 Sets the sound timer to VX
			sound_timer = V[(opcode&0x0F00)>>8]
			pc = pc + 2
		case 0x001E: // 0xFX1E Adds VX to I
			if I+uint16(V[(opcode&0x0F00)>>8]) > 0xFFF {
				V[0xF] = 1
			} else {
				V[0xF] = 0
			}
			I = I + uint16(V[(opcode&0x0F00)>>8])
			pc = pc + 2
		case 0x0029: // 0xFX29 Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
			I = uint16(V[(opcode&0x0F00)>>8]) * 0x5
			pc = pc + 2
		case 0x0033: // 0xFX33 Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2
			memory[I] = V[(opcode&0x0F00)>>8] / 100
			memory[I+1] = (V[(opcode&0x0F00)>>8] / 10) % 10
			memory[I+2] = (V[(opcode&0x0F00)>>8] % 100) / 10
			pc = pc + 2
		case 0x0055: // 0xFX55 Stores V0 to VX (including VX) in memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((opcode&0x0F00)>>8)+1; i++ {
				memory[uint16(i)+I] = V[i]
			}
			I = ((opcode & 0x0F00) >> 8) + 1
			pc = pc + 2
		case 0x0065: // 0xFX65 Fills V0 to VX (including VX) with values from memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((opcode&0x0F00)>>8)+1; i++ {
				V[i] = memory[I+uint16(i)]
			}
			I = ((opcode & 0x0F00) >> 8) + 1
			pc = pc + 2
		default:
			fmt.Printf("Invalid opcode %X\n", opcode)
		}
	default:
		fmt.Printf("Invalid opcode %X\n", opcode)
	}

	if delay_timer > 0 {
		delay_timer = delay_timer - 1
	}
	if sound_timer > 0 {
		if sound_timer == 1 {
			//beeper()
		}
		sound_timer = sound_timer - 1
	}
}

func Key(num uint8, down bool) {
	if down {
		key[num] = 1
	} else {
		key[num] = 0
	}
}
func GetDrawFlag() bool {
	return drawFlag

}


func GetGFX() [32][64]uint8 {
	return gfx
}