# Use the official Golang image as the base image
FROM golang:1.19 as builder

# Set environment variables
ARG VERSION
ARG GIT_COMMITSHA

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application code
COPY . ./

# Build the Go binary with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -o app -tags netgo -installsuffix netgo .

# Start a new stage
FROM alpine:3.16

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage to the final image
COPY --from=builder /app/app .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
