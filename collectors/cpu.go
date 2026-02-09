package collectors

import (
	"encoding/json"
	logger "log/slog"
	"os/exec"
	"strconv"
)

type CPUInfo struct {
	ModelName          string
	CPUFamily          string
	ThreadsPerCore     string
	CoresPerSocket     string
	Sockets            string
	NUMANodes          string
	Model              string
	CpuFrequencyMinMHz string
	CpuFrequencyMaxMHz string
}

type LscpuOutput struct {
	Lscpu []LscpuField `json:"lscpu"`
}

type LscpuField struct {
	Field    string       `json:"field"`
	Data     string       `json:"data"`
	Children []LscpuField `json:"children,omitempty"`
}

func getCPUInfo(cpuInfo *CPUInfo) {
	cmd := exec.Command("lscpu", "-J")
	out, err := cmd.Output()

	if err != nil {
		logger.Error(
			"Error running lscpu command",
			"error", err,
		)
		return
	}

	var lscpuOutput LscpuOutput
	err = json.Unmarshal(out, &lscpuOutput)
	if err != nil {
		logger.Error(
			"Error parsing lscpu JSON output",
			"error", err,
		)
		return
	}

	extractFields(lscpuOutput.Lscpu, cpuInfo)
	logger.Debug("lscpu collection successful")
}

func extractFields(fields []LscpuField, cpuInfo *CPUInfo) {
	if len(fields) == 0 {
		return
	}

	for _, field := range fields {
		switch field.Field {
		case "Model name:":
			cpuInfo.ModelName = field.Data
		case "CPU family:":
			cpuInfo.CPUFamily = field.Data
		case "Model:":
			cpuInfo.Model = field.Data
		case "CPU min MHz:":
			cpuInfo.CpuFrequencyMinMHz = field.Data
		case "CPU max MHz:":
			cpuInfo.CpuFrequencyMaxMHz = field.Data
		case "Thread(s) per core:":
			cpuInfo.ThreadsPerCore = field.Data
		case "Core(s) per socket:":
			cpuInfo.CoresPerSocket = field.Data
		case "Socket(s):":
			cpuInfo.Sockets = field.Data
		case "NUMA node(s):":
			cpuInfo.NUMANodes = field.Data
		}

		if len(field.Children) > 0 {
			extractFields(field.Children, cpuInfo)
		}
	}
}

func (cpuInfo *CPUInfo) calculateTotalThreads() string {
	threadPerCode, err := strconv.Atoi(cpuInfo.ThreadsPerCore)
	if err != nil {
		logger.Error(
			"Error converting threads per core to integer",
			"error", err,
			"value", cpuInfo.ThreadsPerCore,
		)
		return "unknown"
	}

	corePerSocket, err := strconv.Atoi(cpuInfo.CoresPerSocket)
	if err != nil {
		logger.Error(
			"Error converting cores per socket to integer",
			"error", err,
			"value", cpuInfo.CoresPerSocket,
		)
		return "unknown"
	}

	sockets, err := strconv.Atoi(cpuInfo.Sockets)
	if err != nil {
		logger.Error(
			"Error converting sockets to integer",
			"error", err,
			"value", cpuInfo.Sockets,
		)
		return "unknown"
	}

	return strconv.Itoa(threadPerCode * corePerSocket * sockets)
}

func CPU() {
	var cpuInfo CPUInfo
	getCPUInfo(&cpuInfo)

	logger.Info(
		"Collection complete",
		"component", "cpu",
		"model_name", cpuInfo.ModelName,
		"cpu_family", cpuInfo.CPUFamily,
		"threads_per_core", cpuInfo.ThreadsPerCore,
		"cores_per_socket", cpuInfo.CoresPerSocket,
		"sockets", cpuInfo.Sockets,
		"numa_nodes", cpuInfo.NUMANodes,
		"model", cpuInfo.Model,
		"cpu_frequency_min_mhz", cpuInfo.CpuFrequencyMinMHz,
		"cpu_frequency_max_mhz", cpuInfo.CpuFrequencyMaxMHz,
		"total_threads", cpuInfo.calculateTotalThreads(),
	)
}
