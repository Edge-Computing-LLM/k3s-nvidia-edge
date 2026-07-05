# k3s-nvidia-edge Helm Chart

This chart wraps NVIDIA GPU Operator with defaults for a local Ubuntu/Xubuntu k3s edge node:

- k3s-managed containerd paths
- host-managed NVIDIA driver, `driver.enabled=false`
- NVIDIA Container Toolkit enabled
- NVIDIA Device Plugin enabled
- DCGM exporter enabled
- Node Feature Discovery enabled through GPU Operator
- standalone GPU Feature Discovery disabled
- optional k3s CoreDNS and local-path-provisioner chart dependencies, disabled by default
- optional standalone Node Feature Discovery chart dependency, disabled by default

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

The wrapper chart carries packaged dependencies under `charts/k3s-nvidia-edge/charts/` so the local chart path remains usable after cloning the repository. CoreDNS and local-path-provisioner are disabled by default because k3s already deploys those components in `kube-system`; standalone NFD is disabled because the GPU Operator chart deploys its own NFD dependency.

Verify the bundled chart set:

```bash
bin/k3s-nvidia-edge charts
```
