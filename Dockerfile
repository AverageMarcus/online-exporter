FROM golang:1.25-alpine AS builder
WORKDIR /app/
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o online-exporter main.go

FROM golang:1.25-alpine
WORKDIR /app/
COPY --from=builder /app/online-exporter /app/online-exporter
ENTRYPOINT ["/app/online-exporter"]
