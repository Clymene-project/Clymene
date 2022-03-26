package kafka

import (
	"bytes"
	"encoding/json"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/client"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
)

// Marshaller encodes a metric into a byte array to be sent to Kafka
type Marshaller interface {
	MarshalMetric([]prompb.TimeSeries) ([]byte, error)
	MarshalLog(*client.ProducerBatch) ([]byte, error)
}

type protobufMarshaller struct{}

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

func (h *protobufMarshaller) MarshalLog(batch *client.ProducerBatch) ([]byte, error) {
	panic("not supported")
}

func NewProtobufMarshaller() *protobufMarshaller {
	return &protobufMarshaller{}
}

type jsonMarshaller struct {
	pbMarshaller *jsonpb.Marshaler
}

func (h *jsonMarshaller) MarshalMetric(ts []prompb.TimeSeries) ([]byte, error) {
	out := new(bytes.Buffer)
	req := &prompb.WriteRequest{
		Timeseries: ts,
	}
	err := h.pbMarshaller.Marshal(out, req)
	return out.Bytes(), err
}

func (h *jsonMarshaller) MarshalLog(batch *client.ProducerBatch) ([]byte, error) {
	out, err := json.Marshal(batch)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewJSONMarshaller() *jsonMarshaller {
	return &jsonMarshaller{&jsonpb.Marshaler{}}
}
