# Stage 1: Build
FROM golang:1-bookworm AS builder

ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o filectl .

# Stage 2: Runtime (Debian-based for .deb compatibility testing)
FROM debian:bookworm-slim AS runtime

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/filectl /usr/local/bin/filectl
COPY testdata/ /app/testdata/

RUN chmod +x /usr/local/bin/filectl

ENTRYPOINT ["filectl"]
