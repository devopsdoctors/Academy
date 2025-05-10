import http.client
import threading
import time

def make_request(server_port, path, method="GET", body=None):
    try:
        conn = http.client.HTTPConnection("127.0.0.1", server_port)
        if method == "POST":
            conn.request(method, path, body=body, headers={"Content-Type": "text/plain"})
        else:
            conn.request(method, path)
        response = conn.getresponse()
        print(f"Request to http://127.0.0.1:{server_port}{path} -> Status: {response.status}, Response: {response.read().decode()}")
        conn.close()
    except Exception as e:
        print(f"Error making request to http://127.0.0.1:{server_port}{path}: {e}")

def stress_test(server_port, num_requests=50):
    threads = []
    start_time = time.time()
    
    # Test all endpoints for both servers
    paths = ["/", "/hello", "/post"]
    methods = ["GET", "GET", "POST"]
    
    for i in range(num_requests):
        path = paths[i % len(paths)]
        method = methods[i % len(paths)]
        body = "Test data" if method == "POST" else None
        thread = threading.Thread(target=make_request, args=(server_port, path, method, body))
        threads.append(thread)
        thread.start()
    
    for thread in threads:
        thread.join()
    
    return f"Completed {num_requests} requests to port {server_port} in {time.time() - start_time:.7f} seconds"

if __name__ == "__main__":
    print("Testing Threaded Server (port 8081)")
    result1=stress_test(8081, num_requests=50)
    print("\nTesting Non-Threaded Server (port 8082)")
    result2=stress_test(8082, num_requests=50)

    print("Testing 1 No Threaded Server (port 8082)")
    print(result1)
    print("Testing result for with Thread Server (port 8081)")
    print(result2)
    