FROM golang:1.21-alpine AS build

RUN apk add --no-cache git wget
RUN apk --no-cache add ca-certificates

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.19 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

RUN addgroup -S appuser \
    && adduser -S -u 10000 -g appuser appuser
 
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
 
COPY . .
 
RUN CGO_ENABLED=0 GOOS=linux go build -o /app -ldflags="-s -w -a" ./cmd/app
 
FROM alpine:3.18.4

COPY --from=build /app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/passwd
COPY --from=build /src/config /config
COPY --from=build /bin/grpc_health_probe /bin/grpc_health_probe

USER appuser
 
ENTRYPOINT ["/app"]