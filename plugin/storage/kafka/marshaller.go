package kafka

import (
	"bytes"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/golang/snappy"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

// Marshaller encodes a span into a byte array to be sent to Kafka
type Marshaller interface {
	MarshalMetric([]prompb.TimeSeries) ([]byte, error)
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

func newProtobufMarshaller() *protobufMarshaller {
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

func newJSONMarshaller() *jsonMarshaller {
	return &jsonMarshaller{&jsonpb.Marshaler{}}
}
