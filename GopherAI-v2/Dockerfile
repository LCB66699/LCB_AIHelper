FROM golang:1.24-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}
RUN go mod download

COPY . ./
RUN go build -trimpath -o /out/gopherai .

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /out/gopherai /app/gopherai
COPY config/config.toml /app/config/config.toml

RUN mkdir -p /app/uploads

ENV LD_LIBRARY_PATH=/opt/onnxruntime

EXPOSE 9090

CMD ["/app/gopherai"]
