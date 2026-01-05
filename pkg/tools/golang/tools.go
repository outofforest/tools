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
		Version: "1.25.5",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://go.dev/dl/go1.25.5.linux-amd64.tar.gz",
				Hash: "sha256:9e9b755d63b36acf30c12a9a3fc379243714c1c6d3dd72861da637f336ebb35b",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://go.dev/dl/go1.25.5.darwin-amd64.tar.gz",
				Hash: "sha256:b69d51bce599e5381a94ce15263ae644ec84667a5ce23d58dc2e63e2c12a9f56",
				Links: map[string]string{
					"bin/go":    "go/bin/go",
					"bin/gofmt": "go/bin/gofmt",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://go.dev/dl/go1.25.5.darwin-arm64.tar.gz",
				Hash: "sha256:bed8ebe824e3d3b27e8471d1307f803fc6ab8e1d0eb7a4ae196979bd9b801dd3",
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
		Version: "2.7.2",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v2.7.2/golangci-lint-2.7.2-linux-amd64.tar.gz",
				Hash: "sha256:ce46a1f1d890e7b667259f70bb236297f5cf8791a9b6b98b41b283d93b5b6e88",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-2.7.2-linux-amd64/golangci-lint",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v2.7.2/golangci-lint-2.7.2-darwin-amd64.tar.gz", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:6966554840a02229a14c52641bc38c2c7a14d396f4c59ba0c7c8bb0675ca25c9",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-2.7.2-darwin-amd64/golangci-lint",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://github.com/golangci/golangci-lint/releases/download/v2.7.2/golangci-lint-2.7.2-darwin-arm64.tar.gz", //nolint:lll // breaking down urls is not beneficial
				Hash: "sha256:6ce86a00e22b3709f7b994838659c322fdc9eae09e263db50439ad4f6ec5785c",
				Links: map[string]string{
					"bin/golangci-lint": "golangci-lint-2.7.2-darwin-arm64/golangci-lint",
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
