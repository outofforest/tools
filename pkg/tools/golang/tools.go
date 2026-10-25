package golang

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/outofforest/build/v2/pkg/tools"
	"github.com/outofforest/build/v2/pkg/types"
	"github.com/outofforest/libexec"
	"github.com/outofforest/logger"
)

// Tool names.
const (
	Go        tools.Name = "go"
	GolangCI  tools.Name = "golangci"
	LibEVMOne tools.Name = "libevmone"
)

var t = []tools.Tool{
	// https://go.dev/dl/
	tools.BinaryTool{
		Name:    Go,
		Version: "1.27.3",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://go.dev/dl/go1.27.3.linux-amd64.tar.gz",
				Hash: "sha256:65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://go.dev/dl/go1.27.3.darwin-amd64.tar.gz",
				Hash: "sha256:b6f8e931e4ab497fcde6da8ae4499a58f9d258e55c8a80824a516eae5950a750",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://go.dev/dl/go1.27.3.darwin-arm64.tar.gz",
				Hash: "sha256:dd3908cf20959d798c1e06ebe5b72213f0505ce01373b40ce3110975a6cd99c2",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
		},
	},

	// https://github.com/golangci/golangci-lint/releases/
	tools.BinaryTool{
		Name:    GolangCI,
		Version: "2.15.0",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v2.15.0/golangci-lint-2.15.0-linux-amd64.tar.gz",
				Hash: "sha256:9b7fe6a6bf6334b98decf8d86b11415ea7889afdabdd1959fdce23af2a9ebcea",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-2.15.0-linux-amd64/golangci-lint",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v2.15.0/golangci-lint-2.15.0-darwin-amd64.tar.gz", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:767b899a28285e88ded2823d5dcc4bc5ef681173eb2d80f029a7b17ee4ece3b8",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-2.15.0-darwin-amd64/golangci-lint",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v2.15.0/golangci-lint-2.15.0-darwin-arm64.tar.gz", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:d1e409481a8ca846c4bc97e781eb5cb1bd5c23353a0b3977a1cbda1b1b6447d5",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-2.15.0-darwin-arm64/golangci-lint",
				},
			},
		},
	},

	// https://github.com/ethereum/evmone/releases
	tools.BinaryTool{
		Name:    LibEVMOne,
		Version: "0.12.0",
		Sources: tools.Sources{
			tools.PlatformDockerAMD64: {
				URL:  "https://github.com/ethereum/evmone/releases/download/v0.12.0/evmone-0.12.0-linux-x86_64.tar.gz",
				Hash: "sha256:1c7b5eba0c8c3b3b2a7a05101e2d01a13a2f84b323989a29be66285dba4136ce",
				Links: map[string]string{
					"lib/libevmone.so": "lib/libevmone.so",
				},
			},
		},
	},
}

// GoPackageTool is the tool installed using go install command.
type GoPackageTool struct {
	Name    tools.Name
	Version string
	Package string
}

// GetName returns the name of the tool.
func (gpt GoPackageTool) GetName() tools.Name {
	return gpt.Name
}

// GetVersion returns the version of the tool.
func (gpt GoPackageTool) GetVersion() string {
	return gpt.Version
}

// IsCompatible tells if tool is defined for the platform.
func (gpt GoPackageTool) IsCompatible(platform tools.Platform) (bool, error) {
	golang, err := tools.Get(Go)
	if err != nil {
		return false, err
	}
	return golang.IsCompatible(platform)
}

// Verify verifies the cheksums.
func (gpt GoPackageTool) Verify(ctx context.Context) ([]error, error) {
	return nil, nil
}

// Ensure ensures that tool is installed.
func (gpt GoPackageTool) Ensure(ctx context.Context, platform tools.Platform) error {
	binName := filepath.Base(gpt.Package)
	downloadDir := tools.ToolDownloadDir(ctx, platform, gpt)
	dst := filepath.Join("bin", binName)

	//nolint:nestif // complexity comes from trivial error-handling ifs.
	if tools.ShouldReinstall(ctx, platform, gpt, dst, binName) {
		if err := tools.Ensure(ctx, Go, platform); err != nil {
			return errors.Wrapf(err, "ensuring go failed")
		}

		cmd := exec.Command(tools.Bin(ctx, "bin/go", platform), "install", gpt.Package+"@"+gpt.Version)
		cmd.Env = append(env(ctx), "GOBIN="+downloadDir)

		if err := libexec.Exec(ctx, cmd); err != nil {
			return err
		}

		srcPath := filepath.Join(downloadDir, binName)

		binChecksum, err := tools.Checksum(srcPath)
		if err != nil {
			return err
		}

		linksDir := tools.ToolLinksDir(ctx, platform, gpt)
		dstPath := filepath.Join(linksDir, dst)
		dstPathChecksum := dstPath + ":" + binChecksum

		if err := os.Remove(dstPath); err != nil && !os.IsNotExist(err) {
			panic(err)
		}
		if err := os.Remove(dstPathChecksum); err != nil && !os.IsNotExist(err) {
			return errors.WithStack(err)
		}

		if err := os.MkdirAll(filepath.Dir(dstPath), 0o700); err != nil {
			return errors.WithStack(err)
		}
		if err := os.Chmod(srcPath, 0o700); err != nil {
			return errors.WithStack(err)
		}
		srcLinkPath, err := filepath.Rel(filepath.Dir(dstPathChecksum), filepath.Join(downloadDir, binName))
		if err != nil {
			return errors.WithStack(err)
		}
		if err := os.Symlink(srcLinkPath, dstPathChecksum); err != nil {
			return errors.WithStack(err)
		}
		if err := os.Symlink(filepath.Base(dstPathChecksum), dstPath); err != nil {
			return errors.WithStack(err)
		}
		if _, err := filepath.EvalSymlinks(dstPath); err != nil {
			return errors.WithStack(err)
		}

		logger.Get(ctx).Info("Binary installed to path", zap.String("path", dstPath))
	}

	return tools.LinkFiles(ctx, platform, gpt, []string{dst})
}

// EnsureGo ensures that go is available.
func EnsureGo(ctx context.Context, _ types.DepsFunc) error {
	return tools.Ensure(ctx, Go, tools.PlatformLocal)
}

// EnsureGolangCI ensures that go linter is available.
func EnsureGolangCI(ctx context.Context, _ types.DepsFunc) error {
	return tools.Ensure(ctx, GolangCI, tools.PlatformLocal)
}

// EnsureLibEVMOne ensures that libevmone is available.
func EnsureLibEVMOne(ctx context.Context, _ types.DepsFunc) error {
	return tools.Ensure(ctx, LibEVMOne, tools.PlatformDockerAMD64)
}

func init() {
	tools.Add(t...)
}
