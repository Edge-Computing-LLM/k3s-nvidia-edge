# Installation Guide

This guide targets Ubuntu/Xubuntu 22.04 or newer with an NVIDIA GPU.

## Host Prerequisites

Install the NVIDIA display/compute driver and CUDA Toolkit 12.8 or newer before running the full installer.

Verify:

```bash
nvidia-smi
nvcc --version
```

The CLI also accepts systems where CUDA version metadata is available at:

```text
/usr/local/cuda/version.json
/usr/local/cuda/version.txt
```

## Preferred Installation Through edge-cli

For normal operator workflows, install and run the unified CLI:

```bash
git clone https://github.com/Edge-Computing-LLM/edge-cli.git
cd edge-cli
go build -o edge ./cmd/edge
sudo install -m 0755 edge /usr/local/bin/edge
edge install infra --yes
edge validate infra
```

`edge-cli` uses this repository as the infrastructure layer and keeps
organization-wide install/status/validation commands in one place.

## Build The Legacy CLI

```bash
git clone https://github.com/Edge-Computing-LLM/k3s-nvidia-edge.git
cd k3s-nvidia-edge
make check
```

Optional:

```bash
make install-local
```

## Preflight

Preferred:

```bash
edge doctor
```

Legacy direct command:

```bash
bin/k3s-nvidia-edge doctor
```

On a fresh machine where k3s is not installed yet, run a dry-run install instead:

```bash
bin/k3s-nvidia-edge install
```

## Install

Preferred:

```bash
edge install infra --yes
```

Legacy direct command:

```bash
bin/k3s-nvidia-edge install --yes
```

Default values:

```text
k3s channel: stable
GPU Operator: v26.3.3
CUDA test image: nvidia/cuda:12.8.1-base-ubuntu24.04
minimum host CUDA: 12.8
GPU Operator driver: disabled
GPU Operator GFD: disabled
```

The direct upstream GPU Operator chart is the default. The repository also includes a local wrapper chart:

```bash
helm dependency update charts/k3s-nvidia-edge
helm upgrade --install k3s-nvidia-edge charts/k3s-nvidia-edge \
  -n gpu-operator --create-namespace --wait
```

Use it through the CLI with:

```bash
bin/k3s-nvidia-edge install --yes --use-local-chart
```

The repository carries only the wrapper chart and the packaged GPU Operator dependency required by the local NVIDIA profile. k3s already owns CoreDNS and local-path-provisioner, and GPU Operator deploys its own NFD dependency. Verify chart availability with:

```bash
bin/k3s-nvidia-edge charts
```

On an already prepared local machine, deploy only the local Helm profile:

```bash
bin/k3s-nvidia-edge install --yes --sudo=false --use-local-chart --skip-base-package-install --skip-toolkit-install --skip-k3s-install
```

## Validate

Preferred:

```bash
edge validate infra
```

Legacy direct command:

```bash
bin/k3s-nvidia-edge validate --yes
```

Expected result:

```text
NVIDIA-SMI output from inside a Kubernetes pod
pod "cuda-test" deleted
```

## Existing Cluster

For an already working k3s + NVIDIA GPU Operator setup:

```bash
bin/k3s-nvidia-edge cleanup-legacy --yes
bin/k3s-nvidia-edge status
bin/k3s-nvidia-edge validate --yes
```
