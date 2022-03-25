package prometheus

import (
	"github.com/Clymene-project/Clymene/cmd/promtail/app/client"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/Clymene-project/Clymene/storage/logstore"
	"github.com/golang/snappy"

	"github.com/gogo/protobuf/proto"
)

// Marshaller encodes a metric into a byte array to be sent to Kafka
type Marshaller interface {
	MarshalMetric([]prompb.TimeSeries) ([]byte, error)
}

type protobufMarshaller struct{}

func (h *protobufMarshaller) MarshalBatch(logs logstore.Batch) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (h *protobufMarshaller) MarshalLog(logs *client.ProducerBatch) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (h *protobufMarshaller) MarshalMetric(ts []prompb.TimeSeries) ([]byte, error) {
	req := &prompb.WriteRequest{
		Timeseries: ts,
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	compressed := snappy.Encode(nil, data)
	return compressed, nil
}

func newProtobufMarshaller() *protobufMarshaller {
	return &protobufMarshaller{}
}
