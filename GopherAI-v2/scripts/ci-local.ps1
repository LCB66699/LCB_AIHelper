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

if (Get-Command go -ErrorAction SilentlyContinue) {
    go test ./...
}
else {
    docker run --rm -v "${PWD}:/src" -w /src golang:1.24-bookworm go test ./...
}
docker run --rm -v "${PWD}/common/mcp:/src" -w /src golang:1.25-bookworm go test ./...

docker compose --env-file .env config --quiet
if (-not $SkipDockerBuild) {
    docker compose --env-file .env build api frontend mcp
}
