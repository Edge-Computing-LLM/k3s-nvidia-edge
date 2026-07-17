# K3S NVIDIA GPU Edge Setup

`k3s-nvidia-edge` is a reusable Go package and Helm profile for installing, configuring, validating, and cleaning up a local Ubuntu 22+ k3s cluster with NVIDIA GPU support and CUDA Toolkit 12.8+.

The reusable base workflows live in `pkg/edgebase`. New operator workflows should use the unified organization CLI, [`edge-cli`](https://github.com/Edge-Computing-LLM/edge-cli), with commands such as `edge install infra`, `edge validate infra`, and `edge install all`. The legacy `k3s-nvidia-edge` binary remains available during migration.

This is Layer 1 of the `Edge-Computing-LLM` platform. It owns the local Linux + k3s + NVIDIA substrate. `llm-observability-stack` is Layer 2 and must not install GPU Operator, NVIDIA device plugin, or DCGM exporter in the main local NVIDIA path.

[`qwen-gguf-observability`](https://github.com/Edge-Computing-LLM/qwen-gguf-observability)
may read the resulting node, RuntimeClass, and GPU-capacity status as a runtime
evidence companion. It does not install or manage this infrastructure layer.

This layer is conditional. `edge install all --accelerator auto` invokes the
NVIDIA setup only when `nvidia-smi` confirms working host hardware. On CPU-only
hosts, `edge-cli` installs or validates basic k3s without deploying this chart and
then selects the CPU observability profile. Use `--accelerator nvidia` when the
GPU layer is mandatory and fallback would hide an infrastructure problem.

## Documentation

- [edge-cli migration and repo role](docs/edge-cli-migration.md)
- [Installation guide](docs/installation.md)
- [Commands](docs/commands.md)
- [Architecture](docs/architecture.md)
- [Live validation - 2026-07-08](docs/live-validation-2026-07-08.md)
- [Live validation - 2026-07-17](docs/live-validation-2026-07-17.md)
- [Reusable edgebase Go package](docs/edgebase-package.md)
- [Production readiness](docs/production-readiness.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Reference repository analysis](docs/reference-repos.md)
- [Example edge profile](configs/edge-profile.env.example)
- [Bundled Helm chart](charts/k3s-nvidia-edge/README.md)
- [Contributing](CONTRIBUTING.md)
- [Security](SECURITY.md)
- [Qwen GGUF runtime evidence companion](https://github.com/Edge-Computing-LLM/qwen-gguf-observability)

The default profile matches the working local Xubuntu 24 setup:

- k3s with containerd
- host NVIDIA driver already installed
- host CUDA Toolkit 12.8 or newer installed
- NVIDIA Container Toolkit installed on the host
- NVIDIA GPU Operator installed by Helm
- GPU Operator `driver.enabled=false`
- GPU Operator `toolkit.enabled=true`
- GPU Operator `gfd.enabled=false`
- bundled local Helm wrapper for the NVIDIA GPU Operator only
- CUDA validation with an `nvidia-smi` pod
- optional local wrapper Helm chart at `charts/k3s-nvidia-edge`

## Build

Prerequisite:

```bash
go version
```

Build and verify:

```bash
make check
```

The binary is written to:

```bash
bin/k3s-nvidia-edge
```

Optional local install:

```bash
make install-local
```

## Reusable Go Package

Downstream projects should import:

```go
import "github.com/Edge-Computing-LLM/k3s-nvidia-edge/pkg/edgebase"
```

The package exposes `DefaultOptions`, `Runner`, and workflows such as `Doctor`, `Install`, `Status`, `Validate`, `CleanupLegacy`, and `Uninstall`. It is the supported API for reusing the base-layer logic. Do not import from `internal/...`.

## Commands

New operator workflows should normally use `edge-cli`:

```bash
edge doctor
edge install infra --yes
edge validate infra
edge status
```

The legacy binary remains available for compatibility and package development:

```bash
bin/k3s-nvidia-edge doctor
bin/k3s-nvidia-edge status
bin/k3s-nvidia-edge charts
bin/k3s-nvidia-edge repos
bin/k3s-nvidia-edge print-commands
```

Mutating commands are dry-run by default. Add `--yes` to execute them:

```bash
bin/k3s-nvidia-edge install --yes
bin/k3s-nvidia-edge cleanup-legacy --yes
bin/k3s-nvidia-edge validate --yes
```

Important production flags:

```bash
--min-cuda-version 12.8
--require-host-cuda=true
--gpu-operator-version v26.3.3
--cuda-test-image nvidia/cuda:12.8.1-base-ubuntu24.04
--operator-driver-enabled=false
--use-local-chart=false
--local-chart ./charts/k3s-nvidia-edge
--skip-base-package-install=false
```

## End-To-End Install

Preferred through `edge-cli`:

```bash
edge install infra --yes
edge validate infra
```

This command sequence supports an otherwise empty local k3s cluster that only
has default k3s system components such as CoreDNS and local-path-provisioner.
After validation passes, install the LLM observability layer with
`edge install observability --yes` or use `edge install all --yes` for the full
ordered flow.

Legacy direct command:

```bash
bin/k3s-nvidia-edge install --yes
```

This performs:

1. Ubuntu 22+, NVIDIA driver, and CUDA Toolkit 12.8+ preflight checks.
2. Base package installation.
3. NVIDIA Container Toolkit package setup.
4. Removal of superseded `nvidia-container-runtime` / `nvidia-docker*` packages.
5. Helm installation if missing.
6. k3s installation.
7. user-owned kubeconfig setup.
8. NVIDIA Helm repo setup.
9. GPU Operator install/upgrade.
10. CUDA pod validation.

To install through the bundled wrapper chart instead of installing `nvidia/gpu-operator` directly:

```bash
bin/k3s-nvidia-edge install --yes --use-local-chart
```

For an already prepared local machine where k3s, Helm, CUDA, and NVIDIA Container Toolkit are installed:

```bash
bin/k3s-nvidia-edge install --yes --sudo=false --use-local-chart --skip-base-package-install --skip-toolkit-install --skip-k3s-install
```

## Existing k3s Cluster Cleanup

```bash
bin/k3s-nvidia-edge cleanup-legacy --yes
```

This removes legacy host packages and disables the GPU Operator GFD workload:

```yaml
gfd:
  enabled: false
```

It intentionally keeps:

- `nvidia-container-toolkit`
- `libnvidia-container`
- `nvidia-device-plugin`
- `dcgm-exporter`
- GPU Operator managed `node-feature-discovery`
- `RuntimeClass/nvidia`
- `RuntimeClass/nvidia-legacy` generated by current GPU Operator runtime behavior

## Validation

```bash
bin/k3s-nvidia-edge validate --yes
```

The CLI creates this short-lived pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cuda-test
spec:
  restartPolicy: Never
  runtimeClassName: nvidia
  nodeSelector:
    nvidia.com/gpu.present: "true"
  containers:
  - name: cuda-test
    image: nvidia/cuda:12.8.1-base-ubuntu24.04
    command: ["nvidia-smi"]
    resources:
      limits:
        nvidia.com/gpu: 1
```

The pod is deleted after logs are printed.

## Reference Repositories

The CLI was built against the local reference repositories under:

```text
/media/waqasm86/External1/Waqas-Projects/Project-Kubernetes-Sigs/
/media/waqasm86/External1/Waqas-Projects/Project-CoreDNS/
/media/waqasm86/External1/Waqas-Projects/Project-Rancher-K3S/
/media/waqasm86/External1/Waqas-Projects/Project-Nvidia/
/media/waqasm86/External1/Waqas-Projects/Project-Cloudflare/
```

Main tool mapping:

| Component | Repository |
|---|---|
| k3s | `https://github.com/k3s-io/k3s` |
| local-path-provisioner | `https://github.com/rancher/local-path-provisioner` |
| CoreDNS | `https://github.com/coredns/coredns` |
| Node Feature Discovery | `https://github.com/kubernetes-sigs/node-feature-discovery` |
| NVIDIA GPU Operator | `https://github.com/NVIDIA/gpu-operator` |
| NVIDIA Kubernetes Device Plugin | `https://github.com/NVIDIA/k8s-device-plugin` |
| NVIDIA Container Toolkit | `https://github.com/NVIDIA/nvidia-container-toolkit` |
| libnvidia-container | `https://github.com/NVIDIA/libnvidia-container` |
| DCGM Exporter | `https://github.com/NVIDIA/dcgm-exporter` |
| DCGM | `https://github.com/NVIDIA/DCGM` |
| NVIDIA DRA Driver | `https://github.com/kubernetes-sigs/dra-driver-nvidia-gpu` |
| cloudflared | `https://github.com/cloudflare/cloudflared` |

## Notes

This tool assumes the NVIDIA display/compute driver is already installed on the host unless `--operator-driver-enabled=true` is explicitly passed.

For laptop/GeForce edge setups, keep `driver.enabled=false` and use the host driver. That is the least disruptive path for a local workstation.
