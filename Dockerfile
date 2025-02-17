# Step 1: Use the official Go image to build the application
FROM golang:1.21

# Step 2: Set the working directory for your Go app inside the container
WORKDIR /app

# Step 3: Copy Go module files (go.mod and go.sum) to the container
COPY go.mod go.sum ./

# Step 4: Download Go dependencies
RUN go mod tidy

# Step 5: Copy your Go application code to the container
COPY . .

# Step 6: Build the Go application binary
RUN go build -o go-app .

# Step 11: Expose the port the Go app will run on
EXPOSE 8080

# Step 12: Command to run the application
CMD ["./go-app"]