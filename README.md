# gox: Go version of npx

![gox](_assets/gox.webp)

`gox` is a Go tool that allows you to execute Go binaries directly from their source repositories, without global installation. It provides a similar experience to `npx` in the Node.js ecosystem.

## Features
- Runs Go binaries without global installation.
- Caches built binaries in `~/.cache/gox/bin/` to speed up subsequent runs.
- Avoids name collisions by using unique directories for each package.
- Supports versioning (e.g., `@v1.2.3` or `@latest`).

## Usage
```bash
# Run a tool at the latest version
gox <package-path>[@version] [args...]

# Example: Run golangci-lint
gox github.com/golangci/golangci-lint/cmd/golangci-lint@latest --version
```

## How it works
`gox` uses `go install <package>@<version>` to download and build the tool into a managed cache directory (`os.UserCacheDir() + "/gox/bin/"`). Subsequent runs of the same package and version will use the cached binary.

## Development and Contributions
AI agents and contributors should refer to [AGENTS.md](file:///Users/suapapa/ws_suapapa/gox/AGENTS.md) for development guidelines, project structure, and mandatory documentation updates.
