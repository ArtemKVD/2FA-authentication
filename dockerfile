FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /2fa-app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /2fa-app /app/2fa-app
COPY web/templates ./web/templates
COPY migrations ./migrations

RUN apk add --no-cache tzdata

EXPOSE 8080

CMD ["/app/2fa-app"]