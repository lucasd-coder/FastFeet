
FIND ?= $(shell find proto/ -iname "*.proto")

GOPATH ?= go env GOPATH

.PHONY: proto-gen
proto-gen:	
	protoc --proto_path=proto/ $(FIND) \
		--plugin=$(GOPATH)/bin/protoc-gen-go-grpc \
		--go-grpc_out=. --go_out=.;


	
