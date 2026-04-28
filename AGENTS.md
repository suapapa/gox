# AGENTS.md

This document provides guidelines and context for AI agents working on the `gox` project.

> [!IMPORTANT]
> **Automatic Update Requirement**
> All AI agents MUST ensure that both `README.md` and `AGENTS.md` are kept up-to-date with any significant changes to the project structure, features, or workflows. These documents are intended to be living resources for both humans and AI agents.

## Project Overview
`gox` is a Go-based utility designed to provide a streamlined, `npx`-like experience for the Go ecosystem. It allows users to execute Go binaries directly from their source repositories without global installation.

### Key Goals
- Short, intuitive CLI: `gox <package> [args...]`.
- Support for versioning: `gox <package>@<version>`.
- Support GitHub shorthand resolution: `owner/repo` and `owner/repo/...` expand to `github.com/owner/repo...`.
- Intelligent caching of built binaries in `os.UserCacheDir() + "/gox/bin/"`.
- Automatic installation of missing tools into the managed cache.
- Forced reinstall/update flow with `-u` / `--update` for refreshing cached tools.
- Verbose execution tracing with `-v` / `--verbose` for cache inspection and install/update step visibility.

## Development Guidelines for Agents
- **Go Best Practices**: Follow idiomatic Go patterns. Use `golang-pro` skill instructions for reference.
- **CLI Framework**: The project uses `spf13/cobra` for the CLI structure.
- **Error Handling**: Ensure robust error handling and clear user feedback. Prefer actionable remediation hints when failures are caused by missing prerequisites such as the Go toolchain or a missing `go` binary in `PATH`.
- **Concurrency**: Use goroutines and channels where appropriate for performance (e.g., parallel downloads/builds in future phases).
- **Testing**: Maintain high test coverage with unit and integration tests.
- **Documentation parity**: When changing CLI flags or runtime output, update examples and behavioral docs in README immediately.

## File Structure
- `_assets/`: README media (e.g. hero image, before/after demo GIFs).
- `cmd/root.go`: CLI command definitions (Cobra).
- `pkg/`: Core logic and packages.
  - `runner/`: Logic for resolving, building, caching, updating, verbose logging, and running Go packages.
- `main.go`: Entry point for the application, delegates to `cmd`.
- `LICENSE`: Project license (MIT).
- `Makefile`: Standard build, test, and local install tasks.
- `gox_plan.md`: The initial project implementation plan.

## Workflow
1.  **Analyze**: Understand the current state of the codebase.
2.  **Plan**: Draft an implementation plan if the task is complex.
3.  **Implement**: Apply changes using the defined tools.
4.  **Verify**: Run tests and verify the changes manually.
5.  **Document**: Update `README.md` and `AGENTS.md` as required.

## Distribution
End users install from source with the Go toolchain, for example `go install github.com/suapapa/gox@latest` (see `README.md`). Versioning follows Go module tags on the default branch; there is no separate binary release pipeline in this repository.
