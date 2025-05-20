FROM golang:1.21.6-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o main ./cmd/server

FROM alpine:3.19

WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

RUN adduser -D appuser
USER appuser

EXPOSE 3000

CMD ["./main"]
