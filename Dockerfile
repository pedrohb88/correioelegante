# Use the official Go image as the base image
FROM golang:1.20 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o app

# Use a smaller base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the previous build stage
COPY --from=build /app/app .

# Expose the port that the application will listen on
EXPOSE 8080

# Run the application
CMD ["./app"]
