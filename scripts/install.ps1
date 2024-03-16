# This script installs 'gdbuild' by downloading prebuilt binaries from the
# project's GitHub releases page. By default the latest version is installed,
# but a different release can be used instead by setting $GDBUILD_VERSION.
#
# The script will set up a 'gdbuild' cache at '%LOCALAPPDATA%/gdbuild'. This
# behavior can be customized by setting '$GDBUILD_HOME' prior to running the
# script. Existing Godot artifacts cached in a 'gdbuild' store won't be lost, but
# this script will overwrite any 'gdbuild' binary artifacts in '$GDBUILD_HOME/bin'.
#
# NOTE: Unlike the 'install.sh' counterpart, this script exclusively installs
# the 'gdbuild' binary for 64-bit Windows. If an alternative 'gdbuild' binary is
# required, follow the documentation for an alternative means of installation:
# https://github.com/coffeebeats/gdbuild/blob/v0.1.0/docs/installation.md # x-release-please-version

<#
.SYNOPSIS
  Install 'gdbuild' for compiling and exporting Godot projects.

.DESCRIPTION
  This script downloads the specified version of 'gdbuild' from GitHub, extracts
  its artifacts to the 'gdbuild' store ('$GDBUILD_HOME' or a default path), and then
  updates environment variables as needed.

.PARAMETER NoModifyPath
  Do not modify the $PATH environment variable.

.PARAMETER Version
  Install the specified version of 'gdbuild'.

.INPUTS
  None

.OUTPUTS
  $env:GDBUILD_HOME\bin\gdbuild.exe

.NOTES
  Version:        0.1.0 # x-release-please-version
  Author:         https://github.com/coffeebeats

.LINK
  https://github.com/coffeebeats/gdbuild
#>

# ------------------------------ Define: Params ------------------------------ #

Param (
  # NoModifyPath - if set, the user's $PATH variable won't be updated
  [Switch] $NoModifyPath = $False,

  # Version - override the specific version of 'gdbuild' to install
  [String] $Version = "v0.1.0" # x-release-please-version
)

# ------------------------- Function: Get-GdBuildHome ------------------------ #

# Returns the current value of the 'GDBUILD_HOME' environment variable or a
# default if unset.
Function Get-GdBuildHome() {
  if ([string]::IsNullOrEmpty($env:GDBUILD_HOME)) {
    return Join-Path -Path $env:LOCALAPPDATA -ChildPath "gdbuild"
  }

  return $env:GDBUILD_HOME
}

# ----------------------- Function: Get-GdBuildVersion ----------------------- #

Function Get-GdBuildVersion() {
  return "v" + $Version.TrimStart("v")
}

# --------------------- Function: Create-Temporary-Folder -------------------- #

# Creates a new temporary directory for downloading and extracting 'gdbuild'. The
# returned directory path will have a randomized suffix.
Function New-TemporaryFolder() {
  # Make a new temporary folder with a randomized suffix.
  return New-Item `
    -ItemType Directory `
    -Name "gdbuild-$([System.IO.Path]::GetFileNameWithoutExtension([System.IO.Path]::GetRandomFileName()))"`
    -Path $env:temp
}

# ------------------------------- Define: Store ------------------------------ #

$GdBuildHome = Get-GdBuildHome

Write-Host "info: setting 'GDBUILD_HOME' environment variable: ${GdBuildHome}"

[System.Environment]::SetEnvironmentVariable("GDBUILD_HOME", $GdBuildHome, "User")

# ------------------------------ Define: Version ----------------------------- #
  
$GdBuildVersion = Get-GdBuildVersion

$GdBuildArchive = "gdbuild-${GdBuildVersion}-windows-x86_64.zip"

# ----------------------------- Execute: Install ----------------------------- #
  
$GdBuildRepositoryURL = "https://github.com/coffeebeats/gdbuild"

# Install downloads 'gdbuild' and extracts its binaries into the store. It also
# updates environment variables as needed.
Function Install() {
  $GdBuildTempFolder = New-TemporaryFolder

  $GdBuildArchiveURL = "${GdBuildRepositoryURL}/releases/download/${GdBuildVersion}/${GdBuildArchive}"
  $GdBuildDownloadTo = Join-Path -Path $GdBuildTempFolder -ChildPath $GdBuildArchive

  $GdBuildHomeBinPath = Join-Path -Path $GdBuildHome -ChildPath "bin"

  try {
    Write-Host "info: installing version: '${GdBuildVersion}'"

    Invoke-WebRequest -URI $GdBuildArchiveURL -OutFile $GdBuildDownloadTo

    Microsoft.PowerShell.Archive\Expand-Archive `
      -Force `
      -Path $GdBuildDownloadTo `
      -DestinationPath $GdBuildHomeBinPath
  
    if (!($NoModifyPath)) {
      $PathParts = [System.Environment]::GetEnvironmentVariable("PATH", "User").Trim(";") -Split ";"
      $PathParts = $PathParts.where{ $_ -ne $GdBuildHomeBinPath }
      $PathParts = $PathParts + $GdBuildHomeBinPath

      Write-Host "info: updating 'PATH' environment variable: ${GdBuildHomeBinPath}"

      [System.Environment]::SetEnvironmentVariable("PATH", $($PathParts -Join ";"), "User")
    }

    Write-Host "info: sucessfully installed executables:`n"
    Write-Host "  gdbuild.exe: $(Join-Path -Path $GdBuildHomeBinPath -ChildPath "gdbuild.exe")"
    Write-Host "  godot.exe: $(Join-Path -Path $GdBuildHomeBinPath -ChildPath "godot.exe")"
  }
  catch {
    Write-Host "error: failed to install 'gdbuild': ${_}"
  }
  finally {
    Write-Host "`ninfo: cleaning up downloads: ${GdBuildTempFolder}"

    Remove-Item -Recurse $GdBuildTempFolder
  }
}

Install
