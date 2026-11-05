# Offline integration tests. Run with both Windows PowerShell 5.1 and PowerShell 7.
$ErrorActionPreference = 'Stop'
$installer = Join-Path $PSScriptRoot '..\install.ps1'
$testRoot = Join-Path ([IO.Path]::GetTempPath()) ('gocsv tests ' + [Guid]::NewGuid())
$originalUserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$originalPath = $env:Path
$originalInstallDir = $env:GOCSV_INSTALL_DIR
$originalVersion = $env:GOCSV_VERSION
$originalArchitecture = $env:PROCESSOR_ARCHITECTURE
$originalWowArchitecture = $env:PROCESSOR_ARCHITEW6432

# A file invocation gives the installer a separate script scope. Share state
# through an object that the mock can read and update in either invocation.
$installerTestState = @{
    Fixture = Join-Path $testRoot 'fixture.zip'
    ExpectedReleasePath = 'latest/download'
    Mode = 'success'
    DownloadDir = $null
}

function Assert($condition, $message) {
    if (-not $condition) { throw $message }
}

# Only the network boundary is mocked; ZIP extraction, installation, and PATH
# persistence use the real Windows APIs. No Pester dependency is required.
function Invoke-WebRequest {
    param([switch]$UseBasicParsing, $Uri, $OutFile)
    Assert $UseBasicParsing 'UseBasicParsing is required for Windows PowerShell.'
    $expectedUrl = "https://github.com/aotimme/gocsv/releases/$($installerTestState.ExpectedReleasePath)/gocsv-windows-$($installerTestState.ExpectedArch).zip"
    Assert ($Uri -eq $expectedUrl) "Expected download URL '$expectedUrl', got '$Uri'."
    $installerTestState.DownloadDir = Split-Path $OutFile
    if ($installerTestState.Mode -eq 'download') { throw 'Simulated download failure' }
    if ($installerTestState.Mode -eq 'corrupt') {
        Set-Content -LiteralPath $OutFile -Value 'not a zip'
    } else {
        Copy-Item -LiteralPath $installerTestState.Fixture -Destination $OutFile
    }
}

function Expect-Failure {
    $failed = $false
    try { & $installer } catch { $failed = $true }
    Assert $failed 'The installer should have failed.'
    Assert ((Get-Content -LiteralPath $installerTestState.Installed -Raw).Trim() -eq 'existing') 'Existing installation was changed.'
    if ($installerTestState.DownloadDir) {
        Assert (-not (Test-Path -LiteralPath $installerTestState.DownloadDir)) 'Temporary files were not cleaned up.'
    }
}

try {
    New-Item -ItemType Directory -Path $testRoot | Out-Null
    $env:GOCSV_INSTALL_DIR = Join-Path $testRoot 'install with spaces'
    $installerTestState.Installed = Join-Path $env:GOCSV_INSTALL_DIR 'gocsv.exe'
    $env:PROCESSOR_ARCHITEW6432 = ''
    $env:GOCSV_VERSION = $null
    foreach ($arch in @('amd64', 'arm64')) {
        $installerTestState.ExpectedArch = $arch
        $env:PROCESSOR_ARCHITECTURE = $arch.ToUpperInvariant()
        $folder = Join-Path $testRoot "gocsv-windows-$arch"
        New-Item -ItemType Directory -Path $folder | Out-Null
        Set-Content -LiteralPath (Join-Path $folder 'gocsv.exe') -Value 'fixture executable'
        Compress-Archive -LiteralPath $folder -DestinationPath $installerTestState.Fixture -Force
        # Exercise the README's download-and-evaluate execution style too.
        Invoke-Expression (Get-Content -LiteralPath $installer -Raw)
        Assert ((Get-Content -LiteralPath $installerTestState.Installed -Raw).Trim() -eq 'fixture executable') 'Binary not installed.'
        Assert (-not (Test-Path -LiteralPath $installerTestState.DownloadDir)) 'Temporary files were not cleaned up.'
        Assert (@($env:Path -split ';') -contains $env:GOCSV_INSTALL_DIR) 'Session PATH was not updated.'
        $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        Assert (@($userPath -split ';') -contains $env:GOCSV_INSTALL_DIR) 'User PATH was not updated.'
        & $installer
        Assert ([Environment]::GetEnvironmentVariable('Path', 'User') -eq $userPath) 'Repeated install changed user PATH.'
    }
    $env:GOCSV_VERSION = 'v1.0.6'
    $installerTestState.ExpectedReleasePath = 'download/v1.0.6'
    # PowerShell's location can differ from the process working directory.
    $originalWorkingDirectory = [Environment]::CurrentDirectory
    $absoluteInstallDir = $env:GOCSV_INSTALL_DIR
    Push-Location -LiteralPath $testRoot
    try {
        [Environment]::CurrentDirectory = Join-Path $testRoot 'gocsv-windows-arm64'
        $env:GOCSV_INSTALL_DIR = 'relative install'
        $expectedInstallDir = Join-Path $testRoot 'relative install'
        & $installer
        Assert (Test-Path -LiteralPath (Join-Path $expectedInstallDir 'gocsv.exe') -PathType Leaf) 'Relative install did not use the PowerShell location.'
        $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        Assert (@($userPath -split ';') -contains $expectedInstallDir) 'User PATH is missing the resolved install directory.'
    }
    finally {
        [Environment]::CurrentDirectory = $originalWorkingDirectory
        Pop-Location
        $env:GOCSV_INSTALL_DIR = $absoluteInstallDir
    }
    $env:GOCSV_VERSION = $null
    $installerTestState.ExpectedReleasePath = 'latest/download'
    # A 32-bit PowerShell process on 64-bit Windows must use the native binary.
    $env:PROCESSOR_ARCHITECTURE = 'x86'
    $env:PROCESSOR_ARCHITEW6432 = 'ARM64'
    & $installer
    $env:PROCESSOR_ARCHITEW6432 = ''
    Set-Content -LiteralPath $installerTestState.Installed -Value 'existing'
    Expect-Failure # Unsupported 32-bit Windows.
    $env:PROCESSOR_ARCHITECTURE = 'ARM64'
    foreach ($failure in @('download', 'corrupt', 'missing')) {
        $installerTestState.Mode = $failure
        if ($failure -eq 'missing') {
            $folder = Join-Path $testRoot 'unexpected'
            New-Item -ItemType Directory -Path $folder | Out-Null
            Set-Content -LiteralPath (Join-Path $folder 'other.txt') -Value 'missing executable'
            Compress-Archive -LiteralPath $folder -DestinationPath $installerTestState.Fixture -Force
        }
        Expect-Failure
    }
    Write-Host 'All Windows installer tests passed.'
}
finally {
    [Environment]::SetEnvironmentVariable('Path', $originalUserPath, 'User')
    $env:Path = $originalPath
    $env:GOCSV_INSTALL_DIR = $originalInstallDir
    $env:GOCSV_VERSION = $originalVersion
    $env:PROCESSOR_ARCHITECTURE = $originalArchitecture
    $env:PROCESSOR_ARCHITEW6432 = $originalWowArchitecture
    Remove-Item Function:\Invoke-WebRequest
    if (Test-Path -LiteralPath $testRoot) { Remove-Item -LiteralPath $testRoot -Recurse -Force }
}
