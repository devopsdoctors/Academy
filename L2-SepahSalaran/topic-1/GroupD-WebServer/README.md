# Go Web Server

A lightweight, feature-rich HTTP server written in Go with support for file operations, session management, API endpoints, and security features.

## Features

- **Basic HTTP functionality** - Serves static content with proper headers
- **File operations** - Upload, download, and delete files through API endpoints
- **Session management** - Track user sessions with cookies
- **API endpoints** - JSON-based API for status, time, and echo functionality
- **Security features** - Protection against path traversal attacks, secure headers
- **Compression** - Gzip support for bandwidth optimization
- **Directory listing** - Browse files in the server's storage directory

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/mehrdadbn9/codecrafters-http-server-go
   cd go-web-server
   ```

2. Build the server:
   ```
   go build -o server main.go
   ```

### Usage

Run the server with optional configuration flags:

```
./server [--port PORT] [--directory DIRECTORY]
```

Parameters:
- `--port` - TCP port to listen on (default: 8080)
- `--directory` - Base directory for file storage (default: current directory)

Example:
```
./server --port 9000 --directory /var/www
```

## API Documentation

### Basic Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Returns a welcome message |
| `/echo/{string}` | GET | Echoes the provided string |
| `/user-agent` | GET | Returns the client's user agent |

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/status` | GET | Returns server status in JSON format |
| `/api/time` | GET | Returns current server time in JSON format |
| `/api/echo` | POST/PUT | Echoes the request body |
| `/api/session` | GET | Returns current session information |

### File Operations

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/files/` | GET | Lists all files in the files directory |
| `/files/{filename}` | GET | Downloads the specified file |
| `/files/{filename}` | POST | Creates or updates a file |
| `/files/{filename}` | DELETE | Deletes the specified file |

## Testing

A comprehensive test script is included to verify all server functionality.

### Running Tests

1. Make the test script executable:
   ```
   chmod +x webserver-test.sh
   ```

2. Run the tests:
   ```
   ./webserver-test.sh [host] [port]
   ```

   Default host is "localhost"
   Default port is 8080

### What the Tests Cover

The script tests 25 different aspects of the web server:

#### Basic Functionality
- Root endpoint (/)
- Echo endpoint (/echo/hello-world)
- User-agent endpoint (/user-agent)

#### API Endpoints
- Status (/api/status)
- Time (/api/time)
- Echo (/api/echo)
- Session (/api/session)

#### File Operations
- Create a file (POST to /files/test.txt)
- Get a file (GET from /files/test.txt)
- Delete a file (DELETE /files/test.txt)
- Get non-existent file (GET /files/nonexistent.txt)

#### Session Management
- Session cookie handling (using cookie jar)
- Session persistence

#### Security Tests
- Path traversal attempt prevention
- Security headers validation

#### Content Features
- Gzip compression support
- Directory listing
- JSON content type verification

#### Edge Cases
- Large request body handling
- Multiple concurrent requests
- Very long URLs
- Long headers

## Security Features

- Protection against path traversal attacks
- Secure headers (X-Content-Type-Options, X-Frame-Options, X-XSS-Protection)
- Session expiration and cleanup
- Input validation

## Performance

The server is designed for concurrent operation and can handle multiple simultaneous connections efficiently.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
