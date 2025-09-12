FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY . .

RUN go build -o gateway ./infra/gateway/cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway .

EXPOSE 8080

CMD ["./gateway"]
