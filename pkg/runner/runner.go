package runner

import (
	"context"
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
}

// Runner implements the Go package runner logic.
type Runner struct {
	cfg Config
}

// Option is a functional option for configuring a Runner.
type Option func(*Config)

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

	// Check if already installed in the cache.
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		if err := r.install(ctx, pkg, version, binDir); err != nil {
			return err
		}
	}

	return r.execute(ctx, binPath, args)
}

func (r *Runner) install(ctx context.Context, pkg, version, binDir string) error {
	fmt.Fprintf(r.cfg.Output, "gox: installing %s@%s...\n", pkg, version)

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

	// Using exec.CommandContext to allow for context cancellation/timeout.
	cmd := exec.CommandContext(ctx, "go", "install", pkg+"@"+version)
	cmd.Env = append(os.Environ(), "GOBIN="+binDir)
	cmd.Stderr = r.cfg.Output
	cmd.Stdout = r.cfg.Output

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install %s@%s: %w", pkg, version, err)
	}
	return nil
}

func (r *Runner) generateAndBuild(ctx context.Context, pkg, version, binPath string) error {
	// For now, simpler implementation: run go generate in the current directory
	// if it's a local path, or try to download and run if it's remote.
	
	// TODO: fully support remote go generate by downloading source to temp
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

	fmt.Fprintf(r.cfg.Output, "gox: building %s...\n", pkg)
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binPath, pkg)
	buildCmd.Stderr = r.cfg.Output
	buildCmd.Stdout = r.cfg.Output
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build %s: %w", pkg, err)
	}
	
	return nil
}

func isLocalPath(path string) bool {
	return strings.HasPrefix(path, ".") || strings.HasPrefix(path, "/") || strings.HasPrefix(path, "~")
}

func (r *Runner) execute(ctx context.Context, binPath string, args []string) error {
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Sub-process exit error is propagated up as-is.
			return exitErr
		}
		return fmt.Errorf("failed to execute binary %s: %w", binPath, err)
	}
	return nil
}

func parsePackage(pkgWithVersion string) (string, string) {
	pkg := pkgWithVersion
	version := "latest"
	if idx := strings.LastIndex(pkgWithVersion, "@"); idx != -1 {
		pkg = pkgWithVersion[:idx]
		version = pkgWithVersion[idx+1:]
	}
	return pkg, version
}

func getBinaryName(pkg string) string {
	parts := strings.Split(pkg, "/")
	return parts[len(parts)-1]
}
