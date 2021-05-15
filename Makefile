proto-metric:
	go install github.com/gogo/protobuf/protoc-gen-gogo
	go install github.com/gogo/protobuf/protoc-gen-gofast
	go install github.com/gogo/protobuf/protoc-gen-gogofast
	go install github.com/gogo/protobuf/protoc-gen-gogofaster
	go get golang.org/x/tools/cmd/goimports
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

	(protoc -I=./prompb -I=./prompb/googleapis -I=. --gofast_out=plugins=grpc:output/prompb ./prompb/types.proto)
	(protoc -I=./prompb -I=./prompb/googleapis -I=. --gofast_out=plugins=grpc:output/prompb ./prompb/remote.proto)
	(protoc -I=./prompb -I=./prompb/googleapis -I=. --gofast_out=plugins=grpc:output/prompb --grpc-gateway_out=logtostderr=true:output/prompb ./prompb/rpc.proto)