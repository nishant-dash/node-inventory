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

type NetworkDevice struct {
	Device
	numVFs      int
	maxVFs      int
	isPF        bool
	devlinkMode string
}

type DevlinkOutput struct {
	Dev map[string]DevlinkInfo `json:"dev"`
}

type DevlinkInfo struct {
	Mode       string `json:"mode"`
	InlineMode string `json:"inline_mode"`
	EncapMode  string `json:"encap_mode"`
}

func (nd *NetworkDevice) getVFInfoByInterface() (err error) {
	basePath := filepath.Join("/sys/class/net", nd.LogicalName, "device")

	if data, err := os.ReadFile(filepath.Join(basePath, "sriov_totalvfs")); err == nil {
		nd.isPF = true
		nd.maxVFs, _ = strconv.Atoi(strings.TrimSpace(string(data)))
	}

	if data, err := os.ReadFile(filepath.Join(basePath, "sriov_numvfs")); err == nil {
		nd.numVFs, _ = strconv.Atoi(strings.TrimSpace(string(data)))
	}

	return err
}

func (nd *NetworkDevice) getDevLinkMode() (err error) {
	formattedPCI := strings.Replace(nd.BusInfo, "@", "/", 1)
	cmd := exec.Command("devlink", "-j", "dev", "eswitch", "show", formattedPCI)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(
			"Could not get devlink information",
			"error", err,
			"output", string(out),
			"businfo", formattedPCI,
		)
		nd.devlinkMode = "unknown"
	} else {
		var dlOutput DevlinkOutput
		if err := json.Unmarshal(out, &dlOutput); err != nil {
			logger.Error(
				"Error parsing devlink json output",
				"error", err,
				"output", string(out),
			)
			return err
		}
		nd.devlinkMode = dlOutput.Dev[formattedPCI].Mode
	}

	return err
}

func NIC() {
	devices, err := getDevices("net")
	if err != nil {
		logger.Info("No network devices found exiting NIC collection", "error", err)
		return
	}

	networkDevices := make([]NetworkDevice, len(devices))
	for i, device := range devices {
		networkDevices[i] = NetworkDevice{Device: device}
	}

	for _, nd := range networkDevices {
		if nd.Configuration.Link == "yes" {
			err := nd.getDevLinkMode()
			if err != nil {
				logger.Error(
					"Error getting devlink information",
					"interface", nd.LogicalName,
					"error", err,
				)
			}

			err = nd.getVFInfoByInterface()
			if err != nil {
				logger.Error(
					"Error getting VF information",
					"interface", nd.LogicalName,
					"error", err,
				)
			}

			logger.Info(
				"Collection complete",
				"component", "nic",
				"id", nd.ID,
				"product", nd.Product,
				"vendor", nd.Vendor,
				"description", nd.Description,
				"businfo", nd.BusInfo,
				"logicalname", nd.LogicalName,
				"driver", nd.Configuration.Driver,
				"driverversion", nd.Configuration.DriverVersion,
				"firmware", nd.Configuration.FirmwareVersion,
				"mode", nd.devlinkMode,
				"num_vfs", nd.numVFs,
				"max_vfs", nd.maxVFs,
				"is_pf", nd.isPF,
			)
		}
	}
}
