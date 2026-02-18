# Promdigger

Promdigger (Prometheus Digger) is a Go-based CLI tool for executing PromQL queries and exporting the results to CSV files. It interacts directly with the Prometheus HTTP API to run and manage queries efficiently.

## Features

* Execute large-scale PromQL queries efficiently
* Automatically optimize query execution based on query range and size
* Concurrent query processing for improved performance
* Batched query execution
* Export query results to CSV format

## Installation

After cloning the repository, run the installation script:

```sh
./scripts/install.sh
```

> ⚠️ The installation script currently supports Unix-based systems only.

Verify the installation:

```sh
promdigger --help
```
