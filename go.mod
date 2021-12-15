module main

go 1.17

require github.com/veandco/go-sdl2 v0.4.10

require emulator/chip8 v0.0.0-00010101000000-000000000000 // indirect

replace emulator/chip8 => ./chip8
