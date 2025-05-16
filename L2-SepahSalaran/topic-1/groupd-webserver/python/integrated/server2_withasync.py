import asyncio
import os
import mimetypes

HOST = '127.0.0.1'
PORT = 8083  

def get_mime_type(file_path):
    mime_type, _ = mimetypes.guess_type(file_path)
    return mime_type or 'application/octet-stream'


async def handle_client(reader, writer):
    addr = writer.get_extra_info('peername') or ('unknown', 0)
    print(f"New connection from {addr[0]}:{addr[1]}")

    try:
        request = await reader.read(1024)
        request = request.decode('utf-8')
        if not request:
            writer.close()
            await writer.wait_closed()
            return

        lines = request.split('\n')
        request_line = lines[0].strip()
        parts = request_line.split()

        if len(parts) < 2:
            writer.close()
            await writer.wait_closed()
            return

        method, path = parts[0], parts[1]

        if method == 'GET' and path == '/hello':
            response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
            writer.write(response.encode())

        elif method == 'POST' and path == '/post':
            body = request.split('\r\n\r\n', 1)[1] if '\r\n\r\n' in request else ''
            response = f"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nReceived: {body}"
            writer.write(response.encode())

        else:
            if path == '/':
                path = '/index.html'
            file_path = '.' + path

            if os.path.exists(file_path) and os.path.isfile(file_path):
                with open(file_path, 'rb') as f:
                    content = f.read()
                mime_type = get_mime_type(file_path)
                header = f"HTTP/1.1 200 OK\r\nContent-Type: {mime_type}\r\n\r\n"
                writer.write(header.encode() + content)
            else:
                error = "HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\n\r\n<h1>404 Not Found</h1>"
                writer.write(error.encode())
        
        await writer.drain()

    except Exception as e:
        error = "HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/html\r\n\r\n<h1>500 Server Error</h1>"
        writer.write(error.encode())

    finally:
        writer.close()
        await writer.wait_closed()

async def main():
    server = await asyncio.start_server(handle_client, HOST, PORT)
    print(f"🚀 Asyncio Server running on http://{HOST}:{PORT}")
    async with server:
        await server.serve_forever()

if __name__ == '__main__':
    asyncio.run(main())
