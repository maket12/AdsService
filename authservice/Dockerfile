FROM golang:1.25 AS builder
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/authservice
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/authservice ./cmd/app/main.go

FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/authservice /app/authservice


EXPOSE 50051
CMD ["/app/authservice"]
