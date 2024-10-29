FROM golang:1.22.3-alpine

# Set the Working Directory in the Container
WORKDIR /Forum/

# Install Bash(For Alpine)
RUN apk update && apk add bash

# Copy the Project into the Container
COPY . .

# Build the app
RUN go build -o forum .

# Metadata
LABEL version="0.0.1"
LABEL projectname="FORUM"

# Run the app
CMD ["./forum"]