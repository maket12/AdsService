FROM golang:1.24 AS builder
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway ./cmd/server.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /gateway .

EXPOSE 8080

ENV AUTH_GRPC_ADDR=authservice:50051
ENV USER_GRPC_ADDR=userservice:50052

CMD ["./gateway"]
