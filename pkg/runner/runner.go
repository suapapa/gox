package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PackageRunner defines the interface for running Go packages.
type PackageRunner interface {
	Run(ctx context.Context, pkgWithVersion string, args []string) error
}

// Config represents the configuration for the runner.
type Config struct {
	CacheDir string
	Output   io.Writer // Where to write installation progress
	Generate bool      // Run go generate before build
	Update   bool      // Force reinstall even when cached binary exists
	Verbose  bool      // Print cache/version/binary details and install steps
	Env      []string  // Environment variables for execution
}

// Runner implements the Go package runner logic.
type Runner struct {
	cfg Config
}

// Option is a functional option for configuring a Runner.
type Option func(*Config)

// WithEnv sets the environment variables for execution.
func WithEnv(env []string) Option {
	return func(c *Config) {
		c.Env = env
	}
}

// WithCacheDir sets the cache directory.
func WithCacheDir(dir string) Option {
	return func(c *Config) {
		c.CacheDir = dir
	}
}

// WithOutput sets the output writer.
func WithOutput(w io.Writer) Option {
	return func(c *Config) {
		c.Output = w
	}
}

// WithGenerate sets whether to run go generate before building.
func WithGenerate(g bool) Option {
	return func(c *Config) {
		c.Generate = g
	}
}

// WithUpdate sets whether to force reinstall even if the binary is cached.
func WithUpdate(u bool) Option {
	return func(c *Config) {
		c.Update = u
	}
}

// WithVerbose sets whether to print verbose progress information.
func WithVerbose(v bool) Option {
	return func(c *Config) {
		c.Verbose = v
	}
}

// NewRunner creates a new Runner instance with the provided options.
func NewRunner(opt ...Option) (*Runner, error) {
	cfg := Config{
		Output: os.Stderr,
	}

	for _, o := range opt {
		o(&cfg)
	}

	if cfg.CacheDir == "" {
		homeCache, err := os.UserCacheDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user cache dir: %w", err)
		}
		cfg.CacheDir = filepath.Join(homeCache, "gox")
	}

	return &Runner{cfg: cfg}, nil
}

// Run installs and executes the specified package with arguments.
// pkgWithVersion should be in the format "<package>[@version]".
// args are passed directly to the executed binary.
func (r *Runner) Run(ctx context.Context, pkgWithVersion string, args []string) error {
	pkg, version := parsePackage(pkgWithVersion)

	// Create a unique directory for binary to avoid name collisions.
	// Replacing "/" with "_" is a simple way to create a flat directory.
	pkgDirName := strings.ReplaceAll(pkg, "/", "_")
	binDir := filepath.Join(r.cfg.CacheDir, "bin", pkgDirName+"@"+version)
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	binName := getBinaryName(pkg)
	binPath := filepath.Join(binDir, binName)

	if r.cfg.Verbose {
		fmt.Fprintf(r.cfg.Output, "gox: package=%s version=%s\n", pkg, version)
		fmt.Fprintf(r.cfg.Output, "gox: cache-dir=%s\n", r.cfg.CacheDir)
		fmt.Fprintf(r.cfg.Output, "gox: bin-dir=%s\n", binDir)
		fmt.Fprintf(r.cfg.Output, "gox: binary=%s\n", binPath)
	}

	// Check if already installed in the cache unless update was requested.
	if r.cfg.Update {
		if r.cfg.Verbose {
			fmt.Fprintf(r.cfg.Output, "gox: update requested; refreshing cached binary\n")
		}
		if err := r.install(ctx, pkg, version, binDir); err != nil {
			return err
		}
	} else if _, err := os.Stat(binPath); os.IsNotExist(err) {
		if r.cfg.Verbose {
			fmt.Fprintf(r.cfg.Output, "gox: cache miss; binary not found, installing\n")
		}
		if err := r.install(ctx, pkg, version, binDir); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("failed to inspect cached binary: %w", err)
	} else if r.cfg.Verbose {
		fmt.Fprintf(r.cfg.Output, "gox: cache hit; reusing cached binary\n")
	}

	return r.execute(ctx, binPath, args)
}

func (r *Runner) install(ctx context.Context, pkg, version, binDir string) error {
	action := "installing"
	if r.cfg.Update {
		action = "updating"
	}
	pkgWithVer := pkg
	if version != "" {
		pkgWithVer = pkg + "@" + version
	}
	fmt.Fprintf(r.cfg.Output, "gox: %s %s...\n", action, pkgWithVer)

	binName := getBinaryName(pkg)
	if binName == "." || binName == ".." {
		// Use the name of the current directory if pkg is . or ..
		abs, err := filepath.Abs(pkg)
		if err == nil {
			binName = filepath.Base(abs)
		}
	}
	binPath := filepath.Join(binDir, binName)

	if r.cfg.Generate {
		return r.generateAndBuild(ctx, pkg, version, binPath)
	}

	if r.cfg.Verbose {
		fmt.Fprintf(r.cfg.Output, "gox: step 1/3 ensure cache directory exists\n")
		fmt.Fprintf(r.cfg.Output, "gox: step 2/3 run GOBIN=%s go install %s\n", binDir, pkgWithVer)
	}

	// Using exec.CommandContext to allow for context cancellation/timeout.
	cmd := exec.CommandContext(ctx, "go", "install", pkgWithVer)
	cmd.Env = append(os.Environ(), "GOBIN="+binDir)
	cmd.Stderr = r.cfg.Output
	cmd.Stdout = r.cfg.Output

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("failed to install %s: Go toolchain is not installed or `go` is not in PATH. Install Go first: https://go.dev/doc/install", pkgWithVer)
		}
		return fmt.Errorf("failed to install %s: %w", pkgWithVer, err)
	}

	if r.cfg.Verbose {
		fmt.Fprintf(r.cfg.Output, "gox: step 3/3 binary ready at %s\n", binPath)
	}
	return nil
}

func (r *Runner) generateAndBuild(ctx context.Context, pkg, version, binPath string) error {
	// For now, simpler implementation: run go generate in the current directory
	// if it's a local path, or try to download and run if it's remote.

	// TODO: fully support remote go generate by downloading source to temp
	if r.cfg.Verbose {
		fmt.Fprintf(r.cfg.Output, "gox: step 1/3 running go generate\n")
	}
	fmt.Fprintf(r.cfg.Output, "gox: running go generate...\n")

	genCmd := exec.CommandContext(ctx, "go", "generate", "./...")
	// If it's a local path, we might want to run it in that directory.
	if isLocalPath(pkg) {
		genCmd.Dir = pkg
	}
	genCmd.Stderr = r.cfg.Output
	genCmd.Stdout = r.cfg.Output
	if err := genCmd.Run(); err != nil {
		return fmt.Errorf("failed to run go generate: %w", err)
	}

	if r.cfg.Verbose {
		fmt.Fprintf(r.cfg.Output, "gox: step 2/3 building binary\n")
	}
	fmt.Fprintf(r.cfg.Output, "gox: building %s...\n", pkg)
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binPath, pkg)
	buildCmd.Stderr = r.cfg.Output
	buildCmd.Stdout = r.cfg.Output
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build %s: %w", pkg, err)
	}

	if r.cfg.Verbose {
		fmt.Fprintf(r.cfg.Output, "gox: step 3/3 binary ready at %s\n", binPath)
	}
	return nil
}

func isLocalPath(path string) bool {
	return strings.HasPrefix(path, ".") || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "~")
}

func (r *Runner) execute(ctx context.Context, binPath string, args []string) error {
	if r.cfg.Verbose {
		if len(args) > 0 {
			fmt.Fprintf(r.cfg.Output, "gox: executing %s %s\n", binPath, strings.Join(args, " "))
		} else {
			fmt.Fprintf(r.cfg.Output, "gox: executing %s\n", binPath)
		}
	}

	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if len(r.cfg.Env) > 0 {
		cmd.Env = append(os.Environ(), r.cfg.Env...)
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Sub-process exit error is wrapped to indicate it's a command failure
			return fmt.Errorf("command failed: %w", exitErr)
		}
		return fmt.Errorf("failed to execute binary %s: %w", binPath, err)
	}
	return nil
}

func parsePackage(pkgWithVersion string) (string, string) {
	pkg := pkgWithVersion
	version := ""
	if idx := strings.LastIndex(pkgWithVersion, "@"); idx != -1 {
		pkg = pkgWithVersion[:idx]
		version = pkgWithVersion[idx+1:]
	}

	if version == "" && !isLocalPath(pkg) {
		version = "latest"
	}

	pkg = normalizePackagePath(pkg)

	return pkg, version
}

func normalizePackagePath(pkg string) string {
	if pkg == "" || isLocalPath(pkg) || strings.Contains(pkg, ".") {
		return pkg
	}

	parts := strings.Split(pkg, "/")
	if len(parts) >= 2 {
		return "github.com/" + pkg
	}

	return pkg
}

func getBinaryName(pkg string) string {
	parts := strings.Split(strings.TrimRight(pkg, "/"), "/")
	return parts[len(parts)-1]
}
