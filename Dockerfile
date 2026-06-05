FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY . .

# go build prints nothing by default — looks “stuck”. -v streams package names.
# GOMAXPROCS=1 lowers RAM spikes on small VPS (1–2 GB) where the compiler seems to hang (swap/OOM).
# Cache mounts need BuildKit (default in Docker 23+ / compose v2). Build: docker compose build --progress=plain backend
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux GOMAXPROCS=1 \
	go build -trimpath -ldflags="-s -w" -v -o /app/server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["./server"]
