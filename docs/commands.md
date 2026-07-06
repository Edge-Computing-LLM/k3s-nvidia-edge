# Commands

All commands are run from the repository root after building:

```bash
make check
```

The CLI commands are thin wrappers over the public `pkg/edgebase` package. Sibling projects should reuse that package instead of copying command logic or importing from `internal/...`.

## doctor

```bash
bin/k3s-nvidia-edge doctor
```

Checks the current host and cluster:

- Ubuntu 22+
- kernel information
- `nvidia-smi`
- CUDA Toolkit version, default minimum `12.8`
- required commands
- NVIDIA package inventory
- k3s service state
- Kubernetes node status
- GPU capacity and allocatable resources
- GPU Operator values

## install

Dry-run:

```bash
bin/k3s-nvidia-edge install
```

Execute:

```bash
bin/k3s-nvidia-edge install --yes
```

Installs or configures:

- base apt packages
- NVIDIA Container Toolkit
- k3s
- Helm, if missing
- kubeconfig for the invoking user
- NVIDIA Helm repository
- NVIDIA GPU Operator
- optional bundled `charts/k3s-nvidia-edge` wrapper chart with `--use-local-chart`
- CUDA validation pod

For an already prepared local workstation, skip host setup and deploy only through the local chart:

```bash
bin/k3s-nvidia-edge install --yes --sudo=false --use-local-chart --skip-base-package-install --skip-toolkit-install --skip-k3s-install
```

## status

```bash
bin/k3s-nvidia-edge status
```

Prints cluster resources, nodes, runtime classes, pod images, Helm releases, GPU Operator values, and GPU capacity.

## validate

Dry-run:

```bash
bin/k3s-nvidia-edge validate
```

Execute:

```bash
bin/k3s-nvidia-edge validate --yes
```

Creates a short-lived CUDA pod, prints `nvidia-smi`, and deletes the pod.
The command waits for pod phase `Succeeded`, which is more reliable for one-shot validation pods than waiting for a long-lived Ready condition.

## cleanup-legacy

Dry-run:

```bash
bin/k3s-nvidia-edge cleanup-legacy
```

Execute:

```bash
bin/k3s-nvidia-edge cleanup-legacy --yes
```

Removes superseded host packages:

```text
nvidia-container-runtime
nvidia-docker2
nvidia-docker
```

Also persists:

```text
gfd.enabled=false
```

in the GPU Operator Helm release.

## repos

```bash
bin/k3s-nvidia-edge repos
```

Inventories the local reference repositories used while designing the CLI.

## charts

```bash
bin/k3s-nvidia-edge charts
```

Verifies that the bundled wrapper chart and packaged GPU Operator dependency are present and renderable:

- `charts/k3s-nvidia-edge`
- `charts/k3s-nvidia-edge/charts/gpu-operator-v26.3.3.tgz`

## print-commands

```bash
bin/k3s-nvidia-edge print-commands
```

Prints the generated install, cleanup, and validation shell commands.

This is also useful for downstream CLI authors who want to show the base-layer operations that `edgebase` would execute.

## uninstall

Dry-run:

```bash
bin/k3s-nvidia-edge uninstall
```

Remove GPU Operator:

```bash
bin/k3s-nvidia-edge uninstall --yes
```

Remove GPU Operator and k3s:

```bash
bin/k3s-nvidia-edge uninstall --yes --k3s
```
