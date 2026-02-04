package main

import (
	"log/slog"
	logger "log/slog"
	"node-inventory/collectors"
	"os"
)

func init() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		),
	)
}

func main() {
	logger.Info("Starting collection")

	logger.Info("Collecting kernel information")
	if err := collectors.Kernel(); err != nil {
		logger.Error("Kernel collection failed", "error", err)
	}

	logger.Info("Collecting NIC information")
	if err := collectors.NIC(); err != nil {
		logger.Error("NIC collection failed", "error", err)
	}

	logger.Info("Finished collection")
}
