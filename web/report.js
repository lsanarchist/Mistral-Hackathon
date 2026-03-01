document.addEventListener('DOMContentLoaded', function() {
    const darkModeToggle = document.getElementById('darkModeToggle');
    let findingsData = null;
    let insightsData = null;
    let allFindings = [];
    let severityChart = null;
    let categoryChart = null;

    // Initialize dark mode
    function initDarkMode() {
        const savedMode = localStorage.getItem('darkMode');
        if (savedMode === 'enabled') {
            document.documentElement.classList.add('dark-mode');
            if (darkModeToggle) {
                darkModeToggle.innerHTML = '<i class="fas fa-sun"></i> Light Mode';
            }
        }
    }

    // Toggle dark mode
    function toggleDarkMode() {
        document.documentElement.classList.toggle('dark-mode');
        const isDark = document.documentElement.classList.contains('dark-mode');
        localStorage.setItem('darkMode', isDark ? 'enabled' : 'disabled');
        
        if (darkModeToggle) {
            darkModeToggle.innerHTML = isDark ? '<i class="fas fa-sun"></i> Light Mode' : '<i class="fas fa-moon"></i> Dark Mode';
        }
        
        // Update charts with new theme
        updateChartThemes();
    }

    // Update chart themes based on current mode
    function updateChartThemes() {
        const isDark = document.documentElement.classList.contains('dark-mode');
        const textColor = isDark ? '#e0e0e0' : '#262626';
        const gridColor = isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)';
        const backgroundColor = isDark ? '#2d3136' : '#ffffff';
        
        const chartOptions = {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    labels: {
                        color: textColor
                    }
                }
            },
            scales: {
                x: {
                    ticks: { color: textColor },
                    grid: { color: gridColor }
                },
                y: {
                    ticks: { color: textColor },
                    grid: { color: gridColor }
                }
            }
        };
        
        if (severityChart) {
            severityChart.options = chartOptions;
            severityChart.update();
        }
        
        if (categoryChart) {
            categoryChart.options = chartOptions;
            categoryChart.update();
        }
    }

    // Set up dark mode toggle
    if (darkModeToggle) {
        darkModeToggle.addEventListener('click', toggleDarkMode);
        initDarkMode();
    }

    // Load data from URL parameters or local storage
    function loadData() {
        const urlParams = new URLSearchParams(window.location.search);
        const findingsParam = urlParams.get('findings');
        const insightsParam = urlParams.get('insights');

        if (findingsParam) {
            try {
                findingsData = JSON.parse(decodeURIComponent(findingsParam));
                loadInsightsData();
            } catch (e) {
                console.error('Error parsing findings data:', e);
                showError('Invalid findings data format');
            }
        } else {
            // Try to load from local files
            loadFromFiles();
        }
    }

    // Load insights data
    function loadInsightsData() {
        const urlParams = new URLSearchParams(window.location.search);
        const insightsParam = urlParams.get('insights');

        if (insightsParam) {
            try {
                insightsData = JSON.parse(decodeURIComponent(insightsParam));
                renderData();
            } catch (e) {
                console.error('Error parsing insights data:', e);
                // Continue without insights
                renderData();
            }
        } else {
            renderData();
        }
    }

    // Load from local files (fallback)
    function loadFromFiles() {
        // This would be used when opening the report directly
        // For now, we'll just show a message
        document.getElementById('overview').textContent = 'No data loaded. This report should be opened through the TriageProf tool.';
    }

    // Show error message
    function showError(message) {
        const overview = document.getElementById('overview');
        if (overview) {
            overview.innerHTML = `<div class="error-message">${message}</div>`;
        }
        console.error(message);
    }

    // Render all data
    function renderData() {
        if (!findingsData) {
            showError('No findings data available');
            return;
        }

        // Store all findings for filtering
        allFindings = findingsData.findings || [];

        // Set report metadata
        document.getElementById('reportDate').textContent = new Date().toLocaleString();
        document.getElementById('findingsCount').textContent = allFindings.length;
        
        if (insightsData) {
            document.getElementById('aiStatus').textContent = 'Enabled';
            document.getElementById('aiSummarySection').style.display = 'block';
        } else {
            document.getElementById('aiStatus').textContent = 'Disabled';
        }

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

        // Add fade-in animation
        document.body.classList.add('fade-in');
    }

    // Render summary
    function renderSummary() {
        const summary = findingsData.summary || {};
        
        // Update summary values
        document.getElementById('overallScore').textContent = summary.overall_score || 'N/A';
        
        const severityElement = document.getElementById('severity');
        const severity = (summary.severity || 'unknown').toLowerCase();
        severityElement.textContent = severity.charAt(0).toUpperCase() + severity.slice(1);
        severityElement.className = 'severity-indicator ' + severity;
        
        document.getElementById('confidence').textContent = summary.confidence || 'N/A';
        document.getElementById('findingsTotal').textContent = allFindings.length;
        
        // Set overview text
        const overview = document.getElementById('overview');
        if (summary.overview) {
            overview.textContent = summary.overview;
        } else {
            overview.textContent = 'Performance analysis completed successfully.';
        }
        
        // Add score indicator
        const scoreElement = document.getElementById('overallScore');
        const score = parseFloat(summary.overall_score) || 0;
        
        scoreElement.classList.add('summary-value');
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

    // Render AI summary
    function renderAISummary() {
        const aiContent = document.getElementById('aiInsightsContent');
        
        if (!aiContent) return;
        
        // Clear existing content
        aiContent.innerHTML = '';
        
        // Add executive summary if available
        if (insightsData.executiveSummary) {
            const summary = insightsData.executiveSummary;
            
            const summaryElement = document.createElement('div');
            summaryElement.className = 'ai-insight';
            summaryElement.innerHTML = `
                <div class="ai-insight-header">
                    <h3><i class="fas fa-chart-line"></i> Executive Summary</h3>
                    <span class="ai-confidence">Confidence: ${summary.confidence || 'N/A'}%</span>
                </div>
                <div class="ai-content">
                    <p><strong>Overview:</strong> ${summary.overview || 'No overview available'}</p>
                    <p><strong>Overall Severity:</strong> ${summary.overallSeverity || 'Unknown'}</p>
                    ${summary.keyThemes && summary.keyThemes.length > 0 ? `
                    <p><strong>Key Themes:</strong> ${summary.keyThemes.join(', ')}</p>
                    ` : ''}
                </div>
            `;
            
            aiContent.appendChild(summaryElement);
        }
        
        // Add top risks if available
        if (insightsData.topRisks && insightsData.topRisks.length > 0) {
            const risksElement = document.createElement('div');
            risksElement.className = 'ai-insight';
            
            let risksHTML = '<div class="ai-insight-header"><h3><i class="fas fa-exclamation-triangle"></i> Top Risks</h3></div><div class="ai-content">';
            
            insightsData.topRisks.slice(0, 3).forEach((risk, index) => {
                risksHTML += `
                <div style="margin-bottom: 15px;">
                    <strong>${index + 1}. ${risk.description}</strong><br>
                    <small>Severity: ${risk.severity} | Impact: ${risk.impact} | Likelihood: ${risk.likelihood}</small>
                </div>
                `;
            });
            
            risksHTML += '</div>';
            risksElement.innerHTML = risksHTML;
            aiContent.appendChild(risksElement);
        }
        
        // Add top actions if available
        if (insightsData.topActions && insightsData.topActions.length > 0) {
            const actionsElement = document.createElement('div');
            actionsElement.className = 'ai-insight';
            
            let actionsHTML = '<div class="ai-insight-header"><h3><i class="fas fa-tasks"></i> Top Action Items</h3></div><div class="ai-content">';
            
            insightsData.topActions.slice(0, 3).forEach((action, index) => {
                actionsHTML += `
                <div style="margin-bottom: 15px;">
                    <strong>${index + 1}. ${action.description}</strong><br>
                    <small>Priority: ${action.priority} | Effort: ${action.estimatedEffort}</small>
                    ${action.categories && action.categories.length > 0 ? `<br><small>Categories: ${action.categories.join(', ')}</small>` : ''}
                </div>
                `;
            });
            
            actionsHTML += '</div>';
            actionsElement.innerHTML = actionsHTML;
            aiContent.appendChild(actionsElement);
        }
    }

    // Render quick stats
    function renderQuickStats() {
        let criticalCount = 0;
        let highCount = 0;
        let mediumCount = 0;
        let lowCount = 0;
        
        allFindings.forEach(finding => {
            const severity = (finding.severity || 'low').toLowerCase();
            if (severity === 'critical') criticalCount++;
            else if (severity === 'high') highCount++;
            else if (severity === 'medium') mediumCount++;
            else lowCount++;
        });
        
        document.getElementById('criticalCount').textContent = criticalCount;
        document.getElementById('highCount').textContent = highCount;
        document.getElementById('mediumCount').textContent = mediumCount;
        document.getElementById('lowCount').textContent = lowCount;
    }

    // Render charts
    function renderCharts() {
        const ctxSeverity = document.getElementById('severityChart');
        const ctxCategory = document.getElementById('categoryChart');
        
        if (!ctxSeverity || !ctxCategory) return;
        
        // Count severities
        const severityCounts = { critical: 0, high: 0, medium: 0, low: 0 };
        allFindings.forEach(finding => {
            const severity = (finding.severity || 'low').toLowerCase();
            if (severityCounts[severity] !== undefined) {
                severityCounts[severity]++;
            }
        });
        
        // Count categories
        const categoryCounts = {};
        allFindings.forEach(finding => {
            const category = finding.category || 'other';
            categoryCounts[category] = (categoryCounts[category] || 0) + 1;
        });
        
        // Create severity chart
        if (severityChart) severityChart.destroy();
        severityChart = new Chart(ctxSeverity, {
            type: 'doughnut',
            data: {
                labels: ['Critical', 'High', 'Medium', 'Low'],
                datasets: [{
                    data: [severityCounts.critical, severityCounts.high, severityCounts.medium, severityCounts.low],
                    backgroundColor: [
                        '#dc3545',
                        '#fd7e14',
                        '#ffc107',
                        '#28a745'
                    ],
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });
        
        // Create category chart
        if (categoryChart) categoryChart.destroy();
        categoryChart = new Chart(ctxCategory, {
            type: 'bar',
            data: {
                labels: Object.keys(categoryCounts),
                datasets: [{
                    label: 'Findings by Category',
                    data: Object.values(categoryCounts),
                    backgroundColor: '#4a6bff',
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                indexAxis: 'y',
                plugins: {
                    legend: {
                        display: false
                    }
                },
                scales: {
                    x: {
                        beginAtZero: true
                    }
                }
            }
        });
        
        updateChartThemes();
    }

    // Render findings list
    function renderFindingsList() {
        const container = document.getElementById('findingsContainer');
        if (!container) return;
        
        // Clear existing content
        container.innerHTML = '';
        
        // Render each finding
        allFindings.forEach((finding, index) => {
            const findingCard = createFindingCard(finding, index);
            container.appendChild(findingCard);
        });
        
        // Set up severity filters
        setupSeverityFilters();
    }

    // Create finding card
    function createFindingCard(finding, index) {
        const card = document.createElement('div');
        card.className = 'finding-card';
        card.setAttribute('data-severity', finding.severity || 'low');
        card.setAttribute('data-category', finding.category || 'other');
        
        const severity = (finding.severity || 'low').toLowerCase();
        
        card.innerHTML = `
            <div class="finding-header">
                <h3 class="finding-title">${index + 1}. ${finding.title || 'Untitled Finding'}</h3>
                <div class="finding-meta">
                    <span class="finding-severity ${severity}">${severity.charAt(0).toUpperCase() + severity.slice(1)}</span>
                    <span class="finding-score">Score: ${finding.score || 'N/A'}</span>
                    <span class="finding-category">${finding.category || 'Other'}</span>
                </div>
            </div>
            <div class="finding-content">
                ${createFindingContent(finding)}
            </div>
        `;
        
        return card;
    }

    // Create finding content
    function createFindingContent(finding) {
        let html = '';
        
        // Impact summary
        if (finding.impactSummary) {
            html += `
                <div class="finding-section">
                    <h3><i class="fas fa-bullseye"></i> Impact Summary</h3>
                    <p>${finding.impactSummary}</p>
                </div>
            `;
        }
        
        // Evidence
        if (finding.evidence && finding.evidence.length > 0) {
            html += `
                <div class="finding-section">
                    <h3><i class="fas fa-microscope"></i> Evidence</h3>
                    <ul class="evidence-list">
            `;
            
            finding.evidence.forEach(evidence => {
                html += `
                        <li class="evidence-item">
                            <span class="evidence-type">${evidence.type}: ${evidence.description}</span>
                            <span class="evidence-weight">${Math.round(evidence.weight * 100)}%</span>
                        </li>
                `;
            });
            
            html += `
                    </ul>
                </div>
            `;
        }
        
        // Deterministic hints
        if (finding.deterministicHints && finding.deterministicHints.length > 0) {
            html += `
                <div class="finding-section">
                    <h3><i class="fas fa-lightbulb"></i> Optimization Hints</h3>
                    <div class="hints-section">
                        <ul class="hints-list">
            `;
            
            finding.deterministicHints.forEach(hint => {
                html += `<li class="hint-tag">${hint}</li>`;
            });
            
            html += `
                        </ul>
                    </div>
                </div>
            `;
        }
        
        // Tags
        if (finding.tags && finding.tags.length > 0) {
            html += `
                <div class="finding-section tags-section">
                    <h3><i class="fas fa-tags"></i> Tags</h3>
                    <ul class="tags-list">
            `;
            
            finding.tags.forEach(tag => {
                html += `<li class="tag">${tag}</li>`;
            });
            
            html += `
                    </ul>
                </div>
            `;
        }
        
        // Top hotspots
        if (finding.top && finding.top.length > 0) {
            html += `
                <div class="finding-section">
                    <h3><i class="fas fa-fire"></i> Top Hotspots</h3>
                    <table class="hotspots-table">
                        <thead>
                            <tr>
                                <th>Function</th>
                                <th>File</th>
                                <th>Line</th>
                                <th>Cumulative</th>
                                <th>Flat</th>
                            </tr>
                        </thead>
                        <tbody>
            `;
            
            finding.top.forEach(frame => {
                html += `
                            <tr>
                                <td>${frame.function || 'Unknown'}</td>
                                <td>${frame.file || 'Unknown'}</td>
                                <td>${frame.line || 'N/A'}</td>
                                <td>${frame.cum ? frame.cum.toFixed(2) : 'N/A'}</td>
                                <td>${frame.flat ? frame.flat.toFixed(2) : 'N/A'}</td>
                            </tr>
                `;
            });
            
            html += `
                        </tbody>
                    </table>
                </div>
            `;
        }
        
        // Callgraph
        if (finding.callgraph && finding.callgraph.length > 0) {
            html += `
                <div class="finding-section">
                    <h3><i class="fas fa-project-diagram"></i> Callgraph Analysis</h3>
                    <div class="callgraph-section">
                        <code>${formatCallgraph(finding.callgraph)}
            `;
            
            // Add callgraph statistics
            const totalNodes = countCallgraphNodes(finding.callgraph);
            const maxDepth = findMaxCallgraphDepth(finding.callgraph);
            
            html += `
                        </code>
                    </div>
                    <p style="margin-top: 10px; font-size: 0.9rem; color: #666;">
                        <strong>Callgraph Statistics:</strong> ${totalNodes} nodes, max depth ${maxDepth}
                    </p>
                </div>
            `;
        }
        
        // Allocation analysis
        if (finding.allocationAnalysis) {
            html += createAllocationAnalysisSection(finding.allocationAnalysis);
        }
        
        // Regression analysis
        if (finding.regression) {
            html += createRegressionAnalysisSection(finding.regression);
        }
        
        // AI insights
        if (insightsData && insightsData.perFinding) {
            const findingInsights = insightsData.perFinding.find(insight => 
                insight.findingID === finding.id || insight.findingID === finding.category
            );
            
            if (findingInsights) {
                html += createAIInsightsSection(findingInsights);
            }
        }
        
        return html;
    }

    // Create allocation analysis section
    function createAllocationAnalysisSection(analysis) {
        return `
            <div class="finding-section">
                <div class="allocation-analysis">
                    <h3 style="color: white; margin-top: 0;"><i class="fas fa-memory"></i> Allocation Analysis</h3>
                    <div class="allocation-stats">
                        <div class="allocation-stat">
                            <div class="allocation-stat-value">${Math.round(analysis.totalAllocations)}</div>
                            <div class="allocation-stat-label">Total Allocations</div>
                        </div>
                        <div class="allocation-stat">
                            <div class="allocation-stat-value">${(analysis.topConcentration * 100).toFixed(1)}%</div>
                            <div class="allocation-stat-label">Top Concentration</div>
                        </div>
                        <div class="allocation-stat">
                            <div class="allocation-stat-value">${analysis.severity}</div>
                            <div class="allocation-stat-label">Severity</div>
                        </div>
                        <div class="allocation-stat">
                            <div class="allocation-stat-value">${analysis.score}</div>
                            <div class="allocation-stat-label">Score</div>
                        </div>
                    </div>
                    ${analysis.topConcentration > 0.5 ? `
                    <div style="background: rgba(255, 255, 255, 0.2); padding: 10px; border-radius: 5px; margin-top: 10px;">
                        <strong>⚠️ High Allocation Concentration:</strong> Top functions account for ${(analysis.topConcentration * 100).toFixed(1)}% of allocations.
                    </div>
                    ` : `
                    <div style="background: rgba(255, 255, 255, 0.2); padding: 10px; border-radius: 5px; margin-top: 10px;">
                        <strong>✅ Balanced Allocation:</strong> Allocations are reasonably distributed.
                    </div>
                    `}
                    ${analysis.hotspots && analysis.hotspots.length > 0 ? `
                    <div style="margin-top: 15px;">
                        <h4 style="color: white; margin-bottom: 10px;">Top Allocation Hotspots</h4>
                        <table style="width: 100%; border-collapse: collapse;">
                            <thead>
                                <tr style="background: rgba(255, 255, 255, 0.3);">
                                    <th style="padding: 8px; text-align: left;">Function</th>
                                    <th style="padding: 8px; text-align: left;">File</th>
                                    <th style="padding: 8px; text-align: left;">Count</th>
                                    <th style="padding: 8px; text-align: left;">Percentage</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${analysis.hotspots.map(hotspot => `
                                <tr style="background: rgba(255, 255, 255, 0.1);">
                                    <td style="padding: 8px;">${hotspot.function}</td>
                                    <td style="padding: 8px;">${hotspot.file}</td>
                                    <td style="padding: 8px;">${Math.round(hotspot.count)}</td>
                                    <td style="padding: 8px;">${hotspot.percent.toFixed(1)}%</td>
                                </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                    ` : ''}
                </div>
            </div>
        `;
    }

    // Create regression analysis section
    function createRegressionAnalysisSection(regression) {
        const severity = regression.severity.toLowerCase();
        let message = '';
        let icon = '⚠️';
        let color = '#fa709a';
        
        if (severity === 'improved') {
            message = 'Performance Improvement Detected';
            icon = '📈';
            color = '#4facfe';
        } else if (severity !== 'none' && severity !== 'low') {
            message = `Potential Regression Detected (${severity})`;
        } else {
            message = 'No Significant Regression';
            icon = '✅';
            color = '#28a745';
        }
        
        return `
            <div class="finding-section">
                <div class="regression-analysis">
                    <h3 style="margin-top: 0;"><i class="fas fa-chart-line"></i> Regression Analysis</h3>
                    <div class="regression-stats">
                        <div class="regression-stat">
                            <div class="regression-stat-value">${regression.baselineScore}</div>
                            <div class="regression-stat-label">Baseline</div>
                        </div>
                        <div class="regression-stat">
                            <div class="regression-stat-value">${regression.currentScore}</div>
                            <div class="regression-stat-label">Current</div>
                        </div>
                        <div class="regression-stat">
                            <div class="regression-stat-value">${regression.delta} (${regression.percentage.toFixed(1)}%)</div>
                            <div class="regression-stat-label">Delta</div>
                        </div>
                        <div class="regression-stat">
                            <div class="regression-stat-value">${regression.severity}</div>
                            <div class="regression-stat-label">Severity</div>
                        </div>
                    </div>
                    <div style="background: rgba(0, 0, 0, 0.1); padding: 10px; border-radius: 5px; margin-top: 10px;">
                        <strong>${icon} ${message}:</strong> ${regression.currentScore} vs ${regression.baselineScore} (${regression.delta} points, ${regression.percentage.toFixed(1)}% change)
                    </div>
                </div>
            </div>
        `;
    }

    // Create AI insights section
    function createAIInsightsSection(insight) {
        const confidenceEmoji = insight.confidence >= 80 ? '🟢' : insight.confidence >= 50 ? '🟡' : '🔴';
        
        return `
            <div class="finding-section">
                <div class="ai-insights-section">
                    <div class="ai-insights-header">
                        <h3><i class="fas fa-robot"></i> AI Insights</h3>
                        <span class="ai-confidence-badge">${confidenceEmoji} ${insight.confidence}% Confidence</span>
                    </div>
                    <div class="ai-content">
                        <p><strong>Narrative:</strong> ${insight.narrative}</p>
                        ${insight.likelyRootCauses && insight.likelyRootCauses.length > 0 ? `
                        <div>
                            <strong>🔍 Likely Root Causes:</strong>
                            <ul class="ai-list">
                                ${insight.likelyRootCauses.map(cause => `<li>${cause}</li>`).join('')}
                            </ul>
                        </div>
                        ` : ''}
                        ${insight.suggestions && insight.suggestions.length > 0 ? `
                        <div>
                            <strong>💡 Suggestions:</strong>
                            <ul class="ai-list">
                                ${insight.suggestions.map(suggestion => `<li>${suggestion}</li>`).join('')}
                            </ul>
                        </div>
                        ` : ''}
                        ${insight.nextMeasurements && insight.nextMeasurements.length > 0 ? `
                        <div>
                            <strong>📊 Next Measurements:</strong>
                            <ul class="ai-list">
                                ${insight.nextMeasurements.map(measurement => `<li>${measurement}</li>`).join('')}
                            </ul>
                        </div>
                        ` : ''}
                        ${insight.caveats && insight.caveats.length > 0 ? `
                        <div>
                            <strong>⚠️ Caveats:</strong>
                            <ul class="ai-list">
                                ${insight.caveats.map(caveat => `<li>${caveat}</li>`).join('')}
                            </ul>
                        </div>
                        ` : ''}
                    </div>
                </div>
            </div>
        `;
    }

    // Format callgraph for display
    function formatCallgraph(nodes, indent = 0) {
        let result = '';
        const indentStr = '  '.repeat(indent);
        
        nodes.forEach((node, index) => {
            const isLast = index === nodes.length - 1;
            const prefix = isLast ? '└── ' : '├── ';
            
            result += `${indentStr}${prefix}${node.function} (cum: ${node.cum.toFixed(1)}, flat: ${node.flat.toFixed(1)})
`;
            
            if (node.children && node.children.length > 0) {
                result += formatCallgraph(node.children, indent + 1);
            }
        });
        
        return result;
    }

    // Count callgraph nodes
    function countCallgraphNodes(nodes) {
        let count = 0;
        nodes.forEach(node => {
            count += countCallgraphNode(node);
        });
        return count;
    }

    function countCallgraphNode(node) {
        let count = 1;
        if (node.children) {
            node.children.forEach(child => {
                count += countCallgraphNode(child);
            });
        }
        return count;
    }

    // Find max callgraph depth
    function findMaxCallgraphDepth(nodes) {
        let maxDepth = 0;
        nodes.forEach(node => {
            const depth = findMaxCallgraphNodeDepth(node);
            if (depth > maxDepth) maxDepth = depth;
        });
        return maxDepth;
    }

    function findMaxCallgraphNodeDepth(node) {
        let maxDepth = node.depth || 0;
        if (node.children) {
            node.children.forEach(child => {
                const childDepth = findMaxCallgraphNodeDepth(child);
                if (childDepth > maxDepth) maxDepth = childDepth;
            });
        }
        return maxDepth;
    }

    // Set up severity filters
    function setupSeverityFilters() {
        const filters = document.querySelectorAll('.severity-filter');
        
        filters.forEach(filter => {
            filter.addEventListener('click', function() {
                // Remove active class from all filters
                filters.forEach(f => f.classList.remove('active'));
                
                // Add active class to clicked filter
                this.classList.add('active');
                
                // Get selected severity
                const severity = this.getAttribute('data-severity');
                
                // Filter findings
                filterFindings(severity);
            });
        });
    }

    // Filter findings by severity
    function filterFindings(severity) {
        const cards = document.querySelectorAll('.finding-card');
        
        cards.forEach(card => {
            if (severity === 'all') {
                card.style.display = 'block';
            } else {
                const cardSeverity = card.getAttribute('data-severity') || 'low';
                if (cardSeverity.toLowerCase() === severity.toLowerCase()) {
                    card.style.display = 'block';
                } else {
                    card.style.display = 'none';
                }
            }
        });
    }

    // Initialize the report
    function init() {
        // Try to load data from URL parameters
        const urlParams = new URLSearchParams(window.location.search);
        
        if (urlParams.has('findings')) {
            try {
                findingsData = JSON.parse(decodeURIComponent(urlParams.get('findings')));
                
                if (urlParams.has('insights')) {
                    insightsData = JSON.parse(decodeURIComponent(urlParams.get('insights')));
                }
                
                renderData();
            } catch (e) {
                console.error('Error loading data from URL:', e);
                showError('Invalid data format in URL parameters');
            }
        } else {
            // No data in URL, show instructions
            const overview = document.getElementById('overview');
            if (overview) {
                overview.innerHTML = `
                    <p style="text-align: center; color: #666;">
                        <i class="fas fa-info-circle" style="font-size: 2rem; color: #4a6bff; display: block; margin-bottom: 10px;"></i>
                        This report should be opened through the TriageProf tool.<br>
                        It will automatically load performance findings and insights.
                    </p>
                `;
            }
        }
    }

    // Start the application
    init();
});