# Production Readiness Notes

This CLI targets a local edge GPU workstation running Ubuntu/Xubuntu 22.04 or newer with CUDA Toolkit 12.8 or newer.

## Supported Default Profile

| Area | Production default |
|---|---|
| Kubernetes | k3s stable channel |
| Container runtime | k3s-managed containerd |
| Host GPU driver | Installed before running the CLI |
| Host CUDA toolkit | Required, minimum `12.8` |
| NVIDIA Kubernetes install | GPU Operator Helm chart |
| GPU Operator driver | `driver.enabled=false` |
| GPU Operator toolkit | `toolkit.enabled=true` |
| GPU Feature Discovery | `gfd.enabled=false` |
| GPU validation | `nvidia/cuda:12.8.1-base-ubuntu24.04` `nvidia-smi` pod |
| Helm packaging | direct `nvidia/gpu-operator` install by default; bundled wrapper chart available |

## Preflight Gates

`doctor` and `install` check:

- Ubuntu base distro and major version `>= 22`
- working `nvidia-smi`
- host CUDA Toolkit `>= 12.8` by `nvcc`, `/usr/local/cuda/version.json`, or `/usr/local/cuda/version.txt`
- required host commands
- NVIDIA package inventory
- k3s service health, for `doctor`
- Kubernetes node readiness, for `doctor/status`
- GPU Operator values
- `nvidia.com/gpu` capacity and allocatable resources

## Fresh Install

Run a dry-run first:

```bash
bin/k3s-nvidia-edge install
```

Then execute:

```bash
bin/k3s-nvidia-edge install --yes
```

For environments that want all Helm values versioned in this repository:

```bash
bin/k3s-nvidia-edge install --yes --use-local-chart
```

## Existing Cluster Hardening

```bash
bin/k3s-nvidia-edge cleanup-legacy --yes
```

This removes superseded host packages and persists `gfd.enabled=false` in Helm.

## Validation

```bash
bin/k3s-nvidia-edge validate --yes
```

Validation creates a short-lived pod with:

- `runtimeClassName: nvidia`
- `nodeSelector: nvidia.com/gpu.present=true`
- `limits.nvidia.com/gpu: 1`
- `allowPrivilegeEscalation: false`

The pod is deleted after logs are collected.

## Operational Notes

- Keep `driver.enabled=false` for workstation/laptop systems where the NVIDIA display driver is managed by Ubuntu.
- Use `--operator-driver-enabled=true` only on dedicated GPU worker nodes where GPU Operator should own driver lifecycle.
- `RuntimeClass/nvidia-legacy` may be created by the current GPU Operator runtime setup. It is not the archived `nvidia-container-runtime` host package.
- k3s uses embedded/containerized components such as CoreDNS and local-path-provisioner; the CLI does not rebuild those from source.

## Failure Handling

If validation fails:

```bash
kubectl describe pod cuda-test
kubectl logs cuda-test
kubectl get nodes -o json | jq '.items[].status.allocatable'
kubectl get pods -n gpu-operator
```

Then remove the test pod:

```bash
kubectl delete pod cuda-test --ignore-not-found
```
