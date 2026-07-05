package edge

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

func Doctor(ctx context.Context, r *Runner, opts Options) error {
	steps := []Step{
		{Name: "Ubuntu 22+ check", Command: UbuntuVersionCheck()},
		{Name: "Kernel", Command: "uname -a"},
		{Name: "NVIDIA driver and GPU", Command: "nvidia-smi"},
		{Name: "CUDA toolkit version", Command: HostCUDACheck(opts.MinCUDAVersion, opts.RequireHostCUDA)},
		{Name: "Required commands", Command: "missing=0; for c in curl apt-get systemctl kubectl helm jq grep awk sort sed; do command -v $c >/dev/null && echo \"$c: $(command -v $c)\" || { echo \"$c: missing\"; missing=1; }; done; exit $missing"},
		{Name: "NVIDIA host packages", Command: PackageInventoryCommand()},
		{Name: "k3s service", Command: "systemctl is-active k3s"},
		{Name: "cluster nodes", Command: "kubectl get nodes -o wide || true"},
		{Name: "GPU capacity", Command: "kubectl get nodes -o json | jq '.items[] | {name:.metadata.name, capacity_gpu:.status.capacity[\"nvidia.com/gpu\"], allocatable_gpu:.status.allocatable[\"nvidia.com/gpu\"]}' || true"},
		{Name: "GPU Operator values", Command: "helm get values gpu-operator -n gpu-operator -o yaml || true"},
	}
	for _, step := range steps {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

func HostPreflight(ctx context.Context, r *Runner, opts Options) error {
	steps := []Step{
		{Name: "Ubuntu 22+ check", Command: UbuntuVersionCheck()},
		{Name: "Kernel", Command: "uname -a"},
		{Name: "NVIDIA driver and GPU", Command: "nvidia-smi"},
		{Name: "CUDA toolkit version", Command: HostCUDACheck(opts.MinCUDAVersion, opts.RequireHostCUDA)},
		{Name: "Required host commands", Command: "missing=0; for c in curl apt-get systemctl grep awk sort sed; do command -v $c >/dev/null && echo \"$c: $(command -v $c)\" || { echo \"$c: missing\"; missing=1; }; done; exit $missing"},
		{Name: "NVIDIA host packages", Command: PackageInventoryCommand()},
	}
	for _, step := range steps {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

func Install(ctx context.Context, r *Runner, opts Options) error {
	if err := HostPreflight(ctx, r, opts); err != nil {
		return err
	}

	for _, step := range InstallSteps(opts) {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return Validate(ctx, r, opts)
}

func Status(ctx context.Context, r *Runner, opts Options) error {
	steps := []Step{
		{Name: "k3s resources", Command: "kubectl get all -A"},
		{Name: "nodes", Command: "kubectl get nodes -o wide"},
		{Name: "runtime classes", Command: "kubectl get runtimeclass || true"},
		{Name: "running images", Command: "kubectl get pods -A -o jsonpath='{range .items[*]}{.metadata.namespace}{\"/\"}{.metadata.name}{\"\\t\"}{range .spec.containers[*]}{.image}{\" \"}{end}{\"\\n\"}{end}'"},
		{Name: "helm releases", Command: "helm list -A"},
		{Name: "GPU Operator values", Command: "helm get values gpu-operator -n gpu-operator -o yaml || true"},
		{Name: "GPU capacity", Command: "kubectl get nodes -o json | jq '.items[] | {name:.metadata.name, capacity_gpu:.status.capacity[\"nvidia.com/gpu\"], allocatable_gpu:.status.allocatable[\"nvidia.com/gpu\"]}'"},
	}
	for _, step := range steps {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

func Validate(ctx context.Context, r *Runner, opts Options) error {
	return r.Run(ctx, Step{
		Name:     "CUDA validation pod",
		Mutating: true,
		Command: fmt.Sprintf(`kubectl delete pod cuda-test --ignore-not-found
kubectl apply -f - <<'EOF'
%s
EOF
kubectl wait --for=jsonpath='{.status.phase}'=Succeeded pod/cuda-test --timeout=180s
kubectl logs cuda-test
kubectl delete pod cuda-test`, CUDATestManifest(opts.CUDATestImage)),
	})
}

func CleanupLegacy(ctx context.Context, r *Runner, opts Options) error {
	for _, step := range CleanupLegacySteps(opts) {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return Status(ctx, r, opts)
}

func Uninstall(ctx context.Context, r *Runner, opts Options) error {
	steps := []Step{
		{Name: "uninstall GPU Operator", Mutating: true, Command: "helm uninstall gpu-operator -n gpu-operator --wait || true"},
		{Name: "delete gpu-operator namespace", Mutating: true, Command: "kubectl delete namespace gpu-operator --ignore-not-found"},
	}
	if opts.UninstallK3s {
		steps = append(steps, Step{Name: "uninstall k3s", Host: true, Mutating: true, Command: "/usr/local/bin/k3s-uninstall.sh"})
	}
	for _, step := range steps {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

func Repos(ctx context.Context, r *Runner, opts Options) error {
	root := opts.ReferenceRoot
	command := fmt.Sprintf(`for d in %s/Project-Kubernetes-Sigs/* %s/Project-CoreDNS/* %s/Project-Rancher-K3S/* %s/Project-Nvidia/* %s/Project-Cloudflare/*; do
  [ -d "$d/.git" ] || continue
  printf 'repo: %%s\n' "$d"
  printf 'origin: '; git -C "$d" remote get-url origin || true
  printf 'branch: '; git -C "$d" rev-parse --abbrev-ref HEAD || true
  printf 'commit: '; git -C "$d" rev-parse --short HEAD || true
  printf 'go modules: '; find "$d" -maxdepth 3 -name go.mod | wc -l
  printf 'helm charts: '; find "$d" -maxdepth 5 -name Chart.yaml | wc -l
  printf '\n'
done`, shellQuote(root), shellQuote(root), shellQuote(root), shellQuote(root), shellQuote(root))
	return r.Run(ctx, Step{Name: "local reference repo inventory", Command: command})
}

func Charts(ctx context.Context, r *Runner, opts Options) error {
	return r.Run(ctx, BundledChartsStep(opts))
}

func PrintCommands(opts Options) {
	fmt.Println("# Install steps")
	for _, step := range InstallSteps(opts) {
		fmt.Printf("\n## %s\n%s\n", step.Name, step.Command)
	}
	fmt.Println("\n# Cleanup legacy steps")
	for _, step := range CleanupLegacySteps(opts) {
		fmt.Printf("\n## %s\n%s\n", step.Name, step.Command)
	}
	fmt.Println("\n# Validation")
	fmt.Printf("kubectl apply -f - <<'EOF'\n%s\nEOF\n", CUDATestManifest(opts.CUDATestImage))
}

func InstallSteps(opts Options) []Step {
	var steps []Step
	if !opts.SkipBasePackageInstall {
		steps = append(steps, Step{
			Name:     "install base packages",
			Host:     true,
			Mutating: true,
			Command:  "apt-get update && apt-get install -y ca-certificates curl gnupg lsb-release jq apt-transport-https software-properties-common",
		})
	}
	if !opts.SkipToolkitInstall {
		steps = append(steps, Step{
			Name:     "install NVIDIA Container Toolkit",
			Host:     true,
			Mutating: true,
			Command: `set -euo pipefail
install -d -m 0755 /usr/share/keyrings /etc/apt/sources.list.d
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' > /etc/apt/sources.list.d/nvidia-container-toolkit.list
apt-get update
apt-get install -y nvidia-container-toolkit libnvidia-container-tools libnvidia-container1
apt-get remove -y nvidia-container-runtime nvidia-docker2 nvidia-docker || true`,
		})
	}
	if !opts.SkipK3sInstall {
		steps = append(steps, Step{
			Name:     "install k3s",
			Host:     true,
			Mutating: true,
			Command:  fmt.Sprintf("curl -sfL https://get.k3s.io | INSTALL_K3S_CHANNEL=%s INSTALL_K3S_EXEC=%s sh -", shellQuote(opts.K3sChannel), shellQuote(opts.K3sExec)),
		})
	}
	steps = append(steps,
		Step{Name: "install Helm if missing", Host: true, Mutating: true, Command: "command -v helm >/dev/null 2>&1 || curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash"},
		BundledChartsStep(opts),
		Step{Name: "prepare kubeconfig", Mutating: true, Command: `set -euo pipefail
mkdir -p "$HOME/.kube"
if [ -r /etc/rancher/k3s/k3s.yaml ]; then
  cp /etc/rancher/k3s/k3s.yaml "$HOME/.kube/config"
else
  sudo cp /etc/rancher/k3s/k3s.yaml "$HOME/.kube/config"
  sudo chown "$(id -u):$(id -g)" "$HOME/.kube/config"
fi
chmod 600 "$HOME/.kube/config"`},
		Step{Name: "wait for k3s node", Mutating: false, Command: "kubectl wait --for=condition=Ready node --all --timeout=180s"},
		Step{Name: "add NVIDIA Helm repo", Mutating: true, Command: "helm repo add nvidia https://helm.ngc.nvidia.com/nvidia || true && helm repo update"},
	)
	if !opts.SkipGPUOperatorInstall {
		gfd := "false"
		if !opts.DisableGFD {
			gfd = "true"
		}
		steps = append(steps, GPUOperatorInstallStep(opts, gfd))
	}
	steps = append(steps,
		Step{Name: "wait for GPU Operator pods", Command: GPUOperatorReadyCheck(), Mutating: true},
		Step{Name: "verify GPU capacity", Command: GPUCapacityCheck(), Mutating: true},
	)
	return steps
}

func BundledChartsStep(opts Options) Step {
	paths := []string{
		opts.LocalChartPath,
		filepath.Join(opts.LocalChartPath, "charts", "gpu-operator-"+opts.GPUOperatorVersion+".tgz"),
	}

	var checks []string
	checks = append(checks, "set -euo pipefail")
	checks = append(checks, "command -v helm >/dev/null")
	checks = append(checks, fmt.Sprintf("test -f %s", shellQuote(filepath.Join(paths[0], "Chart.yaml"))))
	checks = append(checks, fmt.Sprintf("helm lint %s >/dev/null", shellQuote(paths[0])))
	checks = append(checks, fmt.Sprintf("test -f %s", shellQuote(paths[1])))
	checks = append(checks, fmt.Sprintf("tar -tzf %s >/dev/null", shellQuote(paths[1])))
	checks = append(checks, fmt.Sprintf("helm dependency list %s", shellQuote(opts.LocalChartPath)))
	return Step{
		Name:    "verify bundled Helm charts",
		Command: strings.Join(checks, "\n"),
	}
}

func GPUOperatorInstallStep(opts Options, gfd string) Step {
	if opts.UseLocalChart {
		return Step{
			Name:     "install or upgrade k3s NVIDIA edge Helm chart",
			Mutating: true,
			Command: fmt.Sprintf(`helm dependency update %s
helm upgrade --install k3s-nvidia-edge %s \
  -n gpu-operator --create-namespace \
  --set gpu-operator.driver.enabled=%t \
  --set gpu-operator.toolkit.enabled=true \
  --set gpu-operator.gfd.enabled=%s \
  --set gpu-operator.toolkit.env[0].name=CONTAINERD_CONFIG \
  --set gpu-operator.toolkit.env[0].value=/var/lib/rancher/k3s/agent/etc/containerd/config.toml \
  --set gpu-operator.toolkit.env[1].name=CONTAINERD_SOCKET \
  --set gpu-operator.toolkit.env[1].value=/run/k3s/containerd/containerd.sock \
  --set gpu-operator.toolkit.env[2].name=RUNTIME_CONFIG_SOURCE \
  --set-string gpu-operator.toolkit.env[2].value=file=/var/lib/rancher/k3s/agent/etc/containerd/config.toml \
  --wait`, shellQuote(opts.LocalChartPath), shellQuote(opts.LocalChartPath), opts.DriverEnabled, gfd),
		}
	}
	return Step{
		Name:     "install or upgrade GPU Operator",
		Mutating: true,
		Command: fmt.Sprintf(`helm upgrade --install gpu-operator nvidia/gpu-operator \
  -n gpu-operator --create-namespace \
  --version %s \
  --set driver.enabled=%t \
  --set toolkit.enabled=true \
  --set gfd.enabled=%s \
  --set toolkit.env[0].name=CONTAINERD_CONFIG \
  --set toolkit.env[0].value=/var/lib/rancher/k3s/agent/etc/containerd/config.toml \
  --set toolkit.env[1].name=CONTAINERD_SOCKET \
  --set toolkit.env[1].value=/run/k3s/containerd/containerd.sock \
  --set toolkit.env[2].name=RUNTIME_CONFIG_SOURCE \
  --set-string toolkit.env[2].value=file=/var/lib/rancher/k3s/agent/etc/containerd/config.toml \
  --wait`, opts.GPUOperatorVersion, opts.DriverEnabled, gfd),
	}
}

func CleanupLegacySteps(opts Options) []Step {
	steps := []Step{
		{Name: "remove superseded NVIDIA packages", Host: true, Mutating: true, Command: "apt-get remove -y nvidia-container-runtime nvidia-docker2 nvidia-docker || true"},
	}
	if opts.DisableGFD {
		command := fmt.Sprintf("helm upgrade gpu-operator nvidia/gpu-operator -n gpu-operator --version %s --reuse-values --set gfd.enabled=false --wait", opts.GPUOperatorVersion)
		if opts.UseLocalChart {
			command = fmt.Sprintf("helm dependency update %s && helm upgrade k3s-nvidia-edge %s -n gpu-operator --reuse-values --set gpu-operator.gfd.enabled=false --wait", shellQuote(opts.LocalChartPath), shellQuote(opts.LocalChartPath))
		}
		steps = append(steps, Step{
			Name:     "disable GPU Feature Discovery in GPU Operator",
			Mutating: true,
			Command:  command,
		})
	}
	steps = append(steps,
		Step{Name: "verify no standalone GFD workload", Command: "kubectl get ds,pods -n gpu-operator | grep -E 'gpu-feature-discovery|NAME' || true"},
		Step{Name: "verify host packages", Command: PackageInventoryCommand()},
	)
	return steps
}

func LocalRepoPaths(root string) []string {
	return []string{
		filepath.Join(root, "Project-Rancher-K3S", "k3s"),
		filepath.Join(root, "Project-Rancher-K3S", "local-path-provisioner"),
		filepath.Join(root, "Project-CoreDNS", "coredns"),
		filepath.Join(root, "Project-Kubernetes-Sigs", "node-feature-discovery"),
		filepath.Join(root, "Project-Kubernetes-Sigs", "dra-driver-nvidia-gpu"),
		filepath.Join(root, "Project-Nvidia", "gpu-operator"),
		filepath.Join(root, "Project-Nvidia", "k8s-device-plugin"),
		filepath.Join(root, "Project-Nvidia", "nvidia-container-toolkit"),
		filepath.Join(root, "Project-Nvidia", "libnvidia-container"),
		filepath.Join(root, "Project-Nvidia", "dcgm-exporter"),
		filepath.Join(root, "Project-Nvidia", "DCGM"),
		filepath.Join(root, "Project-Cloudflare", "cloudflared"),
	}
}
