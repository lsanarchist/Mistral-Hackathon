document.addEventListener('DOMContentLoaded', function() {
    const loadBtn = document.getElementById('loadBtn');
    const fileInput = document.getElementById('fileInput');
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const content = document.getElementById('content');
    const severityFilter = document.getElementById('severityFilter');
    const refreshControls = document.getElementById('refreshControls');
    const startRefreshBtn = document.getElementById('startRefreshBtn');
    const stopRefreshBtn = document.getElementById('stopRefreshBtn');
    const refreshNowBtn = document.getElementById('refreshNowBtn');
    const refreshStatus = document.getElementById('refreshStatus');
    const lastRefreshTime = document.getElementById('lastRefreshTime');
    const refreshIntervalSelect = document.getElementById('refreshInterval');
    const websocketControls = document.getElementById('websocketControls');
    const connectWsBtn = document.getElementById('connectWsBtn');
    const disconnectWsBtn = document.getElementById('disconnectWsBtn');
    const wsStatus = document.getElementById('wsStatus');
    const wsUrlInput = document.getElementById('wsUrlInput');
    const wsTokenInput = document.getElementById('wsTokenInput');
    const authControls = document.getElementById('authControls');
    const generateTokenBtn = document.getElementById('generateTokenBtn');
    const darkModeToggle = document.getElementById('darkModeToggle');

    let findingsData = null;
    let insightsData = null;
    let allFindings = [];
    let refreshIntervalId = null;
    let currentFiles = null;
    let isRefreshing = false;
    let websocket = null;
    let isWebSocketConnected = false;
    let severityChart = null;
    let categoryChart = null;

    // Initialize dark mode
    function initDarkMode() {
        const savedMode = localStorage.getItem('darkMode');
        if (savedMode === 'enabled') {
            document.documentElement.classList.add('dark-mode');
            if (darkModeToggle) {
                darkModeToggle.querySelector('button').textContent = '🌞 Light Mode';
            }
        }
    }

    // Toggle dark mode
    function toggleDarkMode() {
        document.documentElement.classList.toggle('dark-mode');
        const isDark = document.documentElement.classList.contains('dark-mode');
        localStorage.setItem('darkMode', isDark ? 'enabled' : 'disabled');
        
        if (darkModeToggle) {
            darkModeToggle.querySelector('button').textContent = isDark ? '🌞 Light Mode' : '🌙 Dark Mode';
        }
        
        // Update charts with new theme
        updateChartThemes();
    }

    // Update chart themes based on current mode
    function updateChartThemes() {
        const isDark = document.documentElement.classList.contains('dark-mode');
        const textColor = isDark ? '#e0e0e0' : '#262626';
        const gridColor = isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)';
        
        if (severityChart) {
            severityChart.options.scales.x.ticks.color = textColor;
            severityChart.options.scales.y.ticks.color = textColor;
            severityChart.options.scales.x.grid.color = gridColor;
            severityChart.options.scales.y.grid.color = gridColor;
            severityChart.update();
        }
        
        if (categoryChart) {
            categoryChart.options.scales.x.ticks.color = textColor;
            categoryChart.options.scales.y.ticks.color = textColor;
            categoryChart.options.scales.x.grid.color = gridColor;
            categoryChart.options.scales.y.grid.color = gridColor;
            categoryChart.update();
        }
    }

    // Set up file input trigger
    loadBtn.addEventListener('click', function() {
        fileInput.click();
    });

    // Handle file selection
    fileInput.addEventListener('change', function(e) {
        const files = Array.from(e.target.files);
        
        if (files.length === 0) {
            return;
        }
        
        // Store current files for refresh
        currentFiles = files;
        
        // Show loading state
        loading.style.display = 'block';
        error.style.display = 'none';
        content.style.display = 'none';
        
        // Show refresh controls
        refreshControls.style.display = 'block';
        
        // Stop any existing refresh
        stopAutoRefresh();
        
        // Process files
        processFiles(files);
    });

    // Handle severity filter changes
    severityFilter.addEventListener('change', function() {
        filterFindings();
    });
    
    // Set up refresh controls
    startRefreshBtn.addEventListener('click', startAutoRefresh);
    stopRefreshBtn.addEventListener('click', stopAutoRefresh);
    refreshNowBtn.addEventListener('click', refreshNow);
    
    // Set up WebSocket controls
    if (websocketControls) {
        connectWsBtn.addEventListener('click', connectWebSocket);
        disconnectWsBtn.addEventListener('click', disconnectWebSocket);
        
        // Set up auth controls if available
        if (authControls && generateTokenBtn) {
            generateTokenBtn.addEventListener('click', generateToken);
            authControls.style.display = 'block';
        }
        
        websocketControls.style.display = 'block';
    }
    
    // Set up dark mode toggle
    if (darkModeToggle) {
        darkModeToggle.querySelector('button').addEventListener('click', toggleDarkMode);
        initDarkMode();
    }

    function processFiles(files) {
        const promises = [];
        
        files.forEach(file => {
            const promise = new Promise((resolve) => {
                const reader = new FileReader();
                reader.onload = function(e) {
                    try {
                        const data = JSON.parse(e.target.result);
                        
                        if (file.name.includes('findings')) {
                            findingsData = data;
                        } else if (file.name.includes('insights')) {
                            insightsData = data;
                        }
                        
                        resolve();
                    } catch (err) {
                        resolve(); // Continue even if one file fails
                    }
                };
                reader.onerror = function() {
                    resolve();
                };
                reader.readAsText(file);
            });
            
            promises.push(promise);
        });
        
        Promise.all(promises).then(() => {
            if (!findingsData) {
                showError('No valid findings.json file found. Please select a findings.json file.');
                return;
            }
            
            // Store all findings for filtering
            allFindings = findingsData.findings || [];
            
            // Update last refresh time
            updateLastRefreshTime();
            
            // Render the data
            renderData();
        });
    }

    function updateLastRefreshTime() {
        const now = new Date();
        lastRefreshTime.textContent = `Last refreshed: ${now.toLocaleTimeString()} (${now.toLocaleDateString()})`;
    }

    function startAutoRefresh() {
        if (refreshIntervalId) {
            stopAutoRefresh();
        }
        
        const intervalSeconds = parseInt(refreshIntervalSelect.value);
        
        refreshIntervalId = setInterval(() => {
            if (!isRefreshing && currentFiles && currentFiles.length > 0) {
                isRefreshing = true;
                refreshStatus.textContent = `Auto-refresh: Active (every ${intervalSeconds} seconds)`;
                refreshStatus.className = 'refresh-active';
                
                // Show loading state briefly
                const tempLoading = document.createElement('div');
                tempLoading.className = 'refresh-loading';
                tempLoading.textContent = 'Refreshing data...';
                content.parentNode.insertBefore(tempLoading, content);
                
                // Re-process files
                processFiles(currentFiles);
                
                // Remove temp loading after a short delay
                setTimeout(() => {
                    if (tempLoading.parentNode) {
                        tempLoading.parentNode.removeChild(tempLoading);
                    }
                    isRefreshing = false;
                }, 500);
            }
        }, intervalSeconds * 1000);
        
        // Update UI
        startRefreshBtn.style.display = 'none';
        stopRefreshBtn.style.display = 'inline-block';
        refreshStatus.textContent = `Auto-refresh: Active (every ${intervalSeconds} seconds)`;
        refreshStatus.className = 'refresh-active';
        updateLastRefreshTime();
    }

    function stopAutoRefresh() {
        if (refreshIntervalId) {
            clearInterval(refreshIntervalId);
            refreshIntervalId = null;
        }
        
        // Update UI
        startRefreshBtn.style.display = 'inline-block';
        stopRefreshBtn.style.display = 'none';
        refreshStatus.textContent = 'Auto-refresh: Off';
        refreshStatus.className = '';
    }

    function refreshNow() {
        if (currentFiles && currentFiles.length > 0) {
            isRefreshing = true;
            
            // Show loading state briefly
            const tempLoading = document.createElement('div');
            tempLoading.className = 'refresh-loading';
            tempLoading.textContent = 'Refreshing data...';
            
            content.parentNode.insertBefore(tempLoading, content);
            
            // Re-process files
            processFiles(currentFiles);
            
            // Remove temp loading after a short delay
            setTimeout(() => {
                if (tempLoading.parentNode) {
                    tempLoading.parentNode.removeChild(tempLoading);
                }
                isRefreshing = false;
            }, 500);
        }
    }

    // Token generation function
    function generateToken() {
        const username = prompt('Enter username:', 'demo-user');
        const password = prompt('Enter password:', 'demo-password');
        const role = prompt('Enter role (viewer/admin):', 'viewer');
        
        if (!username || !password) {
            showError('Username and password are required');
            return;
        }
        
        // Extract base URL from WebSocket URL
        let wsUrl = wsUrlInput.value.trim() || 'ws://localhost:8080/ws';
        let httpUrl = wsUrl.replace('ws://', 'http://').replace('/ws', '');
        
        // Generate token via API
        fetch(httpUrl + '/token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: username,
                password: password,
                role: role
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Token generation failed: ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            if (wsTokenInput) {
                wsTokenInput.value = data.token;
            }
            console.log('Token generated successfully:', data.token);
            console.log('Expires in:', data.expires_in, 'seconds');
            console.log('Username:', data.username);
            console.log('Role:', data.role);
            
            // Show success notification
            showNotification('Token generated successfully!', 'success');
        })
        .catch(error => {
            console.error('Error generating token:', error);
            showError('Failed to generate token: ' + error.message);
        });
    }

    // WebSocket connection functions
    function connectWebSocket() {
        const wsUrl = wsUrlInput.value.trim() || 'ws://localhost:8080/ws';
        const token = wsTokenInput ? wsTokenInput.value.trim() : '';
        
        if (isWebSocketConnected) {
            showError('Already connected to WebSocket');
            return;
        }

        try {
            // Add token to URL if provided
            let finalUrl = wsUrl;
            if (token) {
                const separator = wsUrl.includes('?') ? '&' : '?';
                finalUrl = wsUrl + separator + 'token=' + encodeURIComponent(token);
            }

            websocket = new WebSocket(finalUrl);
            
            websocket.onopen = function() {
                isWebSocketConnected = true;
                wsStatus.textContent = 'WebSocket: Connected';
                wsStatus.className = 'ws-connected';
                connectWsBtn.style.display = 'none';
                disconnectWsBtn.style.display = 'inline-block';
                
                // Show refresh controls for WebSocket mode
                refreshControls.style.display = 'block';
                
                // Show WebSocket stats section
                const websocketStatsSection = document.getElementById('websocketStatsSection');
                if (websocketStatsSection) {
                    websocketStatsSection.style.display = 'block';
                }
                
                console.log('WebSocket connected to', finalUrl);
                showNotification('WebSocket connected successfully!', 'success');
            };

            // Set up message handler
            updateWebSocketMessageHandler();

            websocket.onclose = function() {
                isWebSocketConnected = false;
                wsStatus.textContent = 'WebSocket: Disconnected';
                wsStatus.className = 'ws-disconnected';
                connectWsBtn.style.display = 'inline-block';
                disconnectWsBtn.style.display = 'none';
                
                // Hide WebSocket stats section
                const websocketStatsSection = document.getElementById('websocketStatsSection');
                if (websocketStatsSection) {
                    websocketStatsSection.style.display = 'none';
                }
                
                console.log('WebSocket disconnected');
                showNotification('WebSocket disconnected', 'info');
            };

            websocket.onerror = function(error) {
                console.error('WebSocket error:', error);
                showError('WebSocket connection error: ' + error.message);
            };

        } catch (err) {
            showError('Failed to connect to WebSocket: ' + err.message);
        }
    }

    function disconnectWebSocket() {
        if (websocket && isWebSocketConnected) {
            websocket.close();
            websocket = null;
            isWebSocketConnected = false;
        }
    }

    function updateWebSocketStats(stats) {
        if (stats) {
            // Update quick stats from WebSocket data
            document.getElementById('totalFindings').textContent = stats.total_findings || '0';
            document.getElementById('criticalFindings').textContent = stats.critical_count || '0';
            document.getElementById('highFindings').textContent = stats.high_count || '0';
            document.getElementById('avgScore').textContent = (stats.performance_score || 0).toFixed(1);
            
            // Update WebSocket-specific stats
            document.getElementById('criticalFindingsWs').textContent = stats.critical_count || '0';
            document.getElementById('highFindingsWs').textContent = stats.high_count || '0';
            document.getElementById('mediumFindingsWs').textContent = stats.medium_count || '0';
            document.getElementById('lowFindingsWs').textContent = stats.low_count || '0';
            document.getElementById('performanceScoreWs').textContent = stats.performance_score || '0';
            
            // Update client count
            const clientCountEl = document.getElementById('clientCount');
            if (clientCountEl) {
                clientCountEl.textContent = stats.connected_clients || '0';
            }
            
            // Update last refresh time
            if (stats.last_updated) {
                lastRefreshTime.textContent = 'Last updated: ' + stats.last_updated;
            }
        }
    }
    
    // Update WebSocket message handler to process performance history
    function updateWebSocketMessageHandler() {
        websocket.onmessage = function(event) {
            try {
                const data = JSON.parse(event.data);
                
                if (data.type === 'data_update') {
                    // Update data from WebSocket
                    findingsData = data.findings;
                    insightsData = data.insights;
                    allFindings = findingsData.findings || [];
                    
                    // Update last refresh time
                    updateLastRefreshTime();
                    
                    // Render the data
                    renderData();
                    
                    // Update WebSocket stats
                    updateWebSocketStats(data.stats);
                    
                    // Process performance history if available
                    if (data.history && data.history.length > 0) {
                        updatePerformanceHistory(data.history);
                    }
                    
                    // Show notification for live updates
                    showLiveUpdateNotification();
                }
            } catch (err) {
                console.error('Error processing WebSocket message:', err);
            }
        };
    }
    
    // Update performance history visualization
    function updatePerformanceHistory(history) {
        if (!history || history.length === 0) {
            return;
        }
        
        // Update performance metrics dashboard
        updatePerformanceMetricsDashboard(history);
        
        // Update trends charts
        updateTrendsCharts(history);
        
        // Update score breakdown
        updateScoreBreakdown(history);
    }
    
    // Update performance metrics dashboard
    function updatePerformanceMetricsDashboard(history) {
        const latest = history[history.length - 1];
        
        // Update metric cards
        document.getElementById('performanceScore').textContent = latest.overall_score || 'N/A';
        document.getElementById('criticalIssues').textContent = latest.critical_count || '0';
        document.getElementById('highIssues').textContent = latest.high_count || '0';
        document.getElementById('resolvedIssues').textContent = '0'; // Placeholder for resolved issues
        
        // Add visual indicators
        const scoreElement = document.getElementById('performanceScore');
        const score = parseFloat(latest.overall_score) || 0;
        
        // Clear previous classes
        scoreElement.className = 'metric-value';
        
        if (score >= 80) {
            scoreElement.classList.add('score-excellent');
        } else if (score >= 60) {
            scoreElement.classList.add('score-good');
        } else if (score >= 40) {
            scoreElement.classList.add('score-fair');
        } else {
            scoreElement.classList.add('score-poor');
        }
    }
    
    // Update trends charts
    function updateTrendsCharts(history) {
        if (history.length < 2) {
            return;
        }
        
        // Extract data for charts
        const timestamps = history.map(h => h.timestamp);
        const scores = history.map(h => h.overall_score);
        const criticalCounts = history.map(h => h.critical_count);
        const highCounts = history.map(h => h.high_count);
        
        // Update score trend chart
        updateScoreTrendChart(timestamps, scores);
        
        // Update resolution trend chart (placeholder)
        updateResolutionTrendChart(timestamps, criticalCounts, highCounts);
    }
    
    // Update score trend chart
    function updateScoreTrendChart(timestamps, scores) {
        const ctx = document.getElementById('scoreTrendChart');
        if (!ctx) return;
        
        const isDark = document.documentElement.classList.contains('dark-mode');
        const textColor = isDark ? '#e0e0e0' : '#262626';
        const gridColor = isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)';
        
        // Format timestamps for display
        const labels = timestamps.map(ts => {
            const date = new Date(ts);
            return date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'});
        });
        
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Performance Score',
                    data: scores,
                    borderColor: 'rgba(74, 111, 165, 1)',
                    backgroundColor: 'rgba(74, 111, 165, 0.1)',
                    borderWidth: 2,
                    fill: true,
                    tension: 0.4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        ticks: {
                            color: textColor,
                            font: { size: 10 }
                        },
                        grid: { color: gridColor }
                    },
                    y: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                            color: textColor,
                            font: { size: 10 }
                        },
                        grid: { color: gridColor }
                    }
                },
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        backgroundColor: 'rgba(0, 0, 0, 0.8)',
                        titleColor: '#fff',
                        bodyColor: '#fff'
                    }
                }
            }
        });
    }
    
    // Update resolution trend chart
    function updateResolutionTrendChart(timestamps, criticalCounts, highCounts) {
        const ctx = document.getElementById('resolutionTrendChart');
        if (!ctx) return;
        
        const isDark = document.documentElement.classList.contains('dark-mode');
        const textColor = isDark ? '#e0e0e0' : '#262626';
        const gridColor = isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)';
        
        // Format timestamps for display
        const labels = timestamps.map(ts => {
            const date = new Date(ts);
            return date.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'});
        });
        
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'Critical Issues',
                        data: criticalCounts,
                        borderColor: 'rgba(255, 107, 107, 1)',
                        backgroundColor: 'rgba(255, 107, 107, 0.1)',
                        borderWidth: 2,
                        tension: 0.4
                    },
                    {
                        label: 'High Severity Issues',
                        data: highCounts,
                        borderColor: 'rgba(255, 230, 109, 1)',
                        backgroundColor: 'rgba(255, 230, 109, 0.1)',
                        borderWidth: 2,
                        tension: 0.4
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        ticks: {
                            color: textColor,
                            font: { size: 10 }
                        },
                        grid: { color: gridColor }
                    },
                    y: {
                        beginAtZero: true,
                        ticks: {
                            color: textColor,
                            font: { size: 10 }
                        },
                        grid: { color: gridColor }
                    }
                },
                plugins: {
                    legend: {
                        position: 'top',
                        labels: { color: textColor, font: { size: 10 } }
                    },
                    tooltip: {
                        backgroundColor: 'rgba(0, 0, 0, 0.8)',
                        titleColor: '#fff',
                        bodyColor: '#fff'
                    }
                }
            }
        });
    }
    
    // Update score breakdown visualization
    function updateScoreBreakdown(history) {
        const ctx = document.getElementById('scoreBreakdownChart');
        if (!ctx) return;
        
        const isDark = document.documentElement.classList.contains('dark-mode');
        const textColor = isDark ? '#e0e0e0' : '#262626';
        
        // Calculate average scores by category (simplified for demo)
        const latest = history[history.length - 1];
        const score = latest.overall_score || 0;
        
        // Simulated breakdown - in a real implementation, this would come from detailed analysis
        const cpuEfficiency = Math.min(100, score * 0.4);
        const memoryUsage = Math.min(100, score * 0.3);
        const ioPerformance = Math.min(100, score * 0.2);
        const concurrency = Math.min(100, score * 0.1);
        
        new Chart(ctx, {
            type: 'doughnut',
            data: {
                labels: ['CPU Efficiency', 'Memory Usage', 'I/O Performance', 'Concurrency'],
                datasets: [{
                    data: [cpuEfficiency, memoryUsage, ioPerformance, concurrency],
                    backgroundColor: [
                        'rgba(78, 205, 196, 0.8)',
                        'rgba(255, 230, 109, 0.8)',
                        'rgba(255, 107, 107, 0.8)',
                        'rgba(168, 230, 207, 0.8)'
                    ],
                    borderWidth: 2,
                    borderColor: '#fff'
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'right',
                        labels: { color: textColor, font: { size: 11 } }
                    },
                    tooltip: {
                        backgroundColor: 'rgba(0, 0, 0, 0.8)',
                        titleColor: '#fff',
                        bodyColor: '#fff'
                    }
                }
            }
        });
        
        // Update comparative analysis
        updateComparativeAnalysis(history);
    }
    
    // Update comparative analysis charts
    function updateComparativeAnalysis(history) {
        if (history.length < 2) {
            return;
        }
        
        // Calculate baseline (first snapshot) vs current (last snapshot)
        const baseline = history[0];
        const current = history[history.length - 1];
        
        // Update comparison chart
        updateComparisonChart(baseline, current);
        
        // Update category comparison (simplified for demo)
        updateCategoryComparisonChart(current);
    }
    
    // Update baseline vs current comparison chart
    function updateComparisonChart(baseline, current) {
        const ctx = document.getElementById('comparisonChart');
        if (!ctx) return;
        
        const isDark = document.documentElement.classList.contains('dark-mode');
        const textColor = isDark ? '#e0e0e0' : '#262626';
        const gridColor = isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)';
        
        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: ['Performance Score', 'Critical Issues', 'High Issues', 'Total Findings'],
                datasets: [
                    {
                        label: 'Baseline',
                        data: [baseline.overall_score, baseline.critical_count, baseline.high_count, baseline.total_findings],
                        backgroundColor: 'rgba(168, 230, 207, 0.8)',
                        borderColor: 'rgba(168, 230, 207, 1)',
                        borderWidth: 2
                    },
                    {
                        label: 'Current',
                        data: [current.overall_score, current.critical_count, current.high_count, current.total_findings],
                        backgroundColor: 'rgba(74, 111, 165, 0.8)',
                        borderColor: 'rgba(74, 111, 165, 1)',
                        borderWidth: 2
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        stacked: false,
                        ticks: { color: textColor, font: { size: 10 } },
                        grid: { color: gridColor }
                    },
                    y: {
                        beginAtZero: true,
                        ticks: { color: textColor, font: { size: 10 } },
                        grid: { color: gridColor }
                    }
                },
                plugins: {
                    legend: {
                        position: 'top',
                        labels: { color: textColor, font: { size: 10 } }
                    },
                    tooltip: {
                        backgroundColor: 'rgba(0, 0, 0, 0.8)',
                        titleColor: '#fff',
                        bodyColor: '#fff'
                    }
                }
            }
        });
    }
    
    // Update category comparison chart
    function updateCategoryComparisonChart(current) {
        const ctx = document.getElementById('categoryComparisonChart');
        if (!ctx) return;
        
        const isDark = document.documentElement.classList.contains('dark-mode');
        const textColor = isDark ? '#e0e0e0' : '#262626';
        
        // Simulated category data - in a real implementation, this would come from findings analysis
        const score = current.overall_score || 0;
        const categories = ['CPU', 'Memory', 'I/O', 'Concurrency', 'Other'];
        const values = categories.map((_, i) => Math.max(10, score - i * 10));
        
        new Chart(ctx, {
            type: 'polarArea',
            data: {
                labels: categories,
                datasets: [{
                    data: values,
                    backgroundColor: [
                        'rgba(78, 205, 196, 0.6)',
                        'rgba(255, 230, 109, 0.6)',
                        'rgba(255, 107, 107, 0.6)',
                        'rgba(168, 230, 207, 0.6)',
                        'rgba(144, 164, 174, 0.6)'
                    ],
                    borderWidth: 2,
                    borderColor: '#fff'
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    r: {
                        angleLines: { color: gridColor },
                        ticks: { color: textColor, font: { size: 10 } },
                        pointLabels: { color: textColor, font: { size: 11 } }
                    }
                },
                plugins: {
                    legend: {
                        position: 'right',
                        labels: { color: textColor, font: { size: 10 } }
                    },
                    tooltip: {
                        backgroundColor: 'rgba(0, 0, 0, 0.8)',
                        titleColor: '#fff',
                        bodyColor: '#fff'
                    }
                }
            }
        });
    }
    
    function showLiveUpdateNotification() {
        const notification = document.createElement('div');
        notification.className = 'notification success';
        notification.innerHTML = '<i class="fas fa-sync-alt"></i> Live update received';
        
        document.body.appendChild(notification);
        
        // Remove notification after 3 seconds
        setTimeout(function() {
            notification.style.opacity = '0';
            setTimeout(() => notification.remove(), 300);
        }, 3000);
    }

    function showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        
        document.body.appendChild(notification);
        
        // Remove notification after 3 seconds
        setTimeout(function() {
            notification.style.opacity = '0';
            setTimeout(() => notification.remove(), 300);
        }, 3000);
    }

    function showError(message) {
        loading.style.display = 'none';
        error.textContent = message;
        error.style.display = 'block';
        showNotification(message, 'error');
    }

    function renderData() {
        if (!findingsData) {
            showError('No findings data available');
            return;
        }

        // Show content
        loading.style.display = 'none';
        error.style.display = 'none';
        content.style.display = 'block';

        // Render summary
        renderSummary();
        
        // Render AI summary if insights available
        if (insightsData) {
            renderAISummary();
        }
        
        // Render quick stats
        renderQuickStats();
        
        // Render charts
        renderCharts();
        
        // Render findings list
        renderFindingsList();
        
        // Render plugin info
        renderPluginInfo();
        
        // Add fade-in animation
        content.classList.add('fade-in');
    }

    function renderSummary() {
        const summary = findingsData.summary || {};
        
        // Update summary values
        document.getElementById('overallScore').textContent = summary.overall_score || 'N/A';
        document.getElementById('severity').textContent = summary.severity || 'Unknown';
        document.getElementById('confidence').textContent = summary.confidence || 'N/A';
        document.getElementById('overview').textContent = summary.overview || 'No overview available';
        
        // Add severity indicator
        const severityElement = document.getElementById('severity');
        const severity = (summary.severity || '').toLowerCase();
        
        severityElement.className = 'severity-indicator ' + severity;
        
        // Add score indicator
        const scoreElement = document.getElementById('overallScore');
        const score = parseFloat(summary.overall_score) || 0;
        
        if (score >= 80) {
            scoreElement.classList.add('score-excellent');
        } else if (score >= 60) {
            scoreElement.classList.add('score-good');
        } else if (score >= 40) {
            scoreElement.classList.add('score-fair');
        } else {
            scoreElement.classList.add('score-poor');
        }
    }

    function renderAISummary() {
        const aiSummarySection = document.getElementById('aiSummarySection');
        
        if (!aiSummarySection) return;
        
        aiSummarySection.style.display = 'block';
        
        // Add AI insights to the section
        const insights = insightsData.insights || [];
        
        if (insights.length > 0) {
            const aiContent = document.createElement('div');
            aiContent.className = 'ai-insights-content';
            
            insights.forEach((insight, index) => {
                const insightElement = document.createElement('div');
                insightElement.className = 'ai-insight';
                
                insightElement.innerHTML = `
                    <div class="ai-insight-header">
                        <h3><i class="fas fa-lightbulb"></i> ${insight.title || 'Insight ' + (index + 1)}</h3>
                        <span class="ai-confidence">Confidence: ${insight.confidence || 'N/A'}</span>
                    </div>
                    <div class="ai-insight-content">
                        <p>${insight.description || 'No description available'}</p>
                        ${insight.recommendations ? `<div class="ai-recommendations"><strong>Recommendations:</strong> ${insight.recommendations}</div>` : ''}
                    </div>
                `;
                
                aiContent.appendChild(insightElement);
            });
            
            aiSummarySection.appendChild(aiContent);
        }
    }

    function renderQuickStats() {
        // Calculate statistics
        const totalFindings = allFindings.length;
        const criticalCount = allFindings.filter(f => f.severity?.toLowerCase() === 'critical').length;
        const highCount = allFindings.filter(f => f.severity?.toLowerCase() === 'high').length;
        const mediumCount = allFindings.filter(f => f.severity?.toLowerCase() === 'medium').length;
        const lowCount = allFindings.filter(f => f.severity?.toLowerCase() === 'low').length;
        
        // Calculate average score
        const scores = allFindings.map(f => f.score || 0).filter(score => score > 0);
        const avgScore = scores.length > 0 ? (scores.reduce((a, b) => a + b, 0) / scores.length).toFixed(1) : 0;
        
        // Update DOM
        document.getElementById('totalFindings').textContent = totalFindings;
        document.getElementById('criticalFindings').textContent = criticalCount;
        document.getElementById('highFindings').textContent = highCount;
        document.getElementById('avgScore').textContent = avgScore;
        
        // Add visual indicators for severity levels
        const criticalElement = document.getElementById('criticalFindings');
        const highElement = document.getElementById('highFindings');
        
        if (criticalCount > 0) {
            criticalElement.classList.add('critical-highlight');
        }
        if (highCount > 0) {
            highElement.classList.add('high-highlight');
        }
        
        // Add stat cards hover effects
        const statCards = document.querySelectorAll('.stat-card');
        statCards.forEach(card => {
            card.addEventListener('mouseenter', function() {
                this.style.transform = 'translateY(-3px) scale(1.02)';
                this.style.boxShadow = '0 6px 12px rgba(0, 0, 0, 0.15)';
            });
            
            card.addEventListener('mouseleave', function() {
                this.style.transform = '';
                this.style.boxShadow = '';
            });
        });
    }

    function renderCharts() {
        const isDark = document.documentElement.classList.contains('dark-mode');
        const textColor = isDark ? '#e0e0e0' : '#262626';
        const gridColor = isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)';
        
        // Severity distribution chart
        const severityCtx = document.getElementById('severityChart');
        if (severityCtx) {
            const criticalCount = allFindings.filter(f => f.severity?.toLowerCase() === 'critical').length;
            const highCount = allFindings.filter(f => f.severity?.toLowerCase() === 'high').length;
            const mediumCount = allFindings.filter(f => f.severity?.toLowerCase() === 'medium').length;
            const lowCount = allFindings.filter(f => f.severity?.toLowerCase() === 'low').length;
            
            if (severityChart) {
                severityChart.destroy();
            }
            
            severityChart = new Chart(severityCtx, {
                type: 'doughnut',
                data: {
                    labels: ['Critical', 'High', 'Medium', 'Low'],
                    datasets: [{
                        data: [criticalCount, highCount, mediumCount, lowCount],
                        backgroundColor: [
                            'rgba(255, 107, 107, 0.8)',
                            'rgba(255, 230, 109, 0.8)',
                            'rgba(78, 205, 196, 0.8)',
                            'rgba(168, 230, 207, 0.8)'
                        ],
                        borderWidth: 2,
                        borderColor: '#fff'
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            position: 'bottom',
                            labels: {
                                color: textColor,
                                font: {
                                    size: 12
                                }
                            }
                        },
                        tooltip: {
                            backgroundColor: 'rgba(0, 0, 0, 0.8)',
                            titleColor: '#fff',
                            bodyColor: '#fff'
                        }
                    }
                }
            });
        }
        
        // Category distribution chart
        const categoryCtx = document.getElementById('categoryChart');
        if (categoryCtx) {
            const categoryCounts = {};
            allFindings.forEach(finding => {
                const category = finding.category || 'Other';
                categoryCounts[category] = (categoryCounts[category] || 0) + 1;
            });
            
            const categories = Object.keys(categoryCounts);
            const counts = Object.values(categoryCounts);
            
            if (categoryChart) {
                categoryChart.destroy();
            }
            
            categoryChart = new Chart(categoryCtx, {
                type: 'bar',
                data: {
                    labels: categories,
                    datasets: [{
                        label: 'Findings by Category',
                        data: counts,
                        backgroundColor: 'rgba(74, 111, 165, 0.8)',
                        borderColor: 'rgba(74, 111, 165, 1)',
                        borderWidth: 2
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        x: {
                            ticks: {
                                color: textColor,
                                font: {
                                    size: 11
                                }
                            },
                            grid: {
                                color: gridColor
                            }
                        },
                        y: {
                            beginAtZero: true,
                            ticks: {
                                color: textColor,
                                font: {
                                    size: 11
                                }
                            },
                            grid: {
                                color: gridColor
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            display: false
                        },
                        tooltip: {
                            backgroundColor: 'rgba(0, 0, 0, 0.8)',
                            titleColor: '#fff',
                            bodyColor: '#fff'
                        }
                    }
                }
            });
        }
    }

    function renderFindingsList() {
        const findingsList = document.getElementById('findingsList');
        
        if (!findingsList) return;
        
        findingsList.innerHTML = '';
        
        // Filter findings based on severity filter
        const filteredFindings = filterFindings();
        
        if (filteredFindings.length === 0) {
            findingsList.innerHTML = '<p class="no-findings">No findings match the current filter criteria.</p>';
            return;
        }
        
        // Create findings elements
        filteredFindings.forEach((finding, index) => {
            const findingElement = document.createElement('div');
            findingElement.className = 'finding-item';
            
            const severityClass = `severity-${finding.severity?.toLowerCase() || 'low'}`;
            
            findingElement.innerHTML = `
                <div class="finding-header">
                    <div class="finding-title">${finding.title || 'Untitled Finding'}</div>
                    <div class="finding-severity ${severityClass}">${finding.severity || 'Unknown'}</div>
                </div>
                <div class="finding-description">${finding.description || 'No description available'}</div>
                <div class="finding-meta">
                    ${finding.score ? `<div class="finding-meta-item"><i class="fas fa-chart-line"></i> Score: ${finding.score}</div>` : ''}
                    ${finding.category ? `<div class="finding-meta-item"><i class="fas fa-tag"></i> ${finding.category}</div>` : ''}
                    ${finding.location ? `<div class="finding-meta-item"><i class="fas fa-map-marker-alt"></i> ${finding.location}</div>` : ''}
                </div>
            `;
            
            findingsList.appendChild(findingElement);
            
            // Add click event to show more details
            findingElement.addEventListener('click', function() {
                showFindingDetails(finding);
            });
        });
    }

    function filterFindings() {
        const filterValue = severityFilter.value;
        
        if (filterValue === 'all') {
            return allFindings;
        }
        
        return allFindings.filter(finding => 
            finding.severity?.toLowerCase() === filterValue
        );
    }

    function showFindingDetails(finding) {
        const detailsModal = document.createElement('div');
        detailsModal.className = 'finding-details-modal';
        
        const severityClass = `severity-${finding.severity?.toLowerCase() || 'low'}`;
        
        detailsModal.innerHTML = `
            <div class="modal-content">
                <div class="modal-header">
                    <h3>${finding.title || 'Finding Details'}</h3>
                    <button class="close-modal"><i class="fas fa-times"></i></button>
                </div>
                <div class="modal-body">
                    <div class="finding-detail-section">
                        <div class="detail-label">Severity:</div>
                        <div class="detail-value"><span class="finding-severity ${severityClass}">${finding.severity || 'Unknown'}</span></div>
                    </div>
                    <div class="finding-detail-section">
                        <div class="detail-label">Description:</div>
                        <div class="detail-value">${finding.description || 'No description available'}</div>
                    </div>
                    ${finding.score ? `
                    <div class="finding-detail-section">
                        <div class="detail-label">Score:</div>
                        <div class="detail-value">${finding.score}</div>
                    </div>
                    ` : ''}
                    ${finding.category ? `
                    <div class="finding-detail-section">
                        <div class="detail-label">Category:</div>
                        <div class="detail-value">${finding.category}</div>
                    </div>
                    ` : ''}
                    ${finding.location ? `
                    <div class="finding-detail-section">
                        <div class="detail-label">Location:</div>
                        <div class="detail-value">${finding.location}</div>
                    </div>
                    ` : ''}
                    ${finding.evidence ? `
                    <div class="finding-detail-section">
                        <div class="detail-label">Evidence:</div>
                        <div class="detail-value evidence">${finding.evidence}</div>
                    </div>
                    ` : ''}
                </div>
            </div>
        `;
        
        document.body.appendChild(detailsModal);
        
        // Close modal when clicking X or outside
        detailsModal.querySelector('.close-modal').addEventListener('click', function() {
            detailsModal.remove();
        });
        
        detailsModal.addEventListener('click', function(e) {
            if (e.target === detailsModal) {
                detailsModal.remove();
            }
        });
    }

    function renderPluginInfo() {
        const pluginInfoSection = document.getElementById('pluginInfoSection');
        
        if (!pluginInfoSection) return;
        
        // Create plugin info content
        const pluginInfoHTML = `
            <div class="plugin-info-content">
                <div class="plugin-info-item">
                    <div class="plugin-icon"><i class="fas fa-cogs"></i></div>
                    <div class="plugin-info-text">
                        <h3>Modular Plugin Architecture</h3>
                        <p>TriageProf uses a powerful plugin system where each profiler operates as an independent module. Plugins are discovered automatically via JSON manifests and communicate using JSON-RPC 2.0 protocol.</p>
                        <p><strong>Key Benefits:</strong> Easy extensibility, language-agnostic support, and stable API compatibility.</p>
                    </div>
                </div>
                <div class="plugin-info-item">
                    <div class="plugin-icon"><i class="fas fa-folder-open"></i></div>
                    <div class="plugin-info-text">
                        <h3>Plugin Locations</h3>
                        <p>Plugins are organized in the <code>plugins/</code> directory:</p>
                        <ul>
                            <li><code>plugins/manifests/</code> - JSON manifest files for plugin discovery</li>
                            <li><code>plugins/src/</code> - Source code for each plugin</li>
                            <li><code>plugins/bin/</code> - Compiled plugin binaries</li>
                        </ul>
                        <p>Each plugin has its own manifest defining capabilities, SDK version, and metadata.</p>
                    </div>
                </div>
                <div class="plugin-info-item">
                    <div class="plugin-icon"><i class="fas fa-plus-circle"></i></div>
                    <div class="plugin-info-text">
                        <h3>Extensible & Maintainable</h3>
                        <p>Add new profilers without modifying core code. The plugin API remains stable for backward compatibility.</p>
                        <p><strong>Current Plugins:</strong> Go pprof, Node.js inspector, Python cProfile, Ruby stackprof</p>
                        <p>New plugins can be added by implementing the JSON-RPC interface and providing a manifest.</p>
                    </div>
                </div>
            </div>
        `;
        
        pluginInfoSection.innerHTML = `<h2><i class="fas fa-puzzle-piece"></i> Plugin Information</h2>` + pluginInfoHTML;
        pluginInfoSection.style.display = 'block';
        
        // Add hover effects to plugin info items
        const pluginItems = document.querySelectorAll('.plugin-info-item');
        pluginItems.forEach(item => {
            item.addEventListener('mouseenter', function() {
                this.style.transform = 'translateY(-2px)';
                this.style.boxShadow = '0 4px 8px rgba(0, 0, 0, 0.1)';
            });
            
            item.addEventListener('mouseleave', function() {
                this.style.transform = '';
                this.style.boxShadow = '';
            });
        });
    }
});