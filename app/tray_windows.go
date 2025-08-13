//go:build windows

package main

import (
    "github.com/getlantern/systray"
    icon "github.com/axlroden/volume-speeder/app/src/icon"
)

// startTray initializes the Windows system tray using an external library.
func (a *App) startTray() {
    go systray.Run(func() {
        systray.SetIcon(icon.Data)
        systray.SetTooltip("Volume Speeder")
    mShow := systray.AddMenuItem("Settings", "Open settings")
        systray.AddSeparator()
        mQuit := systray.AddMenuItem("Quit", "Quit")
        go func() {
            for {
                select {
                case <-mShow.ClickedCh:
                    go runGUI(a)
                case <-mQuit.ClickedCh:
                    systray.Quit()
                    return
                }
            }
        }()
    }, func() {})
}
