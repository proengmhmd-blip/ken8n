Param(
  [string]$Version = "2.4.5"
)

$ErrorActionPreference = "Stop"

$App = "ken8n-coder"
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "x64" }
$Os = "windows"
$FileName = "$App-$Os-$Arch.zip"
$InstallDir = Join-Path $HOME ".ken8n-coder\\bin"
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

$Url = "https://github.com/kenkaiii/ken8n-coder/releases/download/v$Version/$FileName"
$TmpDir = Join-Path $env:TEMP "ken8ncodertmp"
New-Item -ItemType Directory -Force -Path $TmpDir | Out-Null

$ZipPath = Join-Path $TmpDir $FileName
Write-Host "Downloading ken8n-coder v$Version..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $Url -OutFile $ZipPath -UseBasicParsing

Write-Host "Extracting..." -ForegroundColor Cyan
Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

$Root = Join-Path $TmpDir "$App-$Os-$Arch"
$BinExe = Join-Path $Root "bin\\ken8n-coder.exe"
if (Test-Path $BinExe) {
  Move-Item -Force $BinExe (Join-Path $InstallDir "ken8n-coder.exe")
}

# Optional assets
if (Test-Path (Join-Path $Root "validation-scripts")) {
  $Dest = Join-Path $HOME ".ken8n-coder\\validation-scripts"
  New-Item -ItemType Directory -Force -Path $Dest | Out-Null
  Copy-Item -Recurse -Force (Join-Path $Root "validation-scripts\\*") $Dest
}
if (Test-Path (Join-Path $Root "deploy-script")) {
  $Dest = Join-Path $HOME ".ken8n-coder\\deploy-script"
  New-Item -ItemType Directory -Force -Path $Dest | Out-Null
  Copy-Item -Recurse -Force (Join-Path $Root "deploy-script\\*") $Dest
}

Remove-Item -Recurse -Force $TmpDir

Write-Host "ken8n-coder installed to: $InstallDir" -ForegroundColor Green
Write-Host "Ensure this directory is on your PATH. Example:" -ForegroundColor Yellow
Write-Host "  setx PATH \"%PATH%;$InstallDir\"" -ForegroundColor Yellow

