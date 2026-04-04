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
COPY --from=whisper-builder /build/whisper.cpp/include/ /usr/local/include/
COPY --from=whisper-builder /build/whisper.cpp/ggml/include/ /usr/local/include/
COPY --from=whisper-builder /build/whisper.cpp/build/src/libwhisper.so /usr/local/lib/libwhisper.so
COPY --from=whisper-builder /build/whisper.cpp/build/ggml/src/libggml.so /usr/local/lib/libggml.so
COPY --from=whisper-builder /build/whisper.cpp/build/ggml/src/libggml-base.so /usr/local/lib/libggml-base.so
COPY --from=whisper-builder /build/whisper.cpp/build/ggml/src/libggml-cpu.so /usr/local/lib/libggml-cpu.so

RUN ln -sf /usr/local/lib/libwhisper.so /usr/local/lib/libwhisper.so.1 && \
    ln -sf /usr/local/lib/libggml.so /usr/local/lib/libggml.so.0 && \
    ln -sf /usr/local/lib/libggml-base.so /usr/local/lib/libggml-base.so.0 && \
    ln -sf /usr/local/lib/libggml-cpu.so /usr/local/lib/libggml-cpu.so.0 && \
    printf '/usr/local/lib\n' > /etc/ld.so.conf.d/whisper.conf && \
    ldconfig

# Set CGO flags
ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-L/usr/local/lib"
ENV CGO_CFLAGS="-I/usr/local/include"
ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

# Build the application
RUN go build -o /app/bin/whisper-server cmd/server/main.go

# Stage 3: Runtime image
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    curl \
    ffmpeg \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create app user
RUN useradd -m -u 1000 whisper && \
    mkdir -p /app/uploads /app/outputs /app/models && \
    chown -R whisper:whisper /app

# Copy Go binary
COPY --from=go-builder /app/bin/whisper-server /app/

# Copy whisper runtime libraries
COPY --from=go-builder /usr/local/lib/libwhisper.so /usr/local/lib/libwhisper.so
COPY --from=go-builder /usr/local/lib/libwhisper.so.1 /usr/local/lib/libwhisper.so.1
COPY --from=go-builder /usr/local/lib/libggml.so /usr/local/lib/libggml.so
COPY --from=go-builder /usr/local/lib/libggml.so.0 /usr/local/lib/libggml.so.0
COPY --from=go-builder /usr/local/lib/libggml-base.so /usr/local/lib/libggml-base.so
COPY --from=go-builder /usr/local/lib/libggml-base.so.0 /usr/local/lib/libggml-base.so.0
COPY --from=go-builder /usr/local/lib/libggml-cpu.so /usr/local/lib/libggml-cpu.so
COPY --from=go-builder /usr/local/lib/libggml-cpu.so.0 /usr/local/lib/libggml-cpu.so.0
COPY --from=go-builder /etc/ld.so.conf.d/whisper.conf /etc/ld.so.conf.d/whisper.conf

RUN ldconfig

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
