package edgebase

import "testing"

func TestInstallStepsProductionDefaults(t *testing.T) {
	opts := DefaultOptions()
	steps := InstallSteps(opts)
	all := joinStepCommands(steps)

	for _, want := range []string{
		"nvidia-container-toolkit",
		"libnvidia-container-tools",
		"INSTALL_K3S_CHANNEL",
		"--write-kubeconfig-mode 0644",
		"--disable traefik",
		"--node-label gpu=nvidia",
		"helm upgrade --install gpu-operator",
		"--set driver.enabled=false",
		"--set toolkit.enabled=true",
		"--set gfd.enabled=false",
		"/var/lib/rancher/k3s/agent/etc/containerd/config.toml",
		"/run/k3s/containerd/containerd.sock",
	} {
		if !contains(all, want) {
			t.Fatalf("install steps missing %q", want)
		}
	}
}

func TestInstallStepsCanSkipBasePackages(t *testing.T) {
	opts := DefaultOptions()
	opts.SkipBasePackageInstall = true
	all := joinStepCommands(InstallSteps(opts))
	if contains(all, "apt-get install -y ca-certificates") {
		t.Fatalf("base package install should be skipped")
	}
}

func TestLocalChartInstallSteps(t *testing.T) {
	opts := DefaultOptions()
	opts.UseLocalChart = true
	steps := InstallSteps(opts)
	all := joinStepCommands(steps)

	for _, want := range []string{
		"helm dependency update './charts/k3s-nvidia-edge'",
		"helm upgrade --install k3s-nvidia-edge './charts/k3s-nvidia-edge'",
		"--set gpu-operator.driver.enabled=false",
		"--set gpu-operator.gfd.enabled=false",
	} {
		if !contains(all, want) {
			t.Fatalf("local chart install steps missing %q", want)
		}
	}
}

func TestInstallPostChecksAreDryRunSafe(t *testing.T) {
	opts := DefaultOptions()
	steps := InstallSteps(opts)

	for _, step := range steps {
		if step.Name != "wait for GPU Operator pods" && step.Name != "verify GPU capacity" {
			continue
		}
		if !step.Mutating {
			t.Fatalf("%s should be dry-run gated", step.Name)
		}
	}
}

func TestBundledChartsStepChecksAllCharts(t *testing.T) {
	opts := DefaultOptions()
	step := BundledChartsStep(opts)

	for _, want := range []string{
		"charts/k3s-nvidia-edge/Chart.yaml",
		"gpu-operator-v26.3.3.tgz",
		"helm dependency list './charts/k3s-nvidia-edge'",
	} {
		if !contains(step.Command, want) {
			t.Fatalf("bundled chart check missing %q", want)
		}
	}
	for _, duplicate := range []string{
		"charts/coredns-k3s/Chart.yaml",
		"charts/local-path-provisioner/Chart.yaml",
		"charts/node-feature-discovery/Chart.yaml",
	} {
		if contains(step.Command, duplicate) {
			t.Fatalf("bundled chart check should not require duplicate chart %q", duplicate)
		}
	}
}

func TestValidateWaitsForSucceededPod(t *testing.T) {
	manifest := CUDATestManifest(DefaultOptions().CUDATestImage)
	if !contains(manifest, "runtimeClassName: nvidia") {
		t.Fatalf("validation manifest should use nvidia runtime class")
	}
	cmd := "kubectl wait --for=jsonpath='{.status.phase}'=Succeeded pod/cuda-test --timeout=180s"
	validate := joinStepCommands([]Step{{Command: cmd}})
	if !contains(validate, "Succeeded") {
		t.Fatalf("validation should wait for completion")
	}
}

func TestCleanupStepsAvoidArchivedPackages(t *testing.T) {
	opts := DefaultOptions()
	all := joinStepCommands(CleanupLegacySteps(opts))
	for _, want := range []string{"nvidia-container-runtime", "nvidia-docker2", "gfd.enabled=false"} {
		if !contains(all, want) {
			t.Fatalf("cleanup steps missing %q", want)
		}
	}
	if contains(all, " rg ") {
		t.Fatalf("cleanup steps should not require ripgrep")
	}
}

func joinStepCommands(steps []Step) string {
	out := ""
	for _, step := range steps {
		out += step.Command + "\n"
	}
	return out
}
