package kafka

import (
	"bytes"
	"github.com/Clymene-project/Clymene/prompb"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
)

type Unmarshaller interface {
	Unmarshal([]byte) ([]prompb.TimeSeries, error)
}

// ProtobufUnmarshaller implements Unmarshaller
type ProtobufUnmarshaller struct{}

func (p *ProtobufUnmarshaller) Unmarshal(msg []byte) ([]prompb.TimeSeries, error) {
	req := &prompb.WriteRequest{}
	decodeMsg, err := snappy.Decode(nil, msg)
	err = proto.Unmarshal(decodeMsg, req)
	return req.Timeseries, err
}

// NewProtobufUnmarshaller constructs a ProtobufUnmarshaller
func NewProtobufUnmarshaller() *ProtobufUnmarshaller {
	return &ProtobufUnmarshaller{}
}

// JSONUnmarshaller implements Unmarshaller
type JSONUnmarshaller struct{}

func (J *JSONUnmarshaller) Unmarshal(msg []byte) ([]prompb.TimeSeries, error) {
	req := &prompb.WriteRequest{}
	err := jsonpb.Unmarshal(bytes.NewReader(msg), req)
	return req.Timeseries, err
}

// NewJSONUnmarshaller constructs a JSONUnmarshaller
func NewJSONUnmarshaller() *JSONUnmarshaller {
	return &JSONUnmarshaller{}
}
