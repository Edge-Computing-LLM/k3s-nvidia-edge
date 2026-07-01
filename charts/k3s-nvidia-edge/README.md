# k3s-nvidia-edge Helm Chart

This chart wraps NVIDIA GPU Operator with defaults for a local Ubuntu/Xubuntu k3s edge node:

- k3s-managed containerd paths
- host-managed NVIDIA driver, `driver.enabled=false`
- NVIDIA Container Toolkit enabled
- NVIDIA Device Plugin enabled
- DCGM exporter enabled
- Node Feature Discovery enabled through GPU Operator
- standalone GPU Feature Discovery disabled

Install from the repository root:

```bash
helm dependency update charts/k3s-nvidia-edge
helm upgrade --install k3s-nvidia-edge charts/k3s-nvidia-edge \
  -n gpu-operator --create-namespace --wait
```

The CLI keeps direct upstream GPU Operator install as the default for existing clusters. Use `--use-local-chart` to install this wrapper chart:

```bash
bin/k3s-nvidia-edge install --yes --use-local-chart
```
