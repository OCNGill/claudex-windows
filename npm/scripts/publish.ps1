$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$npmDir = Resolve-Path (Join-Path $scriptDir "..")

npm whoami | Out-Null
if ($LASTEXITCODE -ne 0) {
    Write-Error "Not logged in to npm. Run 'npm login' first."
    exit 1
}

& (Join-Path $scriptDir "sync-version.ps1")
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host ""
Write-Host "Publishing platform packages first (main package depends on them)..."

$platforms = @("darwin-arm64", "darwin-x64", "linux-x64", "linux-arm64", "windows-x64")
foreach ($platform in $platforms) {
    $pkgDir = Join-Path $npmDir "@claudex-windows\$platform"
    if (-not (Test-Path $pkgDir)) {
        Write-Warning "Skipping @claudex-windows/$platform (missing)"
        continue
    }
    Write-Host "Publishing @claudex-windows/$platform..."
    Push-Location $pkgDir
    npm publish --access public
    Pop-Location
}

Write-Host ""
Write-Host "Publishing main package (@claudex-windows/cli)..."
$cliDir = Join-Path $npmDir "@claudex-windows\cli"
if (Test-Path $cliDir) {
    Push-Location $cliDir
    npm publish --access public
    Pop-Location
}

Write-Host ""
Write-Host "✓ All packages published successfully!"
