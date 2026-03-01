package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mistral-hackathon/triageprof/internal/model"
)

// WebSocketServer handles real-time data streaming
type WebSocketServer struct {
	server      *http.Server
	upgrader    websocket.Upgrader
	clients     map[*websocket.Conn]bool
	clientsMu   sync.Mutex
	dataDir     string
	findings    *model.FindingsBundle
	insights    *model.InsightsBundle
	lastUpdate  time.Time
}

// NewWebSocketServer creates a new WebSocket server instance
func NewWebSocketServer(port int, dataDir string, enableAuth bool) *WebSocketServer {
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

	s := &WebSocketServer{
		server:     server,
		upgrader:   upgrader,
		clients:    make(map[*websocket.Conn]bool),
		dataDir:    dataDir,
		lastUpdate: time.Now(),
	}

	// Set up routes
	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/", s.handleRoot)

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

// handleHealth handles health check requests
func (s *WebSocketServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"clients":   len(s.clients),
		"data_loaded": s.findings != nil,
		"auth_enabled": false,
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
