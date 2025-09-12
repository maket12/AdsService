FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o userservice ./userservice

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/userservice .

EXPOSE 50052

CMD ["./userservice"]
