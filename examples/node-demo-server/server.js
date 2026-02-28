// Simple Node.js server for testing the node-inspector plugin
const http = require('http');

const server = http.createServer((req, res) => {
  res.writeHead(200, {'Content-Type': 'text/plain'});
  res.end('Hello from Node.js demo server!\n');
});

const PORT = 3000;
server.listen(PORT, () => {
  console.log(`Node.js demo server running on http://localhost:${PORT}`);
});

// CPU-intensive function for profiling
function cpuIntensive() {
  let result = 0;
  for (let i = 0; i < 10000000; i++) {
    result += Math.sqrt(i) * Math.random();
  }
  return result;
}

// Memory allocation function
function memoryAllocation() {
  const bigArray = new Array(10000).fill({
    data: new Array(100).fill(Math.random())
  });
  return bigArray;
}

// Call these functions periodically to generate profiling data
setInterval(() => {
  cpuIntensive();
  memoryAllocation();
}, 1000);