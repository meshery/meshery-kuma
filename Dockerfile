# Use Alpine Linux with glibc 2.33 as the base image
FROM frolvlad/alpine-glibc:alpine-3.14_glibc-2.33 as builder

ARG VERSION
ARG GIT_COMMITSHA

WORKDIR /build

# Install build-time dependencies
RUN apk --no-cache add curl

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Cache dependencies before building and copying source
RUN GOPROXY=https://proxy.golang.org,direct go mod download

# Copy the go source
COPY main.go main.go
COPY internal/ internal/
COPY kuma/ kuma/
COPY build/ build/

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux GO111MODULE=on go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -a -o meshery-kuma main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/nodejs:latest

WORKDIR /

ENV DISTRO="debian"
ENV SERVICE_ADDR="meshery-kuma"
ENV MESHERY_SERVER="http://meshery:9081"

COPY templates/ ./templates
COPY --from=builder /build/meshery-kuma .

ENTRYPOINT ["/meshery-kuma"]
