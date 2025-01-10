# Step 1: Use the official Go image to build the application
FROM golang:1.21 AS build

# Step 2: Set the working directory for your Go app inside the container
WORKDIR /app

# Step 3: Copy Go module files (go.mod and go.sum) to the container
COPY go.mod go.sum ./

# Step 4: Download Go dependencies
RUN go mod tidy

# Step 5: Copy your Go application code to the container
COPY . .

# Step 6: Build the Go application binary
RUN go build -o /go-app .

# Step 7: Create a new image to run the application
FROM ubuntu:latest

# Step 8: Install required dependencies for running Go binary (if necessary)
RUN apt-get update && \
    apt-get install -y sqlite3

# Step 9: Set the working directory
WORKDIR /go-app

# Step 10: Copy the Go binary from the build stage and frontend files
COPY --from=build /go-app /go-app
COPY ./app /go-app/app
COPY ./forum.db /go-app/forum.db

# Step 11: Expose the port the Go app will run on
EXPOSE 8080

# Step 12: Command to run the application
CMD ["./go-app"]