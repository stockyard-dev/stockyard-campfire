FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/campfire ./cmd/campfire/
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl
COPY --from=builder /bin/campfire /usr/local/bin/campfire
ENV PORT="9030" DATA_DIR="/data"
EXPOSE 9030
HEALTHCHECK --interval=30s --timeout=5s CMD curl -sf http://localhost:9030/health || exit 1
ENTRYPOINT ["campfire"]
