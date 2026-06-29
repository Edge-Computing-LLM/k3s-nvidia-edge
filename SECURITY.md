# Security

This CLI runs host package installation and Kubernetes administration commands. Review dry-run output before using `--yes`.

## Secrets

Do not commit:

- kubeconfig files
- tokens
- private keys
- `.env` files
- generated evidence or logs containing credentials

The repository `.gitignore` excludes common local secret and runtime artifacts.

## Sudo

Host-level commands use `sudo` when required. The CLI does not store passwords.

## Reporting Issues

Open an issue in the GitHub repository with:

- OS version
- NVIDIA driver version
- CUDA Toolkit version
- k3s version
- GPU Operator version
- sanitized command output
