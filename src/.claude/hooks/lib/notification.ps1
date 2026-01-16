function Send-Notification {
    param(
        [string]$Title,
        [string]$Message,
        [string]$Sound = ""
    )

    try {
        $shell = New-Object -ComObject WScript.Shell
        [void]$shell.Popup($Message, 5, $Title, 0x40)
    } catch {
        Write-Warning "Notification unavailable: $($_.Exception.Message)"
    }
}

function Invoke-Speech {
    param(
        [string]$Message,
        [string]$Voice = ""
    )

    try {
        Add-Type -AssemblyName System.Speech | Out-Null
        $synth = New-Object System.Speech.Synthesis.SpeechSynthesizer
        if ($Voice) {
            $synth.SelectVoice($Voice)
        }
        $synth.Speak($Message)
    } catch {
        Write-Warning "Voice synthesis unavailable: $($_.Exception.Message)"
    }
}

function Send-AgentCompleteNotification {
    param(
        [string]$AgentId,
        [string]$SessionName
    )

    $enableNotifications = $env:CLAUDEX_WINDOWS_NOTIFICATIONS_ENABLED
    if (-not $enableNotifications) { $enableNotifications = $env:CLAUDEX_NOTIFICATIONS_ENABLED }
    if (-not $enableNotifications) { $enableNotifications = "true" }

    $enableVoice = $env:CLAUDEX_WINDOWS_VOICE_ENABLED
    if (-not $enableVoice) { $enableVoice = $env:CLAUDEX_VOICE_ENABLED }
    if (-not $enableVoice) { $enableVoice = "false" }

    $formattedName = $SessionName -replace "-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", ""
    $formattedName = ($formattedName -replace "-", " ") -split " " | ForEach-Object { if ($_) { $_.Substring(0,1).ToUpper() + $_.Substring(1).ToLower() } } | ForEach-Object { $_ } -join " "

    if ($enableNotifications -eq "true") {
    Send-Notification -Title $formattedName -Message "Agent complete"
    }

    if ($enableVoice -eq "true") {
    Invoke-Speech -Message "Agent complete"
    }
}
