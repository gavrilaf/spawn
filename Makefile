SHELL := /bin/bash

generate:
	protoc -I=./pkg/backend/pb -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--gogofast_out=./pkg/backend/pb ./pkg/backend/pb/*.proto

	protoc -I=./pkg/backend/pb -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--mqrpc_out=./pkg/backend/pb/ ./pkg/backend/pb/*.proto

rebuild:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml build

run:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up

run-env:
	docker-compose up -d

run-back:
	go run ./cmd/backend/main.go

run-rest:
	go run ./cmd/rest/main.go

stop:
	docker-compose down --remove-orphans

ptest: 
	( \
		source ./py-test/venv/bin/activate; \
    pytest ./py-test/spawn_api; \
	)

	
	
	