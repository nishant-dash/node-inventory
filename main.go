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
	collectors.Kernel()

	logger.Info("Collecting NIC information")
	collectors.NIC()

	logger.Info("Finished collection")
}
