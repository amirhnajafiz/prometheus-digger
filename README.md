# Prometheus Digger

Prometheus Digger is a tool designed to help you extract and analyze metrics from a Prometheus server. This tool allows you to specify a time range and interval for the metrics you want to retrieve.

## Features

- Extract metrics/queries from a Prometheus server
- Specify custom time ranges
- Define intervals for metric retrieval
- Adaptive pulling method based on the given query

## Usage

To build and run Prometheus Digger, follow these steps:

1. Build the project:

```sh
go build -o pdigger
chmod +x ./pdigger
```

2. Run the executable:

```sh
./pdigger -help
```

## Parameters

Copy the example config file from `config/config.example.json` into `config.json`. You the following fields:

- `url`: The URL of the Prometheus server. (default is http://localhost:9090)
- `from`: The start time for metrics. (e.g -2:12:00:00 is 2 days and 12 hours ago from now)
- `to`: The end time for metrics. (e.g -1:12:00:00 is 1 day and 12 hours ago from now)
- `queries`: A list of metrics and queries to dig
    - `queries.interval`: The interval at which to retrieve metrics. ("1m", "1h", "30s")
    - `queries.name`: The output directory of the query, also a label for logs and debugging.
    - `queries.metric`: Input PromQL.

## Example

Here is an example command to retrieve power usage metrics from a Prometheus server:

```sh
{
    "url": "http://127.0.0.1:9090",
    "from": "-1:12:00:00",
    "to": "+2:12:00:00",
    "queries": [
        {
            "name": "GPU_POWER_USAGE",
            "metric": "DCGM_FI_DEV_GPU_TEMP",
            "interval": "1m"
        }
    ]
}
```

This command will retrieve the `DCGM_FI_DEV_POWER_USAGE` metric from the Prometheus server at `http://127.0.0.1:9090` for the time range from 1 and half day ago to 2 and half days in future, with a 1-minute interval.
