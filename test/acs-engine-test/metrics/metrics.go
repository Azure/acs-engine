package metrics

import (
	"encoding/json"
	"net"
)

type mdmBucket struct {
	Namespace string            `json:"Namespace"`
	Metric    string            `json:"Metric"`
	Dims      map[string]string `json:"Dims"`
}

func AddMetric(namespace, metric string, dims map[string]string) error {
	bucket := mdmBucket{
		Namespace: namespace,
		Metric:    metric,
		Dims:      dims}
	data, _ := json.Marshal(bucket)

	conn, err := net.Dial("udp", "localhost:8125")
	if err != nil {
		return err
	}
	defer conn.Close()

	//simple Read
	//buffer := make([]byte, 1024)
	//conn.Read(buffer)

	//simple write
	_, err = conn.Write(data)
	return err
}
