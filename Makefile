
FIND ?= $(shell find proto/ -iname "*.proto")

GOPATH ?= go env GOPATH

.PHONY: proto-gen generate-mocks wire-gen run-application

proto-gen:	
	protoc --proto_path=proto/ $(FIND) \
		--plugin=$(GOPATH)/bin/protoc-gen-go-grpc \
		--go-grpc_out=. --go_out=.;

generate-mocks:
	./scripts/mockery.sh


wire-gen:
	cd internal/app && wire	

run-application:
	go run ./cmd/app/main.go

