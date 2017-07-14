package metrics

import (
	"encoding/json"
	"fmt"

	"github.com/alexcesaro/statsd"
)

// statsdMetricClientWrapper is an interface wrapping
// interactions to write metrics to statsd
type statsdMetricClientWrapper interface {
	Sector() string
	Region() string
	Count(string, int64)
	Close()
}

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
	data, err := json.Marshal(bucket)
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	//statsdClient.Count(string(durationBucketBytes), latency.Nanoseconds()/nanoSecondToMillisecondConversionFactor)
	client, err := statsd.New(statsd.Address("127.0.0.1:8125"), statsd.Network("udp"))
	if err != nil {
		return err
	}
	client.Count(string(data), 1000)
	client.Close()
	return nil

	/*conn, err := net.Dial("udp", "127.0.0.1:8125")
	if err != nil {
		return err
	}
	defer conn.Close()

	//simple Read
	//buffer := make([]byte, 1024)
	//conn.Read(buffer)

	//simple write
	_, err = conn.Write(data)
	fmt.Printf("AddMetric [%s] [%v]\n", string(data), err)*/
	/*sAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8125")
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, sAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	fmt.Println("client: wrote:", string(data[0:n]))*/
	return nil
}
