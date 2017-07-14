package metrics

import (
	"encoding/json"
	"fmt"

	"github.com/alexcesaro/statsd"
)

type mdmBucket struct {
	Namespace string            `json:"Namespace"`
	Metric    string            `json:"Metric"`
	Dims      map[string]string `json:"Dims"`
}

func AddMetric(namespace, metric string, count int64, dims map[string]string) error {
	bucket := mdmBucket{
		Namespace: namespace,
		Metric:    metric,
		Dims:      dims}
	data, err := json.Marshal(bucket)
	if err != nil {
		return err
	}
	client, err := statsd.New(
		statsd.Address(":8125"),
		statsd.Network("udp"),
		statsd.ErrorHandler(
			func(err error) {
				fmt.Println(err.Error())
			}))
	if err != nil {
		return err
	}
	defer client.Close()
	client.Count(string(data), count)
	return nil
}
