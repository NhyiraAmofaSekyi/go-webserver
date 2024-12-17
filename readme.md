# Go Web Server

[![codecov](https://codecov.io/github/NhyiraAmofaSekyi/go-webserver/graph/badge.svg?token=YLNAWU1AL9)](https://codecov.io/github/NhyiraAmofaSekyi/go-webserver)

https://codecov.io/github/NhyiraAmofaSekyi/go-webserver/graphs/tree.svg?token=YLNAWU1AL9

## Overview

This project contains a Go web server that is containerized using Docker. The server responds to web requests and is designed to be lightweight and efficient, suitable for a variety of backend tasks.

## Prerequisites

Before you can build and run this server, you'll need to have the following installed:
- Docker (for containerization)
- Go (optional for local development)

## Structure

The Docker setup is divided into two stages:
1. **Build stage**: This uses `golang:alpine` to build the Go application.
2. **Final stage**: This uses `alpine:latest` to create a lightweight production image containing only the binary executable.

## Getting Started

These instructions will get your copy of the project up and running on your local machine for development and testing purposes.

### Building the Docker Image

Navigate to the directory containing the Dockerfile and run the following command to build the Docker image:

```bash
docker build -t gowebserver:0.0.1 .
```

This command builds the Docker image with the tag `gowebserver:0.0.1`, using the Dockerfile in the current directory.

### Running the Server

After building the image, run the server using:

```bash
docker run -d -p 8080:8080 --name gowebserver gowebserver:0.0.1
```

This command runs the Docker container in detached mode, maps port 8080 on the host to port 8080 in the container, and names the container `gowebserver`.

## Testing

### Health Check Endpoint

The server includes a `/v1/healthz` endpoint for health checks. To test this endpoint, use the following curl command:

```bash
curl http://localhost:8080/v1/healthz
```

This command sends a GET request to the health check endpoint. The server should respond with a status indicating that it is running correctly.

## Stopping the Server

To stop the running container, use:

```bash
docker stop gowebserver
```

And to remove the container:

```bash
docker rm gowebserver
```

## Additional Information

- **Container Configuration**: The Dockerfile sets up the environment, copies the built executable, and specifies the entry point as the application executable.
- **Security**: The final image includes only the necessary certificates and the binary, minimizing potential vulnerabilities.

