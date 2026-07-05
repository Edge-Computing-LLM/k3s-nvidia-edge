# Reference Repository Analysis

This project was built from the local repositories under:

```text
/media/waqasm86/External1/Waqas-Projects/
```

## Runtime Profile

The target profile is a single-node Ubuntu/Xubuntu 22+ edge workstation:

- Kubernetes distribution: k3s
- Container runtime: k3s-managed containerd
- Storage: local-path-provisioner
- DNS: CoreDNS
- GPU stack: NVIDIA GPU Operator
- Host driver: already installed, so GPU Operator `driver.enabled=false`
- NVIDIA runtime plumbing: NVIDIA Container Toolkit
- GPU allocation: NVIDIA Kubernetes Device Plugin
- GPU metrics: DCGM Exporter
- Hardware feature labels: Kubernetes SIGs Node Feature Discovery
- Archived standalone NVIDIA GFD: disabled via GPU Operator `gfd.enabled=false`

## Local Repository Mapping

| Local repo | Origin | Used by setup | Notes |
|---|---|---|---|
| `Project-Rancher-K3S/k3s` | `https://github.com/k3s-io/k3s.git` | Yes | Lightweight Kubernetes distribution. Bundles containerd, CoreDNS, local-path-provisioner, flannel, and core control-plane behavior. |
| `Project-Rancher-K3S/local-path-provisioner` | `https://github.com/rancher/local-path-provisioner.git` | Yes | k3s default dynamic local storage provisioner. |
| `Project-CoreDNS/coredns` | `https://github.com/coredns/coredns.git` | Yes | Cluster DNS component used by k3s as kube-dns/CoreDNS. |
| `Project-Kubernetes-Sigs/node-feature-discovery` | `https://github.com/kubernetes-sigs/node-feature-discovery.git` | Yes | Deployed by GPU Operator chart for generic hardware and kernel feature labels. |
| `Project-Kubernetes-Sigs/dra-driver-nvidia-gpu` | `https://github.com/kubernetes-sigs/dra-driver-nvidia-gpu.git` | Optional/Future | DRA path for Kubernetes 1.32+; not needed for this GeForce 940M profile. |
| `Project-Nvidia/gpu-operator` | `https://github.com/NVIDIA/gpu-operator.git` | Yes | Helm/operator control plane for NVIDIA GPU Kubernetes components. |
| `Project-Nvidia/k8s-device-plugin` | `https://github.com/NVIDIA/k8s-device-plugin.git` | Yes | Exposes `nvidia.com/gpu`; also contains current integrated GPU Feature Discovery implementation. |
| `Project-Nvidia/nvidia-container-toolkit` | `https://github.com/NVIDIA/nvidia-container-toolkit.git` | Yes | Configures containerd/runtime integration for NVIDIA GPUs. |
| `Project-Nvidia/libnvidia-container` | `https://github.com/NVIDIA/libnvidia-container.git` | Yes | Low-level library and CLI consumed by NVIDIA Container Toolkit. |
| `Project-Nvidia/dcgm-exporter` | `https://github.com/NVIDIA/dcgm-exporter.git` | Yes | GPU metrics exporter deployed by GPU Operator. |
| `Project-Nvidia/DCGM` | `https://github.com/NVIDIA/DCGM.git` | Indirect | GPU management/monitoring library used by DCGM Exporter. |
| `Project-Nvidia/go-dcgm` | `https://github.com/NVIDIA/go-dcgm.git` | Indirect | Go bindings used in DCGM-related tooling. |
| `Project-Nvidia/cuda-samples` | `https://github.com/NVIDIA/cuda-samples.git` | Validation reference | CUDA sample source; not required by the installer. |
| `Project-Cloudflare/cloudflared` | `https://github.com/cloudflare/cloudflared.git` | Optional | Local tunnel path for exposing DCGM or observability endpoints during validation. |

## Why The CLI Uses Helm GPU Operator

The local `gpu-operator` repository includes a development chart, but the working machine is running released chart `gpu-operator-v26.3.3`. The CLI therefore installs the released NVIDIA Helm chart by default and uses local source repositories as implementation references.

This is safer for a workstation than installing `main-latest` development manifests from a source tree.

The repository also includes `charts/k3s-nvidia-edge`, a wrapper chart that pins the released GPU Operator dependency and version-controls the k3s-specific values. It intentionally does not vendor CoreDNS, local-path-provisioner, or standalone Node Feature Discovery charts because the live cluster already receives CoreDNS/local-path from k3s and NFD from the GPU Operator chart.

## Archived NVIDIA Tool Handling

The CLI removes superseded host packages:

```bash
apt-get remove -y nvidia-container-runtime nvidia-docker2 nvidia-docker || true
```

It also disables GPU Operator GFD reconciliation:

```bash
helm upgrade gpu-operator nvidia/gpu-operator \
  -n gpu-operator \
  --reuse-values \
  --set gfd.enabled=false \
  --wait
```

It intentionally keeps Kubernetes SIGs Node Feature Discovery. That is not NVIDIA's archived standalone GPU Feature Discovery repo.

## Validation Contract

A prepared cluster must pass:

```bash
kubectl get nodes -o json | jq '.items[].status.allocatable["nvidia.com/gpu"]'
```

Expected value on the current laptop:

```text
"1"
```

It must also pass a short-lived CUDA pod:

```bash
runtimeClassName: nvidia
resources:
  limits:
    nvidia.com/gpu: 1
```
