package main

import (
	"bufio"
	"flag"
	"fmt"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const HOST = "127.0.0.1"

func getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}
	return "application/octet-stream"
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil || len(requestLine) < 2 {
		return
	}

	requestFields := strings.Fields(requestLine)
	if len(requestFields) < 2 {
		return
	}
	method, path := requestFields[0], requestFields[1]

	// Skip the rest of the headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
	}

	switch {
	case method == "GET" && path == "/hello":
		fmt.Fprint(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!")

	case method == "POST" && path == "/post":
		body := make([]byte, 1024)
		n, _ := reader.Read(body)
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nReceived: %s", string(body[:n]))

	default:
		if path == "/" {
			path = "/index.html"
		}
		filePath := "." + path
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprint(conn, "HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\n\r\n<h1>404 Not Found</h1>")
			return
		}
		contentType := getMimeType(filePath)
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: %s\r\n\r\n", contentType)
		conn.Write(data)
	}
}

func startServer(port string) {
	addr := HOST + ":" + port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Server running on http://%s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func main() {
	port := flag.String("port", "8086", "Port to listen on")
	flag.Parse()
	startServer(*port)
}
