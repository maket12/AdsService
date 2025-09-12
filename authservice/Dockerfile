FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY . .

RUN go build -o authservice ./authservice

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/authservice .

EXPOSE 50051

CMD ["./authservice"]
