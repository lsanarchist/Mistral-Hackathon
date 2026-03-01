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
	authEnabled     bool
	jwtSecretKey    string
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
func NewWebSocketServer(port int, dataDir string, pluginDir string, enableAuth bool) *WebSocketServer {
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
		lastUpdate:      time.Now(),
		authEnabled:     enableAuth,
		jwtSecretKey:    jwtSecretKey,
	}

	// Set up routes
	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/plugins", s.handlePlugins)
	mux.HandleFunc("/plugins/capabilities", s.handlePluginCapabilities)
	mux.HandleFunc("/plugins/health", s.handlePluginHealth)
	
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
		},
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
	
	// Broadcast updated data immediately (outside the lock to avoid deadlock)
	s.BroadcastData()
}

// GetClientCount returns the number of connected clients
func (s *WebSocketServer) GetClientCount() int {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	return len(s.clients)
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
