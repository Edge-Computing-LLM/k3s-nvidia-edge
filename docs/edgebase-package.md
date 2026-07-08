# edgebase Go Package

`pkg/edgebase` is the public reusable Go package for the local k3s + NVIDIA GPU base layer.

It is used by the legacy `k3s-nvidia-edge` CLI and can be used by sibling
Edge-Computing-LLM projects that need to check, install, validate, or remove the
base NVIDIA/k3s substrate without importing from `internal/...`.

For operator workflows, prefer [`edge-cli`](https://github.com/Edge-Computing-LLM/edge-cli).

## Import Path

```go
import "github.com/Edge-Computing-LLM/k3s-nvidia-edge/pkg/edgebase"
```

Do not import from `github.com/Edge-Computing-LLM/k3s-nvidia-edge/internal/...`. The `internal` tree is not part of the supported API surface.

## Responsibilities

`edgebase` exposes reusable workflows for:

- host and cluster doctor checks
- k3s/NVIDIA install orchestration
- status reporting
- CUDA validation pod workflow
- GPU Operator cleanup and uninstall
- bundled Helm chart checks
- GPU capacity and RuntimeClass validation helpers
- dry-run protected command execution through `Runner`

The command-specific user interface stays in `cmd/k3s-nvidia-edge`. Shared base-layer logic stays in `pkg/edgebase`.

## Basic Usage

```go
ctx := context.Background()
opts := edgebase.DefaultOptions()
opts.Yes = false

runner := edgebase.NewRunner(opts)
if err := edgebase.Doctor(ctx, runner, opts); err != nil {
    return err
}
```

Mutating workflows remain dry-run by default unless `Options.Yes` is true.

## Downstream Usage

The organization now uses `edge-cli` as the unified CLI control plane:

- `k3s-nvidia-edge` owns Linux + k3s + NVIDIA GPU readiness.
- `llm-observability-stack` owns Ollama, Open WebUI, OpenTelemetry, dashboards, benchmarks, and notebooks.
- `edge-cli` coordinates those layers through commands such as `edge install infra`, `edge install observability`, and `edge install all`.

That keeps the base GPU logic in one maintained package instead of duplicating shell command orchestration across repositories.
