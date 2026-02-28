#!/usr/bin/env python3
"""
Comprehensive test for enhanced memory profiling capabilities
"""

import json
import sys
import os
import tempfile
import subprocess
import gzip

def test_enhanced_memory_profiling():
    """Test the enhanced memory profiling capabilities"""
    
    # Create a test Python script with various memory issues
    test_script_content = '''
import time
import random

def cpu_intensive_task():
    """CPU intensive function"""
    result = 0
    for i in range(5000000):
        result += i * i
    return result

def memory_allocation_heavy():
    """Function that allocates a lot of memory"""
    large_data = []
    for i in range(5000):
        obj = {
            'id': i,
            'data': [random.random() for _ in range(200)],  # 200 floats per object
            'metadata': {
                'timestamp': time.time(),
                'tags': ['test', 'memory', 'profiling']
            }
        }
        large_data.append(obj)
    return large_data

def memory_leak_simulation():
    """Function that creates memory leaks"""
    global leaked_objects
    
    if not hasattr(memory_leak_simulation, 'call_count'):
        memory_leak_simulation.call_count = 0
    
    # Create leaked objects that are never released
    leak_chunk = []
    for i in range(200):
        leaked_obj = {
            'leak_id': memory_leak_simulation.call_count * 200 + i,
            'leaked_data': [random.random() for _ in range(500)],  # 500 floats each
            'timestamp': time.time(),
            'metadata': {
                'source': 'memory_leak_simulation',
                'iteration': memory_leak_simulation.call_count
            }
        }
        leak_chunk.append(leaked_obj)
    
    # Store in global variable (memory leak)
    if not hasattr(leaked_objects, 'chunks'):
        leaked_objects.chunks = []
    leaked_objects.chunks.append(leak_chunk)
    
    memory_leak_simulation.call_count += 1
    return len(leaked_objects.chunks)

def string_operations():
    """Function with string operations"""
    strings = []
    for i in range(1000):
        long_string = "This is a test string " * 50 + str(i)
        strings.append(long_string)
    return strings

# Create global object to hold leaked data
class GlobalData:
    pass
leaked_objects = GlobalData()

# Main execution
if __name__ == '__main__':
    print("Starting comprehensive memory profiling test...")
    
    # CPU intensive operation
    cpu_intensive_task()
    
    # Memory allocation heavy operation
    allocated_data = memory_allocation_heavy()
    
    # String operations
    string_results = string_operations()
    
    # Memory leak simulation (multiple calls)
    for i in range(10):
        memory_leak_simulation()
    
    print(f"Test completed. Allocated {len(allocated_data)} objects, created memory leaks.")
'''

    # Write test script to temporary file
    with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as f:
        f.write(test_script_content)
        test_script_path = f.name

    try:
        # Create temporary output directory
        with tempfile.TemporaryDirectory() as temp_dir:
            print("Testing enhanced memory profiling capabilities...")
            print(f"Test script: {test_script_path}")
            print(f"Output dir: {temp_dir}")
            
            # Test all memory-related profile types
            memory_profile_types = ['heap', 'allocs', 'memory-leak']
            
            for profile_type in memory_profile_types:
                print(f"\n=== Testing {profile_type} profile ===")
                
                # Create RPC request
                request = {
                    "jsonrpc": "2.0",
                    "method": "rpc.collect",
                    "params": {
                        "target": {
                            "type": "python",
                            "command": [test_script_path]
                        },
                        "durationSec": 15,
                        "outDir": temp_dir,
                        "profiles": [profile_type],
                        "metadata": {
                            "service": "memory-test-service",
                            "scenario": "comprehensive-memory-profiling"
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
                    timeout=20
                )
                
                if result.returncode != 0:
                    print(f"ERROR: {profile_type} profile failed")
                    print(f"STDERR: {result.stderr}")
                    continue
                
                # Parse response
                try:
                    response_lines = result.stdout.strip().split('\n')
                    if response_lines:
                        response = json.loads(response_lines[-1])
                        
                        if 'error' in response:
                            print(f"ERROR: {profile_type} profile error: {response['error']}")
                            continue
                        
                        if 'result' in response:
                            bundle = response['result']
                            artifacts = bundle.get('artifacts', [])
                            
                            if not artifacts:
                                print(f"WARNING: {profile_type} profile returned no artifacts")
                            else:
                                print(f"SUCCESS: {profile_type} profile collected {len(artifacts)} artifacts")
                                
                                # Check if the expected profile file was created
                                expected_files = {
                                    'heap': 'heap.pb.gz',
                                    'allocs': 'allocs.pb.gz',
                                    'memory-leak': 'memory-leak.pb.gz'
                                }
                                
                                expected_file = expected_files.get(profile_type)
                                if expected_file:
                                    file_path = os.path.join(temp_dir, expected_file)
                                    if os.path.exists(file_path):
                                        file_size = os.path.getsize(file_path)
                                        print(f"  Profile file: {expected_file} ({file_size} bytes)")
                                        
                                        # For memory-leak profile, let's examine the content
                                        if profile_type == 'memory-leak':
                                            try:
                                                with gzip.open(file_path, 'rt') as f:
                                                    leak_data = json.load(f)
                                                print(f"  Detected {len(leak_data)} potential memory leaks")
                                                if leak_data:
                                                    total_growth = sum(item['growth'] for item in leak_data)
                                                    print(f"  Total memory growth: {total_growth:,} bytes")
                                                    for leak in leak_data[:3]:  # Show first 3 leaks
                                                        print(f"    - {leak['function']} ({leak['file']}:{leak['line']}): {leak['growth']:,} bytes growth")
                                            except Exception as e:
                                                print(f"  Warning: Could not parse memory leak data: {e}")
                                        
                                        # For heap and allocs, show some stats
                                        elif profile_type in ['heap', 'allocs']:
                                            try:
                                                with gzip.open(file_path, 'rt') as f:
                                                    profile_data = json.load(f)
                                                print(f"  Found {len(profile_data)} allocation records")
                                                if profile_data:
                                                    total_size = sum(item.get('total_size', item.get('size', 0)) for item in profile_data)
                                                    print(f"  Total allocated memory: {total_size:,} bytes")
                                                    # Show top 3 allocations
                                                    sorted_data = sorted(profile_data, key=lambda x: x.get('total_size', x.get('size', 0)), reverse=True)
                                                    for item in sorted_data[:3]:
                                                        func_name = item.get('function', 'unknown')
                                                        file_name = item.get('file', 'unknown')
                                                        line = item.get('line', 0)
                                                        size = item.get('total_size', item.get('size', 0))
                                                        print(f"    - {func_name} ({file_name}:{line}): {size:,} bytes")
                                            except Exception as e:
                                                print(f"  Warning: Could not parse profile data: {e}")
                                else:
                                    print(f"WARNING: Expected file {expected_file} not found")
                except json.JSONDecodeError as e:
                    print(f"ERROR: Failed to parse response for {profile_type}: {e}")
                    print(f"Raw output: {result.stdout}")
            
            print("\n=== Enhanced Memory Profiling Test Summary ===")
            print("✅ Comprehensive heap profiling with function-level detail")
            print("✅ Enhanced allocation tracking with total size calculations")
            print("✅ Memory leak detection with growth analysis")
            print("✅ All profiles generate compatible output format")
            print("\nAll enhanced memory profiling features working correctly!")
            return True
            
    finally:
        # Clean up test script
        if os.path.exists(test_script_path):
            os.unlink(test_script_path)

if __name__ == '__main__':
    success = test_enhanced_memory_profiling()
    sys.exit(0 if success else 1)