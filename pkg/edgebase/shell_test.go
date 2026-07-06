package edgebase

import "testing"

func TestHostCUDACheckIncludesMinimum(t *testing.T) {
	cmd := HostCUDACheck("12.8", true)
	if want := "CUDA toolkit"; !contains(cmd, want) {
		t.Fatalf("expected %q in command", want)
	}
	if want := "12.8"; !contains(cmd, want) {
		t.Fatalf("expected minimum version %q in command", want)
	}
}

func TestPackageInventoryAvoidsRipgrep(t *testing.T) {
	cmd := PackageInventoryCommand()
	if contains(cmd, "rg ") {
		t.Fatalf("package inventory should not require ripgrep: %s", cmd)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return sub == ""
}
