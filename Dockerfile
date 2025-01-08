# Use the Go image with Alpine
FROM golang:1.21-alpine

# Install build dependencies for SQLite and CGO
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o main .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
