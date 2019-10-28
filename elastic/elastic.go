package elastic

// ErrMonitor is a structure used for serializing/deserializing data in Elasticsearch.
type ErrMonitor struct {
	Bid  int64 `json:"bid"`
	Soid int64 `json:"soid"`
	Sid  int64 `json:"sid"`
}
