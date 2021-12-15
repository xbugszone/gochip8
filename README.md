# go-chip8

A practise project to learn about emulation, implements CHIP-8 in golang.

## Installing

- Get dependencies

```shell
go get -u github.com/veandco/go-sdl2/sdl
```

- If missing sdl 2 install
```shell
brew install sdl2
or 
sudo apt install libsdl2-dev
```

- Get code

```shell
git clone https://github.com/xbugszone/gochip8.git
```

## INSTALLING PACKAGES

```
go mod tidy
```

## Running


```
go run main.go 
```

## Key Bindings

```
Chip8 keypad         Keyboard mapping
1 | 2 | 3 | C        1 | 2 | 3 | 4
4 | 5 | 6 | D   =>   Q | W | E | R
7 | 8 | 9 | E   =>   A | S | D | F
A | 0 | B | F        Z | X | C | V
```

## Change Game

Go to the file main.go and on the line 10 change this:

```js
const GAME = "pong.rom"
```
to
```js
const GAME = "filename"
```

The games available are in the folder /games/*

## Sources

- [How to write an emulator chip-8 interpreter](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)
- [Cowgod's Chip-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- [Chip-8 opcode table](https://en.wikipedia.org/wiki/CHIP-8)
