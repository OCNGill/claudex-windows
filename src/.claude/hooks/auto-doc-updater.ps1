$hooksBin = $env:CLAUDEX_WINDOWS_HOOKS_BIN
if (-not $hooksBin) { $hooksBin = $env:CLAUDEX_HOOKS_BIN }
if (-not $hooksBin) { $hooksBin = "claudex-windows-hooks" }
& $hooksBin "auto-doc"
