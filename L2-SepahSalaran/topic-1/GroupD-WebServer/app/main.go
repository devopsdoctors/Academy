package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config represents server configuration
type Config struct {
	Port      string
	Directory string
}

// Session represents a user session
type Session struct {
	ID        string
	CreatedAt time.Time
}

// SessionManager handles all session operations
type SessionManager struct {
	sessions map[string]time.Time
	mutex    sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]time.Time),
	}
}

// GetSession returns the session for the given ID or false if not found
func (sm *SessionManager) GetSession(sessionID string) (time.Time, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	timestamp, exists := sm.sessions[sessionID]
	return timestamp, exists
}

// CreateSession creates a new session and returns the ID
func (sm *SessionManager) CreateSession() string {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sessionID := generateSessionID()
	sm.sessions[sessionID] = time.Now()
	return sessionID
}

// UpdateSession updates the timestamp for a session
func (sm *SessionManager) UpdateSession(sessionID string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.sessions[sessionID] = time.Now()
}

// CleanupSessions removes expired sessions
func (sm *SessionManager) CleanupSessions() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	now := time.Now()
	for id, lastAccess := range sm.sessions {
		if now.Sub(lastAccess) > 30*time.Minute {
			delete(sm.sessions, id)
		}
	}
}

// Server represents our HTTP server
type Server struct {
	config         Config
	sessionManager *SessionManager
	listener       net.Listener
}

// NewServer creates a new server with the given config
func NewServer(config Config) *Server {
	return &Server{
		config:         config,
		sessionManager: NewSessionManager(),
	}
}

// Start starts the server
func (s *Server) Start() error {
	log.Printf("Starting web server on port %s...", s.config.Port)
	log.Printf("Serving files from: %s", filepath.Join(s.config.Directory, "files"))
	
	// Ensure the files directory exists
	filesDir := filepath.Join(s.config.Directory, "files")
	os.MkdirAll(filesDir, 0755)
	
	var err error
	s.listener, err = net.Listen("tcp", "0.0.0.0:"+s.config.Port)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s: %v", s.config.Port, err)
	}
	
	// Start session cleanup routine
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			s.sessionManager.CleanupSessions()
		}
	}()
	
	// Accept connections
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// Stop stops the server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// Handle connection processes each incoming connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	
	for {
		// Read request line
		requestLine, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		requestLine = strings.TrimSpace(requestLine)
		if requestLine == "" {
			continue
		}
		
		// Parse request line
		parts := strings.Split(requestLine, " ")
		if len(parts) < 3 {
			sendResponse(conn, 400, "Bad Request", "text/plain", []byte("Bad Request"), nil, false, true)
			break
		}
		method := parts[0]
		path := parts[1]
		
		// Parse headers
		headers, err := parseHeaders(reader)
		if err != nil {
			break
		}
		
		// Read body if present
		var body []byte
		if clStr, ok := headers["Content-Length"]; ok {
			cl, _ := strconv.Atoi(clStr)
			body = make([]byte, cl)
			io.ReadFull(reader, body)
		}
		
		// Determine if connection should close
		closeConn := strings.ToLower(headers["Connection"]) == "close"
		clientSupportsGzip := supportsGzip(headers["Accept-Encoding"])
		
		// Handle session
		responseHeaders := make(map[string]string)
		sessionID := getSessionCookie(headers["Cookie"])
		
		if sessionID == "" {
			sessionID = s.sessionManager.CreateSession()
			responseHeaders["Set-Cookie"] = fmt.Sprintf("session=%s; Path=/", sessionID)
		} else if _, exists := s.sessionManager.GetSession(sessionID); exists {
			// Update session time
			s.sessionManager.UpdateSession(sessionID)
		} else {
			// Invalid session, create new one
			sessionID = s.sessionManager.CreateSession()
			responseHeaders["Set-Cookie"] = fmt.Sprintf("session=%s; Path=/", sessionID)
		}
		
		// Add security headers
		responseHeaders["X-Content-Type-Options"] = "nosniff"
		responseHeaders["X-Frame-Options"] = "DENY"
		responseHeaders["X-XSS-Protection"] = "1; mode=block"
		
		// Log request
		log.Printf("%s - %s %s", conn.RemoteAddr(), method, path)
		
		// Handle the request
		s.handleRequest(conn, method, path, headers, body, responseHeaders, clientSupportsGzip, closeConn)
		
		// Terminate connection if requested
		if closeConn {
			break
		}
	}
}

// Handle request processes the HTTP request
func (s *Server) handleRequest(
	conn net.Conn,
	method string,
	path string,
	headers map[string]string,
	body []byte,
	responseHeaders map[string]string,
	clientSupportsGzip bool,
	closeConn bool,
) {
	switch {
	case path == "/":
		sendResponse(conn, 200, "OK", "text/plain", []byte("Welcome to the Go Web Server"), responseHeaders, clientSupportsGzip, closeConn)
		
	case strings.HasPrefix(path, "/echo/"):
		echoString := strings.TrimPrefix(path, "/echo/")
		sendResponse(conn, 200, "OK", "text/plain", []byte(echoString), responseHeaders, clientSupportsGzip, closeConn)
		
	case path == "/user-agent":
		// FIX 1: Only allow GET method for user-agent endpoint
		if method != "GET" {
			sendResponse(conn, 405, "Method Not Allowed", "text/plain", []byte("Method not allowed"), responseHeaders, clientSupportsGzip, closeConn)
			return
		}
		userAgent := headers["User-Agent"]
		sendResponse(conn, 200, "OK", "text/plain", []byte(userAgent), responseHeaders, clientSupportsGzip, closeConn)
		
	case path == "/api/status":
		status := map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		}
		jsonResponse, _ := json.Marshal(status)
		sendResponse(conn, 200, "OK", "application/json", jsonResponse, responseHeaders, clientSupportsGzip, closeConn)
		
	case path == "/api/time":
		timeData := map[string]string{
			"time": time.Now().Format(time.RFC3339),
		}
		jsonResponse, _ := json.Marshal(timeData)
		sendResponse(conn, 200, "OK", "application/json", jsonResponse, responseHeaders, clientSupportsGzip, closeConn)
		
	case path == "/api/echo":
		if method != "POST" && method != "PUT" {
			sendResponse(conn, 405, "Method Not Allowed", "text/plain", []byte("Method not allowed"), responseHeaders, clientSupportsGzip, closeConn)
			return
		}
		contentType := "application/json"
		sendResponse(conn, 200, "OK", contentType, body, responseHeaders, clientSupportsGzip, closeConn)
		
	case path == "/api/session":
		timestamp, _ := s.sessionManager.GetSession(getSessionCookie(headers["Cookie"]))
		sessionInfo := map[string]interface{}{
			"session_id": getSessionCookie(headers["Cookie"]),
			"created_at": timestamp.Format(time.RFC3339),
			"age":        time.Since(timestamp).String(),
		}
		jsonResponse, _ := json.Marshal(sessionInfo)
		sendResponse(conn, 200, "OK", "application/json", jsonResponse, responseHeaders, clientSupportsGzip, closeConn)
		
	case strings.HasPrefix(path, "/files"):
		s.handleFiles(conn, method, path, body, responseHeaders, clientSupportsGzip, closeConn)
		
	default:
		sendResponse(conn, 404, "Not Found", "text/plain", []byte("Not Found"), responseHeaders, clientSupportsGzip, closeConn)
	}
}

// Handle files processes file-related requests
func (s *Server) handleFiles(
	conn net.Conn,
	method string,
	path string,
	body []byte,
	responseHeaders map[string]string,
	clientSupportsGzip bool,
	closeConn bool,
) {
	// Handle directory listing for /files/ root
	if path == "/files" || path == "/files/" {
		s.handleDirectoryListing(conn, responseHeaders, clientSupportsGzip, closeConn)
		return
	}
	
	// Extract filename from path
	filename := strings.TrimPrefix(path, "/files/")
	
	// FIX 2: URL-decode the filename to handle encoded traversal attempts
	var err error
	filename, err = url.QueryUnescape(filename)
	if err != nil {
		sendResponse(conn, 400, "Bad Request", "text/plain", []byte("Invalid URL encoding"), responseHeaders, clientSupportsGzip, closeConn)
		return
	}
	
	filesDir := filepath.Join(s.config.Directory, "files")
	filePath := filepath.Join(filesDir, filename)
	
	// Critical security check: prevent path traversal
	// Convert both paths to absolute and check if filePath is contained within filesDir
	absFilesDir, _ := filepath.Abs(filesDir)
	absFilePath, _ := filepath.Abs(filePath)
	
	// FIX 3: Better path traversal detection
	// If the file path is not within the files directory, return Forbidden
	if !strings.HasPrefix(absFilePath, absFilesDir) || strings.Contains(filename, "..") {
		sendResponse(conn, 403, "Forbidden", "text/plain", []byte("Path traversal not allowed"), responseHeaders, clientSupportsGzip, closeConn)
		return
	}
	
	switch method {
	case "GET":
		s.handleFileGet(conn, filePath, responseHeaders, clientSupportsGzip, closeConn)
		
	case "POST":
		s.handleFileCreate(conn, filePath, body, responseHeaders, clientSupportsGzip, closeConn)
		
	case "DELETE":
		s.handleFileDelete(conn, filePath, responseHeaders, clientSupportsGzip, closeConn)
		
	default:
		sendResponse(conn, 405, "Method Not Allowed", "text/plain", []byte("Method not allowed"), responseHeaders, clientSupportsGzip, closeConn)
	}
}

// Handle directory listing shows files in the files directory
func (s *Server) handleDirectoryListing(
	conn net.Conn,
	responseHeaders map[string]string,
	clientSupportsGzip bool,
	closeConn bool,
) {
	filesDir := filepath.Join(s.config.Directory, "files")
	files, err := ioutil.ReadDir(filesDir)
	if err != nil {
		sendResponse(conn, 500, "Internal Server Error", "text/plain", []byte("Error reading directory"), responseHeaders, clientSupportsGzip, closeConn)
		return
	}
	
	var fileList bytes.Buffer
	fileList.WriteString("<html><head><title>Directory Listing</title></head><body>")
	fileList.WriteString("<h1>Directory Listing</h1><ul>")
	
	for _, file := range files {
		fileList.WriteString(fmt.Sprintf("<li><a href=\"/files/%s\">%s</a></li>", file.Name(), file.Name()))
	}
	
	fileList.WriteString("</ul></body></html>")
	sendResponse(conn, 200, "OK", "text/html", fileList.Bytes(), responseHeaders, clientSupportsGzip, closeConn)
}

// Handle file get retrieves a file
func (s *Server) handleFileGet(
	conn net.Conn,
	filePath string,
	responseHeaders map[string]string,
	clientSupportsGzip bool,
	closeConn bool,
) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		sendResponse(conn, 404, "Not Found", "text/plain", []byte("File not found"), responseHeaders, clientSupportsGzip, closeConn)
		return
	}
	
	// Try to determine content type
	contentType := "application/octet-stream"
	ext := filepath.Ext(filePath)
	if ext == ".txt" {
		contentType = "text/plain"
	} else if ext == ".html" {
		contentType = "text/html"
	} else if ext == ".json" {
		contentType = "application/json"
	}
	
	sendResponse(conn, 200, "OK", contentType, fileData, responseHeaders, clientSupportsGzip, closeConn)
}

// Handle file create creates or updates a file
func (s *Server) handleFileCreate(
	conn net.Conn,
	filePath string,
	body []byte,
	responseHeaders map[string]string,
	clientSupportsGzip bool,
	closeConn bool,
) {
	err := ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		sendResponse(conn, 500, "Internal Server Error", "text/plain", []byte("Error writing file"), responseHeaders, clientSupportsGzip, closeConn)
		return
	}
	
	sendResponse(conn, 201, "Created", "text/plain", []byte("File created"), responseHeaders, clientSupportsGzip, closeConn)
}

// Handle file delete removes a file
func (s *Server) handleFileDelete(
	conn net.Conn,
	filePath string,
	responseHeaders map[string]string,
	clientSupportsGzip bool,
	closeConn bool,
) {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			sendResponse(conn, 404, "Not Found", "text/plain", []byte("File not found"), responseHeaders, clientSupportsGzip, closeConn)
		} else {
			sendResponse(conn, 500, "Internal Server Error", "text/plain", []byte("Error deleting file"), responseHeaders, clientSupportsGzip, closeConn)
		}
		return
	}
	
	sendResponse(conn, 200, "OK", "text/plain", []byte("File deleted"), responseHeaders, clientSupportsGzip, closeConn)
}

// Helper functions

// Parse headers parses HTTP headers from reader
func parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}
	return headers, nil
}

// Get session cookie extracts session ID from cookie header
func getSessionCookie(cookies string) string {
	if cookies == "" {
		return ""
	}
	
	cookieParts := strings.Split(cookies, ";")
	for _, part := range cookieParts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "session=") {
			return strings.TrimPrefix(part, "session=")
		}
	}
	return ""
}

// Generate session ID creates a random session ID
func generateSessionID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Supports gzip checks if client supports gzip encoding
func supportsGzip(acceptEncoding string) bool {
	if acceptEncoding == "" {
		return false
	}
	encodings := strings.Split(acceptEncoding, ",")
	for _, encoding := range encodings {
		if strings.TrimSpace(encoding) == "gzip" {
			return true
		}
	}
	return false
}

// Send response sends an HTTP response
func sendResponse(
	conn net.Conn,
	statusCode int,
	statusText string,
	contentType string,
	body []byte,
	headers map[string]string,
	supportsGzip bool,
	closeConnection bool,
) {
	responseHeaders := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
	
	if contentType != "" {
		responseHeaders += fmt.Sprintf("Content-Type: %s\r\n", contentType)
	}
	
	// Add Connection: close header if needed
	if closeConnection {
		responseHeaders += "Connection: close\r\n"
	}
	
	// Add any additional headers
	for key, value := range headers {
		responseHeaders += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	
	// Gzip compression
	if supportsGzip && len(body) > 0 {
		var compressed bytes.Buffer
		gz := gzip.NewWriter(&compressed)
		gz.Write(body)
		gz.Close()
		body = compressed.Bytes()
		responseHeaders += "Content-Encoding: gzip\r\n"
	}
	
	responseHeaders += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	responseHeaders += "\r\n"
	
	conn.Write([]byte(responseHeaders))
	if len(body) > 0 {
		conn.Write(body)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	config := Config{
		Port:      "8080",
		Directory: ".",
	}
	
	// Process command line arguments
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--directory" && i+1 < len(os.Args) {
			config.Directory = os.Args[i+1]
			i++
		} else if os.Args[i] == "--port" && i+1 < len(os.Args) {
			config.Port = os.Args[i+1]
			i++
		}
	}
	
	server := NewServer(config)
	
	// Handle graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		// os.Interrupt is equivalent to SIGINT (Ctrl+C)
		// You'll need to import "os/signal" for this
		// signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Shutting down server...")
		server.Stop()
		os.Exit(0)
	}()
	
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
