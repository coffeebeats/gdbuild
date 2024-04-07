#!/bin/sh
set -e

# This script installs 'gdbuild' by downloading prebuilt binaries from the
# project's GitHub releases page. By default the latest version is installed,
# but a different release can be used instead by setting $GDBUILD_VERSION.
#
# The script will set up a 'gdbuild' cache at '$HOME/.gdbuild'. This behavior can
# be customized by setting '$GDBUILD_HOME' prior to running the script. Existing
# Godot artifacts cached in a 'gdbuild' store won't be lost, but this script will
# overwrite any 'gdbuild' binary artifacts in '$GDBUILD_HOME/bin'.

# ------------------------------ Define: Cleanup ----------------------------- #

trap cleanup EXIT

cleanup() {
    if [ -d "${GDBUILD_TMP=}" ]; then
        rm -rf "${GDBUILD_TMP}"
    fi
}

# ------------------------------ Define: Logging ----------------------------- #

info() {
    if [ "$1" != "" ]; then
        echo info: "$@"
    fi
}

warn() {
    if [ "$1" != "" ]; then
        echo warning: "$1"
    fi
}

error() {
    if [ "$1" != "" ]; then
        echo error: "$1" >&2
    fi
}

fatal() {
    error "$1"
    exit 1
}

unsupported_platform() {
    error "$1"
    echo "See https://github.com/coffeebeats/gdbuild/blob/main/docs/installation.md#install-from-source for instructions on compiling from source."
    exit 1
}

# ------------------------------- Define: Usage ------------------------------ #

usage() {
    cat <<EOF
gdbuild-install: Install 'gdbuild' for compiling and exporting Godot projects.

Usage: gdbuild-install [OPTIONS]

NOTE: The following dependencies are required:
    - curl OR wget
    - grep
    - sha256sum OR shasum
    - tar/unzip
    - tr
    - uname

Available options:
    -h, --help          Print this help and exit
    -v, --verbose       Print script debug info (default=false)
    --no-modify-path    Do not modify the \$PATH environment variable
EOF
    exit
}

check_cmd() {
    command -v "$1" >/dev/null 2>&1
}

need_cmd() {
    if ! check_cmd "$1"; then
        fatal "required command not found: '$1'"
    fi
}

# ------------------------------ Define: Params ------------------------------ #

parse_params() {
    MODIFY_PATH=1

    while :; do
        case "${1:-}" in
        -h | --help) usage ;;
        -v | --verbose) set -x ;;

        --no-modify-path) MODIFY_PATH=0 ;;

        -?*) fatal "Unknown option: $1" ;;
        "") break ;;
        esac
        shift
    done

    return 0
}

parse_params "$@"

# ------------------------------ Define: Version ----------------------------- #

GDBUILD_VERSION="${GDBUILD_VERSION:-0.3.5}" # x-release-please-version
GDBUILD_VERSION="v${GDBUILD_VERSION#v}"

# ----------------------------- Define: Platform ----------------------------- #

need_cmd tr
need_cmd uname

GDBUILD_CLI_OS="$(echo "${GDBUILD_CLI_OS=$(uname -s)}" | tr '[:upper:]' '[:lower:]')"
case "$GDBUILD_CLI_OS" in
darwin*) GDBUILD_CLI_OS="macos" ;;
linux*) GDBUILD_CLI_OS="linux" ;;
mac | macos | osx) GDBUILD_CLI_OS="macos" ;;
cygwin*) GDBUILD_CLI_OS="windows" ;;
msys* | mingw64*) GDBUILD_CLI_OS="windows" ;;
uwin* | win*) GDBUILD_CLI_OS="windows" ;;
*) unsupported_platform "no prebuilt binaries available for operating system: $GDBUILD_CLI_OS" ;;
esac

GDBUILD_CLI_ARCH="$(echo ${GDBUILD_CLI_ARCH=$(uname -m)} | tr '[:upper:]' '[:lower:]')"
case "$GDBUILD_CLI_ARCH" in
aarch64 | arm64)
    GDBUILD_CLI_ARCH="arm64"
    if [ "$GDBUILD_CLI_OS" != "macos" ] && [ "$GDBUILD_CLI_OS" != "linux" ]; then
        fatal "no prebuilt '$GDBUILD_CLI_ARCH' binaries available for operating system: $GDBUILD_CLI_OS"
    fi

    ;;
amd64 | x86_64) GDBUILD_CLI_ARCH="x86_64" ;;
*) unsupported_platform "no prebuilt binaries available for CPU architecture: $GDBUILD_CLI_ARCH" ;;
esac

GDBUILD_CLI_ARCHIVE_EXT=""
case "$GDBUILD_CLI_OS" in
windows) GDBUILD_CLI_ARCHIVE_EXT="zip" ;;
*) GDBUILD_CLI_ARCHIVE_EXT="tar.gz" ;;
esac

GDBUILD_CLI_ARCHIVE="gdbuild-$GDBUILD_VERSION-$GDBUILD_CLI_OS-$GDBUILD_CLI_ARCH.$GDBUILD_CLI_ARCHIVE_EXT"

# ------------------------------- Define: Store ------------------------------ #

GDBUILD_HOME_PREV="${GDBUILD_HOME_PREV=}" # save for later in script

GDBUILD_HOME="${GDBUILD_HOME=}"
if [ "$GDBUILD_HOME" = "" ]; then
    if [ "${HOME=}" = "" ]; then
        fatal "both '\$GDBUILD_HOME' and '\$HOME' unset; one must be specified to determine 'gdbuild' installation path"
    fi

    GDBUILD_HOME="$HOME/.gdbuild"
fi

info "using 'gdbuild' store path: '$GDBUILD_HOME'"

# ----------------------------- Define: Download ----------------------------- #

need_cmd grep
need_cmd mktemp

GDBUILD_TMP=$(mktemp -d --tmpdir gdbuild-XXXXXXXXXX)
cd "$GDBUILD_TMP"

GDBUILD_RELEASE_URL="https://github.com/coffeebeats/gdbuild/releases/download/$GDBUILD_VERSION"

download_with_curl() {
    curl \
        --fail \
        --location \
        --parallel \
        --retry 3 \
        --retry-delay 1 \
        --show-error \
        --silent \
        -o "$GDBUILD_CLI_ARCHIVE" \
        "$GDBUILD_RELEASE_URL/$GDBUILD_CLI_ARCHIVE" \
        -o "checksums.txt" \
        "$GDBUILD_RELEASE_URL/checksums.txt"
}

download_with_wget() {
    wget -q -t 4 -O "$GDBUILD_CLI_ARCHIVE" "$GDBUILD_RELEASE_URL/$GDBUILD_CLI_ARCHIVE" 2>&1
    wget -q -t 4 -O "checksums.txt" "$GDBUILD_RELEASE_URL/checksums.txt" 2>&1
}

if check_cmd curl; then
    download_with_curl
elif check_cmd wget; then
    download_with_wget
else
    fatal "missing one of 'curl' or 'wget' commands"
fi

# -------------------------- Define: Verify checksum ------------------------- #

verify_with_sha256sum() {
    cat "checksums.txt" | grep "$GDBUILD_CLI_ARCHIVE" | sha256sum --check --status
}

verify_with_shasum() {
    cat "checksums.txt" | grep "$GDBUILD_CLI_ARCHIVE" | shasum -a 256 -p --check --status
}

if check_cmd sha256sum; then
    verify_with_sha256sum
elif check_cmd shasum; then
    verify_with_shasum
else
    fatal "missing one of 'sha256sum' or 'shasum' commands"
fi

# ------------------------------ Define: Extract ----------------------------- #

case "$GDBUILD_CLI_OS" in
windows)
    need_cmd unzip

    mkdir -p "$GDBUILD_HOME/bin"
    unzip -u "$GDBUILD_CLI_ARCHIVE" -d "$GDBUILD_HOME/bin"
    ;;
*)
    need_cmd tar

    mkdir -p "$GDBUILD_HOME/bin"
    tar -C "$GDBUILD_HOME/bin" --no-same-owner -xzf "$GDBUILD_CLI_ARCHIVE"
    ;;
esac

info "successfully installed 'gdbuild@$GDBUILD_VERSION' to '$GDBUILD_HOME/bin'"

if [ $MODIFY_PATH -eq 0 ]; then
    exit 0
fi

# The $PATH modification and $GDBUILD_HOME export is already done.
if check_cmd gdbuild && [ "$GDBUILD_HOME_PREV" != "" ]; then
    exit 0
fi

# Simplify the exported $GDBUILD_HOME if possible.
if [ "$HOME" != "" ]; then
    case "$GDBUILD_HOME" in
    $HOME*) GDBUILD_HOME="\$HOME${GDBUILD_HOME#$HOME}" ;;
    esac
fi

CMD_EXPORT_HOME="export GDBUILD_HOME=\"$GDBUILD_HOME\""
CMD_MODIFY_PATH="export PATH=\"\$GDBUILD_HOME/bin:\$PATH\""

case $(basename $SHELL) in
sh) OUT="$HOME/.profile" ;;
bash) OUT="$HOME/.bashrc" ;;
zsh) OUT="$HOME/.zshrc" ;;
*)
    echo ""
    echo "Add the following to your shell profile script:"
    echo "    $CMD_EXPORT_HOME"
    echo "    $CMD_MODIFY_PATH"
    ;;
esac

if [ "$OUT" != "" ]; then
    if [ -f "$OUT" ] && $(cat "$OUT" | grep -q 'export GDBUILD_HOME'); then
        info "Found 'GDBUILD_HOME' export in shell Rc file; skipping modification."
        exit 0
    fi

    if [ -f "$OUT" ] && [ "$(tail -n 1 "$OUT")" != "" ]; then
        echo "" >>"$OUT"
    fi

    echo "# Added by 'gdbuild' install script." >>"$OUT"
    echo "$CMD_EXPORT_HOME" >>"$OUT"
    echo "$CMD_MODIFY_PATH" >>"$OUT"

    info "Updated shell Rc file: $OUT\n      Open a new terminal to start using 'gdbuild'."
fi
