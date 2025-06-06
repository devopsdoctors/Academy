package main

import (
	"flag"
	"fmt"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	HOST = "127.0.0.1"
)

func getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("📥 New connection from %s\n", clientAddr)

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("❌ Error reading from %s: %v\n", clientAddr, err)
		return
	}

	request := string(buffer[:n])
	lines := strings.Split(request, "\n")
	if len(lines) == 0 {
		fmt.Printf("⚠️ Empty request from %s\n", clientAddr)
		return
	}

	requestLine := strings.Fields(lines[0])
	if len(requestLine) < 2 {
		fmt.Printf("⚠️ Invalid request line from %s: %s\n", clientAddr, requestLine)
		return
	}

	method := requestLine[0]
	path := requestLine[1]

	fmt.Printf("➡️ %s %s from %s\n", method, path, clientAddr)

	switch {
	case method == "GET" && path == "/hello":
		fmt.Println("✅ Handling GET /hello")
		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
		conn.Write([]byte(response))

	case method == "POST" && path == "/post":
		fmt.Println("✅ Handling POST /post")
		parts := strings.SplitN(request, "\r\n\r\n", 2)
		body := ""
		if len(parts) == 2 {
			body = parts[1]
		}
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nReceived: %s", body)
		conn.Write([]byte(response))

	default:
		if path == "/" {
			path = "/index.html"
		}
		filePath := "." + path
		if fileData, err := os.ReadFile(filePath); err == nil {
			fmt.Printf("📄 Serving file: %s\n", filePath)
			contentType := getMimeType(filePath)
			header := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\n\r\n", contentType)
			conn.Write([]byte(header))
			conn.Write(fileData)
		} else {
			fmt.Printf("🚫 File not found: %s\n", filePath)
			response := "HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\n\r\n<h1>404 Not Found</h1>"
			conn.Write([]byte(response))
		}
	}
}

func startServer(port string) {
	address := HOST + ":" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("❌ Error starting server on port %s: %v\n", port, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("🚀 Server running on http://%s\n", address)
	fmt.Println("🔒 Press Ctrl+C to stop the server")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("⚠️ Connection error:", err)
			continue
		}
		go handleClient(conn)
	}
}

func main() {
	port := flag.String("port", "8082", "Port number to listen on")
	flag.Parse()

	startServer(*port)
}

