//go:build windows || darwin

package main

import (
    "fmt"

    fyne "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

// runGUI starts a small settings window to control the multiplier.
func runGUI(state *App) {
    a := app.NewWithID("volume-speeder")
    w := a.NewWindow("Volume Speeder")

    cur := state.GetInitialMultiplier()
    valLabel := widget.NewLabel(fmt.Sprintf("%d", cur))
    slider := widget.NewSlider(1, 10)
    slider.Step = 1
    slider.Value = float64(cur)
    slider.OnChanged = func(v float64) {
        valLabel.SetText(fmt.Sprintf("%d", int(v)))
    }

    save := widget.NewButton("Save", func() {
        state.SetMultiplier(int(slider.Value))
    })
    quit := widget.NewButton("Quit", func() { a.Quit() })

    w.SetContent(container.NewVBox(
        widget.NewLabel("Volume step multiplier"),
        slider,
        valLabel,
        container.NewHBox(save, quit),
    ))
    w.Resize(fyne.NewSize(320, 160))
    w.ShowAndRun()
}
