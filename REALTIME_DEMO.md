# Real-time Monitoring Demo

## Overview

TriageProf now supports real-time monitoring with auto-refresh capability in the web viewer! This enables continuous performance analysis without manual file re-uploads.

## Features

### Auto-Refresh Controls
- **Start Auto-Refresh**: Begin continuous monitoring with configurable intervals
- **Stop Auto-Refresh**: Pause monitoring when needed
- **Refresh Now**: Manual refresh at any time
- **Interval Selection**: Choose refresh frequency (5s, 10s, 30s, 1m, 5m)

### Visual Indicators
- **Refresh Status**: Shows current monitoring state
- **Last Refresh Time**: Displays timestamp of last data update
- **Loading States**: Visual feedback during refresh operations

## Usage

### Basic Workflow

1. **Load your analysis files**
   ```bash
   # Run your analysis first
   bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 30 --outdir results/
   ```

2. **Open the web viewer**
   ```bash
   # Serve the web directory
   cd web && python3 -m http.server 8000
   
   # Open in browser
   open http://localhost:8000
   ```

3. **Load your findings.json and insights.json files**
   - Click "Load Results" button
   - Select your findings.json and optionally insights.json
   - Files will be processed and displayed

4. **Start real-time monitoring**
   - Click "Start Auto-Refresh" button
   - Select your preferred refresh interval
   - Watch as data updates automatically!

### Monitoring Scenarios

#### Production Monitoring
```bash
# Continuous monitoring of production application
while true; do
    bin/triageprof run --plugin go-pprof-http --target-url http://production-app:6060 --duration 60 --outdir monitoring/
    sleep 300  # Wait 5 minutes between collections
    cp monitoring/findings.json monitoring/findings_latest.json
    cp monitoring/insights.json monitoring/insights_latest.json
done
```

Then load `findings_latest.json` in the web viewer with auto-refresh enabled!

#### Live Debugging
```bash
# Rapid monitoring during debugging session
while true; do
    bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 10 --outdir debug/
    sleep 15  # Quick refresh for debugging
    cp debug/findings.json debug/findings_current.json
    cp debug/insights.json debug/insights_current.json
done
```

Load `findings_current.json` with 10-second auto-refresh for real-time debugging!

## Screenshots

### Refresh Controls
![Refresh Controls](https://via.placeholder.com/600x200/4a6fa5/ffffff?text=Refresh+Controls+Section)

### Active Monitoring
![Active Monitoring](https://via.placeholder.com/600x200/4ecdc4/ffffff?text=Auto-refresh+Active+with+Timestamp)

### Loading State
![Loading State](https://via.placeholder.com/600x200/ffe66d/ffffff?text=Refreshing+Data...)

## Technical Details

### Implementation
- **JavaScript**: Pure vanilla JS with no external dependencies
- **State Management**: Clean separation of refresh state
- **Error Handling**: Graceful degradation on failures
- **Performance**: Minimal overhead during refresh operations

### Configuration Options

| Interval | Use Case |
|----------|----------|
| 5 seconds | Rapid debugging, immediate feedback |
| 10 seconds | Active monitoring, quick updates |
| 30 seconds | Balanced monitoring, moderate load |
| 1 minute | Production monitoring, reduced overhead |
| 5 minutes | Long-term trends, minimal impact |

### Browser Compatibility
- ✅ Chrome 90+
- ✅ Firefox 88+
- ✅ Safari 14+
- ✅ Edge 90+

## Benefits

### For Developers
- **Immediate Feedback**: See performance changes instantly
- **Continuous Monitoring**: Track optimizations in real-time
- **Debugging Efficiency**: Rapid iteration during development

### For Operations
- **Production Visibility**: Monitor live applications continuously
- **Trend Analysis**: Track performance over time
- **Alerting**: Quickly identify performance regressions

### For Management
- **Dashboard Views**: Always-up-to-date performance metrics
- **Decision Support**: Real-time data for capacity planning
- **ROI Tracking**: Monitor optimization impact immediately

## Future Enhancements

The real-time monitoring feature provides a foundation for future improvements:

1. **Live API Integration**: Direct connection to running applications
2. **WebSocket Support**: Push-based updates without polling
3. **Alerting System**: Threshold-based notifications
4. **Historical Trends**: Time-series data visualization
5. **Multi-source Monitoring**: Aggregate data from multiple applications

## Try It Now!

1. Build the latest version: `make build`
2. Run your analysis: `bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 30 --outdir results/`
3. Open the web viewer: `cd web && python3 -m http.server 8000`
4. Load your files and start auto-refresh!

Experience the power of real-time performance monitoring with TriageProf!