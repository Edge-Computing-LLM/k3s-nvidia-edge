package edgebase

type Options struct {
	Yes                    bool
	Sudo                   bool
	Verbose                bool
	ReferenceRoot          string
	GPUOperatorVersion     string
	CUDATestImage          string
	MinCUDAVersion         string
	RequireHostCUDA        bool
	DisableGFD             bool
	UninstallK3s           bool
	SkipBasePackageInstall bool
	SkipToolkitInstall     bool
	SkipK3sInstall         bool
	SkipGPUOperatorInstall bool
	DriverEnabled          bool
	UseLocalChart          bool
	LocalChartPath         string
	K3sChannel             string
	K3sExec                string
}

func DefaultOptions() Options {
	return Options{
		Sudo:                   true,
		ReferenceRoot:          "/media/waqasm86/External1/Waqas-Projects/Project-Linux-Kubernetes-Nvidia",
		GPUOperatorVersion:     "v26.3.3",
		CUDATestImage:          "nvidia/cuda:12.8.1-base-ubuntu24.04",
		MinCUDAVersion:         "12.8",
		RequireHostCUDA:        true,
		DisableGFD:             true,
		DriverEnabled:          false,
		UseLocalChart:          false,
		LocalChartPath:         "./charts/k3s-nvidia-edge",
		K3sChannel:             "stable",
		K3sExec:                "server --write-kubeconfig-mode 0644 --disable traefik --disable servicelb --disable metrics-server --node-label gpu=nvidia --node-label workload=edge-ai",
		SkipBasePackageInstall: false,
		SkipToolkitInstall:     false,
		SkipK3sInstall:         false,
		SkipGPUOperatorInstall: false,
	}
}
