FROM golang:1.17-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY cmd/ /app/cmd/
COPY internal/transport/grpc/proto /app/internal/transport/grpc/proto/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o application ./cmd/analytics/main.go

FROM alpine:3.15.4
COPY --from=builder /app/application /app/application
WORKDIR /app
CMD ["./application"]