# Start from the latest golang base image
FROM golang:1.25-alpine3.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /src/app

# Installing Git for alpine
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /promdigger

# Second stage
FROM alpine:3.20

# Set the working directory
WORKDIR /usr/local/bin/app

# Copy the binary from the builder stage
COPY --from=builder /promdigger .

# Set the entrypoint and default command
ENTRYPOINT ["./promdigger"]
