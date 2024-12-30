# Step 1: Use the official Go image to build the application
FROM golang:1.21 AS build

# Step 2: Set the working directory for your Go app inside the container
WORKDIR /app

# Step 3: Copy Go module files (go.mod and go.sum) to the container
COPY go.mod go.sum ./

# Step 4: Download Go dependencies
RUN go mod download

# Step 5: Copy your Go application code to the container
COPY . .

# Step 6: Build the Go application binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-app .

# Step 7: Create a new, minimal image to run the application (using scratch for minimal size)
FROM alpine:latest

# Step 8: Install only necessary dependencies for running the Go binary
RUN apk add --no-cache sqlite

# Step 9: Set the working directory for your Go app inside the container
WORKDIR /app

# Step 10: Copy the built Go binary and SQLite database from the build stage
COPY --from=build /app/go-app /app/
COPY ./database.db /app/database.db

# Step 11: Expose the port your Go app will use
EXPOSE 8080

# Step 12: Command to run the application
CMD ["./go-app"]
