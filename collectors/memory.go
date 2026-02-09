package collectors

import (
	logger "log/slog"
)

type MemoryDevice struct {
	Device
	Size     uint64 `json:"size"`
	Capacity uint64 `json:"capacity"`
	Units    string `json:"units"`
	Version  string `json:"version"`
}

func Memory() {
	devices, err := getDevicesTyped[MemoryDevice]("memory")
	if err != nil {
		return
	}

	// Find the system memory device
	for _, dev := range devices {
		if dev.ID == "firmware" {
			logger.Info(
				"Collection complete",
				"component", "memory",
				"description", dev.Description,
				"physid", dev.PhysID,
				"vendor", dev.Vendor,
				"version", dev.Version,
				"size_bytes", dev.Size,
				"capacity_bytes", dev.Capacity,
				"units", dev.Units,
			)
			break
		}
	}
}
