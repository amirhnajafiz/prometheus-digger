package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/query_range", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"resultType": "matrix",
				"result": []map[string]interface{}{
					{
						"metric": map[string]string{
							"__name__": "node_cpu_seconds_total",
							"instance": "localhost:9100",
							"job":      "node-exporter",
							"mode":     "idle",
						},
						"values": [][]interface{}{
							{1706520000, "12345.67"},
							{1706520060, "12346.02"},
							{1706520120, "12346.41"},
						},
					},
				},
			},
		}

		_ = json.NewEncoder(w).Encode(resp)
	})

	http.ListenAndServe(":9090", nil)
}
