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

## Build The CLI

```bash
git clone https://github.com/Edge-Computing-LLM/k3s-nvidia-gpu.git
cd k3s-nvidia-gpu
make check
```

Optional:

```bash
make install-local
```

## Preflight

```bash
bin/k3s-nvidia-edge doctor
```

On a fresh machine where k3s is not installed yet, run a dry-run install instead:

```bash
bin/k3s-nvidia-edge install
```

## Install

```bash
bin/k3s-nvidia-edge install --yes
```

Default values:

```text
k3s channel: stable
GPU Operator: v26.3.3
CUDA test image: nvidia/cuda:12.8.0-base-ubuntu24.04
minimum host CUDA: 12.8
GPU Operator driver: disabled
GPU Operator GFD: disabled
```

## Validate

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
