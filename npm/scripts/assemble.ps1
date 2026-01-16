$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Resolve-Path (Join-Path $scriptDir "..\..")
$platforms = @("darwin-arm64", "darwin-x64", "linux-x64", "linux-arm64", "windows-x64")

Write-Host "Assembling npm packages..."

foreach ($platform in $platforms) {
    $srcDir = Join-Path $rootDir "dist\$platform"
    $destDir = Join-Path $rootDir "npm\@claudex-windows\$platform\bin"

    if (Test-Path $srcDir) {
        New-Item -ItemType Directory -Force -Path $destDir | Out-Null
        Copy-Item (Join-Path $srcDir "claudex-windows*") $destDir -Force
        Write-Host "✓ Assembled $platform"
    } else {
        Write-Warning "Skipped $platform (not built)"
    }
}

Write-Host "Assembly complete!"
