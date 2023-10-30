# Build stage
FROM golang:latest as builder
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/proxy-go .

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /out/proxy-go /proxy-go
EXPOSE 8080
ENTRYPOINT ["/proxy-go"]