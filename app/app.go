//go:build windows || darwin

package main

import (
	"context"
	"log"
	"sync/atomic"
)

// App holds application state
type App struct {
	ctx              context.Context
	volumeMultiplier atomic.Int32
}

// NewApp creates a new App
func NewApp() *App { return &App{} }

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Load multiplier from config or default to 3.
	if cfg, err := loadConfig(); err == nil && cfg.Multiplier >= 1 {
		a.volumeMultiplier.Store(int32(cfg.Multiplier))
	} else {
		a.volumeMultiplier.Store(3)
	}
	// Start tray (no-op on non-Windows)
	a.startTray()
}

// domReady is called after the DOM has been loaded.
func (a *App) domReady(ctx context.Context) {
	// Start the keyboard hook in a separate goroutine
	go a.startKeyboardHook()
}

// showWindow is called when the systray icon is left-clicked.
func (a *App) showWindow() {}

// GetInitialMultiplier returns the current multiplier value to the frontend.
func (a *App) GetInitialMultiplier() int {
	return int(a.volumeMultiplier.Load())
}

// SetMultiplier is called from the frontend to update the multiplier.
func (a *App) SetMultiplier(value int) {
	log.Printf("Multiplier set to: %d\n", value)
	a.volumeMultiplier.Store(int32(value))
	_ = saveConfig(Config{Multiplier: int(a.volumeMultiplier.Load())})
}

// Quit closes the application via Wails runtime
func (a *App) Quit() {}
