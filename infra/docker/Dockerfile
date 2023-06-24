FROM golang:1.20-alpine AS build

RUN apk add --no-cache git wget
RUN apk --no-cache add ca-certificates

RUN addgroup -S appuser \
    && adduser -S -u 10000 -g appuser appuser
 
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
 
COPY . .
 
RUN CGO_ENABLED=0 GOOS=linux go build -o /app -ldflags="-s -w -a" ./cmd/app
 
FROM scratch

COPY --from=build /app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/passwd
COPY --from=build /src/config /config

USER appuser
 
ENTRYPOINT ["/app"]