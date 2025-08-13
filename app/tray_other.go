//go:build !windows && !darwin

package main

// startTray is a no-op on platforms without tray support.
func (a *App) startTray() {}
