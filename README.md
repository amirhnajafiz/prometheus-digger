# Prometheus Digger

You can use this tool to dig your metrics inside a prometheus server.

```sh
go build -o main
./main --prometheus-url http://127.0.0.1:9090 --from 2025-03-07T00:00:00Z --to 2025-03-08T00:00:00Z --interval 1m --metrics DCGM_FI_DEV_POWER_USAGE
```
