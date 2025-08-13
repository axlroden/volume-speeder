//go:build windows

package main

func main() {
    app := NewApp()
    app.startup(nil)
    // start keyboard hook
    app.domReady(nil)
    // No GUI at launch; tray menu opens settings window when needed
    select {}
}
