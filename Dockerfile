FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY metrics/ metrics/

RUN CGO_ENABLED=0 GOOS=linux go build -o /discord-prometheus

FROM alpine:latest
COPY --from=builder /app/discord-prometheus .

EXPOSE 9090
CMD ["/discord-prometheus", "live"]
LABEL org.opencontainers.image.source="https://github.com/AkovStohs/discord-prometheus"
