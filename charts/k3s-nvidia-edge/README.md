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

The preferred operator path is now:

```bash
edge install infra --yes
```

The legacy repository CLI can still install this wrapper chart with
`--use-local-chart`:

```bash
bin/k3s-nvidia-edge install --yes --use-local-chart
```

The wrapper chart carries the packaged GPU Operator dependency under `charts/k3s-nvidia-edge/charts/` so the local chart path remains usable after cloning the repository. CoreDNS and local-path-provisioner are intentionally left to k3s, and NFD is intentionally left to the GPU Operator chart.

Verify the bundled chart set:

```bash
bin/k3s-nvidia-edge charts
```
