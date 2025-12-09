package nats

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// EmbeddedServer wraps an embedded NATS server
type EmbeddedServer struct {
	opts     *server.Options
	srv      *server.Server
	conn     *nats.Conn
	js       nats.JetStreamContext
	port     int
	host     string
	mu       sync.RWMutex
	running  bool
	jetStream bool
}

// NewEmbeddedServer creates and starts an embedded NATS server
// natsURL should be in format "host:port" (e.g., "0.0.0.0:4222")
// storeDir is the directory for JetStream data (should be relative to PocketBase data directory)
func NewEmbeddedServer(natsURL string, jetStream bool, storeDir string) (*EmbeddedServer, error) {
	// Parse natsURL (format: "host:port")
	host, port, err := ParseNATSAddress(natsURL)
	if err != nil {
		return nil, fmt.Errorf("invalid natsURL format: %w", err)
	}

	// Default store directory if not provided
	if storeDir == "" {
		storeDir = "./nats_data"
	}

	// Create NATS data directory if it doesn't exist (for JetStream)
	if jetStream {
		if err := os.MkdirAll(storeDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create NATS store directory: %w", err)
		}
		// Use absolute path to avoid issues
		absStoreDir, err := filepath.Abs(storeDir)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for store directory: %w", err)
		}
		storeDir = absStoreDir
	}

	// Create NATS server options
	// Set MaxPayload to 100MB to support large snapshots (default is 1MB)
	// MaxPending must be >= MaxPayload, so set it to the same value
	maxPayload := int32(100 * 1024 * 1024) // 100MB
	maxPending := int64(100 * 1024 * 1024) // 100MB
	opts := &server.Options{
		Host:           host,
		Port:           port,
		JetStream:      jetStream,
		StoreDir:       storeDir,
		NoLog:          false,
		NoSigs:         true,
		MaxControlLine: 4096,
		MaxPayload:     maxPayload,
		MaxPending:     maxPending,
	}

	// Create and start the server
	log.Printf("Creating NATS server on port %d (JetStream: %v, StoreDir: %s)", port, jetStream, storeDir)
	srv, err := server.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS server: %w", err)
	}

	// Start the server in a goroutine
	log.Printf("Starting NATS server...")
	go srv.Start()

	// Wait for server to be ready (with longer timeout for JetStream initialization)
	log.Printf("Waiting for NATS server to be ready (timeout: 30s)...")
	if !srv.ReadyForConnections(30 * time.Second) {
		log.Printf("NATS server failed to become ready, shutting down...")
		srv.Shutdown()
		return nil, fmt.Errorf("NATS server failed to start within timeout")
	}

	// Additional wait for JetStream to be ready if enabled
	if jetStream {
		log.Printf("Waiting for JetStream to initialize...")
		// Give JetStream a moment to initialize
		time.Sleep(1 * time.Second)
	}

	log.Printf("Embedded NATS server started on %s:%d", host, port)

	// Connect to the embedded server
	connectURL := fmt.Sprintf("nats://%s:%d", host, port)
	conn, err := nats.Connect(connectURL)
	if err != nil {
		srv.Shutdown()
		return nil, fmt.Errorf("failed to connect to embedded NATS server: %w", err)
	}

	es := &EmbeddedServer{
		opts:      opts,
		srv:       srv,
		conn:      conn,
		host:      host,
		port:      port,
		running:   true,
		jetStream: jetStream,
	}

	// Get JetStream context if enabled
	if jetStream {
		js, err := conn.JetStream()
		if err != nil {
			conn.Close()
			srv.Shutdown()
			return nil, fmt.Errorf("failed to get JetStream context: %w", err)
		}
		es.js = js
	}

	return es, nil
}

// Port returns the port the server is running on
func (es *EmbeddedServer) Port() int {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return es.port
}

// URL returns the connection URL for this server
func (es *EmbeddedServer) URL() string {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return fmt.Sprintf("nats://%s:%d", es.host, es.port)
}

// ParseNATSAddress parses a NATS address in format "host:port" and returns host and port
func ParseNATSAddress(addr string) (string, int, error) {
	if addr == "" {
		return "0.0.0.0", 4222, nil
	}

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid address format: %w", err)
	}

	if host == "" {
		host = "0.0.0.0"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %w", err)
	}

	if port < 1 || port > 65535 {
		return "", 0, fmt.Errorf("port out of range: %d", port)
	}

	return host, port, nil
}

// Connection returns the NATS connection to the embedded server
func (es *EmbeddedServer) Connection() *nats.Conn {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return es.conn
}

// JetStream returns the JetStream context
func (es *EmbeddedServer) JetStream() nats.JetStreamContext {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return es.js
}

// IsRunning returns whether the server is running
func (es *EmbeddedServer) IsRunning() bool {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return es.running
}

// Shutdown stops the embedded NATS server
func (es *EmbeddedServer) Shutdown() error {
	es.mu.Lock()
	defer es.mu.Unlock()

	if !es.running {
		return nil
	}

	if es.conn != nil {
		es.conn.Close()
	}

	if es.srv != nil {
		es.srv.Shutdown()
	}

	es.running = false
	log.Printf("Embedded NATS server stopped")
	return nil
}

