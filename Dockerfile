FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o  /app/main ./cmd/room-service

FROM alpine:latest
COPY --from=builder /app/main /main
USER nobody:nobody
EXPOSE 80
CMD ["/main"]