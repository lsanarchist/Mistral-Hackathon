#!/usr/bin/env python3

import json
import sys
import os
import subprocess
import tempfile
import time
from typing import Dict, Any, List

def eprint(*args, **kwargs):
    """Print to stderr"""
    print(*args, file=sys.stderr, **kwargs)

class PythonCProfilePlugin:
    def __init__(self):
        self.info = {
            "name": "python-cprofile",
            "version": "0.5.0",
            "sdkVersion": "1.0",
            "capabilities": {
                "targets": ["python"],
                "profiles": ["cpu", "heap", "allocs", "memory-leak"]
            }
        }

    def handle_rpc(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """Handle RPC requests"""
        method = request.get("method")
        params = request.get("params", {})
        rpc_id = request.get("id", 0)
        
        result = None
        error_obj = None
        
        try:
            if method == "rpc.info":
                result = self.info
            elif method == "rpc.validateTarget":
                result = self.validate_target(params)
            elif method == "rpc.collect":
                result = self.collect(params)
            else:
                error_obj = {
                    "code": -32601,
                    "message": f"Method not found: {method}"
                }
        except Exception as e:
            error_obj = {
                "code": -32603,
                "message": str(e)
            }
        
        response = {
            "jsonrpc": "2.0",
            "id": rpc_id
        }
        
        if error_obj:
            response["error"] = error_obj
        else:
            response["result"] = result
        
        return response

    def validate_target(self, target: Dict[str, Any]) -> bool:
        """Validate target configuration"""
        target_type = target.get("type")
        if target_type != "python":
            raise ValueError(f"Unsupported target type: {target_type}. Expected 'python'")
        
        # Check if required fields are present
        if "command" not in target:
            raise ValueError("Target must include 'command' field for Python execution")
        
        return True

    def collect(self, collect_request: Dict[str, Any]) -> Dict[str, Any]:
        """Collect Python profile data - enhanced to support multiple profiles like Go plugin"""
        target = collect_request.get("target", {})
        duration_sec = collect_request.get("durationSec", 10)
        out_dir = collect_request.get("outDir", ".")
        profiles = collect_request.get("profiles", ["cpu", "allocs"])
        metadata = collect_request.get("metadata", {})
        
        command = target.get("command", [])
        if not command:
            raise ValueError("No command specified in target")
        
        # Create output directory if it doesn't exist
        os.makedirs(out_dir, exist_ok=True)
        
        artifacts = []
        
        # Collect each requested profile type
        for profile_type in profiles:
            if profile_type == "cpu":
                artifact = self._collect_cpu_profile(target, duration_sec, out_dir, command)
                if artifact:
                    artifacts.append(artifact)
            elif profile_type == "heap":
                artifact = self._collect_heap_profile(target, duration_sec, out_dir, command)
                if artifact:
                    artifacts.append(artifact)
            elif profile_type == "allocs":
                artifact = self._collect_allocs_profile(target, duration_sec, out_dir, command)
                if artifact:
                    artifacts.append(artifact)
            elif profile_type == "memory-leak":
                artifact = self._collect_memory_leak_profile(target, duration_sec, out_dir, command)
                if artifact:
                    artifacts.append(artifact)
        
        if not artifacts:
            raise RuntimeError("Failed to collect any profiles")
        
        # Create artifact bundle matching the Go plugin format
        bundle = {
            "metadata": {
                "timestamp": int(time.time()),
                "durationSec": duration_sec,
                "service": metadata.get("service", "python-cprofile"),
                "scenario": metadata.get("scenario", "profiling"),
                "gitSha": metadata.get("gitSha", "")
            },
            "target": target,
            "artifacts": artifacts
        }
        
        return bundle
    
    def _collect_cpu_profile(self, target: Dict[str, Any], duration_sec: int, out_dir: str, 
                           command: List[str]) -> Dict[str, Any]:
        """Collect CPU profile using cProfile"""
        cpu_profile_path = os.path.join(out_dir, "cpu.pb.gz")
        
        try:
            # Run cProfile on the Python command
            cprofile_cmd = [
                "python3", "-m", "cProfile", "-o", cpu_profile_path
            ] + command
            
            # Run the command with timeout
            result = subprocess.run(
                cprofile_cmd,
                timeout=duration_sec,
                capture_output=True,
                text=True
            )
            
            # Check if the profile file was created
            if not os.path.exists(cpu_profile_path):
                raise RuntimeError("cProfile failed to create output file")
            
            # Return artifact in standard format
            return {
                "kind": "cprofile",
                "profileType": "cpu",
                "path": cpu_profile_path,
                "contentType": "application/octet-stream"
            }
            
        except subprocess.TimeoutExpired:
            eprint(f"CPU profiling timed out after {duration_sec} seconds")
            return None
        except Exception as e:
            eprint(f"Failed to collect CPU profile: {str(e)}")
            return None
    
    def _collect_heap_profile(self, target: Dict[str, Any], duration_sec: int, out_dir: str, 
                            command: List[str]) -> Dict[str, Any]:
        """Collect comprehensive heap profile using tracemalloc with detailed memory analysis"""
        heap_profile_path = os.path.join(out_dir, "heap.pb.gz")
        
        try:
            # Create a Python script that runs the target with comprehensive tracemalloc heap tracking
            script_content = f"""
import tracemalloc
import json
import sys
import time
import gzip
import inspect
import linecache

# Start tracemalloc with heap tracking
tracemalloc.start()

# Run the target command
start_time = time.time()

try:
    # Execute the target command as a Python script
    exec(open('{command[0]}').read())
except Exception as e:
    print("Error executing script:", e, file=sys.stderr)
    sys.exit(1)

# Get current memory snapshot
snapshot = tracemalloc.take_snapshot()

# Get comprehensive memory statistics
top_stats = snapshot.statistics('traceback')  # Get full traceback for detailed analysis

# Prepare comprehensive heap allocation data
heap_data = []
for stat in top_stats[:500]:  # Top 500 allocations for comprehensive analysis
    if stat.traceback:
        # Get the most relevant frame (usually the first user frame)
        relevant_frame = None
        for frame in stat.traceback:
            filename = frame.filename
            # Skip standard library and internal frames
            if filename and ('site-packages' not in filename and 'lib/python' not in filename):
                relevant_frame = frame
                break
        
        if relevant_frame:
            filename = relevant_frame.filename
            lineno = relevant_frame.lineno
            
            # Try to get function name from source
            function_name = 'unknown'
            try:
                source_line = linecache.getline(filename, lineno)
                if source_line and 'def ' in source_line:
                    # Extract function name from def statement
                    parts = source_line.strip().split()
                    if len(parts) > 1:
                        function_name = parts[1].split('(')[0]
            except:
                pass
            
            heap_data.append({{
                'function': function_name,
                'file': filename,
                'line': lineno,
                'size': stat.size,
                'count': stat.count,
                'total_size': stat.size * stat.count
            }})

# Write comprehensive heap data to gzipped file
with gzip.open('{heap_profile_path}', 'wt') as f:
    json.dump(heap_data, f)

print("Comprehensive heap profile saved to", '{heap_profile_path}')
"""
            
            # Write the script to a temporary file
            with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as script_file:
                script_file.write(script_content)
                script_path = script_file.name
            
            # Run the script with timeout
            result = subprocess.run(
                ["python3", script_path],
                timeout=duration_sec,
                capture_output=True,
                text=True
            )
            
            # Clean up the temporary script
            os.unlink(script_path)
            
            # Check if the profile file was created
            if not os.path.exists(heap_profile_path):
                raise RuntimeError("tracemalloc heap profiling failed to create output file")
            
            # Return artifact in standard format
            return {
                "kind": "tracemalloc",
                "profileType": "heap",
                "path": heap_profile_path,
                "contentType": "application/octet-stream"
            }
            
        except subprocess.TimeoutExpired:
            eprint(f"Heap profiling timed out after {duration_sec} seconds")
            return None
        except Exception as e:
            eprint(f"Failed to collect heap profile: {str(e)}")
            return None

    def _collect_allocs_profile(self, target: Dict[str, Any], duration_sec: int, out_dir: str, 
                              command: List[str]) -> Dict[str, Any]:
        """Collect comprehensive allocation profile using tracemalloc with function-level detail"""
        allocs_profile_path = os.path.join(out_dir, "allocs.pb.gz")
        
        try:
            # Create a Python script that runs the target with comprehensive tracemalloc
            script_content = f"""
import tracemalloc
import json
import sys
import time
import gzip
import linecache

# Start tracemalloc
tracemalloc.start()

# Run the target command
start_time = time.time()

try:
    # Execute the target command as a Python script
    exec(open('{command[0]}').read())
except Exception as e:
    print("Error executing script:", e, file=sys.stderr)
    sys.exit(1)

# Get memory allocation statistics
snapshot = tracemalloc.take_snapshot()

# Get comprehensive allocation statistics with full traceback
top_stats = snapshot.statistics('traceback')

# Prepare comprehensive allocation data with function names
allocation_data = []
for stat in top_stats[:300]:  # Top 300 allocations for comprehensive analysis
    if stat.traceback:
        # Get the most relevant frame (usually the first user frame)
        relevant_frame = None
        for frame in stat.traceback:
            filename = frame.filename
            # Skip standard library and internal frames
            if filename and ('site-packages' not in filename and 'lib/python' not in filename):
                relevant_frame = frame
                break
        
        if relevant_frame:
            filename = relevant_frame.filename
            lineno = relevant_frame.lineno
            
            # Try to get function name from source
            function_name = 'unknown'
            try:
                source_line = linecache.getline(filename, lineno)
                if source_line and 'def ' in source_line:
                    # Extract function name from def statement
                    parts = source_line.strip().split()
                    if len(parts) > 1:
                        function_name = parts[1].split('(')[0]
            except:
                pass
            
            allocation_data.append({{
                'function': function_name,
                'file': filename,
                'line': lineno,
                'size': stat.size,
                'count': stat.count,
                'total_size': stat.size * stat.count
            }})

# Write comprehensive allocation data to gzipped file
with gzip.open('{allocs_profile_path}', 'wt') as f:
    json.dump(allocation_data, f)

print("Comprehensive allocation profile saved to", '{allocs_profile_path}')
"""
            
            # Write the script to a temporary file
            with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as script_file:
                script_file.write(script_content)
                script_path = script_file.name
            
            # Run the script with timeout
            result = subprocess.run(
                ["python3", script_path],
                timeout=duration_sec,
                capture_output=True,
                text=True
            )
            
            # Clean up the temporary script
            os.unlink(script_path)
            
            # Check if the profile file was created
            if not os.path.exists(allocs_profile_path):
                raise RuntimeError("tracemalloc failed to create output file")
            
            # Return artifact in standard format
            return {
                "kind": "tracemalloc",
                "profileType": "allocs",
                "path": allocs_profile_path,
                "contentType": "application/octet-stream"
            }
            
        except subprocess.TimeoutExpired:
            eprint(f"Allocation profiling timed out after {duration_sec} seconds")
            return None
        except Exception as e:
            eprint(f"Failed to collect allocation profile: {str(e)}")
            return None
    
    def _collect_memory_leak_profile(self, target: Dict[str, Any], duration_sec: int, out_dir: str, 
                                    command: List[str]) -> Dict[str, Any]:
        """Collect memory leak profile using tracemalloc with multiple snapshots"""
        memory_leak_profile_path = os.path.join(out_dir, "memory-leak.pb.gz")
        
        try:
            # Create a Python script that runs the target with memory leak detection
            script_content = f"""
import tracemalloc
import json
import sys
import time
import gzip
import linecache

# Start tracemalloc
tracemalloc.start()

# Take initial snapshot
initial_snapshot = tracemalloc.take_snapshot()

# Run the target command
start_time = time.time()

try:
    # Execute the target command as a Python script
    exec(open('{command[0]}').read())
except Exception as e:
    print("Error executing script:", e, file=sys.stderr)
    sys.exit(1)

# Take final snapshot
final_snapshot = tracemalloc.take_snapshot()

# Compare snapshots to detect memory leaks
def compare_snapshots(old_snapshot, new_snapshot):
    old_stats = old_snapshot.statistics('traceback')
    new_stats = new_snapshot.statistics('traceback')
    
    leak_data = []
    
    # Create mapping of old allocations
    old_allocations = {{}}
    for stat in old_stats:
        if stat.traceback:
            frame = stat.traceback[0]
            key = (frame.filename, frame.lineno)
            old_allocations[key] = stat.size * stat.count
    
    # Find allocations that grew significantly
    for stat in new_stats:
        if stat.traceback:
            frame = stat.traceback[0]
            key = (frame.filename, frame.lineno)
            
            old_size = old_allocations.get(key, 0)
            new_size = stat.size * stat.count
            growth = new_size - old_size
            
            # Consider significant growth as potential leak
            if growth > 100000:  # More than 100KB growth
                # Try to get function name
                function_name = 'unknown'
                try:
                    source_line = linecache.getline(frame.filename, frame.lineno)
                    if source_line and 'def ' in source_line:
                        parts = source_line.strip().split()
                        if len(parts) > 1:
                            function_name = parts[1].split('(')[0]
                except:
                    pass
                
                leak_data.append({{
                    'function': function_name,
                    'file': frame.filename,
                    'line': frame.lineno,
                    'initial_size': old_size,
                    'final_size': new_size,
                    'growth': growth,
                    'count': stat.count
                }})
    
    return leak_data

# Detect memory leaks
leak_data = compare_snapshots(initial_snapshot, final_snapshot)

# Write memory leak data to gzipped file
with gzip.open('{memory_leak_profile_path}', 'wt') as f:
    json.dump(leak_data, f)

print("Memory leak profile saved to", '{memory_leak_profile_path}')
"""
            
            # Write the script to a temporary file
            with tempfile.NamedTemporaryFile(mode='w', suffix='.py', delete=False) as script_file:
                script_file.write(script_content)
                script_path = script_file.name
            
            # Run the script with timeout
            result = subprocess.run(
                ["python3", script_path],
                timeout=duration_sec,
                capture_output=True,
                text=True
            )
            
            # Clean up the temporary script
            os.unlink(script_path)
            
            # Check if the profile file was created
            if not os.path.exists(memory_leak_profile_path):
                raise RuntimeError("memory leak detection failed to create output file")
            
            # Return artifact in standard format
            return {
                "kind": "tracemalloc",
                "profileType": "memory-leak",
                "path": memory_leak_profile_path,
                "contentType": "application/octet-stream"
            }
            
        except subprocess.TimeoutExpired:
            eprint(f"Memory leak profiling timed out after {duration_sec} seconds")
            return None
        except Exception as e:
            eprint(f"Failed to collect memory leak profile: {str(e)}")
            return None

def main():
    plugin = PythonCProfilePlugin()
    
    # Read from stdin line by line (JSON-RPC)
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
            
        try:
            request = json.loads(line)
            response = plugin.handle_rpc(request)
            
            # Write response to stdout
            print(json.dumps(response))
            sys.stdout.flush()
            
        except json.JSONDecodeError as e:
            eprint(f"Error parsing JSON: {e}")
            continue
        except Exception as e:
            eprint(f"Error processing request: {e}")
            continue

if __name__ == "__main__":
    main()