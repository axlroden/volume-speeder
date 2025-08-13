//go:build !windows

package main

// On non-Windows platforms, do nothing so the project still builds.
func (a *App) startKeyboardHook() {}
func (a *App) amplifyKeyPress(_ uint32) {}
