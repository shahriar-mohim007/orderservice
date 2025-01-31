# Builder Stage
FROM golang:1.23.5-alpine AS builder

WORKDIR /app


RUN apk add --no-cache git upx


COPY go.mod go.sum ./
RUN go mod download && go mod verify


COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o api cmd/api/main.go


RUN upx --best --lzma /app/api

# Final Stage
FROM scratch

WORKDIR /app
COPY --from=builder /app/api /usr/local/bin/api


CMD ["/usr/local/bin/api"]
