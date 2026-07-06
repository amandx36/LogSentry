FROM golang:1.22-alpine

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o server ./cmd/server

# Expose port
EXPOSE 8080

# Start container
CMD ["./server"]