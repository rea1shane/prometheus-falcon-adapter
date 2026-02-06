package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"

	"github.com/rea1shane/prometheus-falcon-adapter/config"
	"github.com/rea1shane/prometheus-falcon-adapter/falcon"
	"github.com/rea1shane/prometheus-falcon-adapter/prometheus"
	"github.com/rea1shane/prometheus-falcon-adapter/transform"
)

func main() {
	var (
		configPath = kingpin.Flag(
			"config.path",
			"Path to configuration file.",
		).Default("adapter.yaml").String()
	)

	promslogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promslog.New(promslogConfig)

	configLogger := logger.With("file", *configPath)
	configLogger.Info("Reading config")
	cfg, err := config.New(*configPath)
	if err != nil {
		configLogger.Error("Failed to new config", "file", *configPath, "err", err)
		os.Exit(1)
	}
	configLogger.Info("Config contents", cfg.Content()...)

	exporters := make(map[string]string)
	for _, target := range cfg.Exporter.Targets {
		hostname, _ := config.ExtractHostname(target)
		exporters[target] = hostname
	}

	logger.Info("Scheduling jobs")
	run(logger, cfg.Global.ScrapeInterval, cfg.Global.ScrapeTimeout, cfg.Falcon.Target, exporters)
}

func run(logger *slog.Logger, scrapeInterval, scrapeTimeout time.Duration, falconURL string, exporters map[string]string) {
	ticker := time.NewTicker(scrapeInterval)
	defer ticker.Stop()

	job(logger, scrapeInterval, scrapeTimeout, falconURL, exporters)

	for {
		select {
		case <-ticker.C:
			job(logger, scrapeInterval, scrapeTimeout, falconURL, exporters)
		}
	}
}

func job(logger *slog.Logger, step, timeout time.Duration, falconURL string, exporters map[string]string) {
	start := time.Now()
	defer logger.Info("Job finished", "cost", time.Since(start).String())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	for url, hostname := range exporters {
		wg.Add(1)
		go func() {
			defer wg.Done()

			start := time.Now()
			taskLogger := logger.With("exporter", url)
			defer taskLogger.Debug("Task finished", "cost", time.Since(start).String())

			if err := task(ctx, taskLogger, step, falconURL, url, hostname); err != nil {
				taskLogger.Error("Task failed", "err", err)
			}
		}()
	}
	wg.Wait()
}

func task(ctx context.Context, logger *slog.Logger, step time.Duration, falconURL, exporterURL, exporterHostname string) error {
	metricFamilies, err := prometheus.Pull(ctx, exporterURL)
	if err != nil {
		return fmt.Errorf("failed to pull metrics: %w", err)
	}

	data := transform.PrometheusToFalcon(exporterHostname, time.Now().Unix(), int(step.Seconds()), metricFamilies)

	if err := falcon.Push(ctx, falconURL, data); err != nil {
		return fmt.Errorf("failed to push data: %w", err)
	}

	return nil
}
