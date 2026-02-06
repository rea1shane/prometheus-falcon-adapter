package transform

var input = []byte(`
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{instance="server-1", method="GET", status="200"} 1234
http_requests_total{instance="server-1", method="POST", status="200"} 567

# HELP cpu_usage_percent CPU usage percentage
# TYPE cpu_usage_percent gauge
cpu_usage_percent{instance="server-2", core="0"} 45.5
cpu_usage_percent{instance="server-2", core="1"} 62.3

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{instance="server-3", le="0.1"} 1000
http_request_duration_seconds_bucket{instance="server-3", le="0.5"} 1450
http_request_duration_seconds_bucket{instance="server-3", le="1.0"} 1480
http_request_duration_seconds_bucket{instance="server-3", le="+Inf"} 1500
http_request_duration_seconds_sum{instance="server-3"} 523.7
http_request_duration_seconds_count{instance="server-3"} 1500

# HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 3.0917e-05
go_gc_duration_seconds{quantile="0.25"} 5.0376e-05
go_gc_duration_seconds{quantile="0.5"} 7.8417e-05
go_gc_duration_seconds{quantile="0.75"} 8.9458e-05
go_gc_duration_seconds{quantile="1"} 0.000131374
go_gc_duration_seconds_sum 0.663532651
go_gc_duration_seconds_count 8832
`)
