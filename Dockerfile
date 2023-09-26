FROM golang:1.19 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./


RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM alpine:3.18.3

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080


CMD ["./main"]


