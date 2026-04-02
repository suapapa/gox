package runner

import (
	"testing"
)

func TestParsePackage(t *testing.T) {
	tests := []struct {
		name           string
		pkgWithVersion string
		wantPkg        string
		wantVersion    string
	}{
		{
			name:           "Latest version",
			pkgWithVersion: "github.com/suapapa/gox",
			wantPkg:        "github.com/suapapa/gox",
			wantVersion:    "latest",
		},
		{
			name:           "Specific version",
			pkgWithVersion: "github.com/suapapa/gox@v1.0.0",
			wantPkg:        "github.com/suapapa/gox",
			wantVersion:    "v1.0.0",
		},
		{
			name:           "Version with @latest",
			pkgWithVersion: "github.com/suapapa/gox@latest",
			wantPkg:        "github.com/suapapa/gox",
			wantVersion:    "latest",
		},
		{
			name:           "Package path with @ (imaginary case)",
			pkgWithVersion: "v1.2.3@latest",
			wantPkg:        "v1.2.3",
			wantVersion:    "latest",
		},
		{
			name:           "Implicit GitHub owner repo path",
			pkgWithVersion: "suapapa/gox",
			wantPkg:        "github.com/suapapa/gox",
			wantVersion:    "latest",
		},
		{
			name:           "Implicit GitHub owner repo cmd path with version",
			pkgWithVersion: "suapapa/gox/cmd/gox@v1.2.3",
			wantPkg:        "github.com/suapapa/gox/cmd/gox",
			wantVersion:    "v1.2.3",
		},
		{
			name:           "Keep explicit host path unchanged",
			pkgWithVersion: "golang.org/x/tools/cmd/goimports@latest",
			wantPkg:        "golang.org/x/tools/cmd/goimports",
			wantVersion:    "latest",
		},
		{
			name:           "Keep local path unchanged",
			pkgWithVersion: "./cmd/foo@latest",
			wantPkg:        "./cmd/foo",
			wantVersion:    "latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPkg, gotVer := parsePackage(tt.pkgWithVersion)
			if gotPkg != tt.wantPkg {
				t.Errorf("%s: parsePackage() gotPkg = %v, want %v", tt.name, gotPkg, tt.wantPkg)
			}
			if gotVer != tt.wantVersion {
				t.Errorf("%s: parsePackage() gotVer = %v, want %v", tt.name, gotVer, tt.wantVersion)
			}
		})
	}
}

func TestGetBinaryName(t *testing.T) {
	tests := []struct {
		name    string
		pkg     string
		wantBin string
	}{
		{
			name:    "Standard package",
			pkg:     "github.com/suapapa/gox",
			wantBin: "gox",
		},
		{
			name:    "Nested package",
			pkg:     "golang.org/x/tools/cmd/goimports",
			wantBin: "goimports",
		},
		{
			name:    "Single word package",
			pkg:     "cowsay",
			wantBin: "cowsay",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBinaryName(tt.pkg); got != tt.wantBin {
				t.Errorf("%s: getBinaryName() = %v, want %v", tt.name, got, tt.wantBin)
			}
		})
	}
}
