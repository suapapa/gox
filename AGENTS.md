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

## Development Guidelines for Agents
- **Go Best Practices**: Follow idiomatic Go patterns. Use `golang-pro` skill instructions for reference.
- **CLI Framework**: The project uses `spf13/cobra` for the CLI structure.
- **Error Handling**: Ensure robust error handling and clear user feedback.
- **Concurrency**: Use goroutines and channels where appropriate for performance (e.g., parallel downloads/builds in future phases).
- **Testing**: Maintain high test coverage with unit and integration tests.

## File Structure
- `cmd/root.go`: CLI command definitions (Cobra).
- `pkg/`: Core logic and packages.
  - `runner/`: Logic for resolving, building, and running Go packages.
- `main.go`: Entry point for the application, delegates to `cmd`.
- `LICENSE`: Project license (MIT).
- `Makefile`: Standard build, test, and release tasks.
- `.goreleaser.yaml`: GoReleaser configuration for GitHub releases.
- `gox_plan.md`: The initial project implementation plan.

## Workflow
1.  **Analyze**: Understand the current state of the codebase.
2.  **Plan**: Draft an implementation plan if the task is complex.
3.  **Implement**: Apply changes using the defined tools.
4.  **Verify**: Run tests and verify the changes manually.
5.  **Document**: Update `README.md` and `AGENTS.md` as required.

## Release Workflow
The project uses `goreleaser` for automated GitHub releases and Homebrew cask management.

- **Check config**: `make release-check`
- **Snapshot release (local test)**: `make snapshot`
- **Full release (requires tag and GITHUB_TOKEN)**: `make release`

> [!NOTE]
> The release process automatically pushes an updated Homebrew cask to [suapapa/homebrew-tools](https://github.com/suapapa/homebrew-tools). The generated cask conflicts with the legacy `gox` formula during the migration away from GoReleaser `brews`.
