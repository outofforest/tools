package buildgo

import (
	"context"

	"github.com/outofforest/build"
)

var tools = map[string]build.Tool{
	"go": {
		Name:    "go",
		Version: "1.18.1",
		URL:     "https://go.dev/dl/go1.18.1.linux-amd64.tar.gz",
		Hash:    "sha256:b3b815f47ababac13810fc6021eb73d65478e0b2db4b09d348eefad9581a2334",
		Binaries: []string{
			"go/bin/go",
			"go/bin/gofmt",
		},
	},
	"golangci": {
		Name:    "golangci",
		Version: "1.45.2",
		URL:     "https://github.com/golangci/golangci-lint/releases/download/v1.45.2/golangci-lint-1.45.2-linux-amd64.tar.gz",
		Hash:    "sha256:595ad6c6dade4c064351bc309f411703e457f8ffbb7a1806b3d8ee713333427f",
		Binaries: []string{
			"golangci-lint-1.45.2-linux-amd64/golangci-lint",
		},
	},
}

// InstallAll installs all go tools
func InstallAll(ctx context.Context) error {
	return build.InstallTools(ctx, tools)
}

// EnsureGo ensures that go is installed
func EnsureGo(ctx context.Context) error {
	return build.EnsureTool(ctx, tools["go"])
}

// EnsureGolangCI ensures that golangci is installed
func EnsureGolangCI(ctx context.Context) error {
	return build.EnsureTool(ctx, tools["golangci"])
}