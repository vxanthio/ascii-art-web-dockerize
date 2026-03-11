# syntax=docker/dockerfile:1

# =============================================================================
# Stage 1 — builder
# Build the Go binary inside a full Go toolchain image.
# Nothing from this stage ends up in the final image except the binary we copy.
# =============================================================================
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy dependency manifest first so Docker can cache the module download layer.
# If go.mod/go.sum don't change, this layer is reused on the next build.
COPY go.mod ./

# Download modules (cached separately from source code).
RUN go mod download

# Copy the rest of the source tree.
COPY . .

# Compile the web server.
# CGO_ENABLED=0  — produce a fully static binary (no C runtime dependency).
# -o server      — the audit expects the binary to be named "server".
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/ascii-art-web

# =============================================================================
# Stage 2 — final
# Minimal runtime image: only the binary + static assets, nothing else.
# =============================================================================
FROM alpine:3.19

# ---------------------------------------------------------------------------
# Metadata — the audit requires metadata applied to Docker objects.
# ---------------------------------------------------------------------------
LABEL maintainer="teovaira"
LABEL org.opencontainers.image.title="ascii-art-web"
LABEL org.opencontainers.image.description="ASCII art generator web application written in Go"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.source="https://github.com/teovaira/ascii-art-web-dockerize"
LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /app

# Install bash so auditors can run: docker exec -it dockerize /bin/bash
RUN apk add --no-cache bash

# Copy the compiled binary from the builder stage.
COPY --from=builder /build/server ./server

# Copy static assets and HTML templates needed at runtime.
COPY --from=builder /build/static    ./static
COPY --from=builder /build/templates ./templates

# Expose the port the server listens on.
EXPOSE 8080

# Run the server binary.
# Using exec form (JSON array) — not shell form — so signals (SIGTERM) reach
# the process directly, enabling graceful shutdown.
CMD ["./server"]
