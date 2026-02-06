package transform

import (
	"fmt"
	"log/slog"
	"strings"

	dto "github.com/prometheus/client_model/go"

	"github.com/rea1shane/prometheus-falcon-adapter/falcon"
)

func PrometheusToFalcon(endpoint string, timestamp int64, step int, metricFamilies map[string]*dto.MetricFamily) []falcon.Data {
	var out []falcon.Data

	for name, metricFamily := range metricFamilies {
		for _, metric := range metricFamily.GetMetric() {
			var tags []string
			for _, labelPair := range metric.GetLabel() {
				tags = append(tags, fmt.Sprintf("%s=%s", labelPair.GetName(), labelPair.GetValue()))
			}
			tagString := strings.Join(tags, falcon.DataTagDelimiter)

			switch metricFamily.GetType() {
			case dto.MetricType_COUNTER:
				out = append(out, falcon.Data{
					Metric:      name,
					Endpoint:    endpoint,
					Timestamp:   timestamp,
					Step:        step,
					Value:       metric.GetCounter().GetValue(),
					CounterType: falcon.DataCounterTypeCounter,
					Tags:        tagString,
				})

			case dto.MetricType_GAUGE:
				out = append(out, falcon.Data{
					Metric:      name,
					Endpoint:    endpoint,
					Timestamp:   timestamp,
					Step:        step,
					Value:       metric.GetGauge().GetValue(),
					CounterType: falcon.DataCounterTypeGauge,
					Tags:        tagString,
				})

			case dto.MetricType_UNTYPED:
				slog.Debug("The prometheus metric is of type untyped", "metric", name)
				out = append(out, falcon.Data{
					Metric:      name,
					Endpoint:    endpoint,
					Timestamp:   timestamp,
					Step:        step,
					Value:       metric.GetUntyped().GetValue(),
					CounterType: falcon.DataCounterTypeGauge,
					Tags:        tagString,
				})

			case dto.MetricType_SUMMARY:
				summary := metric.GetSummary()

				out = append(out, falcon.Data{
					Metric:      name + "_sum",
					Endpoint:    endpoint,
					Timestamp:   timestamp,
					Step:        step,
					Value:       summary.GetSampleSum(),
					CounterType: falcon.DataCounterTypeCounter,
					Tags:        tagString,
				})

				out = append(out, falcon.Data{
					Metric:      name + "_count",
					Endpoint:    endpoint,
					Timestamp:   timestamp,
					Step:        step,
					Value:       float64(summary.GetSampleCount()),
					CounterType: falcon.DataCounterTypeCounter,
					Tags:        tagString,
				})

				for _, quantile := range summary.GetQuantile() {
					out = append(out, falcon.Data{
						Metric:      name,
						Endpoint:    endpoint,
						Timestamp:   timestamp,
						Step:        step,
						Value:       quantile.GetValue(),
						CounterType: falcon.DataCounterTypeGauge,
						Tags: strings.Join([]string{
							tagString,
							fmt.Sprintf(`quantile="%f"`, quantile.GetQuantile()),
						}, falcon.DataTagDelimiter),
					})
				}

			case dto.MetricType_HISTOGRAM:
				histogram := metric.GetHistogram()

				out = append(out, falcon.Data{
					Metric:      name + "_sum",
					Endpoint:    endpoint,
					Timestamp:   timestamp,
					Step:        step,
					Value:       histogram.GetSampleSum(),
					CounterType: falcon.DataCounterTypeCounter,
					Tags:        tagString,
				})

				out = append(out, falcon.Data{
					Metric:      name + "_count",
					Endpoint:    endpoint,
					Timestamp:   timestamp,
					Step:        step,
					Value:       float64(histogram.GetSampleCount()),
					CounterType: falcon.DataCounterTypeCounter,
					Tags:        tagString,
				})

				for _, bucket := range histogram.GetBucket() {
					out = append(out, falcon.Data{
						Metric:      name + "_bucket",
						Endpoint:    endpoint,
						Timestamp:   timestamp,
						Step:        step,
						Value:       float64(bucket.GetCumulativeCount()),
						CounterType: falcon.DataCounterTypeCounter,
						Tags: strings.Join([]string{
							tagString,
							fmt.Sprintf(`le="%f"`, bucket.GetUpperBound()),
						}, falcon.DataTagDelimiter),
					})
				}
			}
		}
	}

	return out
}
