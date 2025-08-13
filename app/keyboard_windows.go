//go:build windows

package main

import (
    "syscall"
    "time"
    "unsafe"
)

// Windows API definitions
var (
    user32              = syscall.NewLazyDLL("user32.dll")
    setWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
    callNextHookEx      = user32.NewProc("CallNextHookEx")
    unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
    getMessage          = user32.NewProc("GetMessageW")
    sendInput           = user32.NewProc("SendInput")
)

const (
    whKeyboardLL = 13
    wmKeydown    = 0x0100
    vkVolumeUp   = 0xAF
    vkVolumeDown = 0xAE
)

type kbdllhookstruct struct {
    VkCode      uint32
    ScanCode    uint32
    Flags       uint32
    Time        uint32
    DwExtraInfo uintptr
}

type input struct {
    Type uint32
    Ki   keybdinput
}

type keybdinput struct {
    WVk         uint16
    WScan       uint16
    DwFlags     uint32
    Time        uint32
    DwExtraInfo uintptr
}

var hHook syscall.Handle

func (a *App) startKeyboardHook() {
    // The keyboard hook callback needs to be a global function or a static method.
    // We make it a closure that captures our App instance `a`.
    hookProc := func(nCode int, wParam uintptr, lParam uintptr) uintptr {
        if nCode >= 0 && wParam == wmKeydown {
            kbd := (*kbdllhookstruct)(unsafe.Pointer(lParam))
            vkCode := kbd.VkCode
            if vkCode == vkVolumeUp || vkCode == vkVolumeDown {
                go a.amplifyKeyPress(vkCode)
                return 1 // Block original key press
            }
        }
    ret, _, _ := callNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
        return ret
    }

    hHook, _, _ = setWindowsHookEx.Call(whKeyboardLL, syscall.NewCallback(hookProc), 0, 0)
    if hHook == 0 {
        return
    }
    defer unhookWindowsHookEx.Call(hHook)

    // This message loop is necessary for the hook to receive events.
    var msg struct{}
    for getMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0) != 0 {
    }
}

func (a *App) amplifyKeyPress(vkCode uint32) {
    multiplier := a.volumeMultiplier.Load()
    if multiplier <= 0 {
        return
    }

    press := input{Type: 1, Ki: keybdinput{WVk: uint16(vkCode)}}
    release := input{Type: 1, Ki: keybdinput{WVk: uint16(vkCode), DwFlags: 2}}

    for i := 0; i < int(multiplier); i++ {
        sendInput.Call(1, uintptr(unsafe.Pointer(&press)), unsafe.Sizeof(press))
        time.Sleep(1 * time.Millisecond)
        sendInput.Call(1, uintptr(unsafe.Pointer(&release)), unsafe.Sizeof(release))
    }
}
