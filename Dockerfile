# RVParkBackend/Dockerfile

FROM golang:1.21

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Make sure init.sql is in the right place and accessible
RUN mkdir -p /app/docker
COPY docker/init.sql /app/docker/init.sql
RUN chmod 644 /app/docker/init.sql

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]