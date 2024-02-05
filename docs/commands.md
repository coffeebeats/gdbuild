# Commands

## **gdbuild `template`**

TODO

### Usage

`gdbuild template [OPTIONS] <PLATFORM>`

### Options

- `-f`, `--force` — forcibly overwrite an existing cache entry
- `-g`, `--global` — update the global pin (if `VERSION` is specified) or resolve `VERSION` from the global pin
- `-p`, `--path <PATH>` — resolve the pinned `VERSION` at `PATH`
- `-s`, `--src`, `--source` — install source code instead of an executable (cannot be used with `-g`)

### Arguments

- `<PLATFORM>` — the specific template platform to build
  - Example values:
    - `TODO`

## **gdbuild `project`**

TODO

### Usage

`gdbuild project [OPTIONS] <VERSION>`

### Options

- `-g`, `--global` — pin the system version (cannot be used with `-p`)
- `-i`, `--install` — install the specified version of _Godot_ if missing
- `-f`, `--force` — forcibly overwrite an existing cache entry (only used with `-i`)
- `-p`, `--path <PATH>` — pin the specified path (cannot be used with `-g`)
  - Default value: `$PWD` (current working directory)

### Arguments

- `<VERSION>` — the specific version string to install (must be exact)
  - Example values:
    - `3.5.1` (if missing, the label will default to `stable`)
    - `4.0.4-stable`
    - `4.2-beta2`
