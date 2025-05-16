import socket
import os

HOST = '127.0.0.1'
PORT = 8081

MIME_TYPES = {
    '.html': 'text/html',
    '.css': 'text/css',
    '.js': 'application/javascript',
    '.jpg': 'image/jpeg',
    '.jpeg': 'image/jpeg',
    '.png': 'image/png',
    '.gif': 'image/gif',
    '.txt': 'text/plain'
}

def get_mime_type(filename):
    extension = os.path.splitext(filename)[1].lower()
    return MIME_TYPES.get(extension, 'application/octet-stream')

def send_response(client_socket, status_code, content_type, content):
    header = f"HTTP/1.1 {status_code}\r\nContent-Type: {content_type}\r\n\r\n"
    client_socket.send(header.encode() + content)

def handle_client(client_socket, client_address):
    print(f"New connection from {client_address}")
    
    try:
        request = client_socket.recv(1024).decode('utf-8')
        if not request:
            return

        request_lines = request.split('\n')
        request_parts = request_lines[0].split()
        if len(request_parts) < 2:
            return

        method, requested_path, _ = request_parts
        
        if method == 'GET' and requested_path == '/hello':
            response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
            client_socket.send(response.encode('utf-8'))
        elif method == 'POST' and requested_path == '/post':
            headers, body = request.split('\r\n\r\n', 1)
            response = f"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nReceived: {body}"
            client_socket.send(response.encode('utf-8'))
        else:
            if requested_path == '/':
                requested_path = '/index.html'
            
            file_path = requested_path[1:]
            
            if os.path.exists(file_path):
                with open(file_path, 'rb') as file:
                    file_content = file.read()
                content_type = get_mime_type(file_path)
                send_response(client_socket, "200 OK", content_type, file_content)
            else:
                error_message = "<h1>404 Not Found</h1>"
                send_response(client_socket, "404 Not Found", "text/html", error_message.encode())
    
    except Exception as error:
        print(f"Error handling request: {error}")
        error_message = "<h1>500 Internal Server Error</h1>"
        send_response(client_socket, "500 Internal Server Error", "text/html", error_message.encode())
    
    finally:
        client_socket.close()

def start_server():
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((HOST, PORT))
    server_socket.listen(5)
    
    print(f"🚀 No thred Server running on http://{HOST}:{PORT}")
    print("Press Ctrl+C to stop the server")
    
    try:
        while True:
            client_socket, client_address = server_socket.accept()
            handle_client(client_socket, client_address)
    
    except KeyboardInterrupt:
        print("\nShutting down server...")
    finally:
        server_socket.close()

if __name__ == "__main__":
    start_server()