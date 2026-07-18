# Repository instructions

This repository owns only Ubuntu/k3s/NVIDIA substrate concerns: containerd
runtime wiring, RuntimeClass, GPU Operator, NFD, device plugin, DCGM, and CUDA
validation. Do not add Ollama, OpenAI, Open WebUI, application dashboards, or
model lifecycle behavior here.

Before completing a change run `gofmt`, `go test ./...`, `go vet ./...`,
`go build -o /tmp/k3s-nvidia-edge ./cmd/k3s-nvidia-edge`, and Helm lint for
`charts/k3s-nvidia-edge`. Preserve `--yes` gates and avoid destructive host or
cluster actions unless the operator explicitly requested them.
