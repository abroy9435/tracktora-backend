# 1. Use the official Golang image
FROM golang:1.22-alpine

# 2. Set the working directory inside the container
WORKDIR /app

# 3. Copy dependency files (go.mod and go.sum)
COPY go.mod go.sum ./

# 4. Download all dependencies
RUN go mod download

# 5. Copy the rest of your source code
COPY . .

# 6. Build the binary targeting your main file
# This assumes your main.go is located at cmd/api/main.go
RUN go build -o tracktora-api cmd/api/main.go

# 7. Expose the port Hugging Face expects
EXPOSE 7860

# 8. Run the executable
CMD ["./tracktora-api"]