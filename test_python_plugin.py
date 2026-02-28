#!/usr/bin/env python3
"""
Test script for the enhanced Python cProfile plugin
"""

import json
import sys
import os
import tempfile
import subprocess
import time

def test_python_plugin():
    """Test the enhanced Python cProfile plugin"""
    
    # Create a simple test Python script with performance issues
    test_script_content = '''
import time
import random

def cpu_intensive():
    """CPU intensive function"""
    result = 0
    for i in range(1000000):
        result += i * i
    return result

def memory_allocation():
    """Memory allocation function"""
    data = []
    for i in range(1000):
        obj = {
            'id': i,
            'data': [random.random() for _ in range(100)],
            'metadata': {'timestamp': time.time()}
        }
        data.append(obj)
    return data

def memory_leak():
    """Function that creates memory leak"""
    global leaked_data
    if not hasattr(memory_leak, 'counter'):
        memory_leak.counter = 0
    
    # Create and keep large data structures
    leak_chunk = []
    for i in range(100):
        leak_chunk.append({
            'id': memory_leak.counter + i,
            'data': [random.random() for _ in range(1000)],
            'timestamp': time.time()
        })
    
    if not hasattr(leaked_data, 'chunks'):
        leaked_data.chunks = []
    leaked_data.chunks.append(leak_chunk)
    memory_leak.counter += 100
    
    return len(leaked_data.chunks)

# Create global object to hold leaked data
class GlobalData:
    pass
leaked_data = GlobalData()

# Main execution
if __name__ == '__main__':
    print("Starting performance test...")
    
    # CPU intensive operation
    cpu_intensive()
    
    # Memory allocation
    memory_allocation()
    
    # Memory leak simulation
    for i in range(5):
        memory_leak()
    
    print("Performance test completed")
'''

    # Write test script to temporary file
    with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as f:
        f.write(test_script_content)
        test_script_path = f.name

    try:
        # Create temporary output directory
        with tempfile.TemporaryDirectory() as temp_dir:
            print(f"Testing Python cProfile plugin...")
            print(f"Test script: {test_script_path}")
            print(f"Output dir: {temp_dir}")
            
            # Test each profile type
            profile_types = ['cpu', 'heap', 'allocs', 'memory-leak']
            
            for profile_type in profile_types:
                print(f"\nTesting {profile_type} profile...")
                
                # Create RPC request
                request = {
                    "jsonrpc": "2.0",
                    "method": "rpc.collect",
                    "params": {
                        "target": {
                            "type": "python",
                            "command": [test_script_path]
                        },
                        "durationSec": 10,
                        "outDir": temp_dir,
                        "profiles": [profile_type],
                        "metadata": {
                            "service": "test-service",
                            "scenario": "test-scenario"
                        }
                    },
                    "id": 1
                }
                
                # Run the plugin and capture output
                result = subprocess.run(
                    ['python3', 'plugins/src/python-cprofile/main.py'],
                    input=json.dumps(request) + '\n',
                    capture_output=True,
                    text=True,
                    timeout=15
                )
                
                if result.returncode != 0:
                    print(f"ERROR: {profile_type} profile failed")
                    print(f"STDERR: {result.stderr}")
                    return False
                
                # Parse response
                try:
                    response_lines = result.stdout.strip().split('\n')
                    if response_lines:
                        response = json.loads(response_lines[-1])
                        
                        if 'error' in response:
                            print(f"ERROR: {profile_type} profile error: {response['error']}")
                            return False
                        
                        if 'result' in response:
                            bundle = response['result']
                            artifacts = bundle.get('artifacts', [])
                            
                            if not artifacts:
                                print(f"WARNING: {profile_type} profile returned no artifacts")
                            else:
                                print(f"SUCCESS: {profile_type} profile collected {len(artifacts)} artifacts")
                                
                                # Check if the expected profile file was created
                                expected_files = {
                                    'cpu': 'cpu.pb.gz',
                                    'heap': 'heap.pb.gz', 
                                    'allocs': 'allocs.pb.gz',
                                    'memory-leak': 'memory-leak.pb.gz'
                                }
                                
                                expected_file = expected_files.get(profile_type)
                                if expected_file:
                                    file_path = os.path.join(temp_dir, expected_file)
                                    if os.path.exists(file_path):
                                        file_size = os.path.getsize(file_path)
                                        print(f"  Profile file created: {expected_file} ({file_size} bytes)")
                                    else:
                                        print(f"WARNING: Expected file {expected_file} not found")
                except json.JSONDecodeError as e:
                    print(f"ERROR: Failed to parse response for {profile_type}: {e}")
                    print(f"Raw output: {result.stdout}")
                    return False
            
            print("\nAll profile types tested successfully!")
            return True
            
    finally:
        # Clean up test script
        if os.path.exists(test_script_path):
            os.unlink(test_script_path)

if __name__ == '__main__':
    success = test_python_plugin()
    sys.exit(0 if success else 1)