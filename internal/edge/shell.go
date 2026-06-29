package edge

import "fmt"

func UbuntuVersionCheck() string {
	return `set -euo pipefail
. /etc/os-release
if [ "${ID:-}" != "ubuntu" ]; then
  echo "unsupported distro: ${ID:-unknown}; expected ubuntu/xubuntu base"
  exit 1
fi
major="${VERSION_ID%%.*}"
if [ "$major" -lt 22 ]; then
  echo "unsupported Ubuntu version: ${VERSION_ID}; expected 22.04 or newer"
  exit 1
fi
echo "Ubuntu ${VERSION_ID} supported"`
}

func HostCUDACheck(minVersion string, required bool) string {
	if !required {
		return `if command -v nvcc >/dev/null 2>&1; then nvcc --version | tail -n 1; else echo "host CUDA toolkit check skipped"; fi`
	}
	return fmt.Sprintf(`set -euo pipefail
min=%s
version=""
if command -v nvcc >/dev/null 2>&1; then
  version="$(nvcc --version | sed -n 's/.*release \([0-9][0-9.]*\).*/\1/p' | head -n 1)"
elif [ -f /usr/local/cuda/version.json ]; then
  version="$(grep -o '"version"[[:space:]]*:[[:space:]]*"[^"]*"' /usr/local/cuda/version.json | head -n 1 | sed 's/.*"version"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/' | cut -d. -f1,2)"
elif [ -f /usr/local/cuda/version.txt ]; then
  version="$(sed -n 's/.*CUDA Version \([0-9][0-9.]*\).*/\1/p' /usr/local/cuda/version.txt | head -n 1)"
fi
if [ -z "$version" ]; then
  echo "CUDA toolkit not found; install CUDA toolkit ${min}+ or run with --require-host-cuda=false"
  exit 1
fi
lowest="$(printf '%%s\n%%s\n' "$min" "$version" | sort -V | head -n 1)"
if [ "$lowest" != "$min" ]; then
  echo "CUDA toolkit $version is below required $min"
  exit 1
fi
echo "CUDA toolkit $version satisfies >= $min"`, shellQuote(minVersion))
}

func GPUOperatorReadyCheck() string {
	return `set -euo pipefail
kubectl wait --for=condition=Ready pod -n gpu-operator -l app.kubernetes.io/component=gpu-operator --timeout=180s || true
kubectl get pods -n gpu-operator
bad="$(kubectl get pods -n gpu-operator --no-headers | awk '$3!="Running" && $3!="Completed" {print}')"
if [ -n "$bad" ]; then
  echo "$bad"
  exit 1
fi
notready="$(kubectl get pods -n gpu-operator --no-headers | awk '$3=="Running" {split($2,a,"/"); if (a[1] != a[2]) print}')"
if [ -n "$notready" ]; then
  echo "$notready"
  exit 1
fi`
}

func GPUCapacityCheck() string {
	return `set -euo pipefail
gpu="$(kubectl get nodes -o jsonpath='{.items[0].status.allocatable.nvidia\.com/gpu}')"
if [ -z "$gpu" ] || [ "$gpu" = "<no value>" ] || [ "$gpu" = "0" ]; then
  echo "nvidia.com/gpu is not allocatable"
  kubectl describe node | grep -A5 -B5 'nvidia.com/gpu' || true
  exit 1
fi
echo "nvidia.com/gpu allocatable: $gpu"`
}

func PackageInventoryCommand() string {
	return "dpkg -l | grep -E 'nvidia-docker|nvidia-container-runtime|nvidia-container-toolkit|libnvidia-container|cuda-toolkit|cuda-compiler|cuda-runtime' || true"
}

func CUDATestManifest(image string) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: cuda-test
  labels:
    app.kubernetes.io/name: k3s-nvidia-edge
    app.kubernetes.io/component: cuda-validation
spec:
  restartPolicy: Never
  runtimeClassName: nvidia
  nodeSelector:
    nvidia.com/gpu.present: "true"
  containers:
  - name: cuda-test
    image: %s
    command: ["nvidia-smi"]
    securityContext:
      allowPrivilegeEscalation: false
    resources:
      limits:
        nvidia.com/gpu: 1`, image)
}
