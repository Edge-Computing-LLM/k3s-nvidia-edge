# Contributing

## Development

Run the full local check before committing:

```bash
make check
```

This runs:

```text
gofmt
go vet
go test
go build
```

## Runtime Safety

Mutating CLI commands must remain dry-run by default. A command that changes host packages, k3s, Helm releases, Kubernetes resources, or local kubeconfig must require `--yes`.

## Compatibility

The default profile targets:

- Ubuntu/Xubuntu 22.04+
- CUDA Toolkit 12.8+
- k3s with containerd
- NVIDIA GPU Operator
- host-managed NVIDIA driver

Keep workstation/laptop safety in mind. Do not enable GPU Operator driver management by default.
