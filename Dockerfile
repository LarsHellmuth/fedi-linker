# syntax=docker/dockerfile:1.5.2

FROM golang:latest AS golang
WORKDIR /server
COPY --link ./server .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download && go build -o ./bin/run

FROM redis:latest AS redis
WORKDIR /cache
COPY --link ./cache .

FROM busybox:latest
WORKDIR /app
VOLUME /app/data
COPY --link --from=golang /server/bin/run ./server/bin/run
COPY --link --from=redis /usr/local/bin/redis-server ./cache/redis-server
COPY --link --from=redis /cache ./cache

# redis-server dependencies
COPY --link --from=redis \
/lib/x86_64-linux-gnu/libdl.so.2 \
/lib/x86_64-linux-gnu/libdl.so.2

COPY --link --from=redis \
/usr/lib/x86_64-linux-gnu/libssl.so.1.1 \
/usr/lib/x86_64-linux-gnu/libssl.so.1.1

COPY --link --from=redis \
/usr/lib/x86_64-linux-gnu/libcrypto.so.1.1 \
/usr/lib/x86_64-linux-gnu/libcrypto.so.1.1

ENTRYPOINT \
/app/cache/redis-server ./cache/redis.conf && \
/app/server/bin/run
