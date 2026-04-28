# gox: Go version of npx

[![Website](https://img.shields.io/badge/website-gox-green)](https://suapapa.github.io/gox/)


![gox](_assets/gox.webp)

`gox` is a Go tool that allows you to execute Go binaries directly from their source repositories, without global installation. It provides a similar experience to `npx` in the Node.js ecosystem.

## Before / after

| Before gox | With gox |
| ---------- | -------- |
| ![Before gox](_assets/before_gox.gif) | ![With gox](_assets/gox.gif) |

## Features
- Runs Go binaries without global installation.
- Caches built binaries in `~/.cache/gox/bin/` to speed up subsequent runs.
- Avoids name collisions by using unique directories for each package.
- Supports versioning (e.g., `@v1.2.3` or `@latest`).
- Supports forced reinstall/update with `-u` / `--update`.
- Supports verbose logging with `-v` / `--verbose` to show version, cache paths, binary path, cache hits, and install/update steps.
- Prints a clearer help message when the Go toolchain is not installed or `go` is missing from `PATH`.

## Prerequisites

`gox` requires the [Go toolchain](https://go.dev/doc/install) to be installed and available in your `PATH`.

## Installation

### Scripted Install

```bash
curl -fsSL https://raw.githubusercontent.com/suapapa/gox/main/scripts/install.sh | sh
```

### Manual Install

```bash
go install github.com/suapapa/gox@latest
```

Make sure `$(go env GOPATH)/bin` is in your `PATH` to execute `gox`. You can add it by adding the following line to your shell configuration file (e.g., `.zshrc` or `.bashrc`):

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

## Usage
```bash
# Run a tool at the latest version
gox <package-path>[@version] [args...]

# Force reinstall/update even if it is already cached
gox -u <package-path>[@version] [args...]

# Print version/cache/binary details and install steps
gox -v <package-path>[@version] [args...]

# GitHub shortcuts: owner/repo paths implicitly expand to github.com/owner/repo
gox suapapa/gox --help

# Full import paths still work as-is
gox github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest --version
```

## Verbose output
With `-v`, `gox` prints:
- Resolved package path and version
- Cache root and package-specific cache directory
- Final cached binary path
- Whether the run is a cache hit, cache miss, or forced update
- Per-step install/build progress for first-time installs and updates

Example:
```bash
gox -v golang.org/x/tools/cmd/goimports@latest --version
```

Example output:
```text
gox: package=golang.org/x/tools/cmd/goimports version=latest
gox: cache-dir=/home/user/.cache/gox
gox: bin-dir=/home/user/.cache/gox/bin/golang.org_x_tools_cmd_goimports@latest
gox: binary=/home/user/.cache/gox/bin/golang.org_x_tools_cmd_goimports@latest/goimports
gox: cache miss; binary not found, installing
gox: installing golang.org/x/tools/cmd/goimports@latest...
gox: step 1/3 ensure cache directory exists
gox: step 2/3 run GOBIN=/home/user/.cache/gox/bin/golang.org_x_tools_cmd_goimports@latest go install golang.org/x/tools/cmd/goimports@latest
gox: step 3/3 binary ready at /home/user/.cache/gox/bin/golang.org_x_tools_cmd_goimports@latest/goimports
gox: executing /home/user/.cache/gox/bin/golang.org_x_tools_cmd_goimports@latest/goimports --version
```

## Update semantics
- Without `-u`, `gox` reuses the cached binary when available.
- With `-u`, `gox` rebuilds/reinstalls the requested package into the managed cache before executing it.
- `-u` is useful for refreshing `@latest` tools or forcing a clean rebuild of a pinned version.

### Package path resolution
- `owner/repo` or `owner/repo/...` is treated as `github.com/owner/repo`.
- Full import paths such as `github.com/...` or `golang.org/...` are used unchanged.
- Local paths such as `./cmd/tool` are used unchanged.

## How it works
`gox` uses `go install <package>@<version>` to download and build the tool into a managed cache directory (`os.UserCacheDir() + "/gox/bin/"). Subsequent runs of the same package and version will use the cached binary unless `-u` is specified.

## Development and Contributions
AI agents and contributors should refer to [AGENTS.md](file:///Users/suapapa/ws_suapapa/gox/AGENTS.md) for development guidelines, project structure, and mandatory documentation updates.

You can use the provided `Makefile` for common tasks:
- `make build`: Build the binary locally.
- `make test`: Run tests.
- `make snapshot`: Create a local snapshot release using `goreleaser`.
