//go:build windows

package main

import (
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/wailsapp/wails/v2/pkg/menu"
    icon "github.com/axlroden/volume-speeder/app/src/icon"
)

func getAppOptions(app *App) *options.App {
    return &options.App{
        Title:            "Volume Speeder",
        Width:            280,
        Height:           120,
        HideWindowOnClose: true,
        AssetServer: &assetserver.Options{ Assets: assets },
        Frameless:        true,
        BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
        OnStartup:        app.startup,
        OnDomReady:       app.domReady,
        Bind: []interface{}{ app },
        Systray: &options.Systray{
            Icon:       icon.Data,
            Title:      "Volume Speeder",
            Tooltip:    "Volume Speeder",
            OnLeftClick: app.showWindow,
            Menu: menu.NewMenuFromItems(
                menu.Separator(),
                menu.Text("Quit", app.Quit),
            ),
        },
    }
}
