FROM golang:1.22.3-alpine

# Set the Working Directory in the Container
WORKDIR /Forum/

# Install dependencies for go-sqlite3 and Bash
RUN apk update && \
    apk add --no-cache bash gcc musl-dev

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Enable CGO and download dependencies
ENV CGO_ENABLED=1
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the app
RUN go build -o forum .

# Metadata
LABEL version="0.0.1"
LABEL projectname="FORUM"


# Run the app
CMD ["./forum"]
