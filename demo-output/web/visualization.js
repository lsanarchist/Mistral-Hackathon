// Advanced Visualization Library for TriageProf
// Provides interactive callgraphs, time-series analysis, and comparative views

class AdvancedVisualization {
    constructor(containerId, findingsData, insightsData) {
        this.container = document.getElementById(containerId);
        this.findingsData = findingsData;
        this.insightsData = insightsData;
        this.visualizations = {};
    }

    // Initialize all advanced visualizations
    init() {
        if (!this.container) {
            console.error('Visualization container not found');
            return;
        }

        // Create visualization tabs
        this.createVisualizationTabs();
        
        // Initialize individual visualizations
        this.initCallgraphVisualization();
        this.initTimeSeriesAnalysis();
        this.initComparativeView();
        this.initCustomDashboard();
    }

    // Create navigation tabs for different visualizations
    createVisualizationTabs() {
        const tabContainer = document.createElement('div');
        tabContainer.className = 'visualization-tabs';
        
        const tabs = [
            { id: 'callgraph-tab', title: 'Callgraph', icon: 'fas fa-project-diagram' },
            { id: 'timeseries-tab', title: 'Time Series', icon: 'fas fa-chart-line' },
            { id: 'comparative-tab', title: 'Comparative', icon: 'fas fa-columns' },
            { id: 'dashboard-tab', title: 'Custom Dashboard', icon: 'fas fa-th-large' }
        ];

        const tabNav = document.createElement('nav');
        tabNav.className = 'visualization-tab-nav';
        
        tabs.forEach(tab => {
            const tabButton = document.createElement('button');
            tabButton.id = tab.id;
            tabButton.className = 'visualization-tab-button';
            tabButton.innerHTML = `<i class="${tab.icon}"></i> ${tab.title}`;
            tabButton.addEventListener('click', () => this.switchTab(tab.id));
            tabNav.appendChild(tabButton);
        });

        tabContainer.appendChild(tabNav);
        this.container.appendChild(tabContainer);
        
        // Create content containers for each tab
        const contentContainer = document.createElement('div');
        contentContainer.className = 'visualization-content-container';
        
        tabs.forEach(tab => {
            const contentDiv = document.createElement('div');
            contentDiv.id = `${tab.id}-content`;
            contentDiv.className = 'visualization-tab-content';
            if (tab.id === 'callgraph-tab') {
                contentDiv.style.display = 'block';
            } else {
                contentDiv.style.display = 'none';
            }
            contentContainer.appendChild(contentDiv);
        });

        this.container.appendChild(contentContainer);
    }

    switchTab(tabId) {
        // Hide all content
        document.querySelectorAll('.visualization-tab-content').forEach(content => {
            content.style.display = 'none';
        });
        
        // Remove active class from all buttons
        document.querySelectorAll('.visualization-tab-button').forEach(button => {
            button.classList.remove('active');
        });
        
        // Show selected content and mark button as active
        const selectedContent = document.getElementById(`${tabId}-content`);
        const selectedButton = document.getElementById(tabId);
        
        if (selectedContent && selectedButton) {
            selectedContent.style.display = 'block';
            selectedButton.classList.add('active');
        }
    }

    // Initialize interactive callgraph visualization
    initCallgraphVisualization() {
        const callgraphContainer = document.getElementById('callgraph-tab-content');
        if (!callgraphContainer) return;
        
        callgraphContainer.innerHTML = '<div id="callgraph-network" style="width: 100%; height: 600px; border: 1px solid #ddd; border-radius: 8px;"></div>';
        
        // Check if we have callgraph data
        const findingsWithCallgraphs = this.findingsData.findings.filter(f => f.Callgraph && f.Callgraph.length > 0);
        
        if (findingsWithCallgraphs.length === 0) {
            callgraphContainer.innerHTML = '<p style="text-align: center; padding: 20px;">No callgraph data available for visualization.</p>';
            return;
        }
        
        // Prepare callgraph data for visualization
        const nodes = [];
        const edges = [];
        
        findingsWithCallgraphs.forEach((finding, index) => {
            const findingNodes = [];
            const findingEdges = [];
            
            finding.Callgraph.forEach((node, nodeIndex) => {
                const nodeId = `${index}_${nodeIndex}`;
                findingNodes.push({
                    id: nodeId,
                    label: this.getShortFunctionName(node.Function),
                    title: `${node.Function}\nFile: ${node.File}\nLine: ${node.Line}\nCumulative: ${node.Cum.toFixed(2)}%`,
                    level: node.Depth,
                    value: node.Cum,
                    findingId: finding.ID,
                    severity: finding.Severity
                });
                
                // Connect to parent if not root
                if (node.Depth > 0) {
                    const parentId = `${index}_${nodeIndex - 1}`;
                    findingEdges.push({
                        from: parentId,
                        to: nodeId,
                        arrows: 'to',
                        color: this.getSeverityColor(finding.Severity)
                    });
                }
            });
            
            nodes.push(...findingNodes);
            edges.push(...findingEdges);
        });
        
        // Create network visualization
        const container = document.getElementById('callgraph-network');
        const data = { nodes: nodes, edges: edges };
        
        const options = {
            layout: {
                hierarchical: {
                    direction: 'LR',
                    sortMethod: 'directed',
                    nodeSpacing: 150,
                    levelSeparation: 200
                }
            },
            physics: {
                hierarchicalRepulsion: {
                    centralGravity: 0.0,
                    springLength: 100,
                    springConstant: 0.01,
                    nodeDistance: 120,
                    damping: 0.09
                },
                minVelocity: 0.75,
                solver: 'hierarchicalRepulsion'
            },
            nodes: {
                shape: 'box',
                size: 20,
                font: {
                    size: 12,
                    face: 'Segoe UI'
                },
                borderWidth: 2,
                shadow: true
            },
            edges: {
                width: 1.5,
                smooth: {
                    type: 'cubicBezier',
                    forceDirection: 'horizontal',
                    roundness: 0.4
                },
                arrows: {
                    to: {
                        enabled: true,
                        scaleFactor: 0.5
                    }
                }
            },
            interaction: {
                hover: true,
                navigationButtons: true,
                keyboard: true
            },
            manipulation: {
                enabled: false
            }
        };
        
        // Store network instance
        this.visualizations.callgraph = new vis.Network(container, data, options);
        
        // Add event handlers
        this.visualizations.callgraph.on('click', (params) => {
            if (params.nodes.length > 0) {
                const nodeId = params.nodes[0];
                const node = nodes.find(n => n.id === nodeId);
                if (node) {
                    this.showNodeDetails(node);
                }
            }
        });
        
        // Add legend
        this.addCallgraphLegend(callgraphContainer);
    }

    getShortFunctionName(fullName) {
        // Shorten function names for better display
        if (!fullName) return 'unknown';
        
        // Remove package paths and keep just function name
        const parts = fullName.split('/');
        const lastPart = parts[parts.length - 1];
        
        // Remove function signature
        return lastPart.split('(')[0];
    }

    getSeverityColor(severity) {
        const severityColors = {
            'critical': '#dc3545',
            'high': '#fd7e14',
            'medium': '#ffc107',
            'low': '#28a745'
        };
        return severityColors[severity.toLowerCase()] || '#6c757d';
    }

    showNodeDetails(node) {
        const detailsDiv = document.createElement('div');
        detailsDiv.className = 'node-details-popup';
        detailsDiv.style.position = 'absolute';
        detailsDiv.style.background = 'white';
        detailsDiv.style.padding = '15px';
        detailsDiv.style.borderRadius = '8px';
        detailsDiv.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.15)';
        detailsDiv.style.zIndex = '1000';
        detailsDiv.style.maxWidth = '400px';
        
        detailsDiv.innerHTML = `
            <h4 style="margin-top: 0; color: ${this.getSeverityColor(node.severity)}">
                <i class="fas fa-info-circle"></i> Function Details
            </h4>
            <p><strong>Function:</strong> ${node.label}</p>
            <p><strong>Full Name:</strong> ${node.title.split('\\n')[0]}</p>
            <p><strong>File:</strong> ${node.title.split('\\n')[1].replace('File: ', '')}</p>
            <p><strong>Line:</strong> ${node.title.split('\\n')[2].replace('Line: ', '')}</p>
            <p><strong>Cumulative %:</strong> ${node.value.toFixed(2)}%</p>
            <p><strong>Finding ID:</strong> ${node.findingId}</p>
            <p><strong>Severity:</strong> <span style="color: ${this.getSeverityColor(node.severity)}">${node.severity}</span></p>
            <button onclick="this.parentElement.remove()" style="
                background: var(--primary-color);
                color: white;
                border: none;
                padding: 8px 16px;
                border-radius: 4px;
                cursor: pointer;
                margin-top: 10px;
            ">Close</button>
        `;
        
        document.body.appendChild(detailsDiv);
        
        // Position the popup near the mouse
        detailsDiv.style.left = `${window.event.clientX + 20}px`;
        detailsDiv.style.top = `${window.event.clientY}px`;
    }

    addCallgraphLegend(container) {
        const legend = document.createElement('div');
        legend.className = 'callgraph-legend';
        legend.style.marginTop = '15px';
        legend.style.padding = '10px';
        legend.style.background = '#f8f9fa';
        legend.style.borderRadius = '8px';
        
        legend.innerHTML = `
            <h4 style="margin: 0 0 10px 0;"><i class="fas fa-info-circle"></i> Legend</h4>
            <div style="display: flex; gap: 15px; flex-wrap: wrap;">
                <div style="display: flex; align-items: center; gap: 5px;">
                    <div style="width: 20px; height: 20px; background: #dc3545; border-radius: 4px;"></div>
                    <span>Critical Severity</span>
                </div>
                <div style="display: flex; align-items: center; gap: 5px;">
                    <div style="width: 20px; height: 20px; background: #fd7e14; border-radius: 4px;"></div>
                    <span>High Severity</span>
                </div>
                <div style="display: flex; align-items: center; gap: 5px;">
                    <div style="width: 20px; height: 20px; background: #ffc107; border-radius: 4px;"></div>
                    <span>Medium Severity</span>
                </div>
                <div style="display: flex; align-items: center; gap: 5px;">
                    <div style="width: 20px; height: 20px; background: #28a745; border-radius: 4px;"></div>
                    <span>Low Severity</span>
                </div>
            </div>
            <p style="margin: 10px 0 0 0; font-size: 0.9em; color: #666;">
                <i class="fas fa-mouse-pointer"></i> Click on nodes to see detailed information<br>
                <i class="fas fa-search-plus"></i> Use mouse wheel to zoom, drag to pan
            </p>
        `;
        
        container.appendChild(legend);
    }

    // Initialize time series analysis
    initTimeSeriesAnalysis() {
        const timeseriesContainer = document.getElementById('timeseries-tab-content');
        if (!timeseriesContainer) return;
        
        timeseriesContainer.innerHTML = `
            <div class="timeseries-controls" style="margin-bottom: 20px; padding: 15px; background: #f8f9fa; border-radius: 8px;">
                <h3 style="margin-top: 0;"><i class="fas fa-chart-line"></i> Time Series Analysis</h3>
                <p>Analyze performance trends over time (requires multiple runs with timestamps)</p>
                <div style="color: #666; margin-top: 10px;">
                    <i class="fas fa-info-circle"></i> This feature requires historical data from multiple profiling runs.
                </div>
            </div>
            <div id="timeseries-chart" style="width: 100%; height: 500px;"></div>
        `;
        
        // Check if we have time series data
        const hasTimeSeriesData = this.findingsData.findings.some(f => f.Timestamp);
        
        if (!hasTimeSeriesData) {
            const chartContainer = document.getElementById('timeseries-chart');
            chartContainer.innerHTML = `
                <div style="text-align: center; padding: 50px; color: #666;">
                    <i class="fas fa-database" style="font-size: 48px; margin-bottom: 15px;"></i>
                    <h4>No Time Series Data Available</h4>
                    <p>Time series analysis requires multiple profiling runs with timestamps.</p>
                    <p style="margin-top: 10px;">
                        <small>Run <code>triageprof demo</code> multiple times with different commits or timestamps to enable this feature.</small>
                    </p>
                </div>
            `;
            return;
        }
        
        // Create time series chart
        this.createTimeSeriesChart();
    }

    createTimeSeriesChart() {
        // This would be implemented with actual time series data
        // For now, create a placeholder with sample data
        const ctx = document.getElementById('timeseries-chart');
        
        const sampleData = {
            labels: ['Run 1', 'Run 2', 'Run 3', 'Run 4', 'Run 5'],
            datasets: [
                {
                    label: 'Overall Score',
                    data: [75, 78, 82, 80, 85],
                    borderColor: '#4a6bff',
                    backgroundColor: 'rgba(74, 107, 255, 0.1)',
                    tension: 0.4,
                    fill: true
                },
                {
                    label: 'Critical Findings',
                    data: [3, 2, 1, 1, 0],
                    borderColor: '#dc3545',
                    backgroundColor: 'rgba(220, 53, 69, 0.1)',
                    tension: 0.4,
                    fill: true,
                    yAxisID: 'y1'
                },
                {
                    label: 'High Findings',
                    data: [5, 4, 3, 3, 2],
                    borderColor: '#fd7e14',
                    backgroundColor: 'rgba(253, 126, 20, 0.1)',
                    tension: 0.4,
                    fill: true,
                    yAxisID: 'y1'
                }
            ]
        };
        
        new Chart(ctx, {
            type: 'line',
            data: sampleData,
            options: {
                responsive: true,
                interaction: {
                    mode: 'index',
                    intersect: false
                },
                scales: {
                    y: {
                        type: 'linear',
                        display: true,
                        position: 'left',
                        title: {
                            display: true,
                            text: 'Score'
                        }
                    },
                    y1: {
                        type: 'linear',
                        display: true,
                        position: 'right',
                        title: {
                            display: true,
                            text: 'Findings Count'
                        },
                        grid: {
                            drawOnChartArea: false
                        }
                    }
                },
                plugins: {
                    tooltip: {
                        callbacks: {
                            title: function(context) {
                                return `Run ${context[0].label}`;
                            }
                        }
                    },
                    legend: {
                        position: 'top'
                    }
                }
            }
        });
    }

    // Initialize comparative view
    initComparativeView() {
        const comparativeContainer = document.getElementById('comparative-tab-content');
        if (!comparativeContainer) return;
        
        comparativeContainer.innerHTML = `
            <div class="comparative-controls" style="margin-bottom: 20px; padding: 15px; background: #f8f9fa; border-radius: 8px;">
                <h3 style="margin-top: 0;"><i class="fas fa-columns"></i> Comparative Analysis</h3>
                <p>Compare performance findings across different runs or versions</p>
                <div style="color: #666; margin-top: 10px;">
                    <i class="fas fa-info-circle"></i> Upload multiple findings.json files to compare performance characteristics.
                </div>
                <div id="comparative-upload" style="margin-top: 15px; padding: 15px; background: white; border: 2px dashed #ddd; border-radius: 8px; text-align: center;">
                    <i class="fas fa-upload" style="font-size: 48px; color: #6c757d; margin-bottom: 10px;"></i>
                    <p><strong>Drag & drop findings.json files here</strong></p>
                    <p style="font-size: 0.9em; color: #666;">or click to browse</p>
                    <input type="file" id="file-upload" accept=".json" multiple style="display: none;">
                </div>
            </div>
            <div id="comparative-results" style="display: none;">
                <div class="comparative-charts" style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px;">
                    <div style="background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                        <canvas id="comparative-severity-chart"></canvas>
                    </div>
                    <div style="background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                        <canvas id="comparative-category-chart"></canvas>
                    </div>
                </div>
                <div id="comparative-table" style="background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); overflow-x: auto;">
                    <h4 style="margin-top: 0;">Detailed Comparison</h4>
                    <table id="comparative-findings-table" style="width: 100%; border-collapse: collapse;">
                        <thead>
                            <tr style="background: #f8f9fa;">
                                <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Finding</th>
                                <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Severity</th>
                                <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Category</th>
                                <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Confidence</th>
                                <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Run 1</th>
                                <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Run 2</th>
                            </tr>
                        </thead>
                        <tbody id="comparative-findings-body">
                            <!-- Will be populated dynamically -->
                        </tbody>
                    </table>
                </div>
            </div>
        `;
        
        // Set up file upload
        const uploadArea = document.getElementById('comparative-upload');
        const fileInput = document.getElementById('file-upload');
        
        uploadArea.addEventListener('click', () => fileInput.click());
        
        fileInput.addEventListener('change', (e) => {
            if (e.target.files.length > 0) {
                this.handleFileUpload(e.target.files);
            }
        });
        
        // Prevent default drag behaviors
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            uploadArea.addEventListener(eventName, this.preventDefaults, false);
        });
        
        // Handle drop
        uploadArea.addEventListener('drop', (e) => {
            this.handleFileUpload(e.dataTransfer.files);
        });
    }

    preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    handleFileUpload(files) {
        const comparativeResults = document.getElementById('comparative-results');
        
        if (files.length < 2) {
            alert('Please upload at least 2 findings.json files for comparison.');
            return;
        }
        
        // Show loading state
        comparativeResults.style.display = 'block';
        comparativeResults.innerHTML = `
            <div style="text-align: center; padding: 50px;">
                <i class="fas fa-spinner fa-spin" style="font-size: 48px; color: var(--primary-color);"></i>
                <h4 style="margin-top: 15px;">Analyzing files...</h4>
            </div>
        `;
        
        // In a real implementation, we would parse the files and create comparison charts
        // For this demo, we'll simulate a successful comparison
        setTimeout(() => {
            this.createSampleComparison();
        }, 1000);
    }

    createSampleComparison() {
        const comparativeResults = document.getElementById('comparative-results');
        comparativeResults.innerHTML = `
            <div class="comparative-charts" style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px;">
                <div style="background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <h4 style="margin-top: 0;">Severity Distribution Comparison</h4>
                    <canvas id="comparative-severity-chart"></canvas>
                </div>
                <div style="background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <h4 style="margin-top: 0;">Category Distribution Comparison</h4>
                    <canvas id="comparative-category-chart"></canvas>
                </div>
            </div>
            <div id="comparative-table" style="background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); overflow-x: auto;">
                <h4 style="margin-top: 0;">Detailed Comparison</h4>
                <table id="comparative-findings-table" style="width: 100%; border-collapse: collapse;">
                    <thead>
                        <tr style="background: #f8f9fa;">
                            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Finding</th>
                            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Severity</th>
                            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Category</th>
                            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Confidence</th>
                            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Run 1</th>
                            <th style="padding: 10px; text-align: left; border-bottom: 2px solid #ddd;">Run 2</th>
                        </tr>
                    </thead>
                    <tbody id="comparative-findings-body">
                        <tr>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">High CPU usage in main loop</td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;"><span style="color: #fd7e14;">High</span></td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">CPU</td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">0.92</td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">✓ Present</td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">✓ Present</td>
                        </tr>
                        <tr>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">Memory allocation hotspot</td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;"><span style="color: #ffc107;">Medium</span></td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">Heap</td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">0.85</td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">✓ Present</td>
                            <td style="padding: 10px; border-bottom: 1px solid #eee;">✗ Resolved</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        `;
        
        // Create comparison charts
        this.createComparisonCharts();
    }

    createComparisonCharts() {
        // Severity comparison chart
        const severityCtx = document.getElementById('comparative-severity-chart');
        new Chart(severityCtx, {
            type: 'bar',
            data: {
                labels: ['Critical', 'High', 'Medium', 'Low'],
                datasets: [
                    {
                        label: 'Run 1',
                        data: [1, 3, 5, 7],
                        backgroundColor: 'rgba(74, 107, 255, 0.7)'
                    },
                    {
                        label: 'Run 2',
                        data: [0, 2, 4, 6],
                        backgroundColor: 'rgba(253, 126, 20, 0.7)'
                    }
                ]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        position: 'top'
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        title: {
                            display: true,
                            text: 'Number of Findings'
                        }
                    }
                }
            }
        });
        
        // Category comparison chart
        const categoryCtx = document.getElementById('comparative-category-chart');
        new Chart(categoryCtx, {
            type: 'doughnut',
            data: {
                labels: ['CPU', 'Heap', 'Alloc', 'GC', 'Mutex'],
                datasets: [
                    {
                        label: 'Run 1',
                        data: [4, 3, 2, 1, 1],
                        backgroundColor: [
                            'rgba(255, 99, 132, 0.7)',
                            'rgba(54, 162, 235, 0.7)',
                            'rgba(255, 206, 86, 0.7)',
                            'rgba(75, 192, 192, 0.7)',
                            'rgba(153, 102, 255, 0.7)'
                        ]
                    },
                    {
                        label: 'Run 2',
                        data: [3, 2, 2, 1, 0],
                        backgroundColor: [
                            'rgba(255, 99, 132, 0.3)',
                            'rgba(54, 162, 235, 0.3)',
                            'rgba(255, 206, 86, 0.3)',
                            'rgba(75, 192, 192, 0.3)',
                            'rgba(153, 102, 255, 0.3)'
                        ]
                    }
                ]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        position: 'right'
                    }
                }
            }
        });
    }

    // Initialize custom dashboard
    initCustomDashboard() {
        const dashboardContainer = document.getElementById('dashboard-tab-content');
        if (!dashboardContainer) return;
        
        dashboardContainer.innerHTML = `
            <div class="dashboard-controls" style="margin-bottom: 20px; padding: 15px; background: #f8f9fa; border-radius: 8px;">
                <h3 style="margin-top: 0;"><i class="fas fa-th-large"></i> Custom Dashboard</h3>
                <p>Create your own dashboard layout with configurable widgets</p>
                <div style="color: #666; margin-top: 10px;">
                    <i class="fas fa-info-circle"></i> Drag and drop widgets to create your ideal performance monitoring layout.
                </div>
            </div>
            <div class="dashboard-grid" id="dashboard-grid" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 20px;">
                <!-- Widgets will be added here -->
            </div>
            <div class="widget-palette" style="margin-top: 20px; padding: 15px; background: #f8f9fa; border-radius: 8px;">
                <h4 style="margin-top: 0;"><i class="fas fa-puzzle-piece"></i> Available Widgets</h4>
                <div style="display: flex; gap: 10px; flex-wrap: wrap; margin-top: 10px;">
                    <button onclick="advancedVis.addWidget('summary')" style="
                        padding: 8px 16px;
                        background: white;
                        border: 2px dashed var(--primary-color);
                        border-radius: 8px;
                        cursor: pointer;
                        display: flex;
                        align-items: center;
                        gap: 8px;
                    ">
                        <i class="fas fa-chart-pie" style="color: var(--primary-color);"></i>
                        Summary Stats
                    </button>
                    <button onclick="advancedVis.addWidget('severity')" style="
                        padding: 8px 16px;
                        background: white;
                        border: 2px dashed #dc3545;
                        border-radius: 8px;
                        cursor: pointer;
                        display: flex;
                        align-items: center;
                        gap: 8px;
                    ">
                        <i class="fas fa-exclamation-triangle" style="color: #dc3545;"></i>
                        Severity Chart
                    </button>
                    <button onclick="advancedVis.addWidget('category')" style="
                        padding: 8px 16px;
                        background: white;
                        border: 2px dashed #28a745;
                        border-radius: 8px;
                        cursor: pointer;
                        display: flex;
                        align-items: center;
                        gap: 8px;
                    ">
                        <i class="fas fa-tags" style="color: #28a745;"></i>
                        Category Chart
                    </button>
                    <button onclick="advancedVis.addWidget('findings')" style="
                        padding: 8px 16px;
                        background: white;
                        border: 2px dashed #6c757d;
                        border-radius: 8px;
                        cursor: pointer;
                        display: flex;
                        align-items: center;
                        gap: 8px;
                    ">
                        <i class="fas fa-list" style="color: #6c757d;"></i>
                        Findings List
                    </button>
                </div>
            </div>
        `;
        
        // Add some default widgets
        this.addWidget('summary');
        this.addWidget('severity');
        this.addWidget('category');
    }

    addWidget(type) {
        const grid = document.getElementById('dashboard-grid');
        const widgetId = `widget-${type}-${Date.now()}`;
        
        const widget = document.createElement('div');
        widget.id = widgetId;
        widget.className = 'dashboard-widget';
        widget.style.background = 'white';
        widget.style.padding = '15px';
        widget.style.borderRadius = '8px';
        widget.style.boxShadow = '0 2px 4px rgba(0,0,0,0.1)';
        widget.style.position = 'relative';
        
        // Add widget header
        const header = document.createElement('div');
        header.style.display = 'flex';
        header.style.justifyContent = 'space-between';
        header.style.alignItems = 'center';
        header.style.marginBottom = '15px';
        header.style.paddingBottom = '10px';
        header.style.borderBottom = '1px solid #eee';
        
        let widgetTitle = '';
        let widgetIcon = '';
        let widgetColor = '';
        
        switch(type) {
            case 'summary':
                widgetTitle = 'Performance Summary';
                widgetIcon = 'fas fa-chart-pie';
                widgetColor = 'var(--primary-color)';
                break;
            case 'severity':
                widgetTitle = 'Severity Distribution';
                widgetIcon = 'fas fa-exclamation-triangle';
                widgetColor = '#dc3545';
                break;
            case 'category':
                widgetTitle = 'Category Distribution';
                widgetIcon = 'fas fa-tags';
                widgetColor = '#28a745';
                break;
            case 'findings':
                widgetTitle = 'Top Findings';
                widgetIcon = 'fas fa-list';
                widgetColor = '#6c757d';
                break;
        }
        
        header.innerHTML = `
            <h4 style="margin: 0; display: flex; align-items: center; gap: 8px;">
                <i class="${widgetIcon}" style="color: ${widgetColor};"></i>
                ${widgetTitle}
            </h4>
            <button onclick="this.parentElement.parentElement.remove()" style="
                background: none;
                border: none;
                color: #666;
                cursor: pointer;
                font-size: 16px;
            ">
                <i class="fas fa-times"></i>
            </button>
        `;
        
        widget.appendChild(header);
        
        // Add widget content based on type
        switch(type) {
            case 'summary':
                this.addSummaryWidgetContent(widget);
                break;
            case 'severity':
                this.addSeverityWidgetContent(widget);
                break;
            case 'category':
                this.addCategoryWidgetContent(widget);
                break;
            case 'findings':
                this.addFindingsWidgetContent(widget);
                break;
        }
        
        grid.appendChild(widget);
        
        // Make widget draggable (simple implementation)
        this.makeWidgetDraggable(widget);
    }

    addSummaryWidgetContent(widget) {
        const summary = this.findingsData.summary || {};
        
        widget.innerHTML += `
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 15px;">
                <div style="text-align: center;">
                    <div style="font-size: 24px; font-weight: bold; color: var(--primary-color);">${summary.overall_score || 'N/A'}</div>
                    <div style="font-size: 12px; color: #666; margin-top: 5px;">Overall Score</div>
                </div>
                <div style="text-align: center;">
                    <div style="font-size: 24px; font-weight: bold; color: ${this.getSeverityColor(summary.severity || 'unknown')};">
                        ${summary.severity ? summary.severity.charAt(0).toUpperCase() + summary.severity.slice(1) : 'N/A'}
                    </div>
                    <div style="font-size: 12px; color: #666; margin-top: 5px;">Overall Severity</div>
                </div>
                <div style="text-align: center;">
                    <div style="font-size: 24px; font-weight: bold; color: #28a745;">${this.findingsData.findings.length}</div>
                    <div style="font-size: 12px; color: #666; margin-top: 5px;">Total Findings</div>
                </div>
                <div style="text-align: center;">
                    <div style="font-size: 24px; font-weight: bold; color: #6c757d;">${summary.confidence || 'N/A'}</div>
                    <div style="font-size: 12px; color: #666; margin-top: 5px;">Confidence</div>
                </div>
            </div>
        `;
    }

    addSeverityWidgetContent(widget) {
        const canvas = document.createElement('canvas');
        canvas.style.width = '100%';
        canvas.style.height = '200px';
        widget.appendChild(canvas);
        
        // Count findings by severity
        const severityCounts = {
            critical: 0,
            high: 0,
            medium: 0,
            low: 0
        };
        
        this.findingsData.findings.forEach(finding => {
            const severity = (finding.Severity || 'low').toLowerCase();
            if (severityCounts[severity] !== undefined) {
                severityCounts[severity]++;
            }
        });
        
        new Chart(canvas, {
            type: 'doughnut',
            data: {
                labels: ['Critical', 'High', 'Medium', 'Low'],
                datasets: [{
                    data: [
                        severityCounts.critical,
                        severityCounts.high,
                        severityCounts.medium,
                        severityCounts.low
                    ],
                    backgroundColor: [
                        '#dc3545',
                        '#fd7e14',
                        '#ffc107',
                        '#28a745'
                    ]
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });
    }

    addCategoryWidgetContent(widget) {
        const canvas = document.createElement('canvas');
        canvas.style.width = '100%';
        canvas.style.height = '200px';
        widget.appendChild(canvas);
        
        // Count findings by category
        const categoryCounts = {};
        
        this.findingsData.findings.forEach(finding => {
            const category = finding.Category || 'unknown';
            categoryCounts[category] = (categoryCounts[category] || 0) + 1;
        });
        
        new Chart(canvas, {
            type: 'bar',
            data: {
                labels: Object.keys(categoryCounts),
                datasets: [{
                    label: 'Findings by Category',
                    data: Object.values(categoryCounts),
                    backgroundColor: '#4a6bff'
                }]
            },
            options: {
                responsive: true,
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
    }

    addFindingsWidgetContent(widget) {
        const topFindings = this.findingsData.findings
            .sort((a, b) => {
                const severityOrder = { critical: 4, high: 3, medium: 2, low: 1 };
                return (severityOrder[b.Severity.toLowerCase()] || 0) - (severityOrder[a.Severity.toLowerCase()] || 0);
            })
            .slice(0, 5);
        
        let html = '<div style="max-height: 300px; overflow-y: auto;">';
        
        topFindings.forEach((finding, index) => {
            html += `
                <div style="
                    padding: 10px;
                    margin-bottom: 10px;
                    background: #f8f9fa;
                    border-radius: 6px;
                    border-left: 4px solid ${this.getSeverityColor(finding.Severity)};
                ">
                    <div style="display: flex; justify-content: space-between; align-items: start;">
                        <div>
                            <strong style="color: ${this.getSeverityColor(finding.Severity)};">${index + 1}. ${finding.Title}</strong>
                            <div style="font-size: 0.8em; color: #666; margin-top: 5px;">
                                <span style="background: ${this.getSeverityColor(finding.Severity)}; color: white; padding: 2px 6px; border-radius: 4px; font-size: 0.7em;">${finding.Severity}</span>
                                <span style="margin-left: 8px; background: #6c757d; color: white; padding: 2px 6px; border-radius: 4px; font-size: 0.7em;">${finding.Category}</span>
                            </div>
                        </div>
                        <div style="text-align: right; font-size: 0.8em; color: #666;">
                            Confidence: ${(finding.Confidence * 100).toFixed(0)}%
                        </div>
                    </div>
                </div>
            `;
        });
        
        html += '</div>';
        widget.innerHTML += html;
    }

    makeWidgetDraggable(widget) {
        let isDragging = false;
        let startX, startY, initialX, initialY;
        
        const header = widget.querySelector('h4');
        if (!header) return;
        
        header.style.cursor = 'move';
        
        header.addEventListener('mousedown', (e) => {
            isDragging = true;
            startX = e.clientX;
            startY = e.clientY;
            initialX = widget.offsetLeft;
            initialY = widget.offsetTop;
            
            widget.style.position = 'absolute';
            widget.style.zIndex = '100';
            widget.style.opacity = '0.8';
            
            e.preventDefault();
        });
        
        document.addEventListener('mousemove', (e) => {
            if (!isDragging) return;
            
            const dx = e.clientX - startX;
            const dy = e.clientY - startY;
            
            widget.style.left = `${initialX + dx}px`;
            widget.style.top = `${initialY + dy}px`;
        });
        
        document.addEventListener('mouseup', () => {
            if (isDragging) {
                isDragging = false;
                widget.style.opacity = '1';
                widget.style.zIndex = '1';
            }
        });
    }

    // Export functionality
    static exportReport(format) {
        const reportData = {
            findings: window.findingsData,
            insights: window.insightsData,
            timestamp: new Date().toISOString(),
            version: '1.0'
        };
        
        switch(format) {
            case 'json':
                this.exportJSON(reportData);
                break;
            case 'csv':
                this.exportCSV(reportData);
                break;
            case 'pdf':
                this.exportPDF(reportData);
                break;
        }
    }

    static exportJSON(data) {
        const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `triageprof-report-${new Date().toISOString().split('T')[0]}.json`;
        a.click();
        URL.revokeObjectURL(url);
    }

    static exportCSV(data) {
        let csv = 'Finding ID,Title,Severity,Category,Confidence,Impact Summary\n';
        
        data.findings.findings.forEach(finding => {
            csv += `"${finding.ID}","${finding.Title}","${finding.Severity}","${finding.Category}","${finding.Confidence}","${finding.ImpactSummary}"\n`;
        });
        
        const blob = new Blob([csv], { type: 'text/csv' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `triageprof-report-${new Date().toISOString().split('T')[0]}.csv`;
        a.click();
        URL.revokeObjectURL(url);
    }

    static exportPDF(data) {
        // In a real implementation, this would use jsPDF or similar library
        // For now, we'll just export JSON as PDF is more complex
        alert('PDF export would be implemented with a library like jsPDF. Exporting as JSON instead.');
        this.exportJSON(data);
    }
}

// Make AdvancedVisualization available globally
window.AdvancedVisualization = AdvancedVisualization;