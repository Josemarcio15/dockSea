package connection

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"go-walis/internal/core/db"

	"golang.org/x/crypto/ssh"
)

// Client representa uma conexão ativa (SSH remota ou Local)
type Client struct {
	server     db.VpsServer
	sshClient  *ssh.Client
	httpClient *http.Client
	lastPing   time.Time
	startedAt  time.Time
	mu         sync.RWMutex
}

// NewClient instancia um novo cliente de conexão
func NewClient(server db.VpsServer) (*Client, error) {
	now := time.Now()
	c := &Client{server: server, lastPing: now, startedAt: now}

	if isLocal(server) {
		socketPath := server.DockerSocketPath
		if socketPath == "" {
			socketPath = "/var/run/docker.sock"
		}
		c.httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, proto, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
			Timeout: 30 * time.Second,
		}
		LogConnectionStarted(server, now)
		return c, nil
	}

	// Conexão SSH
	sshClient, err := createSshClient(server)
	if err != nil {
		return nil, err
	}
	c.sshClient = sshClient

	// HTTP Client tunelado pelo SSH para o socket Docker remoto
	c.httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, proto, addr string) (net.Conn, error) {
				socketPath := server.DockerSocketPath
				if socketPath == "" {
					socketPath = "/var/run/docker.sock"
				}
				c.mu.RLock()
				sc := c.sshClient
				c.mu.RUnlock()
				if sc == nil {
					return nil, net.ErrClosed
				}
				return sc.Dial("unix", socketPath)
			},
		},
		Timeout: 30 * time.Second,
	}

	LogConnectionStarted(server, now)
	return c, nil
}

// GetHttpClient retorna o cliente HTTP com timeout padrão para requisições rápidas
func (c *Client) GetHttpClient() *http.Client {
	return c.httpClient
}

// GetStreamHttpClient retorna um cliente HTTP com transporte dedicado e sem timeout para streaming longo
func (c *Client) GetStreamHttpClient() *http.Client {
	if isLocal(c.server) {
		socketPath := c.server.DockerSocketPath
		if socketPath == "" {
			socketPath = "/var/run/docker.sock"
		}
		return &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, proto, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
				DisableKeepAlives: true,
			},
			Timeout: 0,
		}
	}

	c.mu.RLock()
	sc := c.sshClient
	c.mu.RUnlock()
	if sc == nil {
		return nil
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, proto, addr string) (net.Conn, error) {
				c.mu.RLock()
				activeClient := c.sshClient
				c.mu.RUnlock()
				if activeClient == nil {
					return nil, net.ErrClosed
				}
				socketPath := c.server.DockerSocketPath
				if socketPath == "" {
					socketPath = "/var/run/docker.sock"
				}
				return activeClient.Dial("unix", socketPath)
			},
			DisableKeepAlives: true,
		},
		Timeout: 0,
	}
}

// GetSSHClient retorna uma referência segura ao cliente SSH
func (c *Client) GetSSHClient() *ssh.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sshClient
}

// IsAlive verifica se a conexão ainda está saudável (com cache de 10s para evitar RTT excessivo)
func (c *Client) IsAlive() bool {
	if isLocal(c.server) {
		socketPath := c.server.DockerSocketPath
		if socketPath == "" {
			socketPath = "/var/run/docker.sock"
		}
		_, err := os.Stat(socketPath)
		return err == nil
	}

	c.mu.RLock()
	sc := c.sshClient
	last := c.lastPing
	c.mu.RUnlock()

	if sc == nil {
		return false
	}

	// Se testou a menos de 10 segundos, assume saudável
	if time.Since(last) < 10*time.Second {
		return true
	}

	_, _, err := sc.SendRequest("keepalive@openssh.com", true, nil)
	if err == nil {
		c.mu.Lock()
		c.lastPing = time.Now()
		c.mu.Unlock()
		return true
	}
	return false
}

// Close é um no-op para conexões gerenciadas pelo pool (para manter a conexão aberta entre requisições)
func (c *Client) Close() {
	// Mantido aberto para reúso no pool
}

// ForceClose encerra definitivamente a conexão SSH ativa de forma thread-safe
func (c *Client) ForceClose() {
	c.mu.Lock()
	sc := c.sshClient
	c.sshClient = nil
	c.mu.Unlock()

	if sc != nil {
		_ = sc.Close()
	}
	LogConnectionClosed(c.server, c.startedAt, time.Now())
}
