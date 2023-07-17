ARG GOLANG_VERSION=1.20
ARG ALPINE_VERSION=3.18

# Build section
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/bot main.go

# Release
FROM alpine:${ALPINE_VERSION}
WORKDIR /app
COPY --from=builder /app/bot /app/bot
CMD ["/app/bot"]