# Stage 1: Build Stage
FROM golang:1.23-alpine as builder

# Add Maintainer Info
LABEL maintainer="kietnguyen17052001@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR /app

# Install dependencies required for Go build
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Verify the files copied
RUN ls -l /app

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Stage 2: Run Stage (Production Image)
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the build stage
COPY --from=builder app .

# Expose the application port (update as per your app's requirement)
EXPOSE 7777

# Default command to run the application
CMD ./app