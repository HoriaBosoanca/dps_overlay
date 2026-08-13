package overlay

import (
	"errors"
	"strings"
	"syscall"
	"time"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	user32dll               = syscall.NewLazyDLL("user32.dll")
	funcEnumWindows         = user32dll.NewProc("EnumWindows")
	funcGetWindowText       = user32dll.NewProc("GetWindowTextW")
	funcGetWindowRect       = user32dll.NewProc("GetWindowRect")
	funcGetForegroundWindow = user32dll.NewProc("GetForegroundWindow")
	funcIsWindow            = user32dll.NewProc("IsWindow")
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

func RunOverlay() {
	leagueHandle := waitForLeagueHandle()
	rl.SetConfigFlags(
		rl.FlagWindowTransparent |
			rl.FlagWindowTopmost |
			rl.FlagWindowUnfocused |
			rl.FlagWindowMousePassthrough,
	)
	rect := RECT{}
	_, _, _ = funcGetWindowRect.Call(uintptr(leagueHandle), uintptr(unsafe.Pointer(&rect)))
	rl.InitWindow(rect.Right-rect.Left, rect.Bottom-rect.Top, "dps_overlay")
	defer rl.CloseWindow()
	rl.SetWindowPosition(int(rect.Left), int(rect.Top))
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		// if handle is invalid (game ended), wait until next game to get the handle
		if isHandleValid, _, _ := funcIsWindow.Call(uintptr(leagueHandle)); isHandleValid == 0 {
			leagueHandle = waitForLeagueHandle()
		}
		// hide overlay if not tabbed into the game
		if foreground, _, _ := funcGetForegroundWindow.Call(); syscall.Handle(foreground) == leagueHandle {
			rl.ClearWindowState(rl.FlagWindowHidden)
		} else {
			rl.SetWindowState(rl.FlagWindowHidden)
		}
		// raylib drawing stuff
		rl.BeginDrawing()
		rl.ClearBackground(rl.Blank)
		rl.DrawText("Overlay active", 20, 20, 20, rl.Green)
		rl.EndDrawing()
	}
}

func waitForLeagueHandle() (leagueHandle syscall.Handle) {
	for {
		var err error
		leagueHandle, err = findLeagueHandle()
		if err == nil {
			// wait for league to initialize screen size before returning handle
			time.Sleep(2 * time.Second)
			leagueHandle, err = findLeagueHandle()
			break
		} else {
			//fmt.Println(err)
		}
	}
	return leagueHandle
}

func findLeagueHandle() (syscall.Handle, error) {
	var leagueHandle syscall.Handle
	_, _, _ = funcEnumWindows.Call(perWindowCallback, uintptr(unsafe.Pointer(&leagueHandle)))
	if leagueHandle == 0 {
		return 0, errors.New("league window not found")
	}
	return leagueHandle, nil
}

var perWindowCallback = syscall.NewCallback(
	func(hwnd syscall.Handle, lparam uintptr) uintptr {
		//goland:noinspection GoVetUnsafePointer
		leagueHandle := (*syscall.Handle)(unsafe.Pointer(lparam))
		windowTitleUtf16 := make([]uint16, 256)
		_, _, _ = funcGetWindowText.Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(&windowTitleUtf16[0])),
			uintptr(len(windowTitleUtf16)),
		)
		if strings.Contains(syscall.UTF16ToString(windowTitleUtf16), "League of Legends (TM) Client") {
			*leagueHandle = hwnd
			return 0
		}
		return 1
	},
)
