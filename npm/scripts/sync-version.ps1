$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$npmDir = Resolve-Path (Join-Path $scriptDir "..")
$versionFile = Join-Path $npmDir "version.txt"

if (-not (Test-Path $versionFile)) {
    Write-Error "version.txt is missing"
    exit 1
}

$version = (Get-Content $versionFile -Raw).Trim()
if (-not $version) {
    Write-Error "version.txt is empty"
    exit 1
}

Write-Host "Syncing version $version to all packages..."

$mainPackagePath = Join-Path $npmDir "@claudex-windows\cli\package.json"
$mainPackage = Get-Content $mainPackagePath -Raw | ConvertFrom-Json
$mainPackage.version = $version
$mainPackage.optionalDependencies."@claudex-windows/darwin-arm64" = $version
$mainPackage.optionalDependencies."@claudex-windows/darwin-x64" = $version
$mainPackage.optionalDependencies."@claudex-windows/linux-x64" = $version
$mainPackage.optionalDependencies."@claudex-windows/linux-arm64" = $version
$mainPackage.optionalDependencies."@claudex-windows/windows-x64" = $version
$mainPackage | ConvertTo-Json -Depth 10 | Set-Content $mainPackagePath
Write-Host "✓ Updated @claudex-windows/cli"

$platforms = @("darwin-arm64", "darwin-x64", "linux-x64", "linux-arm64", "windows-x64")
foreach ($platform in $platforms) {
    $pkgPath = Join-Path $npmDir "@claudex-windows\$platform\package.json"
    if (-not (Test-Path $pkgPath)) {
        Write-Warning "Skipping $platform (missing package.json)"
        continue
    }
    $pkg = Get-Content $pkgPath -Raw | ConvertFrom-Json
    $pkg.version = $version
    $pkg | ConvertTo-Json -Depth 10 | Set-Content $pkgPath
    Write-Host "✓ Updated @claudex-windows/$platform"
}

Write-Host "Version sync complete: $version"
