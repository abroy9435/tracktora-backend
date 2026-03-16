# Use the official Golang 1.26 image
FROM golang:1.26-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build the binary targeting your main file
RUN go build -o tracktora-api cmd/api/main.go

# Expose the port Hugging Face expects
EXPOSE 7860

# Run the executable
CMD ["./tracktora-api"]