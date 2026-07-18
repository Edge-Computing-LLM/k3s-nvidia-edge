# Clean NVIDIA redeployment — 2026-07-18

The local Ubuntu 24.04 + k3s + NVIDIA layer was diagnosed, removed through
Helm, and installed again from this repository before the LLM stack was
reinstalled.

## Root cause corrected

k3s had retained the former automatically selected node address
`10.165.80.186`, which was no longer assigned to the host. The active USB
Ethernet address was `10.53.163.158`. The k3s journal showed node-IP lookup and
remotedialer failures; those failures caused misleading crash loops in GPU
Operator, Node Feature Discovery, kube-state-metrics, and node-exporter.

Restarting k3s allowed it to select `10.53.163.158`. Cluster networking and the
metrics API recovered before the clean chart installation.

## Final state

- node `waqasm86-thinkpad-t450s` is `Ready` at `10.53.163.158`;
- the `nvidia` RuntimeClass exists;
- the node advertises one allocatable `nvidia.com/gpu`;
- GPU Operator, toolkit, device plugin, DCGM Exporter, Node Feature Discovery,
  and validators are healthy;
- the CUDA validator completed successfully;
- `k3s-nvidia-edge` is deployed as Helm revision 1 after the clean install;
- the toolkit-triggered k3s/containerd restart occurred once and reconciled as
  expected.

The dependent LLM stack subsequently loaded Qwen with 23 of 25 repeating
layers offloaded to the GeForce 940M, confirming the end-to-end runtime path.
