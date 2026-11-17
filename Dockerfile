# === Этап 1: сборка ===
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum отдельно, чтобы кэш слоёв Docker
# переиспользовался, если изменился только код
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# === Этап 2: рантайм ===
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080 9090

ENTRYPOINT ["./server"]
