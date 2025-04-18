FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN \
   --mount=type=cache,target=/go/pkg \
   --mount=type=cache,target=/root/.cache/go-build \
   CGO_ENABLED=0 go build -o discord-prometheus .

FROM scratch
COPY --from=builder /app/discord-prometheus /
ENTRYPOINT ["/discord-prometheus", "live"]
LABEL org.opencontainers.image.source="https://github.com/AkovStohs/discord-prometheus"
