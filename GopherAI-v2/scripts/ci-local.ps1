[CmdletBinding()]
param(
    [switch]$SkipDockerBuild
)

$ErrorActionPreference = "Stop"

function Require-Command([string]$Name) {
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "Required command '$Name' was not found."
    }
}

Require-Command docker

$moduleProxy = $env:GOPROXY
if ([string]::IsNullOrWhiteSpace($moduleProxy)) {
    $moduleProxy = "https://goproxy.cn,direct"
}

if (Get-Command go -ErrorAction SilentlyContinue) {
    $previousProxy = $env:GOPROXY
    $env:GOPROXY = $moduleProxy
    try {
        go test ./...
    }
    finally {
        $env:GOPROXY = $previousProxy
    }
}
else {
    docker run --rm -e "GOPROXY=$moduleProxy" -v "${PWD}:/src" -w /src golang:1.24-bookworm go test ./...
}
docker run --rm -e "GOPROXY=$moduleProxy" -v "${PWD}/common/mcp:/src" -w /src golang:1.25-bookworm go test ./...

docker compose --env-file .env config --quiet
if (-not $SkipDockerBuild) {
    docker compose --env-file .env build api frontend mcp
}
