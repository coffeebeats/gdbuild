# Changelog

## 0.3.9 (2024-04-07)

## What's Changed
* chore(cmd): log cached contents by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/111


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.3.8...v0.3.9

## 0.3.8 (2024-04-07)

## What's Changed
* feat(ci): define a `run-target` action to export targets by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/106
* fix(ci): define `template-archive-path` on correct action by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/108
* fix(ci): exit on print hash error by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/109
* chore(cmd): log the newly computed hash by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/110


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.3.7...v0.3.8

## 0.3.7 (2024-04-07)

## What's Changed
* Revert "fix(exec): quote arguments passed to shell (#102)" by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/103
* feat(ci): define an action to run the `template` command by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/104


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.3.6...v0.3.7

## 0.3.6 (2024-04-07)

## What's Changed
* fix(ci): only cache the `gdbuild` binary directory by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/100
* fix(exec): quote arguments passed to shell by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/102


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.3.5...v0.3.6

## 0.3.5 (2024-04-07)

## What's Changed
* fix(cmd): correctly hash template archives passed via `--template-archive` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/98


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.3.4...v0.3.5

## 0.3.4 (2024-04-06)

## What's Changed
* feat(cmd): allow passing in a template archive to `target` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/96


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.3.3...v0.3.4

## 0.3.3 (2024-04-06)

## What's Changed
* feat(cmd): add `--project` flag to `template` to standardize commands by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/94


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.3.2...v0.3.3

## 0.3.2 (2024-04-06)

## What's Changed
* fix(cmd): ensure hashes are sent to stdout by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/92


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.3.1...v0.3.2

## 0.3.1 (2024-04-06)

## What's Changed
* feat(config): provide top-level target `encrypt` setting; improve encryption validation by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/90


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.3.0...v0.3.1

## 0.3.0 (2024-04-05)

## What's Changed
* chore(deps): bump tj-actions/changed-files from 43 to 44 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/87
* feat(target): implement `target` exporting by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/76
* chore!: update `init` docs; bump minor version by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/89


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.2.5...v0.3.0

## 0.2.5 (2024-04-01)

## What's Changed
* chore: upgrade `gdenv` to `v0.6.16` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/85


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.2.4...v0.2.5

## 0.2.4 (2024-04-01)

## What's Changed
* feat(ci): pre-build `arm64` on `linux` binaries by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/80
* fix(scripts): unblock downloads of new `arm64` on `linux` target by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/82
* fix(scripts): use correct compound condition syntax by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/83
* feat(ci): add support for explicit `--debug` flag by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/84


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.2.3...v0.2.4

## 0.2.3 (2024-03-31)

## What's Changed
* Revert "chore(cmd): disable `target` command (#75)" by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/77
* refactor(cmd): improve code reuse; fix various bugs by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/78


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.2.2...v0.2.3

## 0.2.2 (2024-03-30)

## What's Changed
* chore(deps): bump github.com/coffeebeats/gdenv from 0.6.13 to 0.6.14 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/70
* chore(template): migrate `build` package into `template` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/72
* chore(deps): bump dependabot/fetch-metadata from 1 to 2 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/71
* refactor(run): change `Context.PathBuild` to `Context.PathWorkspace` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/74
* chore(cmd): disable `target` command by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/75


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.2.1...v0.2.2

## 0.2.1 (2024-03-25)

## What's Changed
* feat(target): scaffold support for exporting game projects by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/66
* refactor: move `build.Context` into own `run` package; reorganize `pkg/config` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/68
* feat(cmd): implement export flow (except for actual exporting) by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/69


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.2.0...v0.2.1

## 0.2.0 (2024-03-24)

## What's Changed
* refactor(build): switch `Source.Version` to a `Version` type for improved clarity by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/60
* feat(template): add support for encrypting export templates by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/62
* chore(cmd)!: remove unused `info` command by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/63
* feat(cmd): define an `init` command for creating `gdbuild` manifests by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/64
* fix(cmd): correctly log extra arguments passed to `init` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/65


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.1.3...v0.2.0

## 0.1.3 (2024-03-23)

## What's Changed
* feat(store): add a store package for caching template artifacts by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/44
* fix(ci): fetch full history to enable correct change detection by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/47
* refactor: huge refactoring of package layout; improves organization by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/48
* feat(store): cache built artifacts in the store by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/49
* chore(ci): add temporary logging by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/55
* fix(ci): correctly report whether changes were detected by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/57
* fix(ci): lessen fetch depth; remove test logging by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/58
* fix(build): correctly pass `platform` argument to SCons by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/50
* feat(cmd): update `template` to utilize store cache; add `force` options to `template` and `target` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/59


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.1.2...v0.1.3

## 0.1.2 (2024-03-22)

## What's Changed
* feat(build): implement `Template` hashes for determining export template equivalence by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/39
* feat(template): add support for registering expected template artifacts by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/42
* chore(deps): bump github.com/charmbracelet/log from 0.3.1 to 0.4.0 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/41
* feat(archive): create a package for writing and extracting template archives by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/43


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.1.1...v0.1.2

## 0.1.1 (2024-03-21)

## What's Changed
* refactor(template): overhaul configuration implementation for compiling templates by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/37
* chore(deps): bump tj-actions/changed-files from 42 to 43 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/35
* chore(deps): bump github.com/pelletier/go-toml/v2 from 2.1.1 to 2.2.0 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/36


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.1.0...v0.1.1

## 0.1.0 (2024-03-16)

## What's Changed
* feat(ci): expand environment variables within paths defined in the manifest by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/32
* feat(cmd)!: switch `--path` flag to `--config` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/34


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.0.6...v0.1.0

## 0.0.6 (2024-03-15)

## What's Changed
* chore(deps): bump github.com/coffeebeats/gdenv from 0.6.12 to 0.6.13 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/29
* fix(template): don't vendor source code if build directory is source directory; update `macos.dynamic` to `macos.use_volk` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/30


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.0.5...v0.0.6

## 0.0.5 (2024-03-14)

## What's Changed
* fix: correctly set default version in `install.sh` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/20
* feat(ci): define an action for installing `gdbuild` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/22
* fix(ci): correctly update `PATH` in setup action by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/23
* fix(ci): correctly export environment variable in setup action by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/24
* fix(ci): correctly reference home directory in action by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/25
* fix(ci): correctly use environment variable in cache path by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/26
* fix(ci): use correct path in cache key by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/27
* fix(ci): conditionally check for executable on path during setup by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/28


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.0.4...v0.0.5

## 0.0.4 (2024-03-11)

## What's Changed
* feat(template): add support for building `Linux`, `Windows`, and `MacOS` templates by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/18


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.0.3...v0.0.4

## 0.0.3 (2024-03-10)

## What's Changed
* chore(template): remove implementation of unused interface by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/14
* fix(template): address errors blocking template builds by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/17
* chore(deps): bump github.com/charmbracelet/lipgloss from 0.9.1 to 0.10.0 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/16


**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.0.2...v0.0.3

## 0.0.2 (2024-02-28)

## What's Changed
* feat(manifest): add support for parsing a GDBuild manifest by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/6
* chore(deps): bump golangci/golangci-lint-action from 3 to 4 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/5
* feat(cmd): refactor GDBuild commands to include `template`, `target`, and `info` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/8
* feat: add support for resolving build settings from a `Manifest` and related options by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/9
* refactor(pkg/platform): move `platform` package into `build` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/10
* chore(deps): bump github.com/coffeebeats/gdenv from 0.6.10 to 0.6.11 by @dependabot in https://github.com/coffeebeats/gdbuild/pull/11
* feat(template): expand on template command building (WIP) by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/12
* chore: remove release version pin used for bootstrapping by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/13

## New Contributors
* @dependabot made their first contribution in https://github.com/coffeebeats/gdbuild/pull/5

**Full Changelog**: https://github.com/coffeebeats/gdbuild/compare/v0.0.1...v0.0.2

## 0.0.1 (2024-02-05)

## What's Changed
* chore: set up repository infrastructure by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/1
* chore: pin the release version to `0.0.1` by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/2
* fix: properly set the current version by @coffeebeats in https://github.com/coffeebeats/gdbuild/pull/3

## New Contributors
* @coffeebeats made their first contribution in https://github.com/coffeebeats/gdbuild/pull/1

**Full Changelog**: https://github.com/coffeebeats/gdbuild/commits/v0.0.1
