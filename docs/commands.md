# Commands

## **gdbuild `template`**

Compile an export template for the specified Godot platform `PLATFORM`.

### Usage

`gdbuild template [OPTIONS] <PLATFORM>`

### Options

- `--dry-run` — log the build command without running it
- `--print-hash` — log the unique hash of the export template (skips compilation)

- `-c`, `--config <PATH>` — use the `gdbuild` configuration file found at `PATH`
  - Default value: `$PWD` (current working directory)
- `-o`, `--out <PATH>` — write generated artifacts to `PATH`
  - Default value: `$PWD` (current working directory)
- `--build-dir <PATH>` — build the template within `PATH`
  - Default value: temporary directory

- `--release` — use a release export template (cannot be used with '--release_debug')
- `--release_debug` — use a release export template with debug symbols (cannot be used with '--release')

### Arguments

- `<PLATFORM>` — build for the specified Godot platform `PLATFORM`
  - Default value: `runtime.GOOS` (host platform)

## **gdbuild `target`**

Compile any required export template(s) and then export the specified `TARGET`.

### Usage

`gdbuild target [OPTIONS] <TARGET>`

### Options

- `--dry-run` — log the build command without running it

- `-p`, `--path <PATH>` — use the Godot project found at `PATH`
  - Default value: `$PWD` (current working directory)
- `-o`, `--out <PATH>` — write generated artifacts to `PATH`
  - Default value: `$PWD` (current working directory)

- `-f`, `--feature <FEATURE>` — enable the provided feature tag `FEATURE` (can be specified more than once)
- `-p`, `--platform <PLATFORM>` — build for the specified Godot platform `PLATFORM`
  - Default value: `runtime.GOOS` (host platform)
- `--release` — use a release export template (cannot be used with '--release_debug')
- `--release_debug` — use a release export template with debug symbols (cannot be used with '--release')

### Arguments

- `<TARGET>` — the name of a target specified in the GDBuild manifest (must be exact).
  - Example values (assuming they are present in GDBuild manifest):
    - `client` (define under `target.client` heading)
    - `dlc` (define under `target.dlc` heading; no export template required)

## **gdbuild `info`**

Inspect various properties of the GDBuild manifest.

### Usage

`gdbuild info [OPTIONS] <PROPERTY>`

### Options

- `-p`, `--path <PATH>` — use the Godot project found at `PATH`
  - Default value: `$PWD` (current working directory)

- `--json` — print the property values in JSON format

### Arguments

- `<PROPERTY>` — the name of a manifest property to inspect.
  - Example values (assuming they are present in GDBuild manifest):
    - `target` (lists all targets defined under `target` heading)
