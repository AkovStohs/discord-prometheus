FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN apk --no-cache add --no-check-certificate ca-certificates \
   && update-ca-certificates
RUN \
   --mount=type=cache,target=/go/pkg \
   --mount=type=cache,target=/root/.cache/go-build \
   CGO_ENABLED=0 go build -o discord-prometheus .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /ca-certificates.crt
COPY --from=builder /etc/ssl/certs/ca-bundle.crt /etc/ssl/certs/ca-bundle.crt
COPY --from=builder /app/discord-prometheus /
ENV MECTRICS_PORT=9090
EXPOSE ${MECTRICS_PORT}
ENTRYPOINT ["/discord-prometheus", "live"]
LABEL org.opencontainers.image.source="https://github.com/AkovStohs/discord-prometheus"
