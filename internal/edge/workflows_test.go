package edge

import "testing"

func TestInstallStepsProductionDefaults(t *testing.T) {
	opts := DefaultOptions()
	steps := InstallSteps(opts)
	all := joinStepCommands(steps)

	for _, want := range []string{
		"nvidia-container-toolkit",
		"libnvidia-container-tools",
		"INSTALL_K3S_CHANNEL",
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
