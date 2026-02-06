package falcon

const DataTagDelimiter = ", "

// Data 's 7 fields should be all specified.
type Data struct {
	Metric      string          `json:"metric"`      // Metric indicates the specific measurement of the collecting item.
	Endpoint    string          `json:"endpoint"`    // Endpoint indicates the subject (owner) of Metric.
	Timestamp   int64           `json:"timestamp"`   // Timestamp indicates the unix time when submitting the Data. Unit: seconds.
	Step        int             `json:"step"`        // Step indicates the reporting period of collecting items for the Data. Unit: seconds.
	Value       float64         `json:"value"`       // Value indicates the value of the metric at present time.
	CounterType DataCounterType `json:"counterType"` // CounterType can only choose one from DataCounterTypeCounter or DataCounterTypeGauge.
	Tags        string          `json:"tags"`        // Tags are a group of key values divided by commas which further describes and refines metric.
}

type DataCounterType string

const (
	DataCounterTypeCounter DataCounterType = "COUNTER"
	DataCounterTypeGauge   DataCounterType = "GAUGE"
)
