package edge

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
	SkipToolkitInstall     bool
	SkipK3sInstall         bool
	SkipGPUOperatorInstall bool
	DriverEnabled          bool
	K3sChannel             string
	K3sExec                string
}

func DefaultOptions() Options {
	return Options{
		Sudo:                   true,
		ReferenceRoot:          "/media/waqasm86/External1/Project-Llamatelemetry/Project-Llamatelemetry-End-to-End",
		GPUOperatorVersion:     "v26.3.3",
		CUDATestImage:          "nvidia/cuda:12.8.0-base-ubuntu24.04",
		MinCUDAVersion:         "12.8",
		RequireHostCUDA:        true,
		DisableGFD:             true,
		DriverEnabled:          false,
		K3sChannel:             "stable",
		K3sExec:                "server --write-kubeconfig-mode 644",
		SkipToolkitInstall:     false,
		SkipK3sInstall:         false,
		SkipGPUOperatorInstall: false,
	}
}
