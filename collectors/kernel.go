package collectors

import (
	"bytes"
	logger "log/slog"
	"os/exec"

	"golang.org/x/sys/unix"
)

type KernelInfo struct {
	Sysname     string
	Release     string
	Version     string
	KdumpStatus string
}

func byte65toString(ca [65]byte) string {
	return string(ca[:bytes.IndexByte(ca[:], 0)])
}

func getKernelInfo(kinfo *KernelInfo) {
	var uname unix.Utsname
	err := unix.Uname(&uname)

	if err != nil {
		logger.Error(
			"Error getting uname info",
			"error", err,
		)
	} else {
		logger.Debug("Uname Collection successful")
		kinfo.Sysname = byte65toString(uname.Sysname)
		kinfo.Release = byte65toString(uname.Release)
		kinfo.Version = byte65toString(uname.Version)
	}
}

func getKdumpInfo(kinfo *KernelInfo) {
	cmd := exec.Command("kdump-config", "status")
	out, err := cmd.CombinedOutput()

	if err != nil {
		logger.Warn(
			"Could not get kdump information",
			"error", err,
			"output", string(out),
		)
	} else {
		logger.Debug("Kdump-config status Collection successful")
		kinfo.KdumpStatus = string(out)
	}
}

func Kernel() {
	var kinfo KernelInfo
	getKernelInfo(&kinfo)
	getKdumpInfo(&kinfo)
	logger.Info(
		"Collection complete",
		"component", "kernel",
		"sysname", kinfo.Sysname,
		"release", kinfo.Release,
		"version", kinfo.Version,
		"kdump_status", kinfo.KdumpStatus,
	)
}
