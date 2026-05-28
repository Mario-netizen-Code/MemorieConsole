//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

type _COORD struct {
	X, Y int16
}

type _SMALL_RECT struct {
	Left, Top, Right, Bottom int16
}

type _CONSOLE_SCREEN_BUFFER_INFO struct {
	Size              _COORD
	CursorPosition    _COORD
	Attributes        uint16
	Window            _SMALL_RECT
	MaximumWindowSize _COORD
}

func terminalWidth() int {
	handle, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		return 120
	}

	var info _CONSOLE_SCREEN_BUFFER_INFO
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&info)),
	)
	if ret == 0 {
		return 120
	}
	return int(info.Window.Right - info.Window.Left + 1)
}
