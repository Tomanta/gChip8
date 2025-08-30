# Chip8 Emulator

This is a Chip8 emulator developed in Go using [Ebitengine](https://ebitengine.org/) as the frontend.

Developed using Go version 1.24.6 on Ubuntu running on Windows 11 WSL2. 

# Running

Roms should be located in `./roms`. Run a specific rom using `gchip [ROM_NAME]`. If no rom name supplied will attempt to run `ibm_logo.ch8`.

## Input

The keypad is mapped to the keyboard as:

```
1 2 3 4     1 2 3 C
q w e r     4 5 6 D
a s d f  =  7 8 9 E
z x c v     A 0 B F
```

## Resources:

Most test roms came from: [Timedus' test suite](https://github.com/Timendus/chip8-test-suite/tree/main)

Main programming reference was: [Tobias V. Langhoff's Chip8 guide](https://tobiasvl.github.io/blog/write-a-chip-8-emulator/)