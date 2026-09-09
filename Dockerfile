# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o k8s-mcp ./cmd/main.go

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/k8s-mcp .

EXPOSE 8080

ENTRYPOINT ["/app/k8s-mcp"]
