package collectors

import (
	"bytes"
	logger "log/slog"

	"golang.org/x/sys/unix"
)

func byte65toString(ca [65]byte) string {
	return string(ca[:bytes.IndexByte(ca[:], 0)])
}

func Kernel() error {
	var uname unix.Utsname
	if err := unix.Uname(&uname); err == nil {
		logger.Info(
			"",
			"component", "kernel",
			"sysname", byte65toString(uname.Sysname),
			"release", byte65toString(uname.Release),
			"version", byte65toString(uname.Version))

	} else {
		logger.Error("Error getting uname info", "error", err)
		return err
	}

	return nil
}
