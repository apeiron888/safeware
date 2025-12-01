FROM golang:1.21-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git postgresql-client

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Run the application
CMD ["go", "run", "cmd/server/main.go"]
