package main

import (
	"emulator/chip8"
	sdl "github.com/veandco/go-sdl2/sdl"
	"os"
)
// Default Chip8 resolution
const CHIP_8_WIDTH int32 = 64
const CHIP_8_HEIGHT int32 = 32


func main()  {
	var modifier int32 = 10

	// Create window, chip8 resolution with given modifier
	window, windowErr := sdl.CreateWindow("Chip 8 - ", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, CHIP_8_WIDTH*modifier, CHIP_8_HEIGHT*modifier, sdl.WINDOW_SHOWN)
	if windowErr != nil {
		panic(windowErr)
	}
	defer window.Destroy()

	// Create render surface
	canvas, canvasErr := sdl.CreateRenderer(window, -1, 0)
	if canvasErr != nil {
		panic(canvasErr)
	}
	defer canvas.Destroy()
	// Clear the screen
	canvas.SetDrawColor(255, 0, 0, 255)
	canvas.Clear()

	chip8.Initialize()
	chip8.LoadGame("games/pong.rom")

	for  {
		// Emulate one cycle
		chip8.EmulateCycle()
		if chip8.GetDrawFlag() {
			// Get the display buffer and render
			vector := chip8.GetGFX()
			for j := 0; j < len(vector); j++ {
				for i := 0; i < len(vector[j]); i++ {
					// Values of pixel are stored in 1D array of size 64 * 32
					if vector[j][i] != 0 {
						canvas.SetDrawColor(255, 255, 0, 255)
					} else {
						canvas.SetDrawColor(255, 0, 0, 255)
					}
					canvas.FillRect(&sdl.Rect{
						Y: int32(j) * modifier,
						X: int32(i) * modifier,
						W: modifier,
						H: modifier,
					})
				}
			}

			canvas.Present()
		}
		// Store key press state (Press and Release)
		// Poll for Quit and Keyboard events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch et := event.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			case *sdl.KeyboardEvent:
				if et.Type == sdl.KEYUP {
					switch et.Keysym.Sym {
					case sdl.K_1:
						chip8.Key(0x1, false)
					case sdl.K_2:
						chip8.Key(0x2, false)
					case sdl.K_3:
						chip8.Key(0x3, false)
					case sdl.K_4:
						chip8.Key(0xC, false)
					case sdl.K_q:
						chip8.Key(0x4, false)
					case sdl.K_w:
						chip8.Key(0x5, false)
					case sdl.K_e:
						chip8.Key(0x6, false)
					case sdl.K_r:
						chip8.Key(0xD, false)
					case sdl.K_a:
						chip8.Key(0x7, false)
					case sdl.K_s:
						chip8.Key(0x8, false)
					case sdl.K_d:
						chip8.Key(0x9, false)
					case sdl.K_f:
						chip8.Key(0xE, false)
					case sdl.K_z:
						chip8.Key(0xA, false)
					case sdl.K_x:
						chip8.Key(0x0, false)
					case sdl.K_c:
						chip8.Key(0xB, false)
					case sdl.K_v:
						chip8.Key(0xF, false)
					}
				} else if et.Type == sdl.KEYDOWN {
					switch et.Keysym.Sym {
					case sdl.K_1:
						chip8.Key(0x1, true)
					case sdl.K_2:
						chip8.Key(0x2, true)
					case sdl.K_3:
						chip8.Key(0x3, true)
					case sdl.K_4:
						chip8.Key(0xC, true)
					case sdl.K_q:
						chip8.Key(0x4, true)
					case sdl.K_w:
						chip8.Key(0x5, true)
					case sdl.K_e:
						chip8.Key(0x6, true)
					case sdl.K_r:
						chip8.Key(0xD, true)
					case sdl.K_a:
						chip8.Key(0x7, true)
					case sdl.K_s:
						chip8.Key(0x8, true)
					case sdl.K_d:
						chip8.Key(0x9, true)
					case sdl.K_f:
						chip8.Key(0xE, true)
					case sdl.K_z:
						chip8.Key(0xA, true)
					case sdl.K_x:
						chip8.Key(0x0, true)
					case sdl.K_c:
						chip8.Key(0xB, true)
					case sdl.K_v:
						chip8.Key(0xF, true)
					}
				}
			}
		}
		//wsdl.Delay(1)
	}

}


