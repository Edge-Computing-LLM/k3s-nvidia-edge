# Architecture

`k3s-nvidia-edge` is the infrastructure layer for the `Edge-Computing-LLM` platform. It provides a reusable Go package and local Helm profile for preparing a single-node Ubuntu/Xubuntu edge workstation for NVIDIA GPU workloads. It does not replace k3s, Helm, NVIDIA GPU Operator, or NVIDIA Container Toolkit.

The primary operator CLI is now [`edge-cli`](https://github.com/Edge-Computing-LLM/edge-cli). This repository keeps the base-layer implementation and assets that `edge-cli` coordinates.

This repository is Layer 1 in the organization platform. It owns the NVIDIA GPU
substrate for local k3s. Layer 2, `llm-observability-stack`, may observe DCGM
metrics and schedule Ollama with `RuntimeClass/nvidia`, but it must not install
GPU Operator, NVIDIA device plugin, Node Feature Discovery, or DCGM exporter in
the main local NVIDIA path.

`gguf-observability` is not Layer 3 and owns no infrastructure. It is a
read-only evidence consumer that verifies selected outputs of Layer 1 together
with the selected GGUF model runtime deployed by Layer 2.

## Components

```text
Ubuntu/Xubuntu 22+
  |
  |-- NVIDIA host driver
  |-- CUDA Toolkit 12.8+
  |-- NVIDIA Container Toolkit
  |
  `-- k3s
      |
      |-- containerd
      |-- CoreDNS
      |-- local-path-provisioner
      |
      `-- NVIDIA GPU Operator
          |
          |-- nvidia-container-toolkit daemonset
          |-- nvidia-device-plugin daemonset
          |-- nvidia-dcgm-exporter daemonset
          |-- node-feature-discovery
          `-- validator jobs/pods
```

The repository intentionally vendors only the local wrapper chart and the released GPU Operator dependency archive:

```text
charts/k3s-nvidia-edge
charts/k3s-nvidia-edge/charts/gpu-operator-v26.3.3.tgz
```

CoreDNS and local-path-provisioner are not vendored here because k3s already deploys and owns them in `kube-system`. Standalone Node Feature Discovery is not vendored because the GPU Operator chart deploys the active NFD components needed by the NVIDIA profile.

## Default GPU Policy

The default GPU Operator profile is intentionally conservative for laptops and local edge machines:

```text
driver.enabled=false
toolkit.enabled=true
gfd.enabled=false
```

`driver.enabled=false` keeps the display/compute driver managed by Ubuntu. This avoids the GPU Operator replacing a working laptop driver stack.

`gfd.enabled=false` avoids deploying the standalone GPU Feature Discovery workload. Kubernetes SIGs Node Feature Discovery remains enabled because the GPU Operator chart uses it for generic node feature labels.

## Command Model

The reusable workflow code lives in `pkg/edgebase`; `cmd/k3s-nvidia-edge` is the legacy command-line wrapper around that package. New operator workflows should use `edge-cli`, which exposes module commands such as `edge install infra` and `edge validate infra`.

Other Edge-Computing-LLM repositories can import `pkg/edgebase` to reuse base-layer checks without copying shell orchestration or importing private `internal/...` packages.

The CLI and package use generated shell steps with dry-run protection:

- read-only commands run immediately
- mutating commands print dry-run output unless `--yes` is passed
- host-level commands use `sudo` by default when the current user is not root

The implementation is intentionally dependency-light. Runtime commands avoid non-standard tools such as `rg`; production hosts need ordinary POSIX shell tools plus `curl`, `apt-get`, `systemctl`, `kubectl`, `helm`, and `jq`.

## Public Package Boundary

```text
cmd/k3s-nvidia-edge        CLI parsing and user-facing command dispatch
pkg/edgebase               reusable base workflows, options, runner, and shell command builders
charts/k3s-nvidia-edge     optional local Helm wrapper for GPU Operator
```

`pkg/edgebase` is the supported import surface for sibling projects. The package currently exposes the same workflows used by the existing CLI: doctor, install, status, validate, cleanup legacy resources, uninstall, repository inventory, bundled chart checks, and print-command helpers.

`edge-cli` re-creates the operator workflow in its own Go module while preserving this repository as the infrastructure source of truth.

## Validation Contract

A healthy cluster must satisfy:

```text
kubectl node Ready
nvidia.com/gpu allocatable > 0
RuntimeClass/nvidia exists
GPU Operator pods Running or Completed
CUDA validation pod can run nvidia-smi
```

The validation pod is short-lived and requests exactly one GPU.

When Ollama already holds the single GPU, the evidence companion reads existing
capacity, RuntimeClass, pod, Ollama, and `nvidia-smi` status. It does not launch
a competing CUDA validation pod.

A mostly empty k3s cluster is expected before this layer is installed. CoreDNS
and local-path-provisioner remain owned by k3s; this repository adds the NVIDIA
runtime and GPU readiness layer on top.
