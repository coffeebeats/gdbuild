# **Installation**

The easiest way to install `gdbuild` is by using the pre-built binaries. These can be manually downloaded and configured, but automated installation scripts are provided and recommended.

Alternatively, you can install `gdbuild` from source using the latest supported version of [Go](https://go.dev/). See [Install from source](#install-from-source) for more details.

## **Pre-built binaries (recommended)**

> ⚠️ **WARNING:** It's good practice to inspect an installation script prior to execution. The scripts are included in this repository and can be reviewed prior to use.

### **Linux/MacOS**

```sh
curl https://raw.githubusercontent.com/coffeebeats/gdbuild/main/scripts/install.sh | sh
```

### **Windows**

#### **Git BASH for Windows**

If you're using [Git BASH for Windows](https://gitforwindows.org/) follow the recommended [Linux/MacOS](#linuxmacos) instructions.

#### **Powershell**

> ❕ **NOTE:** In order to run scripts in PowerShell, the [execution policy](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_execution_policies) must _not_ be `Restricted`. Consider running the following command
> if you encounter `UnauthorizedAccess` errors when following these instructions. See [Set-ExecutionPolicy](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.security/set-executionpolicy) documentation for details.
>
> ```sh
> Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope LocalMachine
> ```

```sh
Invoke-WebRequest `
    -UseBasicParsing `
    -Uri "https://raw.githubusercontent.com/coffeebeats/gdbuild/main/scripts/install.ps1" `
    -OutFile "./install-gdbuild.ps1"; `
    &"./install-gdbuild.ps1"; `
    Remove-Item "./install-gdbuild.ps1"
```

### **Manual download**

> ❕ **NOTE:** The instructions below provide `bash`-specific commands for a _Linux_-based system. While these won't work in _Powershell_, the process will be similar.

1. Download a prebuilt binary from the corresponding GitHub release. Set `VERSION`, `OS`, and `ARCH` to the desired values.

    ```sh
    VERSION=0.0.0 OS=linux ARCH=x86_64; \
    curl -LO https://github.com/coffeebeats/gdbuild/releases/download/v$VERSION/gdbuild-$VERSION-$OS-$ARCH.tar.gz
    ```

2. Extract the downloaded archive. To customize the `gdbuild` install location, set `GDBUILD_HOME` to the desired location (defaults to `$HOME/.gdbuild` on Linux/MacOS).

    ```sh
    GDBUILD_HOME=$HOME/.gdbuild; \
    mkdir -p $GDBUILD_HOME/bin && \
    tar -C $GDBUILD_HOME/bin -xf gdbuild-$VERSION-$OS-$ARCH.tar.gz
    ```

3. Export the `GDBUILD_HOME` environment variable and add `$GDBUILD_HOME/bin` to `PATH`. Add the following to your shell profile script (e.g. in `.bashrc`, `.zshenv`, `.profile`, or something similar).

    ```sh
    export GDBUILD_HOME="$HOME/.gdbuild"
    export PATH="$GDBUILD_HOME/bin:$PATH"
    ```

## **Install from source**

`gdbuild` is a Go project and can be installed using `go install`. This option is not recommended as it requires having the Go toolchain installed, it's slower than downloading a prebuilt binary, and there may be instability due to using a different version of Go than it was developed with.

```sh
go install github.com/coffeebeats/gdbuild/cmd/gdbuild@latest
```

Once `gdbuild` is installed a few things need to be configured. Follow the instructions below based on your operating system.

### **Linux/MacOS**

1. Export the `GDBUILD_HOME` environment variable and add `$GDBUILD_HOME/bin` to the `PATH` environment variable.

    Add the following to your shell's profile script/RC file:

    ```sh
    export GDBUILD_HOME="$HOME/.gdbuild"
    export PATH="$GDBUILD_HOME/bin:$PATH"
    ```

### **Windows (Powershell)**

1. Export the `GDBUILD_HOME` environment variable using the following:

    ```sh
    $GdBuildHomePath = "${env:LOCALAPPDATA}\gdbuild" # Replace with whichever path you'd like.
    [System.Environment]::SetEnvironmentVariable("GDBUILD_HOME", $GdBuildHomePath, "User")
    ```

2. Add `$GDBUILD_HOME/bin` to your `PATH` environment variable:

    > ❕ **NOTE:** Make sure to restart your terminal after the previous step so that any changes to `$GDBUILD_HOME` have been updated.

    ```sh
    $PathParts = [System.Environment]::GetEnvironmentVariable("PATH", "User").Trim(";") -Split ";"
    $PathParts = $PathParts.where{ $_ -ne "${env:GDBUILD_HOME}\bin" }
    $PathParts = $PathParts + "${env:GDBUILD_HOME}\bin"

    [System.Environment]::SetEnvironmentVariable("PATH", $($PathParts -Join ";"), "User")
    ```
