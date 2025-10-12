FROM golang:1.24 AS builder
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/userservice
RUN CGO_ENABLED=0 GOOS=linux go build -o /userservice ./cmd/server.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /userservice .

EXPOSE 50052
CMD ["./userservice"]
