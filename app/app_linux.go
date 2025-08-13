//go:build linux

package main

import (
    "context"
    "sync/atomic"
)

// App is a minimal struct to satisfy shared code on Linux without Wails.
type App struct {
    ctx              context.Context
    volumeMultiplier atomic.Int32
}

func NewApp() *App { return &App{} }

// No-ops for Linux CLI build
func (a *App) startup(ctx context.Context)            { a.ctx = ctx; a.volumeMultiplier.Store(3) }
func (a *App) domReady(ctx context.Context)           {}
func (a *App) showWindow()                            {}
func (a *App) GetInitialMultiplier() int              { return int(a.volumeMultiplier.Load()) }
func (a *App) SetMultiplier(value int)                { a.volumeMultiplier.Store(int32(value)) }
func (a *App) Quit()                                  {}
