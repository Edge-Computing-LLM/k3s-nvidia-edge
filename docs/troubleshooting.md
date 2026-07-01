# Troubleshooting

## CUDA Toolkit Not Found

Check:

```bash
nvcc --version
ls -l /usr/local/cuda
cat /usr/local/cuda/version.json
```

If the host intentionally has only the NVIDIA driver and no CUDA Toolkit, use:

```bash
bin/k3s-nvidia-edge doctor --require-host-cuda=false
```

The production profile expects CUDA Toolkit 12.8+.

## GPU Not Allocatable

Check:

```bash
kubectl describe node | grep -A5 -B5 nvidia.com/gpu
kubectl get pods -n gpu-operator
kubectl logs -n gpu-operator -l app=nvidia-device-plugin-daemonset --tail=100
```

Validate host GPU access:

```bash
nvidia-smi
```

## GPU Operator Pods Not Ready

Check:

```bash
kubectl get pods -n gpu-operator -o wide
kubectl describe clusterpolicy cluster-policy
helm get values gpu-operator -n gpu-operator -o yaml
```

For k3s, the toolkit env values should point to:

```text
/var/lib/rancher/k3s/agent/etc/containerd/config.toml
/run/k3s/containerd/containerd.sock
```

If your host changes networks and pods fail with `connect: no route to host` when reaching `10.43.0.1:443`, compare these values:

```bash
ip -4 route get 1.1.1.1
kubectl get endpoints kubernetes -n default -o yaml
kubectl get node -o wide
```

Pin `node-ip` in `/etc/rancher/k3s/config.yaml` when the machine regularly moves between interfaces.

## Validation Pod Stuck

Check:

```bash
kubectl describe pod cuda-test
kubectl logs cuda-test
kubectl get runtimeclass nvidia
kubectl get nodes --show-labels | grep nvidia.com/gpu.present
```

Clean up:

```bash
kubectl delete pod cuda-test --ignore-not-found
```

## Archived NVIDIA Package Still Installed

Run:

```bash
dpkg -l | grep -E 'nvidia-container-runtime|nvidia-docker'
```

Cleanup:

```bash
bin/k3s-nvidia-edge cleanup-legacy --yes
```

This does not remove `nvidia-container-toolkit` or `libnvidia-container`; those are current supported components.
