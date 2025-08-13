# Windows installer one-liner for Volume Speeder
# Usage:
#   powershell -ExecutionPolicy Bypass -NoProfile -Command "iwr -useb https://raw.githubusercontent.com/axlroden/volume-speeder/HEAD/install.ps1 | iex"
# Options via env or params: -Owner, -Repo, -Tag, -InstallDir, -NoStartup

param(
  [string]$Owner = $env:OWNER
, [string]$Repo = $env:REPO
, [string]$Tag = $env:TAG
, [string]$InstallDir = $env:INSTALL_DIR
, [switch]$NoStartup
)

if ([string]::IsNullOrWhiteSpace($Owner)) { $Owner = 'axlroden' }
if ([string]::IsNullOrWhiteSpace($Repo)) { $Repo = 'volume-speeder' }
if ([string]::IsNullOrWhiteSpace($InstallDir)) { $InstallDir = Join-Path $env:LOCALAPPDATA 'Programs\VolumeSpeeder' }

function Get-LatestTag {
  param([string]$Owner,[string]$Repo)
  $url = "https://api.github.com/repos/$Owner/$Repo/releases/latest"
  $headers = @{}
  if ($env:GH_TOKEN) { $headers['Authorization'] = "Bearer $($env:GH_TOKEN)" }
  $res = Invoke-RestMethod -UseBasicParsing -Uri $url -Headers $headers
  return $res.tag_name
}

function Find-AssetUrl {
  param([string]$Owner,[string]$Repo,[string]$Tag,[string]$Pattern)
  $url = "https://api.github.com/repos/$Owner/$Repo/releases/tags/$Tag"
  $headers = @{}
  if ($env:GH_TOKEN) { $headers['Authorization'] = "Bearer $($env:GH_TOKEN)" }
  $res = Invoke-RestMethod -UseBasicParsing -Uri $url -Headers $headers
  foreach ($asset in $res.assets) {
    if ($asset.browser_download_url -match $Pattern) { return $asset.browser_download_url }
  }
  return $null
}

$ErrorActionPreference = 'Stop'
$tag = if ($Tag) { $Tag } else { Get-LatestTag -Owner $Owner -Repo $Repo }
if (-not $tag) { throw 'Could not determine tag' }

  $setupUrl = Find-AssetUrl -Owner $Owner -Repo $Repo -Tag $tag -Pattern '-windows-amd64-setup\.exe$'
$temp = New-Item -ItemType Directory -Force -Path (Join-Path $env:TEMP "volboost.$([guid]::NewGuid().ToString())")
try {
  if ($setupUrl) {
    $setupPath = Join-Path $temp 'setup.exe'
    Write-Host "Downloading $setupUrl"
    Invoke-WebRequest -UseBasicParsing -Uri $setupUrl -OutFile $setupPath
    # Silent install to InstallDir if NSIS; if portable exe, copy instead
    # Try NSIS standard silent flags
    $args = "/VERYSILENT", "/NORESTART"
    try {
      Start-Process -FilePath $setupPath -ArgumentList $args -Wait -NoNewWindow
      Write-Host 'Installed via setup.exe'
    } catch {
      # Fallback: treat as portable exe
      $exeDir = New-Item -ItemType Directory -Force -Path $InstallDir
  Copy-Item $setupPath -Destination (Join-Path $exeDir.FullName 'volume-speeder-gui.exe') -Force
    }
  } else {
    # Fallback to raw exe
  $exeUrl = Find-AssetUrl -Owner $Owner -Repo $Repo -Tag $tag -Pattern '-windows-amd64\.exe$'
    if (-not $exeUrl) { throw 'No Windows asset found' }
    $exeDir = New-Item -ItemType Directory -Force -Path $InstallDir
  $exePath = Join-Path $exeDir.FullName 'volume-speeder-gui.exe'
    Write-Host "Downloading $exeUrl"
    Invoke-WebRequest -UseBasicParsing -Uri $exeUrl -OutFile $exePath
  }

  # Start Menu shortcut
  $shell = New-Object -ComObject WScript.Shell
  $startMenu = [Environment]::GetFolderPath('Programs')
  $lnkDir = Join-Path $startMenu 'Volume Speeder'
  New-Item -ItemType Directory -Force -Path $lnkDir | Out-Null
  $shortcut = $shell.CreateShortcut((Join-Path $lnkDir 'Volume Speeder.lnk'))
  $shortcut.TargetPath = (Join-Path $InstallDir 'volume-speeder-gui.exe')
  $shortcut.WorkingDirectory = $InstallDir
  $shortcut.Save()

  # Startup (optional)
  if (-not $NoStartup) {
  New-ItemProperty -Path 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Run' -Name 'VolumeSpeeder' -Value (Join-Path $InstallDir 'volume-speeder-gui.exe') -PropertyType String -Force | Out-Null
  }

  Write-Host "Install complete."
} finally {
  Remove-Item -Recurse -Force $temp -ErrorAction SilentlyContinue
}
