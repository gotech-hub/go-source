package metric

import "errors"

var (
	DefaultErr = errors.New("ERROR")

	QueryMongoDBMetricHistogram = NewGlobalHistogramInstrument(
		"mongodb", "Time to query data from MongoDB",
	)

	HttpClientMetricHistogram = NewGlobalHistogramInstrument(
		"http_client", "Time to call http client",
	)
)
