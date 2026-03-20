$certPath = "certs/localhost.pem"
$keyPath = "certs/localhost-key.pem"
$certPathCopy="backend/server.crt"
$keyPathCopy="backend/server.key"
$dirs = @("certs", "uploads")

foreach ($dir in $dirs) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir | Out-Null
    }
}

icacls uploads /grant "Everyone:(OI)(CI)F" | Out-Null

New-Item -ItemType Directory -Force certs | Out-Null

if ((Test-Path $certPath) -and (Test-Path $keyPath)) {
    Write-Host "TLS certificate already exists, skipping generation."

    Copy-Item $certPath $certPathCopy -Force
    Copy-Item $keyPath $keyPathCopy -Force
    Write-Host "Certificates copied to backend."

    exit
}

Write-Host ""
Write-Host "Choose TLS setup method:"
Write-Host ""
Write-Host "1) mkcert (recommended, trusted by browser)"
Write-Host "2) self-signed certificate"
Write-Host ""

$choice = Read-Host "Selection [1/2]"

if ($choice -eq "1") {

    if (!(Get-Command mkcert -ErrorAction SilentlyContinue)) {
        Write-Host ""
        Write-Host "mkcert is not installed."
        Write-Host ""
        Write-Host "Install mkcert:"
        Write-Host ""
        Write-Host "Windows (Chocolatey):"
        Write-Host "  choco install mkcert"
        Write-Host ""
        Write-Host "Windows (Scoop):"
        Write-Host "  scoop bucket add extras"
        Write-Host "  scoop install mkcert"
        Write-Host ""
        exit
    }

    Write-Host "Installing local CA..."
    mkcert -install

    Write-Host "Generating trusted certificate..."

    mkcert `
        -cert-file certs/localhost.pem `
        -key-file certs/localhost-key.pem `
        localhost 127.0.0.1 ::1

    Write-Host "Trusted TLS certificate created."

}
elseif ($choice -eq "2") {

    Write-Host "Generating self-signed certificate..."

    openssl req `
        -x509 `
        -nodes `
        -days 3650 `
        -newkey rsa:2048 `
        -keyout certs/localhost-key.pem `
        -out certs/localhost.pem `
        -subj "/CN=localhost"

    Write-Host "Self-signed TLS certificate created."

}
else {
    Write-Host "Invalid selection."
    exit 1
}
if ((Test-Path $certPath) -and (Test-Path $keyPath)) {
    Copy-Item $certPath $certPathCopy -Force
    Copy-Item $keyPath $keyPathCopy -Force
    Write-Host "Certificates copied to backend."
}
