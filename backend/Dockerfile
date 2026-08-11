FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
ENV GOTOOLCHAIN=auto
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
RUN mkdir -p /app/uploads
COPY --from=builder /server .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./server"]