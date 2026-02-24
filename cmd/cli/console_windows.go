//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	enableLineInput      = 0x0002
	enableEchoInput      = 0x0004
	enableProcessedInput = 0x0001
	enableWindowInput    = 0x0008
	enableExtendedFlags  = 0x0080
	enableVTProcessing   = 0x0004

	keyEventType = 0x0001
)

const (
	vkUp     = 0x26
	vkDown   = 0x28
	vkBack   = 0x08
	vkReturn = 0x0D
	vkEscape = 0x1B
)

type coord struct {
	X int16
	Y int16
}

type smallRect struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

type consoleScreenBufferInfo struct {
	Size              coord
	CursorPosition    coord
	Attributes        uint16
	Window            smallRect
	MaximumWindowSize coord
}

type keyEventRecord struct {
	KeyDown          int32
	RepeatCount      uint16
	VirtualKeyCode   uint16
	VirtualScanCode  uint16
	UnicodeChar      uint16
	ControlKeyState  uint32
}

type inputRecord struct {
	EventType uint16
	_         uint16
	Event     [16]byte
}

var (
	kernel32                          = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode                = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode                = kernel32.NewProc("SetConsoleMode")
	procReadConsoleInputW             = kernel32.NewProc("ReadConsoleInputW")
	procGetConsoleScreenBufferInfo    = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procFillConsoleOutputCharacterW   = kernel32.NewProc("FillConsoleOutputCharacterW")
	procFillConsoleOutputAttribute    = kernel32.NewProc("FillConsoleOutputAttribute")
	procSetConsoleCursorPosition      = kernel32.NewProc("SetConsoleCursorPosition")
)

var (
	stdinHandle  syscall.Handle
	stdoutHandle syscall.Handle
	oldInMode    uint32
	oldOutMode   uint32
)

func initWindowsConsole() (func(), error) {
	if stdinHandle == 0 {
		stdinHandle = syscall.Handle(os.Stdin.Fd())
		stdoutHandle = syscall.Handle(os.Stdout.Fd())
	}

	if err := getConsoleMode(stdinHandle, &oldInMode); err != nil {
		return nil, fmt.Errorf("GetConsoleMode stdin: %w", err)
	}
	if err := getConsoleMode(stdoutHandle, &oldOutMode); err != nil {
		return nil, fmt.Errorf("GetConsoleMode stdout: %w", err)
	}

	inMode := oldInMode
	inMode &^= enableLineInput | enableEchoInput
	inMode |= enableExtendedFlags | enableWindowInput | enableProcessedInput
	if err := setConsoleMode(stdinHandle, inMode); err != nil {
		return nil, fmt.Errorf("SetConsoleMode stdin: %w", err)
	}

	outMode := oldOutMode | enableVTProcessing
	_ = setConsoleMode(stdoutHandle, outMode)

	restore := func() {
		_ = setConsoleMode(stdinHandle, oldInMode)
		_ = setConsoleMode(stdoutHandle, oldOutMode)
	}
	return restore, nil
}

func readKeyWindows() (keyEvent, error) {
	var rec inputRecord
	var read uint32
	for {
		if err := readConsoleInput(stdinHandle, &rec, 1, &read); err != nil {
			return keyEvent{}, err
		}
		if rec.EventType != keyEventType {
			continue
		}
		kev := *(*keyEventRecord)(unsafe.Pointer(&rec.Event[0]))
		if kev.KeyDown == 0 {
			continue
		}
		switch kev.VirtualKeyCode {
		case vkUp:
			return keyEvent{kind: keyUp}, nil
		case vkDown:
			return keyEvent{kind: keyDown}, nil
		case vkBack:
			return keyEvent{kind: keyBackspace}, nil
		case vkReturn:
			return keyEvent{kind: keyEnter}, nil
		case vkEscape:
			return keyEvent{kind: keyEsc}, nil
		}
		if kev.UnicodeChar != 0 {
			return keyEvent{kind: keyRune, r: rune(kev.UnicodeChar)}, nil
		}
	}
}

func clearScreenWindows() {
	var info consoleScreenBufferInfo
	if err := getConsoleScreenBufferInfo(stdoutHandle, &info); err != nil {
		return
	}
	size := uint32(info.Size.X) * uint32(info.Size.Y)
	var written uint32
	origin := coord{X: 0, Y: 0}
	_ = fillConsoleOutputCharacter(stdoutHandle, ' ', size, origin, &written)
	_ = fillConsoleOutputAttribute(stdoutHandle, info.Attributes, size, origin, &written)
	_ = setConsoleCursorPosition(stdoutHandle, origin)
}

func getConsoleMode(h syscall.Handle, mode *uint32) error {
	r1, _, e1 := procGetConsoleMode.Call(uintptr(h), uintptr(unsafe.Pointer(mode)))
	if r1 == 0 {
		return e1
	}
	return nil
}

func setConsoleMode(h syscall.Handle, mode uint32) error {
	r1, _, e1 := procSetConsoleMode.Call(uintptr(h), uintptr(mode))
	if r1 == 0 {
		return e1
	}
	return nil
}

func readConsoleInput(h syscall.Handle, rec *inputRecord, length uint32, read *uint32) error {
	r1, _, e1 := procReadConsoleInputW.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(rec)),
		uintptr(length),
		uintptr(unsafe.Pointer(read)),
	)
	if r1 == 0 {
		return e1
	}
	return nil
}

func getConsoleScreenBufferInfo(h syscall.Handle, info *consoleScreenBufferInfo) error {
	r1, _, e1 := procGetConsoleScreenBufferInfo.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(info)),
	)
	if r1 == 0 {
		return e1
	}
	return nil
}

func fillConsoleOutputCharacter(h syscall.Handle, char rune, length uint32, coord coord, written *uint32) error {
	r1, _, e1 := procFillConsoleOutputCharacterW.Call(
		uintptr(h),
		uintptr(char),
		uintptr(length),
		*(*uintptr)(unsafe.Pointer(&coord)),
		uintptr(unsafe.Pointer(written)),
	)
	if r1 == 0 {
		return e1
	}
	return nil
}

func fillConsoleOutputAttribute(h syscall.Handle, attr uint16, length uint32, coord coord, written *uint32) error {
	r1, _, e1 := procFillConsoleOutputAttribute.Call(
		uintptr(h),
		uintptr(attr),
		uintptr(length),
		*(*uintptr)(unsafe.Pointer(&coord)),
		uintptr(unsafe.Pointer(written)),
	)
	if r1 == 0 {
		return e1
	}
	return nil
}

func setConsoleCursorPosition(h syscall.Handle, coord coord) error {
	r1, _, e1 := procSetConsoleCursorPosition.Call(
		uintptr(h),
		*(*uintptr)(unsafe.Pointer(&coord)),
	)
	if r1 == 0 {
		return e1
	}
	return nil
}
