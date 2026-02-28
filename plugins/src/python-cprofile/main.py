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
            "version": "0.4.0",
            "sdkVersion": "1.0",
            "capabilities": {
                "targets": ["python"],
                "profiles": ["cpu", "heap", "allocs"]
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
        """Collect heap profile using tracemalloc with heap snapshot"""
        heap_profile_path = os.path.join(out_dir, "heap.pb.gz")
        
        try:
            # Create a Python script that runs the target with tracemalloc heap tracking
            script_content = f"""
import tracemalloc
import json
import sys
import time
import gzip

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

# Get top memory allocations by size
top_stats = snapshot.statistics('lineno')

# Prepare heap allocation data in pprof-compatible format
heap_data = []
for stat in top_stats[:200]:  # Top 200 allocations for comprehensive analysis
    if stat.traceback:
        frame = stat.traceback[0]
        heap_data.append({{
            'function': 'unknown',  # Simplified - we can't easily get function name from traceback
            'file': frame.filename,
            'line': frame.lineno,
            'size': stat.size,
            'count': stat.count
        }})

# Write heap data to gzipped file (matching Go plugin format)
with gzip.open('{heap_profile_path}', 'wt') as f:
    json.dump(heap_data, f)

print("Heap profile saved to", '{heap_profile_path}')
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
        """Collect allocation profile using tracemalloc"""
        allocs_profile_path = os.path.join(out_dir, "allocs.pb.gz")
        
        try:
            # Create a Python script that runs the target with tracemalloc
            script_content = f"""
import tracemalloc
import json
import sys
import time
import gzip

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

# Get top memory allocations
top_stats = snapshot.statistics('lineno')

# Prepare allocation data
allocation_data = []
for stat in top_stats[:100]:  # Top 100 allocations
    allocation_data.append({{
        'filename': stat.traceback[0].filename,
        'lineno': stat.traceback[0].lineno,
        'size': stat.size,
        'count': stat.count
    }})

# Write allocation data to gzipped file (matching Go plugin format)
with gzip.open('{allocs_profile_path}', 'wt') as f:
    json.dump(allocation_data, f)

print("Allocation profile saved to", '{allocs_profile_path}')
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