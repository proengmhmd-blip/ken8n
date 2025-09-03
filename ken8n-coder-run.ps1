Param(
  [Parameter(ValueFromRemainingArguments = $true)]
  [string[]]$Args
)

$ErrorActionPreference = "Stop"

# Capture the directory where this script is placed
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

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

