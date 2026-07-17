# Live Validation - 2026-07-08

Validated as Layer 1 on local Xubuntu 24, single-node k3s, NVIDIA GeForce 940M.

Tested commands:

```bash
edge install infra --skip-base-package-install --skip-toolkit-install --skip-k3s-install --yes
edge validate infra
bin/k3s-nvidia-edge status
```

Results:

- Installed/upgraded Helm release `k3s-nvidia-edge` in namespace `gpu-operator`.
- GPU Operator, NVIDIA device plugin, DCGM exporter, Node Feature Discovery, and validator pods became healthy.
- `RuntimeClass/nvidia` existed.
- Node advertised `nvidia.com/gpu: 1` capacity and allocatable.
- CUDA validation pod ran `nvidia-smi` successfully.
- Layer 2 was installed and uninstalled without removing this base layer.

This confirms the repo remains independently useful as the NVIDIA/k3s substrate while also working through `edge-cli`.
