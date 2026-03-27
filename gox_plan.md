# gox: Go version of npx

## 1. Background & Motivation
In the Node.js ecosystem, `npx` allows users to execute binaries from npm packages without globally installing them. While Go 1.17+ supports `go run package@version`, the syntax is often verbose (e.g., `go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest`).

`gox` aims to provide a streamlined, `npx`-like experience for the Go ecosystem, making it easier to run Go-based tools with minimal friction.

## 2. Scope & Impact
### Goals
- Provide a short, intuitive CLI: `gox <package> [args...]`.
- Support versioning: `gox <package>@<version>`.
- Intelligent package resolution (optional/future: mapping common tool names to paths).
- Dedicated binary caching to avoid repeated compilation.
- Automatic installation of missing tools into a managed cache.

### Non-Goals
- Replacing `go install` for permanent system-wide installations.
- Managing Go toolchains themselves (use `asdf` or `gvm` for that).

## 3. Proposed Solution
`gox` will be a CLI written in Go.

### Core Workflow
1. **Command Parsing:** Parse `gox <package>[@version] [args...]`.
2. **Resolution:**
   - If `<package>` is a full path (e.g., `github.com/...`), use it.
3. **Cache Check:** Check if the binary for the requested version is already in `$HOME/.cache/gox/bin/<encoded-package-name>`.
   - Path recommendation: `os.UserCacheDir()` + `/gox/bin/`.
4. **Execution/Build:**
   - If not in cache:
     - Run `GOBIN=$CACHE_DIR/bin go install <package>@<version>`.
     - This will build and install the binary into the cache.
5. **Invocation:** Execute the binary from the cache with the provided arguments.

### Technical Stack
- Language: Go
- CLI Library: `spf13/cobra` (standard for Go CLIs)
- Configuration: `spf13/viper` (optional, for cache paths)

## 4. Phased Implementation Plan

### Phase 1: MVP (Minimum Viable Product)
- Basic CLI structure using `cobra`.
- Support for full import paths: `gox github.com/user/repo/cmd/tool@latest`.
- Simple binary caching in `~/.gox/bin`.
- Basic execution logic.

### Phase 2: Enhanced UX
- Short name resolution (mapping common tools).
- Version range support (if feasible).
- `gox --list` to see cached tools.
- `gox --clean` to clear cache.

### Phase 3: Advanced Features
- Parallel downloads/builds.
- Integration with Go proxy for faster discovery.
- Shell completions.

## 5. Verification Plan
- Unit tests for path parsing and version extraction.
- Integration tests:
  - Running a well-known tool (e.g., `github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0`).
  - Verifying cache hits on second run.
  - Verifying argument passing.

## 6. Future Scope
- Support for non-Go binaries (maybe?).
- Cross-platform support (Windows/Linux/macOS).
