# Requires Windows PowerShell 5.1 or PowerShell 7 on Windows.
& {
    $ErrorActionPreference = 'Stop'
    if ($env:OS -ne 'Windows_NT') {
        throw 'This installer supports Windows only. Use install.sh on Linux or macOS.'
    }
    $architecture = $env:PROCESSOR_ARCHITEW6432
    if (-not $architecture) { $architecture = $env:PROCESSOR_ARCHITECTURE }
    switch ($architecture) {
        'AMD64' { $arch = 'amd64' }
        'ARM64' { $arch = 'arm64' }
        default { throw "Unsupported Windows architecture: $architecture. Supported: AMD64 and ARM64." }
    }

    $installDir = $env:GOCSV_INSTALL_DIR
    if (-not $installDir) { $installDir = Join-Path $HOME '.local\bin' }
    $installProvider = $null
    $installDrive = $null
    $installDir = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath($installDir, [ref]$installProvider, [ref]$installDrive)
    if ($installProvider.Name -ne 'FileSystem') {
        throw 'GOCSV_INSTALL_DIR must be a filesystem path.'
    }
    $archive = "gocsv-windows-$arch"
    $version = $env:GOCSV_VERSION
    if (-not $version) { $version = 'latest' }
    if ($version -ceq 'latest') {
        $releasePath = 'latest/download'
    } else {
        if ($version -cnotmatch '\A[A-Za-z0-9._+-]+\z' -or $version -in '.', '..') {
            throw 'Set GOCSV_VERSION to a release tag such as v1.0.7.'
        }
        $releasePath = "download/$version"
    }
    $url = "https://github.com/aotimme/gocsv/releases/$releasePath/$archive.zip"
    $tempDir = Join-Path ([IO.Path]::GetTempPath()) ([Guid]::NewGuid().ToString())
    $oldProtocol = [Net.ServicePointManager]::SecurityProtocol
    try {
        [Net.ServicePointManager]::SecurityProtocol = $oldProtocol -bor [Net.SecurityProtocolType]::Tls12
        New-Item -ItemType Directory -Path $tempDir | Out-Null
        $zipPath = Join-Path $tempDir 'release.zip'
        Write-Host "Downloading $archive ($version)..."
        Invoke-WebRequest -UseBasicParsing -Uri $url -OutFile $zipPath
        Expand-Archive -LiteralPath $zipPath -DestinationPath $tempDir
        $binary = Join-Path (Join-Path $tempDir $archive) 'gocsv.exe'
        if (-not (Test-Path -LiteralPath $binary -PathType Leaf)) {
            throw 'The release archive does not contain the expected gocsv.exe binary.'
        }
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
        $destination = Join-Path $installDir 'gocsv.exe'
        if (Test-Path -LiteralPath $destination -PathType Container) {
            throw "The destination $destination is a directory."
        }
        Copy-Item -LiteralPath $binary -Destination $destination -Force

        $environmentKey = [Microsoft.Win32.Registry]::CurrentUser.CreateSubKey('Environment')
        try {
            if (@($environmentKey.GetValue('Path', '') -split ';') -notcontains $installDir) {
                $userPath = $environmentKey.GetValue('Path', $null, 'DoNotExpandEnvironmentNames')
                $pathKind = if ($null -eq $userPath) { 'String' } else { $environmentKey.GetValueKind('Path') }
                $separator = if ($userPath -and -not $userPath.EndsWith(';')) { ';' } else { '' }

                # A raw registry write preserves expandable entries but does not notify Explorer.
                if (-not ('GoCsv.EnvironmentChange' -as [type])) {
                    Add-Type -Namespace GoCsv -Name EnvironmentChange -MemberDefinition @'
[System.Runtime.InteropServices.DllImport("user32.dll", CharSet = System.Runtime.InteropServices.CharSet.Unicode)]
private static extern System.IntPtr SendMessageTimeout(System.IntPtr window, uint message,
    System.IntPtr wParam, string lParam, uint flags, uint timeout, out System.IntPtr result);

public static void Notify() {
    System.IntPtr result;
    SendMessageTimeout(new System.IntPtr(0xffff), 0x001a, System.IntPtr.Zero, "Environment", 2, 1000, out result);
}
'@
                }
                $environmentKey.SetValue('Path', "$userPath$separator$installDir", $pathKind)
                [GoCsv.EnvironmentChange]::Notify()
            }
        }
        finally {
            $environmentKey.Dispose()
        }
        if (@($env:Path -split ';') -notcontains $installDir) {
            $env:Path = "$installDir;$env:Path"
        }
        Write-Host "GoCSV installed to $installDir\gocsv.exe"
        Write-Host 'Open a new terminal and run: gocsv help'
    }
    finally {
        [Net.ServicePointManager]::SecurityProtocol = $oldProtocol
        if (Test-Path -LiteralPath $tempDir) { Remove-Item -LiteralPath $tempDir -Recurse -Force }
    }
}
