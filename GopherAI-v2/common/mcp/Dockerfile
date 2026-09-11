FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}
RUN go mod download

COPY . ./
RUN go build -trimpath -o /out/gopherai-mcp .

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/gopherai-mcp /usr/local/bin/gopherai-mcp

EXPOSE 8081

ENTRYPOINT ["/usr/local/bin/gopherai-mcp", "-mode", "server", "-http-addr", ":8081"]
