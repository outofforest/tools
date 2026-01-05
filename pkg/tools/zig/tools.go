package zig

import (
	"context"

	"github.com/outofforest/build/v2/pkg/tools"
	"github.com/outofforest/build/v2/pkg/types"
)

// Zig Too name.
const Zig tools.Name = "zig"

var t = []tools.Tool{
	// https://ziglang.org/download/
	tools.BinaryTool{
		Name:    Zig,
		Version: "0.13.0",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://ziglang.org/download/0.13.0/zig-linux-x86_64-0.13.0.tar.xz",
				Hash: "sha256:d45312e61ebcc48032b77bc4cf7fd6915c11fa16e4aad116b66c9468211230ea",
				Links: map[string]string{
					"bin/zig": "zig-linux-x86_64-0.13.0/zig",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://ziglang.org/download/0.13.0/zig-macos-x86_64-0.13.0.tar.xz",
				Hash: "sha256:8b06ed1091b2269b700b3b07f8e3be3b833000841bae5aa6a09b1a8b4773effd",
				Links: map[string]string{
					"bin/zig": "zig-macos-x86_64-0.13.0/zig",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://ziglang.org/download/0.13.0/zig-macos-aarch64-0.13.0.tar.xz",
				Hash: "sha256:46fae219656545dfaf4dce12fb4e8685cec5b51d721beee9389ab4194d43394c",
				Links: map[string]string{
					"bin/zig": "zig-macos-aarch64-0.13.0/zig",
				},
			},
		},
	},
}

// EnsureZig ensures that zig is available.
func EnsureZig(ctx context.Context, _ types.DepsFunc) error {
	return tools.Ensure(ctx, Zig, tools.PlatformLocal)
}

func init() {
	tools.Add(t...)
}
