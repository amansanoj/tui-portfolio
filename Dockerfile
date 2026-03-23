FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o portfolio .

FROM alpine:latest
RUN apk --no-cache add ca-certificates openssh-keygen
WORKDIR /app
COPY --from=builder /app/portfolio .

RUN mkdir -p /data

EXPOSE 22
CMD ["./portfolio"]