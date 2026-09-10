[CmdletBinding()]
param(
    [switch]$SkipBuild,
    [switch]$SkipChecks,
    [switch]$WithImageRecognition,
    [int]$TimeoutSeconds = 120
)

$ErrorActionPreference = "Stop"
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

function New-LocalSecret {
    $bytes = New-Object byte[] 32
    $rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
    try {
        $rng.GetBytes($bytes)
    }
    finally {
        $rng.Dispose()
    }
    return [Convert]::ToBase64String($bytes).TrimEnd('=').Replace('+', '-').Replace('/', '_')
}

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    throw "Docker Desktop and Docker Compose are required for local deployment."
}
if (-not (Test-Path .env)) {
    $localEnv = @(
        "MYSQL_DATABASE=GopherAI",
        "MYSQL_USER=gopherai",
        "MYSQL_PASSWORD=$(New-LocalSecret)",
        "MYSQL_ROOT_PASSWORD=$(New-LocalSecret)",
        "RABBITMQ_USER=gopherai",
        "RABBITMQ_PASSWORD=$(New-LocalSecret)",
        "JWT_KEY=$(New-LocalSecret)",
        "OPENAI_API_KEY=",
        "OPENAI_MODEL_NAME=",
        "OPENAI_BASE_URL="
    )
    Set-Content -LiteralPath .env -Value $localEnv -Encoding ascii
    Write-Host "Created a local .env with generated credentials. Add AI provider settings before using chat or RAG."
}

$requiredKeys = @("MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_ROOT_PASSWORD", "RABBITMQ_USER", "RABBITMQ_PASSWORD", "JWT_KEY")
$environment = @{}
Get-Content .env | ForEach-Object {
    if ($_ -match '^\s*([^#=\s]+)\s*=\s*(.*)\s*$') {
        $environment[$matches[1]] = $matches[2]
    }
}
foreach ($key in $requiredKeys) {
    if (-not $environment.ContainsKey($key) -or [string]::IsNullOrWhiteSpace($environment[$key]) -or $environment[$key] -match '^replace-with') {
        throw "Set a non-placeholder value for $key in .env."
    }
}

docker compose --env-file .env config --quiet
if (-not $SkipChecks) {
    & "$PSScriptRoot\ci-local.ps1" -SkipDockerBuild:$SkipBuild
    if ($LASTEXITCODE -ne 0) {
        throw "Local CI failed; deployment was not started."
    }
}

$composeArgs = @("compose", "--env-file", ".env", "-f", "docker-compose.yml")
if ($WithImageRecognition) {
    $composeArgs += @("-f", "docker-compose.models.yml")
}
$composeArgs += @("up", "--detach")
if (-not $SkipBuild) {
    $composeArgs += "--build"
}
& docker @composeArgs
if ($LASTEXITCODE -ne 0) {
    throw "Docker Compose deployment failed."
}

$deadline = (Get-Date).AddSeconds($TimeoutSeconds)
do {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:9090/healthz" -UseBasicParsing -TimeoutSec 5
        if ($response.StatusCode -eq 200 -and $response.Content -match '"status"\s*:\s*"ok"') {
            Write-Host "GopherAI is ready: http://localhost:8080"
            exit 0
        }
    }
    catch {
        Start-Sleep -Seconds 3
    }
} while ((Get-Date) -lt $deadline)

docker compose --env-file .env logs --tail 100 api frontend
throw "GopherAI did not become healthy within $TimeoutSeconds seconds."
