package overlay

import (
	"dps_overlay/data"
	"errors"
	"fmt"
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
	data.InitGame()
	rl.SetTraceLogLevel(rl.LogNone)
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
	rl.SetTargetFPS(3)
	data.Font = rl.LoadFont("../assets/fonts/JetBrainsMono-Bold.ttf")
	defer rl.UnloadFont(data.Font)
	for !rl.WindowShouldClose() {
		// hide overlay if not tabbed into the game
		if foreground, _, _ := funcGetForegroundWindow.Call(); syscall.Handle(foreground) == leagueHandle {
			rl.ClearWindowState(rl.FlagWindowHidden)
		} else {
			rl.SetWindowState(rl.FlagWindowHidden)
		}
		// if handle is invalid (game ended), wait until next game to get the handle
		if isHandleValid, _, _ := funcIsWindow.Call(uintptr(leagueHandle)); isHandleValid == 0 {
			leagueHandle = waitForLeagueHandle()
			data.InitGame()
		}
		// raylib drawing stuff
		rl.BeginDrawing()
		rl.ClearBackground(rl.Blank)
		data.LoopGame()
		rl.EndDrawing()
	}
}

func waitForLeagueHandle() (leagueHandle syscall.Handle) {
	firstErrPrinted := false
	for {
		var err error
		leagueHandle, err = findLeagueHandle()
		if err == nil {
			break
		}
		if !firstErrPrinted {
			firstErrPrinted = true
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
	// wait for league to initialize screen size and start http server before returning handle
	time.Sleep(10 * time.Second)
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
