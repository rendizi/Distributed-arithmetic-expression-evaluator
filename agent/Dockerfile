FROM golang:latest AS builder

WORKDIR /app

# Copy necessary files into the Docker build context
COPY ../go.mod ../go.sum ./

RUN go mod download

# Copy the entire directory into the image
COPY ../ ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .
COPY --from=builder /app/.env .
COPY --from=builder /app/main.go .

CMD ["./app"]
