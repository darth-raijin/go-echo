# go-echo
go-echo is a simple HTTP server written in Go that echoes back the request details. It is useful for testing HTTP requests, debugging ingress rules, and ensuring that the correct headers and data are being passed through various layers, such as a CDN or a WAF.

# Features

- **Echo Request Data**: Displays the route, method, headers, and origin IP of incoming requests.
- **Configurable Status Code**: Allows setting a custom response status code via the `X-Response-Status` header.
- **Simple Setup**: Easy to run with Go or Docker, in the context of a distributed system.

# Getting Started

## Prerequisites

- **Go**: Ensure that Go is installed on your system, if you wish to run it using Go.
- **Docker**: If you prefer using Docker, make sure Docker is installed.

#### With Go

1. Clone the repository:

    ```bash
    git clone git@github.com:darth-raijin/go-echo.git
    ```

2. Run the server:

    Change your current directory into the freshly baked repository.

    ```bash
    cd ./go-echo
    go run cmd/server/main.go
    ```

3. The server will start on port `8080`. You can send HTTP requests to it using curl or Postman.

    ```bash
    curl -X GET http://localhost:8080/test-route -H "X-Response-Status: 202"
    ```

# Usage

You can send HTTP requests to the server, and it will respond with details of the request. For example:

    ```bash
    curl --location 'http://localhost:8080/test-route' \
    --header 'X-Response-Status: 302' \
    --header 'Accept-Charset: 202' \
    --header 'X-Some-CDN: SOME-value'
    ```

Which will make the server respond with your response status code of choice, indicated by the `X-Response-Status` header.

```
Route: /test-route
Request Method: GET
Request Headers:

- Accept-Charset: 202

- X-Response-Status: 302

- X-Some-Cdn: SOME-value

Origin IP: [::1]:53898
```

Happy debugging! 🎉