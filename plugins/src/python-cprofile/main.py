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
            "version": "0.2.0",
            "sdkVersion": "1.0",
            "capabilities": {
                "targets": ["python"],
                "profiles": ["cpu", "allocs"]
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
        """Collect Python profile data"""
        target = collect_request.get("target", {})
        duration_sec = collect_request.get("durationSec", 10)
        out_dir = collect_request.get("outDir", ".")
        profile_type = collect_request.get("profileType", "cpu")
        
        command = target.get("command", [])
        if not command:
            raise ValueError("No command specified in target")
        
        # Create output directory if it doesn't exist
        os.makedirs(out_dir, exist_ok=True)
        
        # Generate output file paths
        timestamp = int(time.time())
        
        if profile_type == "cpu":
            return self._collect_cpu_profile(target, duration_sec, out_dir, timestamp, command)
        elif profile_type == "allocs":
            return self._collect_allocs_profile(target, duration_sec, out_dir, timestamp, command)
        else:
            raise ValueError(f"Unsupported profile type: {profile_type}")
    
    def _collect_cpu_profile(self, target: Dict[str, Any], duration_sec: int, out_dir: str, 
                           timestamp: int, command: List[str]) -> Dict[str, Any]:
        """Collect CPU profile using cProfile"""
        cpu_profile_path = os.path.join(out_dir, f"cpu_{timestamp}.prof")
        
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
            
            # Create artifact bundle
            bundle = {
                "metadata": {
                    "timestamp": int(time.time()),
                    "durationSec": duration_sec,
                    "service": "python-cprofile",
                    "scenario": "cpu-profiling"
                },
                "target": target,
                "artifacts": [
                    {
                        "kind": "cprofile",
                        "profileType": "cpu",
                        "path": cpu_profile_path,
                        "contentType": "application/octet-stream"
                    }
                ]
            }
            
            return bundle
            
        except subprocess.TimeoutExpired:
            raise RuntimeError(f"Command timed out after {duration_sec} seconds")
        except Exception as e:
            # Clean up profile file if it exists
            if os.path.exists(cpu_profile_path):
                os.remove(cpu_profile_path)
            raise RuntimeError(f"Failed to collect CPU profile: {str(e)}")
    
    def _collect_allocs_profile(self, target: Dict[str, Any], duration_sec: int, out_dir: str, 
                              timestamp: int, command: List[str]) -> Dict[str, Any]:
        """Collect allocation profile using tracemalloc"""
        allocs_profile_path = os.path.join(out_dir, f"allocs_{timestamp}.json")
        
        try:
            # Create a Python script that runs the target with tracemalloc
            script_content = f"""
import tracemalloc
import json
import sys
import time

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

# Write allocation data to file
with open('{allocs_profile_path}', 'w') as f:
    json.dump({{
        'metadata': {{
            'timestamp': time.time(),
            'durationSec': time.time() - start_time,
            'service': 'python-cprofile',
            'scenario': 'allocation-profiling'
        }},
        'allocations': allocation_data
    }}, f)

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
            
            # Create artifact bundle
            bundle = {
                "metadata": {
                    "timestamp": int(time.time()),
                    "durationSec": duration_sec,
                    "service": "python-cprofile",
                    "scenario": "allocation-profiling"
                },
                "target": target,
                "artifacts": [
                    {
                        "kind": "tracemalloc",
                        "profileType": "allocs",
                        "path": allocs_profile_path,
                        "contentType": "application/json"
                    }
                ]
            }
            
            return bundle
            
        except subprocess.TimeoutExpired:
            raise RuntimeError(f"Command timed out after {duration_sec} seconds")
        except Exception as e:
            # Clean up profile file if it exists
            if os.path.exists(allocs_profile_path):
                os.remove(allocs_profile_path)
            raise RuntimeError(f"Failed to collect allocation profile: {str(e)}")

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