# Prometheus Digger

Prometheus Digger is a tool designed to help you extract and analyze metrics from a Prometheus server. This tool allows you to specify a time range and interval for the metrics you want to retrieve.

## Features

- Extract metrics from a Prometheus server
- Specify custom time ranges
- Define intervals for metric retrieval

## Usage

To build and run Prometheus Digger, follow these steps:

1. Build the project:

```sh
go build -o pdigger
chmod +x ./pdigger
```

2. Run the executable with the desired parameters:

```sh
./pdigger --metrics DCGM_FI_DEV_POWER_USAGE,KUBE_SCHEDULER_UP
```

## Parameters

- `--prometheus-url`: The URL of the Prometheus server. (default is http://localhost:9090)
- `--from`: The start time for the metrics in ISO 8601 format. (2025-03-07T00:00:00Z)
- `--to`: The end time for the metrics in ISO 8601 format. (2025-03-07T00:00:00Z)
- `--interval`: The interval at which to retrieve metrics. ("1m", "1h", "30s")
- `--metrics`: The specific metrics to retrieve. (separated using `,`)

## Example

Here is an example command to retrieve power usage metrics from a Prometheus server:

```sh
./main --prometheus-url http://127.0.0.1:9090 --from 2025-03-07T00:00:00Z --to 2025-03-08T00:00:00Z --interval 1m --metrics DCGM_FI_DEV_POWER_USAGE
```

This command will retrieve the `DCGM_FI_DEV_POWER_USAGE` metric from the Prometheus server at `http://127.0.0.1:9090` for the time range from March 7, 2025, to March 8, 2025, with a 1-minute interval.
