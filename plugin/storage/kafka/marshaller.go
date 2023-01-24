package kafka

import (
	"bytes"
	"encoding/json"
	"github.com/Clymene-project/Clymene/cmd/promtail/app/client"
	"github.com/Clymene-project/Clymene/plugin/storage/es/metricstore/dbmodel"
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

func (h *protobufMarshaller) MarshalLog(_ *client.ProducerBatch) ([]byte, error) {
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

type jsonFlattenMarshaller struct {
	converter dbmodel.Converter
}

func (j *jsonFlattenMarshaller) MarshalMetric(metrics []prompb.TimeSeries) ([]byte, error) {
	var metricsMap []byte
	for _, metric := range metrics {
		jsonMetric := j.converter.ConvertTsToJSON(metric)
		marshaledMetric, err := json.Marshal(jsonMetric)
		if err != nil {
			return nil, err
		}
		metricsMap = append(metricsMap, marshaledMetric...)
	}
	return metricsMap, nil
}

func (j *jsonFlattenMarshaller) MarshalLog(_ *client.ProducerBatch) ([]byte, error) {
	panic("not supported")
}

func NewJsonFlattenMarshaller() *jsonFlattenMarshaller {
	return &jsonFlattenMarshaller{
		converter: dbmodel.Converter{},
	}
}
