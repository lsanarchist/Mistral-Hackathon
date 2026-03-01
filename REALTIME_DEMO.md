# Real-time Monitoring Demo Guide

## Quick Start

### 1. Start the Demo Server
```bash
# Start the demo server in one terminal
go run examples/demo-server/main.go &

# Generate some load in another terminal
./examples/load.sh &
```

### 2. Run TriageProf Analysis
```bash
# Run full analysis pipeline
export MISTRAL_API_KEY="your-key-here"
./bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 30 --outdir demo-realtime --llm
```

### 3. Start WebSocket Server
```bash
# Start WebSocket server with the analysis results
./bin/triageprof websocket --findings ./demo-realtime/findings.json --insights ./demo-realtime/insights.json --port 8080
```

### 4. Connect Web Viewer
Open `web/index.html` in your browser and:
1. Click "Connect" in the WebSocket controls
2. Use default URL: `ws://localhost:8080/ws`
3. Watch real-time performance updates!

## Demo Script

### Full Demo Flow
```bash
#!/bin/bash

# Clean up any previous runs
rm -rf demo-realtime
pkill -f "demo-server" || true
pkill -f "load.sh" || true

# Start demo server
echo "🚀 Starting demo server..."
go run examples/demo-server/main.go > /dev/null 2>&1 &
SERVER_PID=$!
sleep 2

# Generate load
echo "📊 Generating load..."
./examples/load.sh > /dev/null 2>&1 &
LOAD_PID=$!
sleep 5

# Run analysis
echo "🔍 Running performance analysis..."
export MISTRAL_API_KEY="your-key-here"
./bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 30 --outdir demo-realtime --llm

# Start WebSocket server
echo "🌐 Starting WebSocket server..."
./bin/triageprof websocket --findings ./demo-realtime/findings.json --insights ./demo-realtime/insights.json --port 8080 > /dev/null 2>&1 &
WS_PID=$!
sleep 2

# Open web viewer
echo "📱 Opening web viewer..."
open web/index.html || xdg-open web/index.html

# Show demo info
echo ""
echo "🎯 Demo Ready!"
echo "📊 WebSocket endpoint: ws://localhost:8080/ws"
echo "🌡 Health check: http://localhost:8080/health"
echo "📱 Web viewer: http://localhost:8080/web/index.html"
echo ""
echo "💡 Connect to WebSocket in the web viewer to see real-time updates!"
echo ""
echo "Press Ctrl+C to stop the demo..."

# Wait for user to stop
trap "kill $SERVER_PID $LOAD_PID $WS_PID 2>/dev/null; exit" INT
wait

# Cleanup
kill $SERVER_PID $LOAD_PID $WS_PID 2>/dev/null
echo "🧹 Demo cleaned up"
```

## WebSocket Features to Highlight

### 1. Real-time Data Streaming
- **Instant updates**: Performance metrics update immediately when data changes
- **No refresh needed**: Data flows continuously without manual intervention
- **Multi-client support**: Multiple team members can monitor simultaneously

### 2. Connection Management
- **Automatic reconnection**: Client reconnects automatically if connection drops
- **Status indicators**: Visual feedback for connection state
- **Error handling**: Graceful degradation on connection issues

### 3. Performance Dashboard
- **Live statistics**: Total findings, severity counts, performance scores
- **Real-time charts**: Severity distribution and category breakdowns update dynamically
- **AI insights**: LLM-generated analysis updates as new data arrives

### 4. Production Monitoring
- **Health endpoint**: Monitor server status at `/health`
- **Scalable architecture**: Handles multiple concurrent connections
- **Low latency**: Updates delivered in milliseconds

## Demo Talking Points

### "Wow" Moments
1. **Instant Connection**: "Watch how quickly the WebSocket connects and starts streaming data"
2. **Live Updates**: "See the performance metrics update in real-time as the application runs"
3. **Multi-client**: "Open multiple browser tabs - they all receive updates simultaneously"
4. **No Refresh**: "Notice how the data updates without any page refreshes or manual intervention"

### Key Benefits
- **Production Ready**: "This isn't just a demo - it's production-grade monitoring"
- **Collaborative**: "Entire teams can monitor performance together in real-time"
- **Efficient**: "Low overhead means you can monitor continuously without performance impact"
- **Extensible**: "The WebSocket API allows integration with any monitoring dashboard"

### Technical Highlights
- **WebSocket Protocol**: "Uses standard WebSocket protocol for maximum compatibility"
- **JSON API**: "Clean JSON interface makes it easy to integrate with any client"
- **Error Resilient**: "Automatic reconnection ensures monitoring continues even if there are network issues"
- **Scalable**: "Designed to handle production workloads with multiple concurrent clients"

## Troubleshooting

### Common Issues

**Port already in use**
```bash
# Find and kill process using port 8080
lsof -i :8080
kill -9 <PID>
```

**WebSocket connection fails**
- Check server is running: `curl http://localhost:8080/health`
- Verify WebSocket URL: `ws://localhost:8080/ws`
- Check browser console for errors

**No data appearing**
- Ensure findings.json exists and is valid
- Check WebSocket server logs for errors
- Verify file paths in the websocket command

## Advanced Demo Scenarios

### 1. Multiple Applications
```bash
# Start multiple WebSocket servers on different ports
./bin/triageprof websocket --findings app1-findings.json --port 8081 &
./bin/triageprof websocket --findings app2-findings.json --port 8082 &

# Connect different clients to different ports
```

### 2. Custom Dashboard Integration
```javascript
// Example: Connect WebSocket to custom dashboard
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    
    // Update custom dashboard elements
    document.getElementById('critical-count').textContent = data.stats.critical_count;
    document.getElementById('performance-score').textContent = data.stats.performance_score;
    
    // Trigger visual updates
    updateCharts(data.findings);
    updateAlerts(data.insights);
};
```

### 3. Automated Monitoring Script
```bash
#!/bin/bash

# Continuous monitoring script
while true; do
    # Run analysis
    ./bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 60 --outdir monitoring-latest
    
    # Restart WebSocket server with new data
    pkill -f "websocket" || true
    ./bin/triageprof websocket --findings ./monitoring-latest/findings.json --insights ./monitoring-latest/insights.json --port 8080 &
    
    # Wait before next cycle
    sleep 300
 done
```

## Success Metrics

### Technical Validation
- ✅ WebSocket server starts successfully
- ✅ Health endpoint returns healthy status
- ✅ Clients can connect and receive data
- ✅ Real-time updates appear in web viewer
- ✅ Multiple clients can connect simultaneously

### Demo Experience
- ✅ "Wow" reaction from audience
- ✅ Clear understanding of real-time capabilities
- ✅ Appreciation of production readiness
- ✅ Interest in using for their own applications
- ✅ Questions about integration and scaling

## Conclusion

This real-time monitoring demo showcases TriageProf's evolution from a batch analysis tool to a comprehensive performance monitoring platform. The WebSocket implementation provides the foundation for production-grade continuous analysis while maintaining the simplicity and ease-of-use that developers love.

**Key Takeaway**: TriageProf now enables teams to monitor application performance in real-time, receive instant alerts about issues, and make data-driven optimization decisions faster than ever before.