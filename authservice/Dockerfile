FROM golang:1.24 AS builder
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/authservice
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/authservice ./cmd/server.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/authservice /app/authservice

EXPOSE 50051
CMD ["/app/authservice"]
