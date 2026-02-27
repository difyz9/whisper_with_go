# Whisper API Server - CPU Version
# Multi-stage build for optimized image size

# Stage 1: Build whisper.cpp
FROM ubuntu:22.04 AS whisper-builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Clone and build whisper.cpp
WORKDIR /build
RUN git clone https://github.com/ggerganov/whisper.cpp.git && \
    cd whisper.cpp && \
    cmake -B build -DCMAKE_BUILD_TYPE=Release && \
    cmake --build build --config Release

# Stage 2: Build Go application
FROM golang:1.23.4-bullseye AS go-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy whisper.cpp from previous stage
COPY --from=whisper-builder /build/whisper.cpp /usr/local/whisper.cpp

# Set CGO flags
ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-L/usr/local/whisper.cpp/build"
ENV CGO_CFLAGS="-I/usr/local/whisper.cpp"

# Build the application
RUN go build -o /app/bin/whisper-server cmd/server/main.go

# Stage 3: Runtime image
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create app user
RUN useradd -m -u 1000 whisper && \
    mkdir -p /app/uploads /app/outputs /app/models && \
    chown -R whisper:whisper /app

# Copy whisper.cpp libraries
COPY --from=whisper-builder /build/whisper.cpp/build/libwhisper.so /usr/local/lib/
COPY --from=whisper-builder /build/whisper.cpp/ggml/src/libggml.so /usr/local/lib/

# Copy Go binary
COPY --from=go-builder /app/bin/whisper-server /app/

# Copy configuration files
COPY --from=go-builder /app/.env.example /app/.env

# Set library path
ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

# Set working directory
WORKDIR /app

# Switch to non-root user
USER whisper

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/whisper-server"]
