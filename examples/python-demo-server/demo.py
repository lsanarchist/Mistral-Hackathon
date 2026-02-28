#!/usr/bin/env python3
"""
Python demo server with intentional performance issues for profiling
"""

import time
import json
import random
from http.server import HTTPServer, BaseHTTPRequestHandler
import threading

class DemoHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/cpu-hotspot':
            self.cpu_hotspot()
        elif self.path == '/alloc-heavy':
            self.alloc_heavy()
        elif self.path == '/memory-leak':
            self.memory_leak()
        else:
            self.send_response(404)
            self.end_headers()
            self.wfile.write(b'Not Found')
    
    def cpu_hotspot(self):
        """Simulate CPU-intensive operation"""
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        
        start = time.time()
        # CPU-bound calculation
        result = 0
        for i in range(10000000):
            result += i * i
            if i % 1000000 == 0:
                # Small delay to make it more realistic
                time.sleep(0.001)
        
        self.wfile.write(f"CPU hotspot completed in {time.time() - start:.3f}s\n".encode())
    
    def alloc_heavy(self):
        """Simulate memory allocation heavy operation"""
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        
        start = time.time()
        # Allocate many small objects
        data = []
        for i in range(10000):
            obj = {
                'id': i,
                'data': [random.random() for _ in range(100)],
                'metadata': {
                    'timestamp': time.time(),
                    'value': random.random()
                }
            }
            data.append(obj)
        
        # Keep the data allocated
        global allocated_data
        allocated_data = data
        
        self.wfile.write(f"Allocation heavy completed in {time.time() - start:.3f}s\n".encode())
    
    def memory_leak(self):
        """Simulate memory leak"""
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        
        # Create a memory leak by keeping references
        global memory_leak_data
        if not hasattr(self, 'memory_leak_counter'):
            self.memory_leak_counter = 0
        
        # Allocate and keep large data structures
        leak_chunk = []
        for i in range(1000):
            leak_chunk.append({
                'id': self.memory_leak_counter + i,
                'data': [random.random() for _ in range(1000)],
                'timestamp': time.time()
            })
        
        if not hasattr(memory_leak_data, 'chunks'):
            memory_leak_data.chunks = []
        memory_leak_data.chunks.append(leak_chunk)
        self.memory_leak_counter += 1000
        
        self.wfile.write(f"Memory leak chunk {self.memory_leak_counter} created\n".encode())

def run_server():
    """Run the demo server"""
    server = HTTPServer(('localhost', 8080), DemoHandler)
    print("Python demo server running on http://localhost:8080")
    print("Endpoints:")
    print("- /cpu-hotspot - CPU intensive endpoint")
    print("- /alloc-heavy - Memory allocation heavy endpoint")
    print("- /memory-leak - Memory leak simulation endpoint")
    
    server.serve_forever()

if __name__ == '__main__':
    # Create a global object to hold allocated data
    class GlobalData:
        pass
    allocated_data = GlobalData()
    memory_leak_data = GlobalData()
    
    run_server()