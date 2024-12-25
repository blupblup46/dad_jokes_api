FROM golang:latest AS builder
WORKDIR /app/hahapi
COPY . .
RUN go mod download
RUN go build -o client 

FROM alpine:latest
WORKDIR /app/hahapi
COPY --from=builder /app/hahapi/ .
RUN apk add libc6-compat
ENTRYPOINT [ "./client" ]