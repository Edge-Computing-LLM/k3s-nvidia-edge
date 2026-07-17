# Live NVIDIA validation — 2026-07-17

The Layer 1 chart was installed on Ubuntu 24.04.3, k3s v1.36.2+k3s1, and a
GeForce 940M using the host driver 580.95.05.

Validated components:

- NVIDIA Container Toolkit 1.19.1
- GPU Operator v26.3.3
- Node Feature Discovery v0.18.3
- NVIDIA device plugin v0.19.3
- NVIDIA DCGM Exporter
- NVIDIA operator and CUDA validators
- `RuntimeClass/nvidia`
- one allocatable `nvidia.com/gpu`

The CUDA validation pod ran `nvidia-smi` successfully inside
`nvidia/cuda:12.8.1-base-ubuntu24.04`. This confirms the complete path from host
driver through k3s containerd, RuntimeClass, device plugin, scheduler, and container.

The toolkit reconciliation restarted k3s once. Early device-plugin events briefly
reported that the NVIDIA runtime was not configured; they resolved after the
toolkit restart and all final daemonsets were healthy.
