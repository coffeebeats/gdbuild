# Commands

## **gdbuild `template`**

Compile an export template for the specified Godot platform `PLATFORM`.

### Usage

`gdbuild template [OPTIONS] <PLATFORM>`

### Options

- `--dry-run` — log the build command without running it
- `--force` — build the export template even if it was cached in the store
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

> [!NOTE]  
> This command has been disabled as the implementation is not complete.

### Options

- `--dry-run` — log the build command without running it
- `--force` - export the target even if it was cached in the store (does not rebuild the export template)
- `--print-hash` — log the unique hash of the game binary (skips exporting)

- `-p`, `--path <PATH>` — use the Godot project found at `PATH`
  - Default value: `$PWD` (current working directory)
- `-o`, `--out <PATH>` — write generated artifacts to `PATH`
  - Default value: `$PWD` (current working directory)
- `--build-dir <PATH>` — build the template within `PATH`
  - Default value: temporary directory

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

## **gdbuild `init`**

Initialize a Godot project with a GDBuild manifest.

### Usage

`gdbuild init [OPTIONS]`

### Options

- `-p`, `--path <PATH>` — use the Godot project found at `PATH`
  - Default value: `$PWD` (current working directory)
