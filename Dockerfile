FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# go build prints nothing by default — looks “stuck”. -v streams package names.
# GOMAXPROCS=1 lowers RAM spikes on small VPS (1–2 GB) where the compiler seems to hang (swap/OOM).
RUN CGO_ENABLED=0 GOOS=linux GOMAXPROCS=1 \
	go build -trimpath -ldflags="-s -w" -v -o /app/server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["./server"]
