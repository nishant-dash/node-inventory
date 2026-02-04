package collectors

import (
	"encoding/json"
	logger "log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type NetworkDeviceConfiguration struct {
	Driver          string `json:"driver"`
	DriverVersion   string `json:"driverversion"`
	FirmwareVersion string `json:"firmware"`
	Link            string `json:"link"`
}

type NetworkDevice struct {
	ID            string                     `json:"id"`
	Product       string                     `json:"product"`
	Vendor        string                     `json:"vendor"`
	Description   string                     `json:"description"`
	PhysID        string                     `json:"physid"`
	BusInfo       string                     `json:"businfo"`
	LogicalName   string                     `json:"logicalname"`
	Configuration NetworkDeviceConfiguration `json:"configuration"`
}

func getVFInfoByInterface(ifName string) (numVFs, maxVFs int, isPF bool, err error) {
	basePath := filepath.Join("/sys/class/net", ifName, "device")

	if data, err := os.ReadFile(filepath.Join(basePath, "sriov_totalvfs")); err == nil {
		isPF = true
		maxVFs, _ = strconv.Atoi(strings.TrimSpace(string(data)))
	}

	if data, err := os.ReadFile(filepath.Join(basePath, "sriov_numvfs")); err == nil {
		numVFs, _ = strconv.Atoi(strings.TrimSpace(string(data)))
	}

	return
}

func getDevLinkMode(devicePCI string) (string, error) {
	formattedPCI := strings.Replace(devicePCI, "@", "/", 1)
	cmd := exec.Command("devlink", "dev", "eswitch", "show", formattedPCI)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(
			"Could not get devlink information",
			"error", err,
			"output", string(out),
			"businfo", formattedPCI,
		)
	}
	return string(out), err
}

func NIC() error {
	out, err := exec.Command("lshw", "-c", "net", "-json").Output()
	if err != nil {
		logger.Debug("Error executing lshw command", "error", err)
		return err
	}

	var devices []NetworkDevice
	if err := json.Unmarshal(out, &devices); err != nil {
		logger.Debug("Error parsing lshw output", "error", err)
		return err
	}
	for _, device := range devices {
		if device.Configuration.Link == "yes" {
			devlinkMode, err := getDevLinkMode(device.BusInfo)
			if err != nil {
				devlinkMode = "unknown"
			}

			logger.Info(
				"",
				"component", "nic",
				"id", device.ID,
				"product", device.Product,
				"vendor", device.Vendor,
				"description", device.Description,
				"businfo", device.BusInfo,
				"logicalname", device.LogicalName,
				"driver", device.Configuration.Driver,
				"driverversion", device.Configuration.DriverVersion,
				"firmware", device.Configuration.FirmwareVersion,
				"mode", devlinkMode,
			)
		}
	}

	return nil
}
