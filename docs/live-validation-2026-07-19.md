# Live validation: 2026-07-19

The repository gates passed on Ubuntu 24.04 with Go 1.26.5, Helm 4.2.3, k3s
v1.36.2+k3s1, and NVIDIA GPU Operator v26.3.3:

- `go mod verify`, formatting, unit tests, race tests, vet, and build
- `govulncheck ./...` with no reachable vulnerabilities
- Helm lint and chart rendering
- host `nvidia-smi`, RuntimeClass discovery, and GPU capacity inspection

The live cluster audit correctly found a post-reboot network fault. k3s still
advertised an InternalIP and Flannel interface that were no longer assigned to
the host. Pod traffic to the Kubernetes API therefore failed with `no route to
host`, causing GPU Operator, Node Feature Discovery, metrics-server,
kube-state-metrics, and node-exporter readiness failures. GPU allocation and
the active Ollama workload remained available, which demonstrated why capacity
alone is not a sufficient health signal.

The CLI now checks both invariants before CUDA validation:

1. The advertised k3s InternalIP must exist on a current host interface.
2. GPU Operator pods must be Running and ready, or successfully Completed.

No host addresses, kubeconfig data, Secrets, prompts, responses, or pod logs
are recorded in this document.
