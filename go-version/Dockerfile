# Start from the official Go image
FROM golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
# The output binary is named after the project
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gmail-ai-console cmd/main.go

# The application reads config files from the current directory.
# The user will need to mount the directory with config files.

# Set the entrypoint
ENTRYPOINT ["/app/gmail-ai-console"]
