//go:build windows

package app

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32                        = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo  = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procFillConsoleOutputCharacterW = kernel32.NewProc("FillConsoleOutputCharacterW")
	procFillConsoleOutputAttribute  = kernel32.NewProc("FillConsoleOutputAttribute")
	procSetConsoleCursorPosition    = kernel32.NewProc("SetConsoleCursorPosition")
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

func openFileShared(path string) (*os.File, error) {
	pathp, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}

	h, err := syscall.CreateFile(
		pathp,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return nil, os.NewSyscallError("CreateFile", err)
	}
	return os.NewFile(uintptr(h), path), nil
}

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
