# Build stage

# https://hub.docker.com/_/golang
FROM golang:1.19-alpine3.17 AS builder

LABEL maintainer="charles schiavinato charles.schiavinato@yahoo.gom.br"

# Create appuser
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  appuser

WORKDIR /app

COPY . .

# Fetch dependencies
RUN go mod download
RUN go mod verify

# Build the binary
RUN go build -o luizalabs-order server.go

# Run stage
FROM alpine:3.17

WORKDIR /app

# Import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy binary from the builder
COPY --from=builder /app/luizalabs-order .
COPY --from=builder /app/swagger.yaml .
COPY --from=builder /app/service/database/migration/scripts ./migration

EXPOSE 9000

# Use an unprivileged user
USER appuser:appuser

# Run the binary
CMD [ "./luizalabs-order" ]