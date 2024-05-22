#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app .

# Run tests
RUN CGO_ENABLED=0 go test -v ./... > test_output.txt && cat test_output.txt

# Build the application.
# RUN CGO_ENABLED=0 go build -o /go/bin/app .

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
# Create the app directory
RUN mkdir -p /app

# Copy the built binary and .env file from the builder stage
COPY --from=builder /go/bin/app /app/app
COPY .env /app/.env

# Set the entrypoint to the binary location
ENTRYPOINT ["/app/app"]
LABEL Name=gowebserver Version=0.0.2
EXPOSE 8080
