#build stage
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git

ENV PATH="/go/bin:${PATH}"

RUN apk update && apk add --no-cache protobuf make protobuf-dev & \
    go install github.com/google/wire/cmd/wire@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    wget https://github.com/ktr0731/evans/releases/download/v0.10.9/evans_linux_amd64.tar.gz && \
    tar -xzvf evans_linux_amd64.tar.gz && \
    mv evans ../bin && rm -f evans_linux_amd64.tar.gz

WORKDIR /go/src

EXPOSE ${GRPC_PORT}

CMD [ "tail", "-f", "/dev/null" ]