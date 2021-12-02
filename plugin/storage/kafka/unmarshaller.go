package kafka

// ProtobufUnmarshaller implements Unmarshaller
type ProtobufUnmarshaller struct{}

// NewProtobufUnmarshaller constructs a ProtobufUnmarshaller
func NewProtobufUnmarshaller() *ProtobufUnmarshaller {
	return &ProtobufUnmarshaller{}
}

// JSONUnmarshaller implements Unmarshaller
type JSONUnmarshaller struct{}

// NewJSONUnmarshaller constructs a JSONUnmarshaller
func NewJSONUnmarshaller() *JSONUnmarshaller {
	return &JSONUnmarshaller{}
}

// ZipkinThriftUnmarshaller implements Unmarshaller
type ZipkinThriftUnmarshaller struct{}

// NewZipkinThriftUnmarshaller constructs a zipkinThriftUnmarshaller
func NewZipkinThriftUnmarshaller() *ZipkinThriftUnmarshaller {
	return &ZipkinThriftUnmarshaller{}
}
