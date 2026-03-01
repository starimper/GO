# ---------- BUILD STAGE ----------
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp ./cmd/api/main.go

# ---------- FINAL STAGE ----------
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp .

EXPOSE 8080

ENTRYPOINT ["./myapp"]