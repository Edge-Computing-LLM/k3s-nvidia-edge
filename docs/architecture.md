# Architecture

`k3s-nvidia-edge` is a small orchestration CLI plus a reusable Go package. It does not replace k3s, Helm, NVIDIA GPU Operator, or NVIDIA Container Toolkit. It codifies a known-good sequence for preparing a single-node Ubuntu/Xubuntu edge workstation for GPU workloads.

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

The reusable workflow code lives in `pkg/edgebase`; `cmd/k3s-nvidia-edge` is the command-line wrapper around that package. Other Edge-Computing-LLM repositories can import `pkg/edgebase` to reuse base-layer checks without copying shell orchestration or importing private `internal/...` packages.

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
