// Plugin Management UI
function initPluginManagement() {
    const pluginManagementSection = document.createElement('div');
    pluginManagementSection.className = 'plugin-management-section';
    pluginManagementSection.id = 'pluginManagementSection';
    
    pluginManagementSection.innerHTML = `
        <div class="plugin-management-header">
            <h2><i class="fas fa-puzzle-piece"></i> Plugin Management</h2>
            <div class="plugin-management-controls">
                <button id="refreshPluginsBtn" class="btn info"><i class="fas fa-sync-alt"></i> Refresh</button>
                <button id="checkHealthBtn" class="btn success"><i class="fas fa-heartbeat"></i> Check Health</button>
            </div>
        </div>
        
        <div class="plugin-management-tabs">
            <button class="tab-btn active" data-tab="plugins-list">Plugin List</button>
            <button class="tab-btn" data-tab="capabilities-matrix">Capabilities Matrix</button>
            <button class="tab-btn" data-tab="health-monitoring">Health Monitoring</button>
            <button class="tab-btn" data-tab="plugin-marketplace">Marketplace</button>
        </div>
        
        <div class="plugin-management-content">
            <!-- Plugin List Tab -->
            <div id="plugins-list" class="plugin-tab-content active">
                <div class="plugin-list-header">
                    <h3><i class="fas fa-list"></i> Available Plugins</h3>
                    <div class="plugin-search">
                        <input type="text" id="pluginSearch" placeholder="Search plugins...">
                        <i class="fas fa-search"></i>
                    </div>
                </div>
                <div class="plugin-list" id="pluginList"></div>
            </div>
            
            <!-- Capabilities Matrix Tab -->
            <div id="capabilities-matrix" class="plugin-tab-content">
                <div class="capabilities-header">
                    <h3><i class="fas fa-table"></i> Plugin Capabilities Matrix</h3>
                    <div class="capability-filters">
                        <label><input type="checkbox" id="filterTargets" checked> Targets</label>
                        <label><input type="checkbox" id="filterProfiles" checked> Profiles</label>
                    </div>
                </div>
                <div class="capabilities-matrix" id="capabilitiesMatrix"></div>
            </div>
            
            <!-- Health Monitoring Tab -->
            <div id="health-monitoring" class="plugin-tab-content">
                <div class="health-header">
                    <h3><i class="fas fa-heartbeat"></i> Plugin Health Monitoring</h3>
                    <div class="health-stats">
                        <span id="healthyCount">0 Healthy</span> | 
                        <span id="unhealthyCount">0 Unhealthy</span>
                    </div>
                </div>
                <div class="health-monitoring" id="healthMonitoring"></div>
            </div>
            
            <!-- Plugin Marketplace Tab -->
            <div id="plugin-marketplace" class="plugin-tab-content">
                <div class="marketplace-header">
                    <h3><i class="fas fa-store"></i> Plugin Marketplace</h3>
                    <div class="marketplace-controls">
                        <button id="installPluginBtn" class="btn success"><i class="fas fa-download"></i> Install Plugin</button>
                        <button id="updateAllBtn" class="btn info"><i class="fas fa-sync-alt"></i> Update All</button>
                    </div>
                </div>
                <div class="marketplace-content" id="pluginMarketplace"></div>
            </div>
        </div>
    `;
    
    // Insert the plugin management section after the summary section
    const summarySection = document.querySelector('.summary-card');
    if (summarySection) {
        summarySection.parentNode.insertBefore(pluginManagementSection, summarySection.nextSibling);
    } else {
        document.getElementById('content').prepend(pluginManagementSection);
    }
    
    // Set up tab switching
    const tabButtons = document.querySelectorAll('.plugin-management-tabs .tab-btn');
    const tabContents = document.querySelectorAll('.plugin-tab-content');
    
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabName = button.getAttribute('data-tab');
            
            // Update active tab button
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
            
            // Update active tab content
            tabContents.forEach(content => content.classList.remove('active'));
            document.getElementById(tabName).classList.add('active');
            
            // Load data for the specific tab
            if (tabName === 'plugins-list') {
                loadPluginList();
            } else if (tabName === 'capabilities-matrix') {
                loadCapabilitiesMatrix();
            } else if (tabName === 'health-monitoring') {
                loadHealthMonitoring();
            } else if (tabName === 'plugin-marketplace') {
                loadPluginMarketplace();
            }
        });
    });
    
    // Set up control buttons
    document.getElementById('refreshPluginsBtn').addEventListener('click', () => {
        loadPluginList();
        showNotification('Plugin list refreshed!', 'success');
    });
    
    document.getElementById('checkHealthBtn').addEventListener('click', () => {
        checkAllPluginHealth();
        showNotification('Checking plugin health...', 'info');
    });
    
    document.getElementById('installPluginBtn').addEventListener('click', () => {
        showInstallPluginDialog();
    });
    
    document.getElementById('updateAllBtn').addEventListener('click', () => {
        updateAllPlugins();
    });
    
    // Set up search functionality
    document.getElementById('pluginSearch').addEventListener('input', (e) => {
        filterPlugins(e.target.value);
    });
    
    // Load initial data
    loadPluginList();
}

// Load plugin list from server
function loadPluginList() {
    const pluginListElement = document.getElementById('pluginList');
    pluginListElement.innerHTML = '<div class="loading-plugins"><i class="fas fa-spinner fa-spin"></i> Loading plugins...</div>';
    
    fetch('/plugins')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to load plugins');
            }
            return response.json();
        })
        .then(data => {
            displayPluginList(data.plugins);
        })
        .catch(error => {
            console.error('Error loading plugins:', error);
            pluginListElement.innerHTML = '<div class="error-message">Failed to load plugins: ' + error.message + '</div>';
        });
}

// Display plugin list
function displayPluginList(plugins) {
    const pluginListElement = document.getElementById('pluginList');
    pluginListElement.innerHTML = '';
    
    if (plugins.length === 0) {
        pluginListElement.innerHTML = '<div class="no-plugins">No plugins found. Make sure plugins are installed in the plugins/ directory.</div>';
        return;
    }
    
    plugins.forEach(plugin => {
        const pluginElement = document.createElement('div');
        pluginElement.className = 'plugin-card';
        
        const healthStatus = plugin.health.status || 'unknown';
        const healthClass = `health-status ${healthStatus}`;
        
        pluginElement.innerHTML = `
            <div class="plugin-card-header">
                <div class="plugin-card-title">
                    <h4>${plugin.name}</h4>
                    <span class="plugin-version">v${plugin.version}</span>
                </div>
                <div class="${healthClass}">
                    <i class="fas fa-${getHealthIcon(healthStatus)}"></i> ${healthStatus}
                </div>
            </div>
            <div class="plugin-card-body">
                <p class="plugin-description">${plugin.description || 'No description available'}</p>
                <div class="plugin-meta">
                    <span class="plugin-author">${plugin.author || 'Unknown'}</span>
                    <span class="plugin-sdk">SDK: ${plugin.sdkVersion}</span>
                </div>
                <div class="plugin-capabilities">
                    <div class="capability-item">
                        <strong>Targets:</strong> ${plugin.capabilities.targets.join(', ')}
                    </div>
                    <div class="capability-item">
                        <strong>Profiles:</strong> ${plugin.capabilities.profiles.join(', ')}
                    </div>
                </div>
            </div>
            <div class="plugin-card-footer">
                <button class="btn info plugin-details-btn" data-plugin="${plugin.name}">
                    <i class="fas fa-info-circle"></i> Details
                </button>
                <button class="btn success plugin-check-health-btn" data-plugin="${plugin.name}">
                    <i class="fas fa-heartbeat"></i> Check Health
                </button>
                <button class="btn warning plugin-update-btn" data-plugin="${plugin.name}">
                    <i class="fas fa-sync-alt"></i> Update
                </button>
            </div>
        `;
        
        pluginListElement.appendChild(pluginElement);
        
        // Add event listeners for buttons
        pluginElement.querySelector('.plugin-details-btn').addEventListener('click', () => {
            showPluginDetails(plugin);
        });
        
        pluginElement.querySelector('.plugin-check-health-btn').addEventListener('click', () => {
            checkPluginHealth(plugin.name);
        });
        
        pluginElement.querySelector('.plugin-update-btn').addEventListener('click', () => {
            updatePlugin(plugin.name);
        });
    });
}

// Get health icon based on status
function getHealthIcon(status) {
    switch (status.toLowerCase()) {
        case 'healthy': return 'check-circle';
        case 'unhealthy': return 'exclamation-triangle';
        case 'warning': return 'exclamation-circle';
        default: return 'question-circle';
    }
}

// Show plugin details dialog
function showPluginDetails(plugin) {
    const dialog = document.createElement('div');
    dialog.className = 'plugin-details-dialog';
    
    const healthStatus = plugin.health.status || 'unknown';
    const healthClass = `health-status ${healthStatus}`;
    
    dialog.innerHTML = `
        <div class="dialog-content">
            <div class="dialog-header">
                <h3>${plugin.name} Details</h3>
                <button class="close-dialog"><i class="fas fa-times"></i></button>
            </div>
            <div class="dialog-body">
                <div class="plugin-detail-section">
                    <div class="detail-label">Version:</div>
                    <div class="detail-value">${plugin.version}</div>
                </div>
                <div class="plugin-detail-section">
                    <div class="detail-label">SDK Version:</div>
                    <div class="detail-value">${plugin.sdkVersion}</div>
                </div>
                <div class="plugin-detail-section">
                    <div class="detail-label">Author:</div>
                    <div class="detail-value">${plugin.author || 'Unknown'}</div>
                </div>
                <div class="plugin-detail-section">
                    <div class="detail-label">Health Status:</div>
                    <div class="detail-value">
                        <span class="${healthClass}">
                            <i class="fas fa-${getHealthIcon(healthStatus)}"></i> ${healthStatus}
                        </span>
                    </div>
                </div>
                ${plugin.health.error ? `
                <div class="plugin-detail-section">
                    <div class="detail-label">Health Error:</div>
                    <div class="detail-value error-text">${plugin.health.error}</div>
                </div>
                ` : ''}
                ${plugin.health.binaryPath ? `
                <div class="plugin-detail-section">
                    <div class="detail-label">Binary Path:</div>
                    <div class="detail-value">${plugin.health.binaryPath}</div>
                </div>
                ` : ''}
                <div class="plugin-detail-section">
                    <div class="detail-label">Description:</div>
                    <div class="detail-value">${plugin.description || 'No description available'}</div>
                </div>
                <div class="plugin-detail-section">
                    <div class="detail-label">Supported Targets:</div>
                    <div class="detail-value">${plugin.capabilities.targets.join(', ')}</div>
                </div>
                <div class="plugin-detail-section">
                    <div class="detail-label">Supported Profiles:</div>
                    <div class="detail-value">${plugin.capabilities.profiles.join(', ')}</div>
                </div>
            </div>
            <div class="dialog-footer">
                <button class="btn secondary close-dialog">Close</button>
            </div>
        </div>
    `;
    
    document.body.appendChild(dialog);
    
    // Close dialog
    dialog.querySelectorAll('.close-dialog').forEach(btn => {
        btn.addEventListener('click', () => {
            dialog.remove();
        });
    });
    
    // Close when clicking outside
    dialog.addEventListener('click', (e) => {
        if (e.target === dialog) {
            dialog.remove();
        }
    });
}

// Check health for a specific plugin
function checkPluginHealth(pluginName) {
    fetch('/plugins/health')
        .then(response => response.json())
        .then(data => {
            const health = data.health[pluginName];
            if (health) {
                showNotification(`Plugin ${pluginName} health: ${health.status}`, 'info');
                // Refresh the plugin list to show updated health
                loadPluginList();
            }
        })
        .catch(error => {
            console.error('Error checking plugin health:', error);
            showNotification('Failed to check plugin health: ' + error.message, 'error');
        });
}

// Check health for all plugins
function checkAllPluginHealth() {
    fetch('/plugins/health')
        .then(response => response.json())
        .then(data => {
            let healthyCount = 0;
            let unhealthyCount = 0;
            
            for (const pluginName in data.health) {
                const health = data.health[pluginName];
                if (health.status === 'healthy') {
                    healthyCount++;
                } else {
                    unhealthyCount++;
                }
            }
            
            showNotification(`Health check complete: ${healthyCount} healthy, ${unhealthyCount} unhealthy`, 'success');
            
            // Update health monitoring tab
            if (document.getElementById('health-monitoring').classList.contains('active')) {
                loadHealthMonitoring();
            }
        })
        .catch(error => {
            console.error('Error checking plugin health:', error);
            showNotification('Failed to check plugin health: ' + error.message, 'error');
        });
}

// Load capabilities matrix
function loadCapabilitiesMatrix() {
    const matrixElement = document.getElementById('capabilitiesMatrix');
    matrixElement.innerHTML = '<div class="loading-matrix"><i class="fas fa-spinner fa-spin"></i> Loading capabilities...</div>';
    
    fetch('/plugins/capabilities')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to load capabilities');
            }
            return response.json();
        })
        .then(data => {
            displayCapabilitiesMatrix(data);
        })
        .catch(error => {
            console.error('Error loading capabilities:', error);
            matrixElement.innerHTML = '<div class="error-message">Failed to load capabilities: ' + error.message + '</div>';
        });
}

// Display capabilities matrix
function displayCapabilitiesMatrix(data) {
    const matrixElement = document.getElementById('capabilitiesMatrix');
    matrixElement.innerHTML = '';
    
    // Create matrix table
    const table = document.createElement('table');
    table.className = 'capabilities-table';
    
    // Create header row
    const headerRow = document.createElement('tr');
    const headerCell = document.createElement('th');
    headerCell.textContent = 'Plugin';
    headerRow.appendChild(headerCell);
    
    // Add target columns if enabled
    if (document.getElementById('filterTargets').checked) {
        data.targets.forEach(target => {
            const th = document.createElement('th');
            th.textContent = target;
            th.title = 'Target: ' + target;
            headerRow.appendChild(th);
        });
    }
    
    // Add profile columns if enabled
    if (document.getElementById('filterProfiles').checked) {
        data.profiles.forEach(profile => {
            const th = document.createElement('th');
            th.textContent = profile;
            th.title = 'Profile: ' + profile;
            headerRow.appendChild(th);
        });
    }
    
    table.appendChild(headerRow);
    
    // Create data rows
    data.matrix.forEach(pluginData => {
        const row = document.createElement('tr');
        
        // Plugin name cell
        const pluginCell = document.createElement('td');
        pluginCell.textContent = pluginData.plugin;
        pluginCell.className = 'plugin-name';
        row.appendChild(pluginCell);
        
        // Target capability cells
        if (document.getElementById('filterTargets').checked) {
            data.targets.forEach(target => {
                const cell = document.createElement('td');
                if (pluginData.targets[target]) {
                    cell.innerHTML = '<i class="fas fa-check-circle" style="color: #4CAF50;"></i>';
                    cell.title = 'Supports ' + target;
                } else {
                    cell.innerHTML = '<i class="fas fa-times-circle" style="color: #F44336;"></i>';
                    cell.title = 'Does not support ' + target;
                }
                cell.className = 'capability-cell';
                row.appendChild(cell);
            });
        }
        
        // Profile capability cells
        if (document.getElementById('filterProfiles').checked) {
            data.profiles.forEach(profile => {
                const cell = document.createElement('td');
                if (pluginData.profiles[profile]) {
                    cell.innerHTML = '<i class="fas fa-check-circle" style="color: #4CAF50;"></i>';
                    cell.title = 'Supports ' + profile;
                } else {
                    cell.innerHTML = '<i class="fas fa-times-circle" style="color: #F44336;"></i>';
                    cell.title = 'Does not support ' + profile;
                }
                cell.className = 'capability-cell';
                row.appendChild(cell);
            });
        }
        
        table.appendChild(row);
    });
    
    matrixElement.appendChild(table);
    
    // Add filter change listeners
    document.getElementById('filterTargets').addEventListener('change', () => {
        loadCapabilitiesMatrix();
    });
    
    document.getElementById('filterProfiles').addEventListener('change', () => {
        loadCapabilitiesMatrix();
    });
}

// Load health monitoring
function loadHealthMonitoring() {
    const healthElement = document.getElementById('healthMonitoring');
    healthElement.innerHTML = '<div class="loading-health"><i class="fas fa-spinner fa-spin"></i> Loading health data...</div>';
    
    fetch('/plugins/health')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to load health data');
            }
            return response.json();
        })
        .then(data => {
            displayHealthMonitoring(data.health);
        })
        .catch(error => {
            console.error('Error loading health data:', error);
            healthElement.innerHTML = '<div class="error-message">Failed to load health data: ' + error.message + '</div>';
        });
}

// Display health monitoring
function displayHealthMonitoring(healthData) {
    const healthElement = document.getElementById('healthMonitoring');
    healthElement.innerHTML = '';
    
    let healthyCount = 0;
    let unhealthyCount = 0;
    
    for (const pluginName in healthData) {
        const health = healthData[pluginName];
        
        const pluginHealthElement = document.createElement('div');
        pluginHealthElement.className = 'plugin-health-item';
        
        const statusClass = `health-status ${health.status || 'unknown'}`;
        
        pluginHealthElement.innerHTML = `
            <div class="plugin-health-header">
                <h4>${pluginName}</h4>
                <div class="${statusClass}">
                    <i class="fas fa-${getHealthIcon(health.status)}"></i> ${health.status}
                </div>
            </div>
            <div class="plugin-health-details">
                <div class="health-detail">
                    <strong>Last Checked:</strong> ${health.LastChecked ? new Date(health.LastChecked).toLocaleString() : 'Never'}
                </div>
                ${health.Error ? `
                <div class="health-detail error">
                    <strong>Error:</strong> ${health.Error}
                </div>
                ` : ''}
                ${health.BinaryPath ? `
                <div class="health-detail">
                    <strong>Binary:</strong> ${health.BinaryPath}
                </div>
                ` : ''}
            </div>
            <div class="plugin-health-actions">
                <button class="btn info check-health-btn" data-plugin="${pluginName}">
                    <i class="fas fa-sync-alt"></i> Recheck
                </button>
            </div>
        `;
        
        healthElement.appendChild(pluginHealthElement);
        
        // Update counters
        if (health.status === 'healthy') {
            healthyCount++;
        } else {
            unhealthyCount++;
        }
        
        // Add event listener for recheck button
        pluginHealthElement.querySelector('.check-health-btn').addEventListener('click', () => {
            checkPluginHealth(pluginName);
        });
    }
    
    // Update health stats
    document.getElementById('healthyCount').textContent = `${healthyCount} Healthy`;
    document.getElementById('unhealthyCount').textContent = `${unhealthyCount} Unhealthy`;
}

// Load plugin marketplace
function loadPluginMarketplace() {
    const marketplaceElement = document.getElementById('pluginMarketplace');
    marketplaceElement.innerHTML = '<div class="loading-marketplace"><i class="fas fa-spinner fa-spin"></i> Loading marketplace...</div>';
    
    fetch('/plugins/marketplace')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to load marketplace');
            }
            return response.json();
        })
        .then(data => {
            displayPluginMarketplace(data.plugins);
        })
        .catch(error => {
            console.error('Error loading marketplace:', error);
            marketplaceElement.innerHTML = '<div class="error-message">Failed to load marketplace: ' + error.message + '</div>';
        });
}

// Display plugin marketplace
function displayPluginMarketplace(plugins) {
    const marketplaceElement = document.getElementById('pluginMarketplace');
    marketplaceElement.innerHTML = '';
    
    if (!plugins || plugins.length === 0) {
        marketplaceElement.innerHTML = '<div class="no-plugins">No plugins available in the marketplace.</div>';
        return;
    }
    
    plugins.forEach(plugin => {
        const pluginElement = document.createElement('div');
        pluginElement.className = 'marketplace-plugin-card';
        
        const installedBadge = plugin.installed ? 
            '<span class="installed-badge"><i class="fas fa-check-circle"></i> Installed</span>' :
            '<span class="not-installed-badge"><i class="fas fa-download"></i> Available</span>';
        
        pluginElement.innerHTML = `
            <div class="marketplace-plugin-header">
                <h4>${plugin.name}</h4>
                ${installedBadge}
                <span class="marketplace-plugin-version">v${plugin.version}</span>
            </div>
            <div class="marketplace-plugin-body">
                <p class="marketplace-plugin-description">${plugin.description}</p>
                <div class="marketplace-plugin-meta">
                    <span class="marketplace-plugin-author">${plugin.author}</span>
                </div>
                <div class="marketplace-plugin-capabilities">
                    <div class="capability-badge">
                        <strong>Targets:</strong> ${plugin.capabilities.targets.join(', ')}
                    </div>
                    <div class="capability-badge">
                        <strong>Profiles:</strong> ${plugin.capabilities.profiles.join(', ')}
                    </div>
                </div>
            </div>
            <div class="marketplace-plugin-footer">
                ${plugin.installed ? `
                    <button class="btn warning marketplace-update-btn" data-plugin="${plugin.name}">
                        <i class="fas fa-sync-alt"></i> Update
                    </button>
                    <button class="btn danger marketplace-uninstall-btn" data-plugin="${plugin.name}">
                        <i class="fas fa-trash"></i> Uninstall
                    </button>
                ` : `
                    <button class="btn success marketplace-install-btn" data-plugin="${plugin.name}">
                        <i class="fas fa-download"></i> Install
                    </button>
                `}
            </div>
        `;
        
        marketplaceElement.appendChild(pluginElement);
        
        // Add event listeners for buttons
        if (plugin.installed) {
            pluginElement.querySelector('.marketplace-update-btn').addEventListener('click', () => {
                updatePlugin(plugin.name);
            });
            
            pluginElement.querySelector('.marketplace-uninstall-btn').addEventListener('click', () => {
                uninstallPlugin(plugin.name);
            });
        } else {
            pluginElement.querySelector('.marketplace-install-btn').addEventListener('click', () => {
                installPlugin(plugin.name);
            });
        }
    });
}

// Show install plugin dialog
function showInstallPluginDialog() {
    const dialog = document.createElement('div');
    dialog.className = 'install-plugin-dialog';
    
    dialog.innerHTML = `
        <div class="dialog-content">
            <div class="dialog-header">
                <h3>Install Plugin</h3>
                <button class="close-dialog"><i class="fas fa-times"></i></button>
            </div>
            <div class="dialog-body">
                <div class="install-option">
                    <h4>Option 1: Install from Marketplace</h4>
                    <p>Browse available plugins in the marketplace tab and click Install.</p>
                </div>
                <div class="install-option">
                    <h4>Option 2: Install from URL</h4>
                    <div class="install-url-form">
                        <input type="text" id="pluginUrlInput" placeholder="https://example.com/plugin.zip">
                        <button id="installFromUrlBtn" class="btn success">
                            <i class="fas fa-download"></i> Install
                        </button>
                    </div>
                </div>
                <div class="install-option">
                    <h4>Option 3: Manual Installation</h4>
                    <p>Place plugin files in the <code>plugins/</code> directory and restart the server.</p>
                </div>
            </div>
            <div class="dialog-footer">
                <button class="btn secondary close-dialog">Close</button>
            </div>
        </div>
    `;
    
    document.body.appendChild(dialog);
    
    // Set up install from URL
    document.getElementById('installFromUrlBtn').addEventListener('click', () => {
        const url = document.getElementById('pluginUrlInput').value.trim();
        if (url) {
            installPluginFromUrl(url);
            dialog.remove();
        } else {
            showNotification('Please enter a valid URL', 'error');
        }
    });
    
    // Close dialog
    dialog.querySelectorAll('.close-dialog').forEach(btn => {
        btn.addEventListener('click', () => {
            dialog.remove();
        });
    });
    
    // Close when clicking outside
    dialog.addEventListener('click', (e) => {
        if (e.target === dialog) {
            dialog.remove();
        }
    });
}

// Install plugin from URL
function installPluginFromUrl(url) {
    showNotification(`Installing plugin from ${url}...`, 'info');
    
    fetch('/plugins/install', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ url: url })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to install plugin');
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            showNotification(data.message, 'success');
            loadPluginMarketplace();
        } else {
            showNotification('Failed to install plugin', 'error');
        }
    })
    .catch(error => {
        console.error('Error installing plugin:', error);
        showNotification('Failed to install plugin: ' + error.message, 'error');
    });
}

// Install plugin
function installPlugin(pluginName) {
    showNotification(`Installing plugin ${pluginName}...`, 'info');
    
    fetch('/plugins/install', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ pluginName: pluginName })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to install plugin');
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            showNotification(data.message, 'success');
            loadPluginMarketplace();
            loadPluginList();
        } else {
            showNotification('Failed to install plugin', 'error');
        }
    })
    .catch(error => {
        console.error('Error installing plugin:', error);
        showNotification('Failed to install plugin: ' + error.message, 'error');
    });
}

// Update plugin
function updatePlugin(pluginName) {
    showNotification(`Updating plugin ${pluginName}...`, 'info');
    
    fetch('/plugins/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ pluginName: pluginName })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to update plugin');
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            showNotification(data.message, 'success');
            loadPluginMarketplace();
            loadPluginList();
        } else {
            showNotification('Failed to update plugin', 'error');
        }
    })
    .catch(error => {
        console.error('Error updating plugin:', error);
        showNotification('Failed to update plugin: ' + error.message, 'error');
    });
}

// Update all plugins
function updateAllPlugins() {
    showNotification('Checking for plugin updates...', 'info');
    
    // In a real implementation, this would call a server endpoint
    setTimeout(() => {
        showNotification('All plugins are up to date!', 'success');
    }, 1500);
}

// Uninstall plugin
function uninstallPlugin(pluginName) {
    if (confirm(`Are you sure you want to uninstall ${pluginName}?`)) {
        showNotification(`Uninstalling plugin ${pluginName}...`, 'info');
        
        fetch('/plugins/uninstall', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ pluginName: pluginName })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to uninstall plugin');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                showNotification(data.message, 'success');
                loadPluginMarketplace();
                loadPluginList();
            } else {
                showNotification('Failed to uninstall plugin', 'error');
            }
        })
        .catch(error => {
            console.error('Error uninstalling plugin:', error);
            showNotification('Failed to uninstall plugin: ' + error.message, 'error');
        });
    }
}

// Filter plugins based on search query
function filterPlugins(query) {
    const pluginCards = document.querySelectorAll('.plugin-card');
    const searchTerm = query.toLowerCase();
    
    pluginCards.forEach(card => {
        const pluginName = card.querySelector('.plugin-card-title h4').textContent.toLowerCase();
        const pluginDescription = card.querySelector('.plugin-description').textContent.toLowerCase();
        
        if (pluginName.includes(searchTerm) || pluginDescription.includes(searchTerm)) {
            card.style.display = 'block';
        } else {
            card.style.display = 'none';
        }
    });
}

// Initialize plugin management when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    // Check if we should initialize plugin management
    if (document.getElementById('content')) {
        // Add a small delay to ensure other UI elements are loaded
        setTimeout(initPluginManagement, 500);
    }
});