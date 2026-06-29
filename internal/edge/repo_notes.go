package edge

type RepoNote struct {
	Name    string
	Repo    string
	Purpose string
}

var RepoNotes = []RepoNote{
	{Name: "k3s", Repo: "https://github.com/k3s-io/k3s", Purpose: "lightweight Kubernetes distribution bundling containerd, CoreDNS, local-path provisioning, and core control-plane components"},
	{Name: "local-path-provisioner", Repo: "https://github.com/rancher/local-path-provisioner", Purpose: "default local persistent-volume provisioner used by k3s"},
	{Name: "CoreDNS", Repo: "https://github.com/coredns/coredns", Purpose: "cluster DNS service used through k3s kube-dns/CoreDNS deployment"},
	{Name: "node-feature-discovery", Repo: "https://github.com/kubernetes-sigs/node-feature-discovery", Purpose: "hardware and kernel feature label discovery deployed by the GPU Operator chart"},
	{Name: "gpu-operator", Repo: "https://github.com/NVIDIA/gpu-operator", Purpose: "operator and Helm chart that reconciles NVIDIA Kubernetes GPU components"},
	{Name: "k8s-device-plugin", Repo: "https://github.com/NVIDIA/k8s-device-plugin", Purpose: "exposes nvidia.com/gpu resources and current integrated GPU Feature Discovery implementation"},
	{Name: "nvidia-container-toolkit", Repo: "https://github.com/NVIDIA/nvidia-container-toolkit", Purpose: "configures container runtimes so containers can access NVIDIA GPUs"},
	{Name: "libnvidia-container", Repo: "https://github.com/NVIDIA/libnvidia-container", Purpose: "low-level NVIDIA container runtime library and CLI used by the toolkit"},
	{Name: "dcgm-exporter", Repo: "https://github.com/NVIDIA/dcgm-exporter", Purpose: "GPU metrics exporter deployed by GPU Operator"},
	{Name: "DCGM", Repo: "https://github.com/NVIDIA/DCGM", Purpose: "NVIDIA GPU monitoring/management library used by dcgm-exporter"},
	{Name: "dra-driver-nvidia-gpu", Repo: "https://github.com/kubernetes-sigs/dra-driver-nvidia-gpu", Purpose: "future/optional DRA path; not required for the GeForce/k3s default profile"},
}
