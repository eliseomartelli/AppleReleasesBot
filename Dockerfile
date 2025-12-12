FROM golang:alpine3.21 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o applereleases ./cmd/applereleasesbot

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/applereleases .

RUN mkdir -p /data && \
	addgroup -S appgroup && \
	adduser -S appuser -G appgroup && \
	chown -R appuser:appgroup /app && \
	chown appuser:appgroup /data

VOLUME /data
WORKDIR /data

USER appuser

CMD ["/app/applereleases"]
