[CmdletBinding()]
param(
    [switch]$RemoveData,
    [switch]$WithImageRecognition,
    [string]$EnvFile = $env:GOPHERAI_LOCAL_ENV_FILE
)

$ErrorActionPreference = "Stop"
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

if ([string]::IsNullOrWhiteSpace($EnvFile)) {
    $EnvFile = ".env"
}
if (-not (Test-Path -LiteralPath $EnvFile)) {
    throw "Environment file not found: $EnvFile"
}

$composeArgs = @("compose", "--env-file", $EnvFile, "-f", "docker-compose.yml")
if ($WithImageRecognition) {
    $composeArgs += @("-f", "docker-compose.models.yml")
}
$composeArgs += "down"
if ($RemoveData) {
    $composeArgs += "--volumes"
}
& docker @composeArgs
