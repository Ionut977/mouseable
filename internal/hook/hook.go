package hook

import (
	"sync"
	"unsafe"

	"github.com/JamesHovious/w32"

	"github.com/wirekang/mouseable/internal/def"
	"github.com/wirekang/mouseable/internal/lg"
)

const activateID = 100

var hHook w32.HHOOK
var state struct {
	sync.RWMutex
	activateKey def.HotKey
	isHooking   bool
}

func SetKey(config def.Config) {
	state.Lock()
	state.activateKey = config.HotKeyMap[def.Activate]
	state.Unlock()
	go messageLoop(config.HotKeyMap[def.Activate])
}

func messageLoop(activateKey def.HotKey) {
	unregisterHotKey(activateID)
	registerHotKey(activateID, activateKey)
	lg.Logf("Start message loop")
	defer func() { lg.Logf("Exit message loop") }()
	var msg w32.MSG
	for {
		r := w32.GetMessage(&msg, 0, 0, 0)
		state.RLock()
		if state.activateKey != activateKey {
			state.RUnlock()
			return
		}
		state.RUnlock()
		lg.Logf("message: %+v", msg)
		if r == 0 {
			return
		}
		if msg.Message == w32.WM_HOTKEY {
			switch msg.WParam {
			case activateID:
				hook()
			}
		}
	}
}

func getMod(h def.HotKey) (mod uint) {
	if h.IsAlt {
		mod = mod | w32.MOD_ALT
	}
	if h.IsControl {
		mod = mod | w32.MOD_CONTROL
	}
	if h.IsWin {
		mod = mod | w32.MOD_WIN
	}
	if h.IsShift {
		mod = mod | w32.MOD_SHIFT
	}
	return
}

func hook() {
	state.Lock()
	if state.isHooking {
		lg.Errorf("Already hooking")
		state.Unlock()
		return
	}
	state.isHooking = true
	state.Unlock()

	lg.Logf("Hook")
	hHook = w32.SetWindowsHookEx(w32.WH_KEYBOARD_LL, hookProc, 0, 0)
	DI.OnHook()
}

func Unhook() {
	state.Lock()
	if !state.isHooking {
		lg.Errorf("Already not hooking")
		state.Unlock()
		return
	}
	state.isHooking = false
	state.Unlock()
	lg.Logf("Unhook")
	w32.UnhookWindowsHookEx(hHook)
	DI.OnUnhook()
}

func registerHotKey(id int, key def.HotKey) {
	err := w32.RegisterHotKey(
		0, id, getMod(key)|w32.MOD_NOREPEAT, uint(key.KeyCode),
	)

	if err != nil {
		lg.Errorf("registerHotKey: %v", err)
	} else {
		lg.Logf("registerHotKey: %d %+v", id, key)
	}
}

func unregisterHotKey(id int) {
	err := w32.UnregisterHotKey(0, id)
	if err != nil {
		lg.Errorf("unregisterHotKey: %v", err)
	} else {
		lg.Logf("unregisterHotKey: %d", id)
	}
}

func SetCursorPos(x, y int) {
	w32.SetCursorPos(x, y)
}

func AddCursorPos(dx, dy int32) {
	sendMouseInput(dx, dy, 0, w32.MOUSEEVENTF_MOVE)
}

func GetCursorPos() (x, y int) {
	x, y, _ = w32.GetCursorPos()
	return
}

// MouseDown send mouse down event.
//
// left = 0
//
// right = 1
//
// middle = 2
func MouseDown(button int) {
	var flag uint32
	switch button {
	case 0:
		flag = w32.MOUSEEVENTF_LEFTDOWN
	case 1:
		flag = w32.MOUSEEVENTF_RIGHTDOWN
	case 2:
		flag = w32.MOUSEEVENTF_MIDDLEDOWN
	}
	sendMouseInput(0, 0, 0, flag)
}

// MouseUp send mouse up event. checkout MouseDown for button.
func MouseUp(button int) {
	var flag uint32
	switch button {
	case 0:
		flag = w32.MOUSEEVENTF_LEFTUP
	case 1:
		flag = w32.MOUSEEVENTF_RIGHTUP
	case 2:
		flag = w32.MOUSEEVENTF_MIDDLEUP
	}
	sendMouseInput(0, 0, 0, flag)
}

func Wheel(amount int, horizontal bool) {
	flag := uint32(w32.MOUSEEVENTF_WHEEL)
	if horizontal {
		flag = w32.MOUSEEVENTF_HWHEEL
	}
	sendMouseInput(0, 0, uint32(amount), flag)
}

func sendMouseInput(dx, dy int32, mouseData uint32, flags ...uint32) {
	var dwFlags uint32
	for _, f := range flags {
		dwFlags |= f
	}

	input := []w32.INPUT{
		{
			Type: w32.INPUT_MOUSE,
			Mi: w32.MOUSEINPUT{
				Dx:          dx,
				Dy:          dy,
				MouseData:   mouseData,
				DwFlags:     dwFlags,
				Time:        0,
				DwExtraInfo: 0,
			},
		},
	}
	w32.SendInput(input)
}

var hookProc w32.HOOKPROC = func(
	code int, wParam w32.WPARAM, lParam w32.LPARAM,
) w32.LRESULT {
	data := *(*w32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	flagMap := map[string]w32.DWORD{
		"Extended":      data.Flags & (w32.KF_EXTENDED >> 8),
		"LowerInjected": data.Flags & 0x00000002,
		"Injected":      data.Flags & 0x00000010,
		"AltDown":       data.Flags & (w32.KF_ALTDOWN >> 8),
		"Up":            data.Flags & (w32.KF_UP >> 8),
	}

	if DI.OnKey != nil {
		isDown := flagMap["Up"] == 0
		preventDefault := DI.OnKey(uint32(data.VkCode), isDown)

		if preventDefault {
			return 1
		}
	}

	return w32.CallNextHookEx(0, code, wParam, lParam)
}
