# Prometheus Falcon adapter

A lightweight adapter that scrapes Prometheus-format metrics from one or more exporters and pushes the converted data to Falcon via the push API.

## Architecture

- Periodically scrape metrics from `exporter.targets` in configuration
- Parse Prometheus exposition format
- Convert metrics to Falcon data model
- Push to `falcon.target` in configuration

## Configuration

The configuration file is written in YAML. Example:

```yaml
global:
  scrape_interval: 15s # Interval between scrapes. Must be > 0.
  scrape_timeout: 10s  # Timeout for each scrape request. Must be > 0 and less than `scrape_interval`.
falcon:
  target: http://localhost:8888/push # Falcon push endpoint. Must be a valid `http://` or `https://` URL with a non-empty hostname.
exporter:
  targets: # List of Prometheus metrics endpoints to scrape. Each entry must be a valid `http://` or `https://` URL with a non-empty hostname.
    - http://foo:9999/metrics
    - http://bar
    - https://foo.bar:7777
```

## Build

```bash
go build -o prometheus-falcon-adapter .
```

## Run

```bash
./prometheus-falcon-adapter --config.path=adapter.yaml
```

If not specified, `--config.path` defaults to `adapter.yaml` in the current working directory.
