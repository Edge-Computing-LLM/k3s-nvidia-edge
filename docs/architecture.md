# Architecture

`k3s-nvidia-edge` is a small orchestration CLI. It does not replace k3s, Helm, NVIDIA GPU Operator, or NVIDIA Container Toolkit. It codifies a known-good sequence for preparing a single-node Ubuntu/Xubuntu edge workstation for GPU workloads.

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

The CLI uses generated shell steps with dry-run protection:

- read-only commands run immediately
- mutating commands print dry-run output unless `--yes` is passed
- host-level commands use `sudo` by default when the current user is not root

The implementation is intentionally dependency-light. Runtime commands avoid non-standard tools such as `rg`; production hosts need ordinary POSIX shell tools plus `curl`, `apt-get`, `systemctl`, `kubectl`, `helm`, and `jq`.

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
