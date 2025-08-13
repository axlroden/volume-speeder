//go:build darwin

package main

import (
    "github.com/getlantern/systray"
    icon "github.com/axlroden/volume-speeder/app/src/icon"
)

// startTray initialises the macOS menu bar icon with a Settings menu.
func (a *App) startTray() {
    go systray.Run(func() {
        systray.SetIcon(icon.Data)
        systray.SetTooltip("Volume Speeder")
        mSettings := systray.AddMenuItem("Settings", "Open settings")
        systray.AddSeparator()
        mQuit := systray.AddMenuItem("Quit", "Quit")
        go func() {
            for {
                select {
                case <-mSettings.ClickedCh:
                    go runGUI(a)
                case <-mQuit.ClickedCh:
                    systray.Quit()
                    return
                }
            }
        }()
    }, func() {})
}
