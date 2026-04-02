# gox: Go version of npx

![gox](_assets/gox.webp)

`gox` is a Go tool that allows you to execute Go binaries directly from their source repositories, without global installation. It provides a similar experience to `npx` in the Node.js ecosystem.

## Features
- Runs Go binaries without global installation.
- Caches built binaries in `~/.cache/gox/bin/` to speed up subsequent runs.
- Avoids name collisions by using unique directories for each package.
- Supports versioning (e.g., `@v1.2.3` or `@latest`).
- Supports forced reinstall/update with `-u` / `--update`.

## Installation

### Homebrew
```bash
brew install --cask suapapa/tools/gox
```
(Alternatively: `brew tap suapapa/tools && brew install --cask gox`)

## Usage
```bash
# Run a tool at the latest version
gox <package-path>[@version] [args...]

# Force reinstall/update even if it is already cached
gox -u <package-path>[@version] [args...]

# GitHub shortcuts: owner/repo paths implicitly expand to github.com/owner/repo
gox suapapa/gox --help

# Full import paths still work as-is
gox github.com/golangci/golangci-lint/cmd/golangci-lint@latest --version
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
