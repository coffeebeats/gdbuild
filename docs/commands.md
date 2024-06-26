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
  - Default value: `<PROJECT>/gdbuild.toml` (`gdbuild.toml` in project directory)
- `-p`, `--project <PATH>` — use the Godot project found at `PATH`
  - Default value: `$PWD` (current working directory)
- `-o`, `--out <PATH>` — write generated artifacts to `PATH`
  - Default value: `$PWD` (current working directory)

- `--release` — use a release export template (cannot be used with `--release_debug` or `--debug`)
- `--release_debug` — use a release export template with debug symbols (cannot be used with `--release` or `--debug`)
- `--debug` — use a debug export template (cannot be used with `--release` or `--release_debug`)

### Arguments

- `<PLATFORM>` — build for the specified Godot platform `PLATFORM`
  - Default value: `runtime.GOOS` (host platform)

## **gdbuild `target`**

Compile any required export template(s) and then export the specified `TARGET`.

### Usage

`gdbuild target [OPTIONS] <TARGET>`

### Options

- `--dry-run` — log the build command without running it
- `--force` - export the target even if it was cached in the store (does not rebuild the export template)
- `--print-hash` — log the unique hash of the game binary (skips exporting)

- `-c`, `--config <PATH>` — use the `gdbuild` configuration file found at `PATH`
  - Default value: `<PROJECT>/gdbuild.toml` (`gdbuild.toml` in project directory)
- `-p`, `--project <PATH>` — use the Godot project found at `PATH`
  - Default value: `$PWD` (current working directory)
- `-o`, `--out <PATH>` — write generated artifacts to `PATH`
  - Default value: `$PWD` (current working directory)
- `--template-archive <PATH>` - extract the template from the archive found at `PATH` (skips template build)

- `-f`, `--feature <FEATURE>` — enable the provided feature tag `FEATURE` (can be specified more than once)
- `-p`, `--platform <PLATFORM>` — build for the specified Godot platform `PLATFORM`
  - Default value: `runtime.GOOS` (host platform)
- `--release` — use a release export template (cannot be used with `--release_debug` or `--debug`)
- `--release_debug` — use a release export template with debug symbols (cannot be used with `--release` or `--debug`)
- `--debug` — use a debug export template (cannot be used with `--release` or `--release_debug`)

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

- `-p`, `--project <PATH>` — use the Godot project found at `PATH`
  - Default value: `$PWD` (current working directory)
