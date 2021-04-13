package chip8

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

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

type Chip8 struct {
	drawFlag bool
	opcode   uint16
	memory   [4096]uint8
	V        [16]uint8
	I        uint16
	pc       uint16

	gfx   [32][64]uint8
	stack [16]uint16
	sp    uint16

	key [16]uint8

	delay_timer uint8
	sound_timer uint8
}

func Init() Chip8 {
	instance := Chip8{
		pc:     0x200,
		opcode: 0,
		I:      0,
		sp:     0,
	}

	//Load fontset
	for i := 1; i < len(chip8_fontset); i++ {
		instance.memory[i] = chip8_fontset[i]
	}

	return instance
}

func (chip8 *Chip8) LoadGame(name string) {

	file, err := os.Open(name)

	if err != nil {

	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {

	}

	var size = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	for i := 0; i < len(bytes); i++ {
		chip8.memory[i+512] = bytes[i]
	}
}

func (chip8 *Chip8) EmulateCycle() {
	chip8.opcode = (uint16(chip8.memory[chip8.pc]) << 8) | uint16(chip8.memory[chip8.pc+1])

	switch chip8.opcode & 0xF000 {
	case 0x0000:
		switch chip8.opcode & 0x000F {
		case 0x0000: // 0x00E0 Clears screen
			for i := 0; i < len(chip8.gfx); i++ {
				for j := 0; j < len(chip8.gfx[i]); j++ {
					chip8.gfx[i][j] = 0x0
				}
			}
			chip8.drawFlag = true
			chip8.pc = chip8.pc + 2
		case 0x000E: // 0x00EE Returns from a subroutine
			chip8.sp = chip8.sp - 1
			chip8.pc = chip8.stack[chip8.sp]
			chip8.pc = chip8.pc + 2
		default:
			fmt.Printf("Invalid chip8.opcode %X\n", chip8.opcode)
		}
	case 0x1000: // 0x1NNN Jump to address NNN
		chip8.pc = chip8.opcode & 0x0FFF
	case 0x2000: // 0x2NNN Calls subroutine at NNN
		chip8.stack[chip8.sp] = chip8.pc // store current program counter
		chip8.sp = chip8.sp + 1          // increment chip8.stack  pointer
		chip8.pc = chip8.opcode & 0x0FFF // jump to address NNN
	case 0x3000: // 0x3XNN Skips the next instruction if VX equals NN
		if uint16(chip8.V[(chip8.opcode&0x0F00)>>8]) == chip8.opcode&0x00FF {
			chip8.pc = chip8.pc + 4
		} else {
			chip8.pc = chip8.pc + 2
		}
	case 0x4000: // 0x4XNN Skips the next instruction if VX doesn't equal NN
		if uint16(chip8.V[(chip8.opcode&0x0F00)>>8]) != chip8.opcode&0x00FF {
			chip8.pc = chip8.pc + 4
		} else {
			chip8.pc = chip8.pc + 2
		}
	case 0x5000: // 0x5XY0 Skips the next instruction if VX equals VY
		if chip8.V[(chip8.opcode&0x0F00)>>8] == chip8.V[(chip8.opcode&0x00F0)>>4] {
			chip8.pc = chip8.pc + 4
		} else {
			chip8.pc = chip8.pc + 2
		}
	case 0x6000: // 0x6XNN Sets VX to NN
		chip8.V[(chip8.opcode&0x0F00)>>8] = uint8(chip8.opcode & 0x00FF)
		chip8.pc = chip8.pc + 2
	case 0x7000: // 0x7XNN Adds NN to VX
		chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x0F00)>>8] + uint8(chip8.opcode&0x00FF)
		chip8.pc = chip8.pc + 2
	case 0x8000:
		switch chip8.opcode & 0x000F {
		case 0x0000: // 0x8XY0 Sets VX to the value of VY
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x00F0)>>4]
			chip8.pc = chip8.pc + 2
		case 0x0001: // 0x8XY1 Sets VX to VX or VY
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x0F00)>>8] | chip8.V[(chip8.opcode&0x00F0)>>4]
			chip8.pc = chip8.pc + 2
		case 0x0002: // 0x8XY2 Sets VX to VX and VY
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x0F00)>>8] & chip8.V[(chip8.opcode&0x00F0)>>4]
			chip8.pc = chip8.pc + 2
		case 0x0003: // 0x8XY3 Sets VX to VX xor VY
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x0F00)>>8] ^ chip8.V[(chip8.opcode&0x00F0)>>4]
			chip8.pc = chip8.pc + 2
		case 0x0004: // 0x8XY4 Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't
			if chip8.V[(chip8.opcode&0x00F0)>>4] > 0xFF-chip8.V[(chip8.opcode&0x0F00)>>8] {
				chip8.V[0xF] = 1
			} else {
				chip8.V[0xF] = 0
			}
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x0F00)>>8] + chip8.V[(chip8.opcode&0x00F0)>>4]
			chip8.pc = chip8.pc + 2
		case 0x0005: // 0x8XY5 VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if chip8.V[(chip8.opcode&0x00F0)>>4] > chip8.V[(chip8.opcode&0x0F00)>>8] {
				chip8.V[0xF] = 0
			} else {
				chip8.V[0xF] = 1
			}
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x0F00)>>8] - chip8.V[(chip8.opcode&0x00F0)>>4]
			chip8.pc = chip8.pc + 2
		case 0x0006: // 0x8XY6 Shifts VY right by one and stores the result to VX (VY remains unchanged). VF is set to the value of the least significant bit of VY before the shift
			chip8.V[0xF] = chip8.V[(chip8.opcode&0x0F00)>>8] & 0x1
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x0F00)>>8] >> 1
			chip8.pc = chip8.pc + 2
		case 0x0007: // 0x8XY7 Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if chip8.V[(chip8.opcode&0x0F00)>>8] > chip8.V[(chip8.opcode&0x00F0)>>4] {
				chip8.V[0xF] = 0
			} else {
				chip8.V[0xF] = 1
			}
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x00F0)>>4] - chip8.V[(chip8.opcode&0x0F00)>>8]
			chip8.pc = chip8.pc + 2
		case 0x000E: // 0x8XYE Shifts VY left by one and copies the result to VX. VF is set to the value of the most significant bit of VY before the shift
			chip8.V[0xF] = chip8.V[(chip8.opcode&0x0F00)>>8] >> 7
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.V[(chip8.opcode&0x0F00)>>8] << 1
			chip8.pc = chip8.pc + 2
		default:
			fmt.Printf("Invalid chip8.opcode %X\n", chip8.opcode)
		}
	case 0x9000: // 9XY0 Skips the next instruction if VX doesn't equal VY
		if chip8.V[(chip8.opcode&0x0F00)>>8] != chip8.V[(chip8.opcode&0x00F0)>>4] {
			chip8.pc = chip8.pc + 4
		} else {
			chip8.pc = chip8.pc + 2
		}
	case 0xA000: // 0xANNN Sets I to the address NNN
		chip8.I = chip8.opcode & 0x0FFF
		chip8.pc = chip8.pc + 2
	case 0xB000: // 0xBNNN Jumps to the address NNN plus V0
		chip8.pc = (chip8.opcode & 0x0FFF) + uint16(chip8.V[0x0])
	case 0xC000: // 0xCXNN Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
		chip8.V[(chip8.opcode&0x0F00)>>8] = uint8(rand.Intn(256)) & uint8(chip8.opcode&0x00FF)
		chip8.pc = chip8.pc + 2
	case 0xD000: // 0xDXYN Draws a sprite at coordinate (VX, VY)
		x := chip8.V[(chip8.opcode&0x0F00)>>8]
		y := chip8.V[(chip8.opcode&0x00F0)>>4]
		h := chip8.opcode & 0x000F
		chip8.V[0xF] = 0
		var j uint16 = 0
		var i uint16 = 0
		for j = 0; j < h; j++ {
			pixel := chip8.memory[chip8.I+j]
			for i = 0; i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					if chip8.gfx[(y + uint8(j))][x+uint8(i)] == 1 {
						chip8.V[0xF] = 1
					}
					chip8.gfx[(y + uint8(j))][x+uint8(i)] ^= 1
				}
			}
		}
		chip8.drawFlag = true
		chip8.pc = chip8.pc + 2
	case 0xE000:
		switch chip8.opcode & 0x00FF {
		case 0x009E: // 0xEX9E Skips the next instruction if the key stored in VX is pressed
			if chip8.key[chip8.V[(chip8.opcode&0x0F00)>>8]] == 1 {
				chip8.pc = chip8.pc + 4
			} else {
				chip8.pc = chip8.pc + 2
			}
		case 0x00A1: // 0xEXA1 Skips the next instruction if the key stored in VX isn't pressed
			if chip8.key[chip8.V[(chip8.opcode&0x0F00)>>8]] == 0 {
				chip8.pc = chip8.pc + 4
			} else {
				chip8.pc = chip8.pc + 2
			}
		default:
			fmt.Printf("Invalid chip8.opcode %X\n", chip8.opcode)
		}
	case 0xF000:
		switch chip8.opcode & 0x00FF {
		case 0x0007: // 0xFX07 Sets VX to the value of the delay timer
			chip8.V[(chip8.opcode&0x0F00)>>8] = chip8.delay_timer
			chip8.pc = chip8.pc + 2
		case 0x000A: // 0xFX0A A key press is awaited, and then stored in VX
			pressed := false
			for i := 0; i < len(chip8.key); i++ {
				if chip8.key[i] != 0 {
					chip8.V[(chip8.opcode&0x0F00)>>8] = uint8(i)
					pressed = true
				}
			}
			if !pressed {
				return
			}
			chip8.pc = chip8.pc + 2
		case 0x0015: // 0xFX15 Sets the delay timer to VX
			chip8.delay_timer = chip8.V[(chip8.opcode&0x0F00)>>8]
			chip8.pc = chip8.pc + 2
		case 0x0018: // 0xFX18 Sets the sound timer to VX
			chip8.sound_timer = chip8.V[(chip8.opcode&0x0F00)>>8]
			chip8.pc = chip8.pc + 2
		case 0x001E: // 0xFX1E Adds VX to I
			if chip8.I+uint16(chip8.V[(chip8.opcode&0x0F00)>>8]) > 0xFFF {
				chip8.V[0xF] = 1
			} else {
				chip8.V[0xF] = 0
			}
			chip8.I = chip8.I + uint16(chip8.V[(chip8.opcode&0x0F00)>>8])
			chip8.pc = chip8.pc + 2
		case 0x0029: // 0xFX29 Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font
			chip8.I = uint16(chip8.V[(chip8.opcode&0x0F00)>>8]) * 0x5
			chip8.pc = chip8.pc + 2
		case 0x0033: // 0xFX33 Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2
			chip8.memory[chip8.I] = chip8.V[(chip8.opcode&0x0F00)>>8] / 100
			chip8.memory[chip8.I+1] = (chip8.V[(chip8.opcode&0x0F00)>>8] / 10) % 10
			chip8.memory[chip8.I+2] = (chip8.V[(chip8.opcode&0x0F00)>>8] % 100) / 10
			chip8.pc = chip8.pc + 2
		case 0x0055: // 0xFX55 Stores V0 to VX (including VX) in chip8.memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((chip8.opcode&0x0F00)>>8)+1; i++ {
				chip8.memory[uint16(i)+chip8.I] = chip8.V[i]
			}
			chip8.I = ((chip8.opcode & 0x0F00) >> 8) + 1
			chip8.pc = chip8.pc + 2
		case 0x0065: // 0xFX65 Fills V0 to VX (including VX) with values from chip8.memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((chip8.opcode&0x0F00)>>8)+1; i++ {
				chip8.V[i] = chip8.memory[chip8.I+uint16(i)]
			}
			chip8.I = ((chip8.opcode & 0x0F00) >> 8) + 1
			chip8.pc = chip8.pc + 2
		default:
			fmt.Printf("Invalid chip8.opcode %X\n", chip8.opcode)
		}
	default:
		fmt.Printf("Invalid chip8.opcode %X\n", chip8.opcode)
	}

}

func (chip8 *Chip8) Key(num uint8, down bool) {
	if down {
		chip8.key[num] = 1
	} else {
		chip8.key[num] = 0
	}
}
func (chip8 *Chip8) GetDrawFlag() bool {
	return chip8.drawFlag
}

func (chip8 *Chip8) GetGFX() [32][64]uint8 {
	return chip8.gfx
}
