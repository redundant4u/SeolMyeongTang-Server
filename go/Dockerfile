FROM golang:1.25-alpine AS builder

ARG ARCHITECTURE

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCHITECTURE} go build -trimpath -o main main.go

FROM alpine:3.23
WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
