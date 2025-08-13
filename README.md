# Volume Speeder

Root contains install scripts and docs. The app lives in `app/`.

Structure:
- app/            Go sources, Wails config, and frontend assets
- installer/      Windows NSIS script
- docs/           GitHub Pages site
- install.sh      Linux/macOS installer
- install.ps1     Windows installer

Build locally:
- Linux CLI: `cd app && go build`
- Windows/macOS GUI: `cd app && wails build`

Config (Linux): `~/.config/volume-speeder/config` with `multiplier=<n>`.

## macOS (no Apple Developer ID)

The app is unsigned, so Gatekeeper may warn on first launch.

- Easiest (auto-trust during install):
	- `sh -c "TRUST=1 $(curl -fsSL https://raw.githubusercontent.com/axlroden/volume-speeder/HEAD/install.sh)"`
- Manual: run the installer normally, then in Finder:
	- Applications → right-click "Volume Speeder" → Open → Open

