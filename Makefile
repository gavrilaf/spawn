generate:

	protoc -I=./pkg/backend/pb -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--gogofast_out=./pkg/backend/pb ./pkg/backend/pb/*.proto

	protoc -I=./pkg/backend/pb -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--mqrpc_out=./pkg/backend/pb/ ./pkg/backend/pb/*.proto