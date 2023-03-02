# Use an official Golang runtime as a parent image
FROM golang:1.20-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go app
RUN go build -o main .

# Set the PORT environment variable
ENV PORT=8080

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go build`
CMD ["./main"]
