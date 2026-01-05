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
		Version: "0.15.2",
		Sources: tools.Sources{
			tools.PlatformLinuxAMD64: {
				URL:  "https://ziglang.org/download/0.15.2/zig-x86_64-linux-0.15.2.tar.xz",
				Hash: "sha256:02aa270f183da276e5b5920b1dac44a63f1a49e55050ebde3aecc9eb82f93239",
				Links: map[string]string{
					"bin/zig": "zig-x86_64-linux-0.15.2/zig",
				},
			},
			tools.PlatformDarwinAMD64: {
				URL:  "https://ziglang.org/download/0.15.2/zig-x86_64-macos-0.15.2.tar.xz",
				Hash: "sha256:375b6909fc1495d16fc2c7db9538f707456bfc3373b14ee83fdd3e22b3d43f7f",
				Links: map[string]string{
					"bin/zig": "zig-x86_64-macos-0.15.2/zig",
				},
			},
			tools.PlatformDarwinARM64: {
				URL:  "https://ziglang.org/download/0.15.2/zig-aarch64-macos-0.15.2.tar.xz",
				Hash: "sha256:3cc2bab367e185cdfb27501c4b30b1b0653c28d9f73df8dc91488e66ece5fa6b",
				Links: map[string]string{
					"bin/zig": "zig-aarch64-macos-0.15.2/zig",
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
