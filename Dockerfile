FROM alpine:latest

WORKDIR /app
COPY dist/juno-agent_linux_amd64_v1/juno-agent /app/bin/
COPY config /app/config

CMD ["/app/bin/juno-agent", "--config", "/app/config/config.toml"]
