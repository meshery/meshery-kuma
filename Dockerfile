# Stage 1: Build the Meshery Kuma adapter
FROM golang:1.19 as builder

# Install necessary dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

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

# Stage 2: Create the final production image
FROM alpine:3.14

# Install glibc
RUN apk --no-cache add ca-certificates wget && \
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.33-r0/glibc-2.33-r0.apk && \
    apk add glibc-2.33-r0.apk && \
    rm glibc-2.33-r0.apk && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.33-r0/glibc-bin-2.33-r0.apk && \
    apk add glibc-bin-2.33-r0.apk && \
    rm glibc-bin-2.33-r0.apk && \
    /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib && \
    echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf

# Set necessary environment variables
ENV DISTRO="alpine"
ENV SERVICE_ADDR="meshery-kuma"
ENV MESHERY_SERVER="http://meshery:9081"

# Copy the built Meshery Kuma adapter binary into the image from the builder stage
COPY --from=builder /meshery-kuma /meshery-kuma

# Set the entry point
ENTRYPOINT ["/meshery-kuma"]
