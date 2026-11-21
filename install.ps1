# ============================================================
# DockSea — Instalador / Atualizador (Windows .exe)
# Uso:  irm https://raw.githubusercontent.com/Josemarcio15/dockSea/dev/install.ps1 | iex
# ============================================================
$ErrorActionPreference = 'Stop'

$Repo = "Josemarcio15/dockSea"
$Api = "https://api.github.com/repos/$Repo/releases"

function Write-Step {
    param([string]$Message)
    Write-Host "[dockSea] " -ForegroundColor Cyan -NoNewline
    Write-Host $Message
}

function Write-Success {
    param([string]$Message)
    Write-Host "[dockSea] " -ForegroundColor Green -NoNewline
    Write-Host $Message -ForegroundColor Green
}

Write-Step "Consultando versao mais recente em github.com/$Repo..."
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

try {
    $releases = Invoke-RestMethod -Uri $Api -Headers @{ "User-Agent" = "docksea-installer" }
} catch {
    Write-Host "Erro ao consultar API do GitHub: $_" -ForegroundColor Red
    exit 1
}

if (-not $releases -or $releases.Count -eq 0) {
    Write-Host "Nenhum release encontrado no repositorio." -ForegroundColor Red
    exit 1
}

$latestRelease = $releases[0]
$tag = $latestRelease.tag_name
Write-Step "Versao mais recente encontrada: $tag"

$asset = $latestRelease.assets | Where-Object { 
    $_.name -match '\.exe$' -and $_.name -notmatch 'sig$' 
} | Select-Object -First 1

if (-not $asset) {
    Write-Host "Nenhum instalador .exe disponivel na versao $tag ainda." -ForegroundColor Yellow
    exit 1
}

$downloadUrl = $asset.browser_download_url
$fileName = $asset.name
$tempPath = Join-Path $env:TEMP $fileName

Write-Step "Baixando $fileName..."
$ProgressPreference = 'SilentlyContinue'
if (Get-Command curl.exe -ErrorAction SilentlyContinue) {
    & curl.exe -fL --retry 3 --retry-delay 2 --progress-bar $downloadUrl -o $tempPath
} else {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempPath -UseBasicParsing
}

Write-Step "Iniciando instalador..."
Start-Process -FilePath $tempPath -Wait

if (Test-Path $tempPath) {
    Remove-Item $tempPath -Force -ErrorAction SilentlyContinue
}

Write-Success "DockSea $tag instalado com sucesso! 🚀"
