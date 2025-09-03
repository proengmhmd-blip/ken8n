Param(
  [Parameter(ValueFromRemainingArguments = $true)]
  [string[]]$Args
)

$ErrorActionPreference = "Stop"

# Launch ken8n-coder from repo root using Bun (dev mode)
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

# Ensure Bun is available on PATH for Windows
if (-not (Get-Command bun -ErrorAction SilentlyContinue)) {
  $bunDefault = Join-Path $HOME ".bun\\bin\\bun.exe"
  if (Test-Path $bunDefault) {
    $env:PATH = (Split-Path $bunDefault) + ";" + $env:PATH
  } else {
    Write-Error "Bun not found. Please install Bun: https://bun.sh"
    exit 1
  }
}

& bun run --conditions=development packages/ken8n-coder/src/index.ts @Args
exit $LASTEXITCODE

