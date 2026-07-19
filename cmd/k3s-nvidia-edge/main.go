package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Edge-Computing-LLM/k3s-nvidia-edge/pkg/edgebase"
)

const usage = `k3s-nvidia-edge is an Ubuntu 22+ CLI for local k3s + NVIDIA GPU setup.

Deprecation notice:
  The unified organization CLI is now edge-cli:
    https://github.com/Edge-Computing-LLM/edge-cli

  Use "edge install infra", "edge validate infra", and "edge status" for new
  workflows. This legacy command remains available during migration.

Usage:
  k3s-nvidia-edge <command> [flags]

Commands:
  doctor          Check OS, GPU, packages, k3s, kubectl, helm, and cluster GPU state
  install         Install/configure k3s + NVIDIA GPU Operator profile
  status          Print live k3s/NVIDIA status
  validate        Run an nvidia-smi CUDA pod and remove it after success
  cleanup-legacy Remove superseded NVIDIA host packages and disable GFD in GPU Operator
  uninstall       Remove GPU Operator, optionally uninstall k3s
  repos           Inventory local reference GitHub repos
  charts          Verify bundled Helm charts are present and renderable
  print-commands  Print the exact shell commands used by install/cleanup/validate

Global flags:
  --yes                 execute mutating commands; otherwise mutating commands are dry-run
  --sudo                use sudo for host-level commands (default true)
  --verbose             print command output while commands run
  --reference-root DIR  root containing Project-Rancher, Project-Nvidia, etc.
  --min-cuda-version    minimum host CUDA toolkit version for doctor/install
  --require-host-cuda   require host CUDA toolkit, not only driver/container CUDA

Examples:
  k3s-nvidia-edge doctor
  k3s-nvidia-edge install --yes
  k3s-nvidia-edge cleanup-legacy --yes
  k3s-nvidia-edge validate --yes
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	opts := edgebase.DefaultOptions()
	fs.BoolVar(&opts.Yes, "yes", false, "execute mutating commands")
	fs.BoolVar(&opts.Sudo, "sudo", true, "use sudo for host-level commands")
	fs.BoolVar(&opts.Verbose, "verbose", false, "print command output while commands run")
	fs.StringVar(&opts.ReferenceRoot, "reference-root", opts.ReferenceRoot, "root containing local cloned reference repos")
	fs.StringVar(&opts.GPUOperatorVersion, "gpu-operator-version", opts.GPUOperatorVersion, "GPU Operator Helm chart version")
	fs.StringVar(&opts.CUDATestImage, "cuda-test-image", opts.CUDATestImage, "CUDA image used for validation")
	fs.StringVar(&opts.MinCUDAVersion, "min-cuda-version", opts.MinCUDAVersion, "minimum host CUDA toolkit version")
	fs.BoolVar(&opts.RequireHostCUDA, "require-host-cuda", opts.RequireHostCUDA, "require host CUDA toolkit in doctor/install")
	fs.BoolVar(&opts.DisableGFD, "disable-gfd", true, "disable GPU Feature Discovery in GPU Operator")
	fs.BoolVar(&opts.UninstallK3s, "k3s", false, "with uninstall: also uninstall k3s")
	fs.BoolVar(&opts.SkipBasePackageInstall, "skip-base-package-install", false, "with install: skip base apt package setup")
	fs.BoolVar(&opts.SkipToolkitInstall, "skip-toolkit-install", false, "with install: skip host NVIDIA Container Toolkit package setup")
	fs.BoolVar(&opts.SkipK3sInstall, "skip-k3s-install", false, "with install: skip k3s install command")
	fs.BoolVar(&opts.SkipGPUOperatorInstall, "skip-gpu-operator-install", false, "with install: skip GPU Operator install")
	fs.BoolVar(&opts.DriverEnabled, "operator-driver-enabled", false, "install GPU Operator driver component instead of using host driver")
	fs.BoolVar(&opts.UseLocalChart, "use-local-chart", false, "with install/cleanup: install the bundled Helm chart instead of the upstream chart directly")
	fs.StringVar(&opts.LocalChartPath, "local-chart", opts.LocalChartPath, "path to bundled k3s-nvidia-edge Helm chart")
	fs.StringVar(&opts.K3sChannel, "k3s-channel", opts.K3sChannel, "k3s install channel")
	fs.StringVar(&opts.K3sExec, "k3s-exec", opts.K3sExec, "INSTALL_K3S_EXEC value")

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	ctx := context.Background()
	r := edgebase.NewRunner(opts)
	var err error

	switch strings.ToLower(cmd) {
	case "doctor":
		err = edgebase.Doctor(ctx, r, opts)
	case "install":
		err = edgebase.Install(ctx, r, opts)
	case "status":
		err = edgebase.Status(ctx, r, opts)
	case "validate":
		err = edgebase.Validate(ctx, r, opts)
	case "cleanup-legacy":
		err = edgebase.CleanupLegacy(ctx, r, opts)
	case "uninstall":
		err = edgebase.Uninstall(ctx, r, opts)
	case "repos":
		err = edgebase.Repos(ctx, r, opts)
	case "charts":
		err = edgebase.Charts(ctx, r, opts)
	case "print-commands":
		edgebase.PrintCommands(opts)
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", cmd, usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
