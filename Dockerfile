# -------- Stage 1: Build the Go binary --------
FROM golang:1.24.4-alpine AS builder

WORKDIR /workspace

# 1) Copy go.mod & go.sum first for layer cache
COPY go.mod go.sum ./
RUN go mod download

# 2) Copy the rest of the source code
COPY . .

# 3) Build statically linked binary
#    -trimpath  : reproducible builds
#    -ldflags   : strip symbol table & DWARF info
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o rlaas ./cmd/api/main.go

# -------- Stage 2: Minimal runtime image --------
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /workspace/rlaas /rlaas

# The app listens on port 8080 by default
EXPOSE 8080

# Healthcheck (optional but handy for Compose / k8s)
HEALTHCHECK --interval=30s --timeout=3s CMD [ "/rlaas", "-health" ]

# Run as non-root user provided by distroless (uid=65532)
USER nonroot:nonroot

ENTRYPOINT ["/rlaas"]
