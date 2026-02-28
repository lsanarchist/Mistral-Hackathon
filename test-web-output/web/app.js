document.addEventListener('DOMContentLoaded', function() {
    const loadBtn = document.getElementById('loadBtn');
    const fileInput = document.getElementById('fileInput');
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const content = document.getElementById('content');
    
    let findingsData = null;
    let insightsData = null;
    
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
            document.getElementById('overallScore').textContent = findingsData.Summary.OverallScore || 'N/A';
            
            // Determine severity
            const score = findingsData.Summary.OverallScore || 0;
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
        
        findingsData.Findings.forEach(finding => {
            const findingCard = document.createElement('div');
            findingCard.className = 'finding-card';
            
            // Determine severity class
            const severityClass = `severity-${finding.Severity.toLowerCase()}`;
            
            // Create finding header
            const findingHeader = document.createElement('div');
            findingHeader.className = 'finding-header';
            
            const findingTitle = document.createElement('div');
            findingTitle.className = 'finding-title';
            findingTitle.textContent = finding.Title;
            
            const severityBadge = document.createElement('div');
            severityBadge.className = `severity-badge ${severityClass}`;
            severityBadge.textContent = finding.Severity;
            
            findingHeader.appendChild(findingTitle);
            findingHeader.appendChild(severityBadge);
            
            // Create finding details
            const findingDetails = document.createElement('div');
            findingDetails.className = 'finding-details';
            
            const details = [
                { label: 'Category', value: finding.Category },
                { label: 'Score', value: finding.Score },
                { label: 'Profile Type', value: finding.Evidence.ProfileType },
                { label: 'Artifact', value: finding.Evidence.ArtifactPath }
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
            if (finding.Top && finding.Top.length > 0) {
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
                finding.Top.forEach(frame => {
                    const row = document.createElement('tr');
                    
                    [frame.Function, frame.File, frame.Line, frame.Cum.toFixed(2), frame.Flat.toFixed(2)].forEach(cellData => {
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
                const findingInsight = insightsData.PerFinding.find(i => i.FindingID === finding.Category);
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
        });
    }
});