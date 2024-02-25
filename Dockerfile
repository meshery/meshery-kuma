# Use Alpine Linux as the base image
FROM golang:1.19 as builder

# Install necessary dependencies
RUN apk --no-cache add curl

# Set necessary environment variables
ENV GO111MODULE=on

# Set up working directory
WORKDIR /build

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Cache dependencies before building and copying source
RUN go mod download

# Update the version of github.com/layer5io/meshkit
RUN go get github.com/layer5io/meshkit@latest

# Copy the go source
COPY . .

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -o /meshery-kuma

# Use Alpine Linux as the final base image
FROM alpine:3.14

# Install necessary dependencies
RUN apk --no-cache add curl

# Set environment variables
ENV DISTRO="alpine"
ENV SERVICE_ADDR="meshery-kuma"
ENV MESHERY_SERVER="http://meshery:9081"

# Copy the built binary from the builder stage
COPY --from=builder /meshery-kuma /meshery-kuma

# Set the entry point
ENTRYPOINT ["/meshery-kuma"]
