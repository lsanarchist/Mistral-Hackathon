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

    let findingsData = null;
    let insightsData = null;
    let allFindings = [];
    let refreshIntervalId = null;
    let currentFiles = null;
    let isRefreshing = false;
    let websocket = null;
    let isWebSocketConnected = false;

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
            };

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
                        
                        // Show notification for live updates
                        showLiveUpdateNotification();
                    }
                } catch (err) {
                    console.error('Error processing WebSocket message:', err);
                }
            };

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
    
    function showLiveUpdateNotification() {
        const notification = document.createElement('div');
        notification.className = 'live-update-notification';
        notification.textContent = '🔄 Live update received';
        
        document.body.appendChild(notification);
        
        // Remove notification after 3 seconds
        setTimeout(function() {
            notification.remove();
        }, 3000);
    }
    }

    function showError(message) {
        loading.style.display = 'none';
        error.textContent = message;
        error.style.display = 'block';
    }

    function renderQuickStats() {
        // Calculate statistics
        const totalFindings = allFindings.length;
        const criticalCount = allFindings.filter(f => f.severity?.toLowerCase() === 'critical').length;
        const highCount = allFindings.filter(f => f.severity?.toLowerCase() === 'high').length;
        const mediumCount = allFindings.filter(f => f.severity?.toLowerCase() === 'medium').length;
        
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
    
    function renderPluginInfo() {
        const pluginInfoSection = document.getElementById('pluginInfoSection');
        
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
                this.style.transform = 'translateX(5px)';
                this.style.transition = 'transform 0.3s ease';
            });
            
            item.addEventListener('mouseleave', function() {
                this.style.transform = '';
            });
        });
    }
    
    function renderData() {
        try {
            // Set up basic info
            document.getElementById('overallScore').textContent = findingsData.summary?.OverallScore || 'N/A';
            
            // Determine severity
            const score = findingsData.summary?.OverallScore || 0;
            let severity = 'Unknown';
            if (score >= 80) severity = 'Critical';
            else if (score >= 60) severity = 'High';
            else if (score >= 40) severity = 'Medium';
            else if (score >= 20) severity = 'Low';
            else severity = 'Info';
            
            document.getElementById('severity').textContent = severity;
            
            // Set confidence if insights available
            if (insightsData && insightsData.ExecutiveSummary) {
                const confidence = insightsData.ExecutiveSummary.Confidence || 0;
                document.getElementById('confidence').textContent = confidence + '%';
                
                // Add confidence indicator
                const confidenceIndicator = document.createElement('div');
                confidenceIndicator.className = 'confidence-indicator';
                
                let confidenceClass = 'confidence-low';
                let confidenceText = 'Low Confidence';
                if (confidence >= 80) {
                    confidenceClass = 'confidence-high';
                    confidenceText = 'High Confidence';
                } else if (confidence >= 50) {
                    confidenceClass = 'confidence-medium';
                    confidenceText = 'Medium Confidence';
                }
                
                confidenceIndicator.className = 'confidence-indicator ' + confidenceClass;
                confidenceIndicator.textContent = confidenceText;
                
                const confidenceElement = document.getElementById('confidence');
                confidenceElement.parentNode.insertBefore(confidenceIndicator, confidenceElement.nextSibling);
                
                // Enhanced overview with AI branding
                let overviewHTML = `<div class="ai-overview">`;
                overviewHTML += `<div class="ai-header"><i class="fas fa-robot"></i> AI Analysis Overview</div>`;
                overviewHTML += `<p class="ai-overview-text">${insightsData.ExecutiveSummary.Overview || 'No overview available'}</p>`;
                overviewHTML += `<div class="ai-metrics">`;
                overviewHTML += `<span class="ai-metric"><strong>Severity:</strong> ${insightsData.ExecutiveSummary.OverallSeverity || 'Unknown'}</span>`;
                overviewHTML += `<span class="ai-metric"><strong>Confidence:</strong> ${confidence}%</span>`;
                
                if (insightsData.ExecutiveSummary.KeyThemes && insightsData.ExecutiveSummary.KeyThemes.length > 0) {
                    overviewHTML += `<span class="ai-metric"><strong>Themes:</strong> ${insightsData.ExecutiveSummary.KeyThemes.join(', ')}</span>`;
                }
                overviewHTML += `</div></div>`;
                
                document.getElementById('overview').innerHTML = overviewHTML;
                
            } else {
                document.getElementById('confidence').textContent = 'N/A';
                document.getElementById('overview').innerHTML = '<p>No LLM insights available</p>';
            }
            
            // Render quick stats
            renderQuickStats();
            
            // Render charts
            renderSeverityChart();
            renderCategoryChart();
            
            // Show AI summary section if insights are available
            if (insightsData) {
                document.getElementById('aiSummarySection').style.display = 'block';
                renderPluginInfo();
            }
            
            // Render top risks if available
            if (insightsData && insightsData.TopRisks && insightsData.TopRisks.length > 0) {
                renderTopRisks();
            }
            
            // Render top actions if available
            if (insightsData && insightsData.TopActions && insightsData.TopActions.length > 0) {
                renderTopActions();
            }
            
            // Render findings
            renderFindings();
            
            // Render metrics dashboard
            renderMetricsDashboard();
            
            // Show content
            loading.style.display = 'none';
            content.style.display = 'block';
            content.classList.add('content-visible');
            
        } catch (err) {
            showError('Error rendering data: ' + err.message);
            console.error('Render error:', err);
        }
    }

    function renderSeverityChart() {
        const ctx = document.getElementById('severityChart').getContext('2d');
        
        // Count findings by severity
        const severityCounts = {
            critical: 0,
            high: 0,
            medium: 0,
            low: 0,
            info: 0
        };
        
        allFindings.forEach(finding => {
            const severity = finding.severity?.toLowerCase();
            if (severityCounts[severity] !== undefined) {
                severityCounts[severity]++;
            }
        });
        
        // Prepare chart data
        const labels = Object.keys(severityCounts);
        const data = Object.values(severityCounts);
        
        const colors = {
            critical: '#ff6b6b',
            high: '#ff8e53',
            medium: '#ffd166',
            low: '#06d6a0',
            info: '#118ab2'
        };
        
        new Chart(ctx, {
            type: 'doughnut',
            data: {
                labels: labels.map(label => label.charAt(0).toUpperCase() + label.slice(1)),
                datasets: [{
                    data: data,
                    backgroundColor: labels.map(label => colors[label]),
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: {
                            padding: 15,
                            font: {
                                size: 12
                            }
                        }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                const label = context.label || '';
                                const value = context.parsed || 0;
                                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                const percentage = Math.round((value / total) * 100);
                                return `${label}: ${value} findings (${percentage}%)`;
                            }
                        }
                    }
                }
            }
        });
    }

    function renderCategoryChart() {
        const ctx = document.getElementById('categoryChart').getContext('2d');
        
        // Count findings by category
        const categoryCounts = {};
        
        allFindings.forEach(finding => {
            const category = finding.category || 'unknown';
            categoryCounts[category] = (categoryCounts[category] || 0) + 1;
        });
        
        // Prepare chart data
        const labels = Object.keys(categoryCounts);
        const data = Object.values(categoryCounts);
        
        // Generate colors
        const generateColors = (count) => {
            const colors = [];
            for (let i = 0; i < count; i++) {
                colors.push(`hsl(${Math.random() * 360}, 70%, 60%)`);
            }
            return colors;
        };
        
        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Findings by Category',
                    data: data,
                    backgroundColor: generateColors(labels.length),
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return `${context.parsed.y} findings in ${context.label}`;
                            }
                        }
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            stepSize: 1
                        }
                    }
                }
            }
        });
    }

    function filterFindings() {
        const filterValue = severityFilter.value;
        const findingsList = document.getElementById('findingsList');
        
        // Clear current findings
        findingsList.innerHTML = '';
        
        // Filter findings
        const filteredFindings = filterValue === 'all' 
            ? allFindings 
            : allFindings.filter(finding => finding.severity?.toLowerCase() === filterValue);
        
        // Render filtered findings
        filteredFindings.forEach(finding => {
            renderFindingCard(finding);
        });
        
        // Update summary
        const filteredCount = filteredFindings.length;
        const totalCount = allFindings.length;
        const summaryText = document.createElement('div');
        summaryText.className = 'filter-summary';
        summaryText.textContent = `Showing ${filteredCount} of ${totalCount} findings`;
        summaryText.style.textAlign = 'right';
        summaryText.style.marginBottom = '10px';
        summaryText.style.color = '#666';
        
        findingsList.prepend(summaryText);
    }

    function renderTopRisks() {
        const risksSection = document.getElementById('topRisksSection');
        const riskCards = document.getElementById('riskCards');
        
        // Show top 3 risks
        const risksToShow = insightsData.TopRisks.slice(0, 3);
        
        risksToShow.forEach((risk, index) => {
            const card = document.createElement('div');
            card.className = 'risk-card';
            
            // Determine risk level class
            let riskLevelClass = 'risk-low';
            if (risk.Severity.toLowerCase().includes('critical') || risk.Severity.toLowerCase().includes('high')) {
                riskLevelClass = 'risk-high';
            } else if (risk.Severity.toLowerCase().includes('medium')) {
                riskLevelClass = 'risk-medium';
            }
            
            card.innerHTML = `
                <div class="risk-header">
                    <div class="risk-number">${index + 1}</div>
                    <div class="risk-title">${risk.Description}</div>
                    <div class="risk-badge ${riskLevelClass}">${risk.Severity}</div>
                </div>
                <div class="risk-details">
                    <div class="risk-detail-item">
                        <i class="fas fa-bullseye"></i>
                        <span><strong>Impact:</strong> ${risk.Impact}</span>
                    </div>
                    <div class="risk-detail-item">
                        <i class="fas fa-chart-line"></i>
                        <span><strong>Likelihood:</strong> ${risk.Likelihood}</span>
                    </div>
                </div>
            `;
            
            riskCards.appendChild(card);
        });

        risksSection.style.display = 'block';
    }

    function renderTopActions() {
        const actionsSection = document.getElementById('topActionsSection');
        const actionCards = document.getElementById('actionCards');
        
        // Show top 3 actions
        const actionsToShow = insightsData.TopActions.slice(0, 3);
        
        actionsToShow.forEach((action, index) => {
            const card = document.createElement('div');
            card.className = 'action-card';
            
            // Determine priority level class
            let priorityLevelClass = 'priority-low';
            if (action.Priority.toLowerCase().includes('high') || action.Priority.toLowerCase().includes('critical')) {
                priorityLevelClass = 'priority-high';
            } else if (action.Priority.toLowerCase().includes('medium')) {
                priorityLevelClass = 'priority-medium';
            }
            
            card.innerHTML = `
                <div class="action-header">
                    <div class="action-number">${index + 1}</div>
                    <div class="action-title">${action.Description}</div>
                    <div class="action-badge ${priorityLevelClass}">${action.Priority}</div>
                </div>
                <div class="action-details">
                    <div class="action-detail-item">
                        <i class="fas fa-clock"></i>
                        <span><strong>Estimated Effort:</strong> ${action.EstimatedEffort}</span>
                    </div>
                    ${action.Categories && action.Categories.length > 0 ? 
                        `<div class="action-detail-item">
                            <i class="fas fa-tags"></i>
                            <span><strong>Categories:</strong> ${action.Categories.join(', ')}</span>
                        </div>` : ''}
                </div>
            `;
            
            actionCards.appendChild(card);
        });
        
        actionsSection.style.display = 'block';
    }

    function renderFindings() {
        const findingsList = document.getElementById('findingsList');
        
        allFindings.forEach(finding => {
            renderFindingCard(finding);
        });
    }

    function renderFindingCard(finding) {
        const findingsList = document.getElementById('findingsList');
        const findingCard = document.createElement('div');
        findingCard.className = 'finding-card';
        
        // Determine severity class
        const severityClass = `severity-${finding.severity?.toLowerCase() || 'info'}`;
        
        // Create finding header
        const findingHeader = document.createElement('div');
        findingHeader.className = 'finding-header';
        
        const findingTitle = document.createElement('div');
        findingTitle.className = 'finding-title';
        findingTitle.textContent = finding.title || 'Untitled Finding';
        
        const severityBadge = document.createElement('div');
        severityBadge.className = `severity-badge ${severityClass}`;
        severityBadge.textContent = finding.severity || 'Unknown';
        
        findingHeader.appendChild(findingTitle);
        findingHeader.appendChild(severityBadge);
        
        // Create finding details
        const findingDetails = document.createElement('div');
        findingDetails.className = 'finding-details';
        
        const details = [
            { label: 'Category', value: finding.category || 'Unknown' },
            { label: 'Score', value: finding.score || 'N/A' },
            { label: 'Profile Type', value: finding.evidence?.ProfileType || 'Unknown' },
            { label: 'Artifact', value: finding.evidence?.ArtifactPath || 'N/A' }
        ];
        
        details.forEach(detail => {
            const detailItem = document.createElement('div');
            detailItem.className = 'detail-item';
            
            const detailLabel = document.createElement('div');
            detailLabel.className = 'detail-label';
            detailLabel.textContent = detail.label;
            
            const detailValue = document.createElement('div');
            detailValue.className = 'detail-value';
            detailValue.textContent = detail.value;
            
            detailItem.appendChild(detailLabel);
            detailItem.appendChild(detailValue);
            findingDetails.appendChild(detailItem);
        });
        
        // Add hotspots table if available
        if (finding.top && finding.top.length > 0) {
            const hotspotsTable = document.createElement('table');
            hotspotsTable.className = 'hotspots-table';
            
            const thead = document.createElement('thead');
            const headerRow = document.createElement('tr');
            
            ['Function', 'File', 'Line', 'Cumulative', 'Flat'].forEach(header => {
                const th = document.createElement('th');
                th.textContent = header;
                headerRow.appendChild(th);
            });
            
            thead.appendChild(headerRow);
            hotspotsTable.appendChild(thead);
            
            const tbody = document.createElement('tbody');

            finding.top.forEach(frame => {
                const row = document.createElement('tr');
                
                [frame.function || 'Unknown', frame.file || 'Unknown', frame.line || 'N/A', 
                 frame.cum?.toFixed(2) || '0.00', frame.flat?.toFixed(2) || '0.00'].forEach(cellData => {
                    const td = document.createElement('td');
                    td.textContent = cellData;
                    row.appendChild(td);
                });
                
                tbody.appendChild(row);
            });
            
            hotspotsTable.appendChild(tbody);
            findingCard.appendChild(hotspotsTable);
        }
        
        // Add LLM insights if available
        if (insightsData && insightsData.PerFinding) {
            const findingInsight = insightsData.PerFinding.find(i => i.FindingID === finding.category);
            if (findingInsight) {
                const insightsSection = document.createElement('div');
                insightsSection.className = 'insights-section';
                
                // Add AI header with confidence
                const insightsHeader = document.createElement('div');
                insightsHeader.className = 'insights-header';
                
                const confidence = findingInsight.Confidence || 0;
                let confidenceEmoji = '🟡';
                if (confidence >= 80) {
                    confidenceEmoji = '🟢';
                } else if (confidence <= 50) {
                    confidenceEmoji = '🔴';
                }
                
                insightsHeader.innerHTML = `
                    <div class="insights-header-title">
                        <i class="fas fa-robot"></i> AI Analysis (${confidenceEmoji} ${confidence}% Confidence)
                    </div>
                `;
                
                insightsSection.appendChild(insightsHeader);
                
                if (findingInsight.Narrative) {
                    const narrative = document.createElement('div');
                    narrative.className = 'insights-narrative';
                    narrative.innerHTML = `
                        <div class="narrative-header"><i class="fas fa-comment-dots"></i> Root Cause Analysis</div>
                        <p>${findingInsight.Narrative}</p>
                    `;
                    insightsSection.appendChild(narrative);
                }
                
                if (findingInsight.LikelyRootCauses && findingInsight.LikelyRootCauses.length > 0) {
                    const rootCausesSection = document.createElement('div');
                    rootCausesSection.className = 'insights-root-causes';
                    
                    const rootCausesTitle = document.createElement('div');
                    rootCausesTitle.className = 'root-causes-title';
                    rootCausesTitle.innerHTML = '<i class="fas fa-search"></i> Likely Root Causes';
                    rootCausesSection.appendChild(rootCausesTitle);
                    
                    const rootCausesList = document.createElement('ul');
                    rootCausesList.className = 'root-causes-list';
                    
                    findingInsight.LikelyRootCauses.forEach((cause, i) => {
                        const li = document.createElement('li');
                        li.innerHTML = `<span class="cause-number">${i + 1}.</span> ${cause}`;
                        rootCausesList.appendChild(li);
                    });
                    
                    rootCausesSection.appendChild(rootCausesList);
                    insightsSection.appendChild(rootCausesSection);
                }
                
                if (findingInsight.Suggestions && findingInsight.Suggestions.length > 0) {
                    const suggestionsSection = document.createElement('div');
                    suggestionsSection.className = 'insights-suggestions';
                    
                    const suggestionsTitle = document.createElement('div');
                    suggestionsTitle.className = 'suggestions-title';
                    suggestionsTitle.innerHTML = '<i class="fas fa-lightbulb"></i> Optimization Recommendations';
                    suggestionsSection.appendChild(suggestionsTitle);
                    
                    const suggestionsList = document.createElement('ul');
                    suggestionsList.className = 'suggestions-list';
                    
                    findingInsight.Suggestions.forEach((suggestion, i) => {
                        const li = document.createElement('li');
                        li.innerHTML = `<span class="suggestion-number">${i + 1}.</span> ${suggestion}`;
                        suggestionsList.appendChild(li);
                    });
                    
                    suggestionsSection.appendChild(suggestionsList);
                    insightsSection.appendChild(suggestionsSection);
                }
                
                if (findingInsight.NextMeasurements && findingInsight.NextMeasurements.length > 0) {
                    const measurementsSection = document.createElement('div');
                    measurementsSection.className = 'insights-measurements';
                    
                    const measurementsTitle = document.createElement('div');
                    measurementsTitle.className = 'measurements-title';
                    measurementsTitle.innerHTML = '<i class="fas fa-chart-line"></i> Validation Metrics';
                    measurementsSection.appendChild(measurementsTitle);
                    
                    const measurementsList = document.createElement('ul');
                    measurementsList.className = 'measurements-list';
                    
                    findingInsight.NextMeasurements.forEach((measurement, i) => {
                        const li = document.createElement('li');
                        li.innerHTML = `<span class="measurement-number">${i + 1}.</span> ${measurement}`;
                        measurementsList.appendChild(li);
                    });
                    
                    measurementsSection.appendChild(measurementsList);
                    insightsSection.appendChild(measurementsSection);
                }
                
                if (findingInsight.Caveats && findingInsight.Caveats.length > 0) {
                    const caveatsSection = document.createElement('div');
                    caveatsSection.className = 'insights-caveats';
                    
                    const caveatsTitle = document.createElement('div');
                    caveatsTitle.className = 'caveats-title';
                    caveatsTitle.innerHTML = '<i class="fas fa-exclamation-triangle"></i> Considerations & Limitations';
                    caveatsSection.appendChild(caveatsTitle);
                    
                    const caveatsList = document.createElement('ul');
                    caveatsList.className = 'caveats-list';
                    
                    findingInsight.Caveats.forEach((caveat, i) => {
                        const li = document.createElement('li');
                        li.innerHTML = `<span class="caveat-number">${i + 1}.</span> ${caveat}`;
                        caveatsList.appendChild(li);
                    });
                    
                    caveatsSection.appendChild(caveatsList);
                    insightsSection.appendChild(caveatsSection);
                }
                
                findingCard.appendChild(insightsSection);
            }
        }
        
        // Build the finding card
        findingCard.appendChild(findingHeader);
        findingCard.appendChild(findingDetails);
        
        findingsList.appendChild(findingCard);
    }

    function renderMetricsDashboard() {
        // Set up metrics data
        const performanceScore = findingsData.summary?.OverallScore || 0;
        const criticalCount = allFindings.filter(f => f.severity?.toLowerCase() === 'critical').length;
        const highCount = allFindings.filter(f => f.severity?.toLowerCase() === 'high').length;
        
        // Update metric values
        document.getElementById('performanceScore').textContent = performanceScore;
        document.getElementById('criticalIssues').textContent = criticalCount;
        document.getElementById('highIssues').textContent = highCount;
        document.getElementById('resolvedIssues').textContent = '0'; // Placeholder for resolved issues
        
        // Set performance score color based on value
        const scoreElement = document.getElementById('performanceScore');
        if (performanceScore >= 80) {
            scoreElement.style.color = '#2e7d32'; // Green for good
        } else if (performanceScore >= 60) {
            scoreElement.style.color = '#e65100'; // Orange for medium
        } else {
            scoreElement.style.color = '#c62828'; // Red for poor
        }
        
        // Render trend charts
        renderScoreTrendChart();
        renderResolutionTrendChart();
        
        // Render comparison charts
        renderComparisonChart();
        renderCategoryComparisonChart();
        
        // Render score breakdown chart
        renderScoreBreakdownChart();
        
        // Add hover effects to metric cards
        const metricCards = document.querySelectorAll('.metric-card');
        metricCards.forEach(card => {
            card.addEventListener('mouseenter', function() {
                this.style.transform = 'translateY(-5px) scale(1.02)';
                this.style.boxShadow = '0 8px 16px rgba(0, 0, 0, 0.2)';
            });
            
            card.addEventListener('mouseleave', function() {
                this.style.transform = '';
                this.style.boxShadow = '';
            });
        });
    }

    function renderScoreTrendChart() {
        const ctx = document.getElementById('scoreTrendChart').getContext('2d');
        
        // Mock data for trend - in a real scenario this would come from historical data
        const labels = ['Week 1', 'Week 2', 'Week 3', 'Week 4', 'Current'];
        const data = [75, 78, 82, 80, findingsData.summary?.OverallScore || 0];
        
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Performance Score',
                    data: data,
                    borderColor: '#4a6fa5',
                    backgroundColor: 'rgba(74, 111, 165, 0.1)',
                    tension: 0.4,
                    fill: true,
                    pointBackgroundColor: '#4a6fa5',
                    pointBorderColor: '#fff',
                    pointBorderWidth: 2,
                    pointRadius: 5
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return 'Score: ' + context.parsed.y;
                            }
                        }
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                            stepSize: 20
                        }
                    }
                }
            }
        });
    }

    function renderResolutionTrendChart() {
        const ctx = document.getElementById('resolutionTrendChart').getContext('2d');
        
        // Mock data for resolution trend
        const labels = ['Week 1', 'Week 2', 'Week 3', 'Week 4', 'Current'];
        const resolvedData = [5, 8, 12, 15, 20]; // Mock resolved issues
        const createdData = [10, 12, 15, 18, allFindings.length]; // Current total + mock
        
        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Resolved',
                    data: resolvedData,
                    backgroundColor: '#4ecdc4',
                    borderWidth: 1
                }, {
                    label: 'Created',
                    data: createdData,
                    backgroundColor: '#ff6b6b',
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        position: 'top',
                        labels: {
                            padding: 15,
                            font: {
                                size: 12
                            }
                        }
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            stepSize: 5
                        }
                    }
                }
            }
        });
    }

    function renderComparisonChart() {
        const ctx = document.getElementById('comparisonChart').getContext('2d');
        
        const currentScore = findingsData.summary?.OverallScore || 0;
        const baselineScore = Math.max(0, currentScore - 15); // Mock baseline
        
        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: ['Performance Score'],
                datasets: [{
                    label: 'Baseline',
                    data: [baselineScore],
                    backgroundColor: '#e0e0e0',
                    borderWidth: 1
                }, {
                    label: 'Current',
                    data: [currentScore],
                    backgroundColor: '#4a6fa5',
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                indexAxis: 'y',
                plugins: {
                    legend: {
                        position: 'top',
                        labels: {
                            padding: 15,
                            font: {
                                size: 12
                            }
                        }
                    }
                },
                scales: {
                    x: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                            stepSize: 20
                        }
                    }
                }
            }
        });
    }

    function renderCategoryComparisonChart() {
        const ctx = document.getElementById('categoryComparisonChart').getContext('2d');
        
        // Count findings by category
        const categoryCounts = {};
        allFindings.forEach(finding => {
            const category = finding.category || 'unknown';
            categoryCounts[category] = (categoryCounts[category] || 0) + 1;
        });
        
        // Prepare chart data
        const labels = Object.keys(categoryCounts);
        const data = Object.values(categoryCounts);
        
        // Generate colors
        const generateColors = (count) => {
            const colors = [];
            const baseColors = ['#4a6fa5', '#4ecdc4', '#ffe66d', '#ff6b6b', '#a8e6cf', '#ff8e53'];
            for (let i = 0; i < count; i++) {
                colors.push(baseColors[i % baseColors.length]);
            }
            return colors;
        };
        
        new Chart(ctx, {
            type: 'doughnut',
            data: {
                labels: labels,
                datasets: [{
                    data: data,
                    backgroundColor: generateColors(labels.length),
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        position: 'right',
                        labels: {
                            padding: 15,
                            font: {
                                size: 11
                            }
                        }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                const label = context.label || '';
                                const value = context.parsed || 0;
                                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                const percentage = Math.round((value / total) * 100);
                                return `${label}: ${value} findings (${percentage}%)`;
                            }
                        }
                    }
                }
            }
        });
    }

    function renderScoreBreakdownChart() {
        const ctx = document.getElementById('scoreBreakdownChart').getContext('2d');
        
        // Mock breakdown data - in a real scenario this would come from detailed analysis
        const breakdownData = [
            { category: 'CPU Efficiency', value: 30, color: '#4ecdc4' },
            { category: 'Memory Usage', value: 25, color: '#ffe66d' },
            { category: 'I/O Performance', value: 20, color: '#ff6b6b' },
            { category: 'Concurrency', value: 25, color: '#a8e6cf' }
        ];
        
        new Chart(ctx, {
            type: 'pie',
            data: {
                labels: breakdownData.map(item => item.category),
                datasets: [{
                    data: breakdownData.map(item => item.value),
                    backgroundColor: breakdownData.map(item => item.color),
                    borderWidth: 1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        display: false // We're using custom legend
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                const label = context.label || '';
                                const value = context.parsed || 0;
                                return `${label}: ${value}%`;
                            }
                        }
                    }
                }
            }
        });
    }
});