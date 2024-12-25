FROM golang:latest AS builder
WORKDIR /app/hahapi
COPY . .
RUN go mod download
RUN go build -o server

FROM alpine:latest
WORKDIR /app/hahapi
COPY --from=builder /app/hahapi/ .
EXPOSE 8080
RUN apk add libc6-compat
ENTRYPOINT ["./server"]
