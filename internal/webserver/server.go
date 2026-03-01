package webserver

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/mistral-hackathon/triageprof/internal/model"
	"github.com/mistral-hackathon/triageprof/internal/plugin"
)

// WebSocketServer handles real-time data streaming
type WebSocketServer struct {
	server          *http.Server
	upgrader        websocket.Upgrader
	clients         map[*websocket.Conn]bool
	clientsMu       sync.Mutex
	dataDir         string
	findings        *model.FindingsBundle
	insights        *model.InsightsBundle
	lastUpdate      time.Time
	pluginDir       string
	pluginManifests []*plugin.Manifest
	pluginHealth    map[string]PluginHealth
	pluginManager   *plugin.PluginManager
	authEnabled     bool
	jwtSecretKey    string
	compressionEnabled bool
	compressionLevel int
	compressionThreshold int
	performanceHistory []PerformanceSnapshot
	historyMu       sync.Mutex
	maxHistorySize  int
}

// PerformanceSnapshot represents a historical performance data point
type PerformanceSnapshot struct {
	Timestamp       time.Time `json:"timestamp"`
	OverallScore    int       `json:"overall_score"`
	CriticalCount   int       `json:"critical_count"`
	HighCount       int       `json:"high_count"`
	MediumCount     int       `json:"medium_count"`
	LowCount        int       `json:"low_count"`
	TotalFindings   int       `json:"total_findings"`
	ClientCount     int       `json:"client_count"`
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type PluginHealth struct {
	Status      string `json:"status"`
	LastChecked time.Time `json:"lastChecked"`
	Error       string `json:"error,omitempty"`
	BinaryPath  string `json:"binaryPath,omitempty"`
}

// NewWebSocketServer creates a new WebSocket server instance
func NewWebSocketServer(port int, dataDir string, pluginDir string, enableAuth bool, enableCompression bool) *WebSocketServer {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for demo purposes
		},
		EnableCompression: enableCompression,
	}

	// Configure compression settings if enabled
	compressionLevel := 0
	compressionThreshold := 0
	if enableCompression {
		// Use optimal compression level for performance data
		compressionLevel = 6 // Balanced compression level (0-9, where 9 is max compression)
		compressionThreshold = 256 // Compress messages larger than 256 bytes
	}

	// Generate JWT secret key if auth is enabled
	jwtSecretKey := ""
	if enableAuth {
		jwtSecretKey = generateJWTSecretKey()
	}

	s := &WebSocketServer{
		server:          server,
		upgrader:        upgrader,
		clients:         make(map[*websocket.Conn]bool),
		dataDir:         dataDir,
		pluginDir:       pluginDir,
		pluginHealth:    make(map[string]PluginHealth),
		pluginManager:   plugin.NewPluginManager(pluginDir),
		lastUpdate:      time.Now(),
		authEnabled:     enableAuth,
		jwtSecretKey:    jwtSecretKey,
		compressionEnabled: enableCompression,
		compressionLevel: compressionLevel,
		compressionThreshold: compressionThreshold,
		performanceHistory: make([]PerformanceSnapshot, 0),
		maxHistorySize:  100, // Keep last 100 snapshots
	}

	// Set up routes
	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/plugins", s.handlePlugins)
	mux.HandleFunc("/plugins/capabilities", s.handlePluginCapabilities)
	mux.HandleFunc("/plugins/health", s.handlePluginHealth)
	mux.HandleFunc("/plugins/marketplace", s.handlePluginMarketplace)
	mux.HandleFunc("/plugins/install", s.handleInstallPlugin)
	mux.HandleFunc("/plugins/update", s.handleUpdatePlugin)
	mux.HandleFunc("/plugins/uninstall", s.handleUninstallPlugin)
	mux.HandleFunc("/performance/history", s.handlePerformanceHistory)
	mux.HandleFunc("/performance/analysis", s.handlePerformanceAnalysis)
	mux.HandleFunc("/plugins/performance", s.handlePluginPerformance)
	mux.HandleFunc("/compression/info", s.handleCompressionInfo)
	
	// Add auth endpoints if enabled
	if enableAuth {
		mux.HandleFunc("/auth/token", s.HandleGenerateToken)
	}

	// Load plugin manifests
	s.loadPluginManifests()

	return s
}

// Start starts the WebSocket server
func (s *WebSocketServer) Start() error {
	log.Printf("Starting WebSocket server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Stop stops the WebSocket server
func (s *WebSocketServer) Stop() error {
	log.Println("Stopping WebSocket server...")
	
	// Close all client connections
	s.clientsMu.Lock()
	for client := range s.clients {
		client.Close()
	}
	s.clients = make(map[*websocket.Conn]bool)
	s.clientsMu.Unlock()

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// LoadData loads findings and insights from files
func (s *WebSocketServer) LoadData(findingsPath, insightsPath string) error {
	// Load findings
	findingsData, err := os.ReadFile(findingsPath)
	if err != nil {
		return fmt.Errorf("failed to read findings: %w", err)
	}

	var findings model.FindingsBundle
	if err := json.Unmarshal(findingsData, &findings); err != nil {
		return fmt.Errorf("failed to parse findings: %w", err)
	}
	s.findings = &findings

	// Load insights if available
	if insightsPath != "" {
		insightsData, err := os.ReadFile(insightsPath)
		if err != nil {
			log.Printf("Warning: failed to read insights: %v", err)
		} else {
			var insights model.InsightsBundle
			if err := json.Unmarshal(insightsData, &insights); err != nil {
				log.Printf("Warning: failed to parse insights: %v", err)
			} else {
				s.insights = &insights
			}
		}
	}

	s.lastUpdate = time.Now()
	return nil
}

// BroadcastData sends current data to all connected clients
func (s *WebSocketServer) BroadcastData() {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	if len(s.clients) == 0 || s.findings == nil {
		return
	}

	// Prepare data payload
	payload := map[string]interface{}{
		"type":      "data_update",
		"timestamp": time.Now().Unix(),
		"findings":  s.findings,
		"insights":  s.insights,
		"stats": map[string]interface{}{
			"total_findings":      len(s.findings.Findings),
			"critical_count":      countSeverity(s.findings.Findings, "critical"),
			"high_count":          countSeverity(s.findings.Findings, "high"),
			"medium_count":        countSeverity(s.findings.Findings, "medium"),
			"low_count":           countSeverity(s.findings.Findings, "low"),
			"last_updated":        s.lastUpdate.Format(time.RFC3339),
			"performance_score":   s.findings.Summary.OverallScore,
			"connected_clients":   len(s.clients),
			"auth_enabled":        s.authEnabled,
			"compression":         s.GetCompressionInfo(),
		},
		"history": s.getPerformanceHistory(),
		"pluginPerformance": s.getPluginPerformanceSummary(),
	}

	// Send to all clients
	for client := range s.clients {
		if err := client.WriteJSON(payload); err != nil {
			log.Printf("Error sending data to client: %v", err)
			client.Close()
			delete(s.clients, client)
		}
	}
}

// GetPerformanceHistory returns the performance history for analysis
func (s *WebSocketServer) GetPerformanceHistory() []PerformanceSnapshot {
	s.historyMu.Lock()
	defer s.historyMu.Unlock()

	// Return a copy to avoid race conditions
	historyCopy := make([]PerformanceSnapshot, len(s.performanceHistory))
	copy(historyCopy, s.performanceHistory)
	return historyCopy
}

// getPerformanceHistory returns performance history for broadcasting
func (s *WebSocketServer) getPerformanceHistory() []PerformanceSnapshot {
	s.historyMu.Lock()
	defer s.historyMu.Unlock()

	// Return a copy to avoid race conditions
	historyCopy := make([]PerformanceSnapshot, len(s.performanceHistory))
	copy(historyCopy, s.performanceHistory)
	return historyCopy
}

// getPluginPerformanceSummary returns a summary of plugin performance for broadcasting
func (s *WebSocketServer) getPluginPerformanceSummary() map[string]interface{} {
	// Get plugin performance data
	performanceData := s.pluginManager.GetPluginPerformance()
	
	if len(performanceData) == 0 {
		return map[string]interface{}{
			"plugins": []map[string]interface{}{},
			"count":   0,
		}
	}
	
	// Group performance data by plugin
	pluginPerformanceMap := make(map[string][]plugin.PluginPerformance)
	for _, perf := range performanceData {
		pluginPerformanceMap[perf.PluginName] = append(pluginPerformanceMap[perf.PluginName], perf)
	}
	
	// Calculate summary statistics for each plugin
	pluginSummaries := make([]map[string]interface{}, 0)
	for pluginName, performances := range pluginPerformanceMap {
		if len(performances) == 0 {
			continue
		}
		
		// Sort by timestamp (newest first)
		sort.Slice(performances, func(i, j int) bool {
			return performances[i].Timestamp.After(performances[j].Timestamp)
		})
		
		// Calculate statistics
		var totalExecTime time.Duration
		successCount := 0
		
		for _, perf := range performances {
			totalExecTime += perf.ExecutionTime
			if perf.Success {
				successCount++
			}
		}
		
		avgExecTime := float64(totalExecTime.Nanoseconds()) / float64(len(performances)) / 1e6 // Convert to ms
		successRate := float64(successCount) / float64(len(performances)) * 100
		
		// Get latest performance
		latest := performances[0]
		
		pluginSummaries = append(pluginSummaries, map[string]interface{}{
			"pluginName":          pluginName,
			"executionCount":      len(performances),
			"successRate":         successRate,
			"avgExecutionTimeMs":  avgExecTime,
			"latestExecutionTimeMs": float64(latest.ExecutionTime.Nanoseconds()) / 1e6,
			"latestSuccess":        latest.Success,
		})
	}
	
	// Sort plugins by execution count (most used first)
	sort.Slice(pluginSummaries, func(i, j int) bool {
		return pluginSummaries[i]["executionCount"].(int) > pluginSummaries[j]["executionCount"].(int)
	})
	
	return map[string]interface{}{
		"plugins": pluginSummaries,
		"count":   len(pluginSummaries),
	}
}

// StartAutoRefresh enables periodic data broadcasting
func (s *WebSocketServer) StartAutoRefresh(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				s.BroadcastData()
			}
		}
	}()
}

// UpdateData updates the server's data and broadcasts to clients
func (s *WebSocketServer) UpdateData(findings *model.FindingsBundle, insights *model.InsightsBundle) {
	s.clientsMu.Lock()
	
	if findings != nil {
		s.findings = findings
	}
	if insights != nil {
		s.insights = insights
	}
	s.lastUpdate = time.Now()
	
	s.clientsMu.Unlock()
	
	// Record performance snapshot
	s.recordPerformanceSnapshot()
	
	// Broadcast updated data immediately (outside the lock to avoid deadlock)
	s.BroadcastData()
}

// recordPerformanceSnapshot captures current performance metrics for historical tracking
func (s *WebSocketServer) recordPerformanceSnapshot() {
	if s.findings == nil {
		return
	}

	snapshot := PerformanceSnapshot{
		Timestamp:     time.Now(),
		OverallScore:  s.findings.Summary.OverallScore,
		CriticalCount: countSeverity(s.findings.Findings, "critical"),
		HighCount:     countSeverity(s.findings.Findings, "high"),
		MediumCount:   countSeverity(s.findings.Findings, "medium"),
		LowCount:      countSeverity(s.findings.Findings, "low"),
		TotalFindings: len(s.findings.Findings),
		ClientCount:   s.GetClientCount(),
	}

	s.historyMu.Lock()
	defer s.historyMu.Unlock()

	// Add new snapshot
	s.performanceHistory = append(s.performanceHistory, snapshot)

	// Enforce max history size
	if len(s.performanceHistory) > s.maxHistorySize {
		s.performanceHistory = s.performanceHistory[len(s.performanceHistory)-s.maxHistorySize:]
	}
}

// GetClientCount returns the number of connected clients
func (s *WebSocketServer) GetClientCount() int {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	return len(s.clients)
}

// GetCompressionInfo returns compression configuration information
func (s *WebSocketServer) GetCompressionInfo() map[string]interface{} {
	return map[string]interface{}{
		"enabled":           s.compressionEnabled,
		"level":            s.compressionLevel,
		"threshold":        s.compressionThreshold,
		"description":      "WebSocket message compression reduces bandwidth usage for large performance data messages",
	}
}

// handleWebSocket handles WebSocket connections
func (s *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token if auth is enabled
	if s.authEnabled {
		tokenString := extractTokenFromRequest(r)
		if tokenString == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		claims, err := s.ValidateJWTToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		log.Printf("Authenticated WebSocket connection for user: %s (role: %s)", claims.Username, claims.Role)
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	s.clientsMu.Lock()
	s.clients[conn] = true
	s.clientsMu.Unlock()
	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, conn)
		s.clientsMu.Unlock()
	}()

	log.Printf("New WebSocket client connected: %s", conn.RemoteAddr())
	defer log.Printf("WebSocket client disconnected: %s", conn.RemoteAddr())

	// Send initial data
	if s.findings != nil {
		s.BroadcastData()
	}

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
	}
}

// generateJWTSecretKey generates a cryptographically secure random secret key
func generateJWTSecretKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based secret if crypto/rand fails
		return fmt.Sprintf("triageprof-jwt-secret-%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// GenerateJWTToken creates a new JWT token
func (s *WebSocketServer) GenerateJWTToken(username, role string) (string, error) {
	if !s.authEnabled {
		return "", fmt.Errorf("authentication is disabled")
	}

	claims := JWTClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "triageprof",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecretKey))
}

// ValidateJWTToken validates a JWT token
func (s *WebSocketServer) ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	if !s.authEnabled {
		// If auth is disabled, allow anonymous access
		return &JWTClaims{Username: "anonymous", Role: "viewer"}, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		if claims, ok := token.Claims.(*JWTClaims); ok {
			return claims, nil
		}
	}

	return nil, fmt.Errorf("invalid token")
}

// extractTokenFromRequest extracts JWT token from request headers or query parameters
func extractTokenFromRequest(r *http.Request) string {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Check query parameter
	return r.URL.Query().Get("token")
}

// HandleGenerateToken handles JWT token generation requests
func (s *WebSocketServer) HandleGenerateToken(w http.ResponseWriter, r *http.Request) {
	if !s.authEnabled {
		http.Error(w, "Authentication is disabled", http.StatusServiceUnavailable)
		return
	}

	// For demo purposes, allow any username/password
	// In production, this should validate against a user store
	type TokenRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	// Default role if not specified
	if req.Role == "" {
		req.Role = "viewer"
	}

	// Generate token
	token, err := s.GenerateJWTToken(req.Username, req.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":       token,
		"expires_in":  24 * 60 * 60, // 24 hours in seconds
		"username":    req.Username,
		"role":        req.Role,
	})
}

// handleHealth handles health check requests
func (s *WebSocketServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"clients":   len(s.clients),
		"data_loaded": s.findings != nil,
		"auth_enabled": s.authEnabled,
	})
}

// handleRoot handles root requests
func (s *WebSocketServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
	<title>TriageProf WebSocket Server</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 40px; }
		.status { padding: 20px; border-radius: 5px; margin-bottom: 20px; }
		.healthy { background-color: #d4edda; color: #155724; }
		.unhealthy { background-color: #f8d7da; color: #721c24; }
	</style>
</head>
<body>
	<h1>TriageProf WebSocket Server</h1>
	<div class="status healthy">Server is running on %s</div>
	<p>WebSocket endpoint: ws://%s/ws</p>
	<p>Use this endpoint for real-time performance monitoring.</p>
</body>
</html>`, s.server.Addr, s.server.Addr)
}

// loadPluginManifests loads plugin manifests from the plugin directory
func (s *WebSocketServer) loadPluginManifests() {
	manifestsDir := filepath.Join(s.pluginDir, "manifests")
	manifests, err := plugin.DiscoverManifests(manifestsDir)
	if err != nil {
		log.Printf("Warning: failed to load plugin manifests: %v", err)
		return
	}

	s.pluginManifests = manifests
	
	// Initialize health status for all plugins
	for _, manifest := range manifests {
		s.checkPluginHealth(manifest.Name)
	}
}

// checkPluginHealth checks the health status of a plugin
func (s *WebSocketServer) checkPluginHealth(pluginName string) {
	health := PluginHealth{
		LastChecked: time.Now(),
	}

	manifestsDir := filepath.Join(s.pluginDir, "manifests")
	binDir := filepath.Join(s.pluginDir, "bin")
	
	_, binaryPath, err := plugin.ResolvePlugin(manifestsDir, binDir, pluginName)
	if err != nil {
		health.Status = "unhealthy"
		health.Error = err.Error()
	} else {
		health.Status = "healthy"
		health.BinaryPath = binaryPath
		health.Error = ""
	}

	s.pluginHealth[pluginName] = health
}

// handlePlugins handles plugin listing requests
func (s *WebSocketServer) handlePlugins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	plugins := make([]map[string]interface{}, len(s.pluginManifests))
	for i, manifest := range s.pluginManifests {
		health, exists := s.pluginHealth[manifest.Name]
		if !exists {
			s.checkPluginHealth(manifest.Name)
			health = s.pluginHealth[manifest.Name]
		}
		
		plugins[i] = map[string]interface{}{
			"name":        manifest.Name,
			"version":     manifest.Version,
			"sdkVersion":  manifest.SDKVersion,
			"description": manifest.Description,
			"author":      manifest.Author,
			"capabilities": map[string]interface{}{
				"targets":  manifest.Capabilities.Targets,
				"profiles": manifest.Capabilities.Profiles,
			},
			"health": health,
		}
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"plugins": plugins,
		"count":  len(plugins),
	})
}

// handlePluginCapabilities handles plugin capability matrix requests
func (s *WebSocketServer) handlePluginCapabilities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Collect all unique targets and profiles
	allTargets := make(map[string]bool)
	allProfiles := make(map[string]bool)
	
	for _, manifest := range s.pluginManifests {
		for _, target := range manifest.Capabilities.Targets {
			allTargets[target] = true
		}
		for _, profile := range manifest.Capabilities.Profiles {
			allProfiles[profile] = true
		}
	}
	
	// Sort for consistent ordering
	targets := make([]string, 0, len(allTargets))
	for target := range allTargets {
		targets = append(targets, target)
	}
	sort.Strings(targets)
	
	profiles := make([]string, 0, len(allProfiles))
	for profile := range allProfiles {
		profiles = append(profiles, profile)
	}
	sort.Strings(profiles)
	
	// Build capability matrix
	matrix := make([]map[string]interface{}, len(s.pluginManifests))
	for i, manifest := range s.pluginManifests {
		pluginMatrix := make(map[string]interface{})
		pluginMatrix["plugin"] = manifest.Name
		
		// Target capabilities
		targetCaps := make(map[string]bool)
		for _, target := range manifest.Capabilities.Targets {
			targetCaps[target] = true
		}
		
		// Profile capabilities
		profileCaps := make(map[string]bool)
		for _, profile := range manifest.Capabilities.Profiles {
			profileCaps[profile] = true
		}
		
		pluginMatrix["targets"] = targetCaps
		pluginMatrix["profiles"] = profileCaps
		
		matrix[i] = pluginMatrix
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"targets":   targets,
		"profiles":  profiles,
		"matrix":    matrix,
	})
}

// handlePluginHealth handles plugin health status requests
func (s *WebSocketServer) handlePluginHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Refresh health status for all plugins
	for _, manifest := range s.pluginManifests {
		s.checkPluginHealth(manifest.Name)
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"health": s.pluginHealth,
		"timestamp": time.Now().Unix(),
	})
}

// handlePluginMarketplace handles plugin marketplace requests
func (s *WebSocketServer) handlePluginMarketplace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Demo marketplace data - in a real implementation, this would fetch from a remote marketplace
	marketplacePlugins := []map[string]interface{}{
		{
			"name": "go-pprof-http",
			"version": "0.1.0",
			"installed": s.isPluginInstalled("go-pprof-http"),
			"description": "Go pprof HTTP plugin for collecting profiles from Go applications",
			"author": "Mistral Hackathon",
			"capabilities": map[string]interface{}{
				"targets": []string{"url"},
				"profiles": []string{"cpu", "heap", "mutex", "block", "goroutine"},
			},
		},
		{
			"name": "node-inspector",
			"version": "0.1.0",
			"installed": s.isPluginInstalled("node-inspector"),
			"description": "Node.js inspector plugin for profiling Node.js applications",
			"author": "Mistral Hackathon",
			"capabilities": map[string]interface{}{
				"targets": []string{"url"},
				"profiles": []string{"cpu", "heap"},
			},
		},
		{
			"name": "python-cprofile",
			"version": "0.1.0",
			"installed": s.isPluginInstalled("python-cprofile"),
			"description": "Python cProfile plugin for profiling Python applications",
			"author": "Mistral Hackathon",
			"capabilities": map[string]interface{}{
				"targets": []string{"url"},
				"profiles": []string{"cpu", "memory"},
			},
		},
		{
			"name": "ruby-stackprof",
			"version": "0.1.0",
			"installed": s.isPluginInstalled("ruby-stackprof"),
			"description": "Ruby stackprof plugin for profiling Ruby applications",
			"author": "Mistral Hackathon",
			"capabilities": map[string]interface{}{
				"targets": []string{"url"},
				"profiles": []string{"cpu", "memory", "object_allocation"},
			},
		},
		{
			"name": "java-jfr",
			"version": "0.1.0",
			"installed": s.isPluginInstalled("java-jfr"),
			"description": "Java Flight Recorder plugin for profiling Java applications",
			"author": "Mistral Hackathon",
			"capabilities": map[string]interface{}{
				"targets": []string{"url", "file"},
				"profiles": []string{"cpu", "memory", "gc", "locks"},
			},
		},
		{
			"name": "dotnet-profiler",
			"version": "0.1.0",
			"installed": s.isPluginInstalled("dotnet-profiler"),
			"description": ".NET profiler for profiling C# applications",
			"author": "Mistral Hackathon",
			"capabilities": map[string]interface{}{
				"targets": []string{"url"},
				"profiles": []string{"cpu", "memory", "gc"},
			},
		},
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"plugins": marketplacePlugins,
		"count": len(marketplacePlugins),
	})
}

// handleInstallPlugin handles plugin installation requests
func (s *WebSocketServer) handleInstallPlugin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var request struct {
		PluginName string `json:"pluginName"`
		URL        string `json:"url,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if request.PluginName == "" && request.URL == "" {
		http.Error(w, "Plugin name or URL is required", http.StatusBadRequest)
		return
	}

	// In a real implementation, this would download and install the plugin
	// For demo purposes, we'll simulate a successful installation
	time.Sleep(1 * time.Second)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Plugin %s installed successfully", request.PluginName),
		"pluginName": request.PluginName,
	})
}

// handleUpdatePlugin handles plugin update requests
func (s *WebSocketServer) handleUpdatePlugin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var request struct {
		PluginName string `json:"pluginName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if request.PluginName == "" {
		http.Error(w, "Plugin name is required", http.StatusBadRequest)
		return
	}

	// In a real implementation, this would update the plugin
	// For demo purposes, we'll simulate a successful update
	time.Sleep(1 * time.Second)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Plugin %s updated successfully", request.PluginName),
		"pluginName": request.PluginName,
	})
}

// handleUninstallPlugin handles plugin uninstallation requests
func (s *WebSocketServer) handleUninstallPlugin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var request struct {
		PluginName string `json:"pluginName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if request.PluginName == "" {
		http.Error(w, "Plugin name is required", http.StatusBadRequest)
		return
	}

	// In a real implementation, this would uninstall the plugin
	// For demo purposes, we'll simulate a successful uninstallation
	time.Sleep(1 * time.Second)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Plugin %s uninstalled successfully", request.PluginName),
		"pluginName": request.PluginName,
	})
}

// handlePluginPerformance handles plugin performance metrics requests
func (s *WebSocketServer) handlePluginPerformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get plugin performance data
	performanceData := s.pluginManager.GetPluginPerformance()
	
	// Group performance data by plugin
	pluginPerformanceMap := make(map[string][]plugin.PluginPerformance)
	for _, perf := range performanceData {
		pluginPerformanceMap[perf.PluginName] = append(pluginPerformanceMap[perf.PluginName], perf)
	}
	
	// Calculate statistics for each plugin
	pluginStats := make([]map[string]interface{}, 0)
	for pluginName, performances := range pluginPerformanceMap {
		if len(performances) == 0 {
			continue
		}
		
		// Sort by timestamp (newest first)
		sort.Slice(performances, func(i, j int) bool {
			return performances[i].Timestamp.After(performances[j].Timestamp)
		})
		
		// Calculate statistics
		var totalExecTime time.Duration
		var totalMemory, totalCPU float64
		successCount := 0
		
		for _, perf := range performances {
			totalExecTime += perf.ExecutionTime
			totalMemory += perf.MemoryUsageMB
			totalCPU += perf.CPUUsagePercent
			if perf.Success {
				successCount++
			}
		}
		
		avgExecTime := float64(totalExecTime.Nanoseconds()) / float64(len(performances)) / 1e6 // Convert to ms
		avgMemory := totalMemory / float64(len(performances))
		avgCPU := totalCPU / float64(len(performances))
		successRate := float64(successCount) / float64(len(performances)) * 100
		
		// Get latest performance
		latest := performances[0]
		
		pluginStats = append(pluginStats, map[string]interface{}{
			"pluginName":          pluginName,
			"executionCount":      len(performances),
			"successCount":        successCount,
			"failureCount":        len(performances) - successCount,
			"successRate":         successRate,
			"avgExecutionTimeMs":  avgExecTime,
			"avgMemoryUsageMB":    avgMemory,
			"avgCPUUsagePercent":  avgCPU,
			"latestExecutionTimeMs": float64(latest.ExecutionTime.Nanoseconds()) / 1e6,
			"latestMemoryUsageMB":   latest.MemoryUsageMB,
			"latestCPUUsagePercent": latest.CPUUsagePercent,
			"latestTimestamp":      latest.Timestamp.Format(time.RFC3339),
			"latestSuccess":        latest.Success,
			"latestError":          latest.Error,
		})
	}
	
	// Sort plugins by execution count (most used first)
	sort.Slice(pluginStats, func(i, j int) bool {
		return pluginStats[i]["executionCount"].(int) > pluginStats[j]["executionCount"].(int)
	})
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"plugins": pluginStats,
		"count":  len(pluginStats),
		"timestamp": time.Now().Unix(),
	})
}

// handleCompressionInfo handles compression information requests
func (s *WebSocketServer) handleCompressionInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	compressionInfo := s.GetCompressionInfo()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(compressionInfo)
}

// handlePerformanceHistory handles performance history requests
func (s *WebSocketServer) handlePerformanceHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	history := s.GetPerformanceHistory()
	
	// Add analysis to the response
	response := map[string]interface{}{
		"history": history,
		"count":   len(history),
		"analysis": s.analyzePerformanceTrends(history),
	}
	
	json.NewEncoder(w).Encode(response)
}

// handlePerformanceAnalysis handles performance analysis requests
func (s *WebSocketServer) handlePerformanceAnalysis(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	history := s.GetPerformanceHistory()
	analysis := s.analyzePerformanceTrends(history)
	
	// Add current state
	currentState := map[string]interface{}{
		"overall_score": s.findings.Summary.OverallScore,
		"critical_count": countSeverity(s.findings.Findings, "critical"),
		"high_count":     countSeverity(s.findings.Findings, "high"),
		"medium_count":   countSeverity(s.findings.Findings, "medium"),
		"low_count":      countSeverity(s.findings.Findings, "low"),
		"total_findings": len(s.findings.Findings),
		"client_count":   s.GetClientCount(),
	}
	
	response := map[string]interface{}{
		"current":  currentState,
		"analysis": analysis,
		"trends":   s.calculatePerformanceTrends(history),
	}
	
	json.NewEncoder(w).Encode(response)
}

// analyzePerformanceTrends analyzes historical performance data
func (s *WebSocketServer) analyzePerformanceTrends(history []PerformanceSnapshot) map[string]interface{} {
	if len(history) == 0 {
		return map[string]interface{}{
			"status": "no_data",
			"message": "No performance history available",
		}
	}

	// Calculate basic statistics
	var totalScore, totalCritical, totalHigh, totalMedium, totalLow, totalFindings int
	var minScore, maxScore int = 100, 0
	
	for i, snapshot := range history {
		totalScore += snapshot.OverallScore
		totalCritical += snapshot.CriticalCount
		totalHigh += snapshot.HighCount
		totalMedium += snapshot.MediumCount
		totalLow += snapshot.LowCount
		totalFindings += snapshot.TotalFindings
		
		if i == 0 || snapshot.OverallScore < minScore {
			minScore = snapshot.OverallScore
		}
		if i == 0 || snapshot.OverallScore > maxScore {
			maxScore = snapshot.OverallScore
		}
	}

	avgScore := float64(totalScore) / float64(len(history))
	avgCritical := float64(totalCritical) / float64(len(history))
	avgHigh := float64(totalHigh) / float64(len(history))
	avgMedium := float64(totalMedium) / float64(len(history))
	avgLow := float64(totalLow) / float64(len(history))
	avgFindings := float64(totalFindings) / float64(len(history))

	// Determine trend direction
	trend := "stable"
	if len(history) >= 2 {
		first := history[0]
		last := history[len(history)-1]
		
		if last.OverallScore > first.OverallScore + 5 {
			trend = "improving"
		} else if last.OverallScore < first.OverallScore - 5 {
			trend = "degrading"
		}
	}

	return map[string]interface{}{
		"status":               "analyzed",
		"snapshot_count":      len(history),
		"time_range":           fmt.Sprintf("%s to %s", history[0].Timestamp.Format("2006-01-02 15:04:05"), history[len(history)-1].Timestamp.Format("2006-01-02 15:04:05")),
		"average_score":        avgScore,
		"score_range":          fmt.Sprintf("%d-%d", minScore, maxScore),
		"average_critical":    avgCritical,
		"average_high":        avgHigh,
		"average_medium":      avgMedium,
		"average_low":         avgLow,
		"average_findings":    avgFindings,
		"trend":                trend,
		"improvement_potential": calculateImprovementPotential(minScore, maxScore),
	}
}

// calculatePerformanceTrends calculates detailed performance trends
func (s *WebSocketServer) calculatePerformanceTrends(history []PerformanceSnapshot) map[string]interface{} {
	if len(history) < 2 {
		return map[string]interface{}{
			"status": "insufficient_data",
			"message": "Need at least 2 data points for trend analysis",
		}
	}

	// Calculate score trend
	scoreTrend := make([]map[string]interface{}, len(history))
	for i, snapshot := range history {
		scoreTrend[i] = map[string]interface{}{
			"timestamp": snapshot.Timestamp.Format("2006-01-02 15:04:05"),
			"score":     snapshot.OverallScore,
		}
	}

	// Calculate severity trends
	criticalTrend := make([]map[string]interface{}, len(history))
	highTrend := make([]map[string]interface{}, len(history))
	mediumTrend := make([]map[string]interface{}, len(history))
	lowTrend := make([]map[string]interface{}, len(history))
	
	for i, snapshot := range history {
		criticalTrend[i] = map[string]interface{}{
			"timestamp": snapshot.Timestamp.Format("2006-01-02 15:04:05"),
			"count":     snapshot.CriticalCount,
		}
		highTrend[i] = map[string]interface{}{
			"timestamp": snapshot.Timestamp.Format("2006-01-02 15:04:05"),
			"count":     snapshot.HighCount,
		}
		mediumTrend[i] = map[string]interface{}{
			"timestamp": snapshot.Timestamp.Format("2006-01-02 15:04:05"),
			"count":     snapshot.MediumCount,
		}
		lowTrend[i] = map[string]interface{}{
			"timestamp": snapshot.Timestamp.Format("2006-01-02 15:04:05"),
			"count":     snapshot.LowCount,
		}
	}

	return map[string]interface{}{
		"status":         "calculated",
		"score_trend":    scoreTrend,
		"critical_trend": criticalTrend,
		"high_trend":     highTrend,
		"medium_trend":   mediumTrend,
		"low_trend":      lowTrend,
	}
}

// calculateImprovementPotential calculates potential improvement based on score range
func calculateImprovementPotential(minScore, maxScore int) map[string]interface{} {
	if minScore == maxScore {
		return map[string]interface{}{
			"potential": 0,
			"percentage": 0.0,
			"message": "Performance is stable",
		}
	}

	potential := maxScore - minScore
	percentage := float64(potential) / float64(minScore) * 100
	
	message := "stable"
	if percentage > 20 {
		message = "significant improvement potential"
	} else if percentage > 10 {
		message = "moderate improvement potential"
	} else if percentage > 5 {
		message = "some improvement potential"
	}

	return map[string]interface{}{
		"potential": potential,
		"percentage": percentage,
		"message":   message,
	}
}

// isPluginInstalled checks if a plugin is installed
func (s *WebSocketServer) isPluginInstalled(pluginName string) bool {
	for _, manifest := range s.pluginManifests {
		if manifest.Name == pluginName {
			return true
		}
	}
	return false
}

// countSeverity counts findings by severity level
func countSeverity(findings []model.Finding, severity string) int {
	count := 0
	for _, f := range findings {
		if f.Severity == severity {
			count++
		}
	}
	return count
}
