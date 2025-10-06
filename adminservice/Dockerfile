FROM golang:1.24 AS builder
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/adminservice
RUN CGO_ENABLED=0 GOOS=linux go build -o /adminservice ./cmd/server.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /adminservice .

EXPOSE 50053
CMD ["./adminservice"]
