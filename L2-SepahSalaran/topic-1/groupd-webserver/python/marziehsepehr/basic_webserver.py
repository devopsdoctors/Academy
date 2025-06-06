import socket
import threading

def handle_client(client_socket):
    data = client_socket.recv(1024).decode('utf-8')
    method, path, _ = data.split(' ', 2)
    body = ''
    if method == 'POST':
        headers, body = data.split('\r\n\r\n', 1)
    
    if method == 'GET' and path == '/hello':
        response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
    elif method == 'POST' and path == '/post':
        response = f"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nReceived: {body}"
    else:
        response = "HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\nNot Found"
    
    client_socket.send(response.encode('utf-8'))
    client_socket.close()

server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind(('localhost', 8080))
server_socket.listen(5)

while True:
    client_socket, address = server_socket.accept()
    print(f"Connection from {address}")
    thread = threading.Thread(target=handle_client, args=(client_socket,))
    thread.start()