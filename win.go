package main

import (
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	procSetWindowsHookEx = user32.NewProc("SetWindowsHookExW")
	procLowLevelKeyboard = user32.NewProc("LowLevelKeyboardProc")
	procCallNextHookEx = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage = user32.NewProc("DispatchMessageW")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetKeyState = user32.NewProc("GetKeyState")
	procGetWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	procGetWindowText = user32.NewProc("GetWindowTextW")

	keyboardHook HHOOK
)

type (
	DWORD uint32
	WPARAM uintptr
	LPARAM uintptr
	LRESULT uintptr
	HANDLE uintptr
	HINSTANCE HANDLE
	HHOOK HANDLE
	HWND HANDLE
)


const (
	WH_KEYBOARD_LL = 13
	WH_KEYBOARD = 2
	WM_KEYDOWN = 256
	WM_SYSKEYDOWN = 260
	WM_KEYUP = 257
	WM_SYSKEYUP = 261
	WM_KEYFIRST = 256
	WM_KEYLAST = 264
	PM_NOREMOVE = 0x000
	PM_REMOVE = 0x001
	PM_NOYIELD = 0x002
	WM_LBUTTONDOWN = 513
	WM_RBUTTONDOWN = 516
	NULL = 0
)

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

type KBDLLHOOKSTRUCT struct {
	VkCode DWORD
	ScanCode DWORD
	Flags DWORD
	Time DWORD
	DwExtraInfo uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd HWND
	Message uint32
	WParam uintptr
	LParam uintptr
	Time uint32
	Pt POINT
}

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

func TranslateMessage(msg *MSG) bool {
	ret, _, _ := procTranslateMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret != 0
}

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret
}

func LowLevelKeyboardProc (nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procLowLevelKeyboard.Call(
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func GetForegroundWindow() HWND {
	ret, _, _ := procGetForegroundWindow.Call(0)

	return HWND(ret)
}

func GetKeyState(nVirtKey int32) int16 {
	ret, _, _ := procGetKeyState.Call(1, uintptr(nVirtKey))

	return int16(ret)
}

func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

// Virtual-Key Codes
const (
	VK_LBUTTON             = 0x01
	VK_RBUTTON             = 0x02
	VK_CANCEL              = 0x03
	VK_MBUTTON             = 0x04
	VK_XBUTTON1            = 0x05
	VK_XBUTTON2            = 0x06
	VK_BACK                = 0x08
	VK_TAB                 = 0x09
	VK_CLEAR               = 0x0C
	VK_RETURN              = 0x0D
	VK_SHIFT               = 0x10
	VK_CONTROL             = 0x11
	VK_MENU                = 0x12
	VK_PAUSE               = 0x13
	VK_CAPITAL             = 0x14
	VK_KANA                = 0x15
	VK_HANGEUL             = 0x15
	VK_HANGUL              = 0x15
	VK_JUNJA               = 0x17
	VK_FINAL               = 0x18
	VK_HANJA               = 0x19
	VK_KANJI               = 0x19
	VK_ESCAPE              = 0x1B
	VK_CONVERT             = 0x1C
	VK_NONCONVERT          = 0x1D
	VK_ACCEPT              = 0x1E
	VK_MODECHANGE          = 0x1F
	VK_SPACE               = 0x20
	VK_PRIOR               = 0x21
	VK_NEXT                = 0x22
	VK_END                 = 0x23
	VK_HOME                = 0x24
	VK_LEFT                = 0x25
	VK_UP                  = 0x26
	VK_RIGHT               = 0x27
	VK_DOWN                = 0x28
	VK_SELECT              = 0x29
	VK_PRINT               = 0x2A
	VK_EXECUTE             = 0x2B
	VK_SNAPSHOT            = 0x2C
	VK_INSERT              = 0x2D
	VK_DELETE              = 0x2E
	VK_HELP                = 0x2F
	VK_LWIN                = 0x5B
	VK_RWIN                = 0x5C
	VK_APPS                = 0x5D
	VK_SLEEP               = 0x5F
	VK_NUMPAD0             = 0x60
	VK_NUMPAD1             = 0x61
	VK_NUMPAD2             = 0x62
	VK_NUMPAD3             = 0x63
	VK_NUMPAD4             = 0x64
	VK_NUMPAD5             = 0x65
	VK_NUMPAD6             = 0x66
	VK_NUMPAD7             = 0x67
	VK_NUMPAD8             = 0x68
	VK_NUMPAD9             = 0x69
	VK_MULTIPLY            = 0x6A
	VK_ADD                 = 0x6B
	VK_SEPARATOR           = 0x6C
	VK_SUBTRACT            = 0x6D
	VK_DECIMAL             = 0x6E
	VK_DIVIDE              = 0x6F
	VK_F1                  = 0x70
	VK_F2                  = 0x71
	VK_F3                  = 0x72
	VK_F4                  = 0x73
	VK_F5                  = 0x74
	VK_F6                  = 0x75
	VK_F7                  = 0x76
	VK_F8                  = 0x77
	VK_F9                  = 0x78
	VK_F10                 = 0x79
	VK_F11                 = 0x7A
	VK_F12                 = 0x7B
	VK_F13                 = 0x7C
	VK_F14                 = 0x7D
	VK_F15                 = 0x7E
	VK_F16                 = 0x7F
	VK_F17                 = 0x80
	VK_F18                 = 0x81
	VK_F19                 = 0x82
	VK_F20                 = 0x83
	VK_F21                 = 0x84
	VK_F22                 = 0x85
	VK_F23                 = 0x86
	VK_F24                 = 0x87
	VK_NUMLOCK             = 0x90
	VK_SCROLL              = 0x91
	VK_OEM_NEC_EQUAL       = 0x92
	VK_OEM_FJ_JISHO        = 0x92
	VK_OEM_FJ_MASSHOU      = 0x93
	VK_OEM_FJ_TOUROKU      = 0x94
	VK_OEM_FJ_LOYA         = 0x95
	VK_OEM_FJ_ROYA         = 0x96
	VK_LSHIFT              = 0xA0
	VK_RSHIFT              = 0xA1
	VK_LCONTROL            = 0xA2
	VK_RCONTROL            = 0xA3
	VK_LMENU               = 0xA4
	VK_RMENU               = 0xA5
	VK_BROWSER_BACK        = 0xA6
	VK_BROWSER_FORWARD     = 0xA7
	VK_BROWSER_REFRESH     = 0xA8
	VK_BROWSER_STOP        = 0xA9
	VK_BROWSER_SEARCH      = 0xAA
	VK_BROWSER_FAVORITES   = 0xAB
	VK_BROWSER_HOME        = 0xAC
	VK_VOLUME_MUTE         = 0xAD
	VK_VOLUME_DOWN         = 0xAE
	VK_VOLUME_UP           = 0xAF
	VK_MEDIA_NEXT_TRACK    = 0xB0
	VK_MEDIA_PREV_TRACK    = 0xB1
	VK_MEDIA_STOP          = 0xB2
	VK_MEDIA_PLAY_PAUSE    = 0xB3
	VK_LAUNCH_MAIL         = 0xB4
	VK_LAUNCH_MEDIA_SELECT = 0xB5
	VK_LAUNCH_APP1         = 0xB6
	VK_LAUNCH_APP2         = 0xB7
	VK_OEM_1               = 0xBA
	VK_OEM_PLUS            = 0xBB
	VK_OEM_COMMA           = 0xBC
	VK_OEM_MINUS           = 0xBD
	VK_OEM_PERIOD          = 0xBE
	VK_OEM_2               = 0xBF
	VK_OEM_3               = 0xC0
	VK_OEM_4               = 0xDB
	VK_OEM_5               = 0xDC
	VK_OEM_6               = 0xDD
	VK_OEM_7               = 0xDE
	VK_OEM_8               = 0xDF
	VK_OEM_AX              = 0xE1
	VK_OEM_102             = 0xE2
	VK_ICO_HELP            = 0xE3
	VK_ICO_00              = 0xE4
	VK_PROCESSKEY          = 0xE5
	VK_ICO_CLEAR           = 0xE6
	VK_PACKET              = 0xE7
	VK_OEM_RESET           = 0xE9
	VK_OEM_JUMP            = 0xEA
	VK_OEM_PA1             = 0xEB
	VK_OEM_PA2             = 0xEC
	VK_OEM_PA3             = 0xED
	VK_OEM_WSCTRL          = 0xEE
	VK_OEM_CUSEL           = 0xEF
	VK_OEM_ATTN            = 0xF0
	VK_OEM_FINISH          = 0xF1
	VK_OEM_COPY            = 0xF2
	VK_OEM_AUTO            = 0xF3
	VK_OEM_ENLW            = 0xF4
	VK_OEM_BACKTAB         = 0xF5
	VK_ATTN                = 0xF6
	VK_CRSEL               = 0xF7
	VK_EXSEL               = 0xF8
	VK_EREOF               = 0xF9
	VK_PLAY                = 0xFA
	VK_ZOOM                = 0xFB
	VK_NONAME              = 0xFC
	VK_PA1                 = 0xFD
	VK_OEM_CLEAR           = 0xFE
)
