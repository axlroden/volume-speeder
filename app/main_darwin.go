//go:build darwin

package main

func main() {
    app := NewApp()
    app.startup(nil)
    runGUI(app)
}
