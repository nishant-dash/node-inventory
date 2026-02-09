package collectors

import (
	"encoding/json"
	logger "log/slog"
	"os/exec"
	"strings"
)

type DeviceConfiguration struct {
	Driver          string `json:"driver"`
	DriverVersion   string `json:"driverversion"`
	FirmwareVersion string `json:"firmware"`
	Link            string `json:"link"`
}

type Device struct {
	ID            string              `json:"id"`
	Product       string              `json:"product"`
	Vendor        string              `json:"vendor"`
	Description   string              `json:"description"`
	PhysID        string              `json:"physid"`
	BusInfo       string              `json:"businfo"`
	LogicalName   string              `json:"logicalname"`
	Configuration DeviceConfiguration `json:"configuration"`
}

// Generic version that can unmarshal into any struct type
func getDevicesTyped[T any](class string) ([]T, error) {
	var devices []T
	cmdArgs := []string{"lshw", "-c", class, "-json"}

	out, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).Output()
	if err != nil {
		logger.Error(
			"Error collecting device information",
			"command", strings.Join(cmdArgs, " "),
			"error", err,
			"output", string(out),
		)
		return []T{}, err
	}
	if err := json.Unmarshal(out, &devices); err != nil {
		logger.Error("Error parsing lshw output", "error", err, "output", string(out))
		return []T{}, err
	}
	return devices, nil
}

// Convenience function that defaults to Device struct
func getDevices(class string) ([]Device, error) {
	return getDevicesTyped[Device](class)
}
