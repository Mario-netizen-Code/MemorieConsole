//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	procFillConsoleOutputCharacterW = kernel32.NewProc("FillConsoleOutputCharacterW")
	procFillConsoleOutputAttribute  = kernel32.NewProc("FillConsoleOutputAttribute")
	procSetConsoleCursorPosition    = kernel32.NewProc("SetConsoleCursorPosition")
)

func clearScreen() {
	handle, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		return
	}

	var csbi _CONSOLE_SCREEN_BUFFER_INFO
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&csbi)),
	)
	if ret == 0 {
		return
	}

	total := uint32(csbi.Size.X) * uint32(csbi.Size.Y)

	var written uint32
	procFillConsoleOutputCharacterW.Call(
		uintptr(handle),
		uintptr(' '),
		uintptr(total),
		0,
		uintptr(unsafe.Pointer(&written)),
	)
	procFillConsoleOutputAttribute.Call(
		uintptr(handle),
		uintptr(csbi.Attributes),
		uintptr(total),
		0,
		uintptr(unsafe.Pointer(&written)),
	)
	procSetConsoleCursorPosition.Call(
		uintptr(handle),
		0,
	)
}
