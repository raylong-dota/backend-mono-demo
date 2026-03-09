FROM golang:1.26.0 AS builder

ARG APP_SVC

RUN apt-get update && apt-get install -y --no-install-recommends make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY . .

RUN make build ${APP_SVC}

# ── Runtime ───────────────────────────────────────────────────────────────────
FROM debian:stable-slim

ARG APP_SVC

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        netbase \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin/orbit-${APP_SVC}-svc /app/server

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["/app/server", "-conf", "/data/conf"]
