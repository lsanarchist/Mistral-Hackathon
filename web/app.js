document.addEventListener('DOMContentLoaded', function() {
    const loadBtn = document.getElementById('loadBtn');
    const fileInput = document.getElementById('fileInput');
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const content = document.getElementById('content');
    const severityFilter = document.getElementById('severityFilter');

    let findingsData = null;
    let insightsData = null;
    let allFindings = [];

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
        
        // Show loading state
        loading.style.display = 'block';
        error.style.display = 'none';
        content.style.display = 'none';
        
        // Process files
        processFiles(files);
    });

    // Handle severity filter changes
    severityFilter.addEventListener('change', function() {
        filterFindings();
    });

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
            
            // Render the data
            renderData();
        });
    }

    function showError(message) {
        loading.style.display = 'none';
        error.textContent = message;
        error.style.display = 'block';
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
                document.getElementById('confidence').textContent = insightsData.ExecutiveSummary.Confidence + '%';
                document.getElementById('overview').innerHTML = 
                    `<p><strong>Overview:</strong> ${insightsData.ExecutiveSummary.Overview || 'No overview available'}</p>` +
                    `<p><strong>Severity:</strong> ${insightsData.ExecutiveSummary.OverallSeverity || 'Unknown'} (${insightsData.ExecutiveSummary.Confidence || 0}% confidence)</p>`;
            } else {
                document.getElementById('confidence').textContent = 'N/A';
                document.getElementById('overview').innerHTML = '<p>No LLM insights available</p>';
            }
            
            // Render charts
            renderSeverityChart();
            renderCategoryChart();
            
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
            
            card.innerHTML = `
                <div class="card-title">${index + 1}. ${risk.Description}</div>
                <div class="card-detail"><strong>Severity:</strong> ${risk.Severity}</div>
                <div class="card-detail"><strong>Impact:</strong> ${risk.Impact}</div>
                <div class="card-detail"><strong>Likelihood:</strong> ${risk.Likelihood}</div>
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
            
            card.innerHTML = `
                <div class="card-title">${index + 1}. ${action.Description}</div>
                <div class="card-detail"><strong>Priority:</strong> ${action.Priority}</div>
                <div class="card-detail"><strong>Estimated Effort:</strong> ${action.EstimatedEffort}</div>
                ${action.Categories && action.Categories.length > 0 ? 
                    `<div class="card-detail"><strong>Categories:</strong> ${action.Categories.join(', ')}</div>` : ''}
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
                
                const insightsTitle = document.createElement('div');
                insightsTitle.className = 'insights-title';
                insightsTitle.innerHTML = '<i class="fas fa-robot"></i> LLM Insights';
                
                insightsSection.appendChild(insightsTitle);
                
                if (findingInsight.Narrative) {
                    const narrative = document.createElement('p');
                    narrative.innerHTML = `<strong>Narrative:</strong> ${findingInsight.Narrative}`;
                    insightsSection.appendChild(narrative);
                }
                
                if (findingInsight.LikelyRootCauses && findingInsight.LikelyRootCauses.length > 0) {
                    const rootCausesTitle = document.createElement('div');
                    rootCausesTitle.innerHTML = '<strong><i class="fas fa-search"></i> Likely Root Causes:</strong>';
                    insightsSection.appendChild(rootCausesTitle);
                    
                    const rootCausesList = document.createElement('ul');
                    rootCausesList.className = 'insights-list';
                    
                    findingInsight.LikelyRootCauses.forEach(cause => {
                        const li = document.createElement('li');
                        li.textContent = cause;
                        rootCausesList.appendChild(li);
                    });
                    
                    insightsSection.appendChild(rootCausesList);
                }
                
                if (findingInsight.Suggestions && findingInsight.Suggestions.length > 0) {
                    const suggestionsTitle = document.createElement('div');
                    suggestionsTitle.innerHTML = '<strong><i class="fas fa-lightbulb"></i> Suggestions:</strong>';
                    insightsSection.appendChild(suggestionsTitle);
                    
                    const suggestionsList = document.createElement('ul');
                    suggestionsList.className = 'insights-list';
                    
                    findingInsight.Suggestions.forEach(suggestion => {
                        const li = document.createElement('li');
                        li.textContent = suggestion;
                        suggestionsList.appendChild(li);
                    });
                    
                    insightsSection.appendChild(suggestionsList);
                }
                
                if (findingInsight.Confidence) {
                    const confidence = document.createElement('p');
                    confidence.innerHTML = `<strong>Confidence:</strong> ${findingInsight.Confidence}%`;
                    insightsSection.appendChild(confidence);
                }
                
                findingCard.appendChild(insightsSection);
            }
        }
        
        // Build the finding card
        findingCard.appendChild(findingHeader);
        findingCard.appendChild(findingDetails);
        
        findingsList.appendChild(findingCard);
    }
});