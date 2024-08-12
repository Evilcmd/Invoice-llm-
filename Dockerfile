# Use an official Go runtime as a parent image
FROM golang:1.22.5

# Install xpdf-reader
RUN apt-get update && apt-get install -y xpdf

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 for the Go server
EXPOSE 8080

# Run the startup script when the container launches
CMD ["./main"]