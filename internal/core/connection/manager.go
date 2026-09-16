package connection

import (
	"sync"

	"go-walis/internal/core/db"
)

// Manager gerencia instâncias de conexões ativas reutilizáveis
type Manager struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var globalManager = &Manager{
	clients: make(map[string]*Client),
}

// GetManager retorna a instância global do gerenciador de conexões
func GetManager() *Manager {
	return globalManager
}

// GetClient obtém ou cria uma nova conexão para o servidor reaproveitando a sessão existente
func (m *Manager) GetClient(server db.VpsServer) (*Client, error) {
	m.mu.Lock()
	client, exists := m.clients[server.ID]
	if exists {
		// Se os dados essenciais de conexão mudaram, fecha e recria
		if client.server.Host == server.Host &&
			client.server.Port == server.Port &&
			client.server.Username == server.Username &&
			client.server.SshPassword == server.SshPassword &&
			client.server.SudoPassword == server.SudoPassword &&
			client.server.SshKeyPath == server.SshKeyPath &&
			client.server.DockerSocketPath == server.DockerSocketPath &&
			client.server.SshKeyPassphrase == server.SshKeyPassphrase &&
			client.server.ConnectionType == server.ConnectionType {
			// Libera o lock global antes de realizar I/O de rede potencialmente bloqueante
			m.mu.Unlock()

			if client.IsAlive() {
				return client, nil
			}

			// Se morreu, adquire o lock novamente para remover e recriar
			m.mu.Lock()
			// Verifica se outro já não substituiu
			if current, stillSame := m.clients[server.ID]; stillSame && current == client {
				client.ForceClose()
				delete(m.clients, server.ID)
			}
		} else {
			client.ForceClose()
			delete(m.clients, server.ID)
		}
	}
	m.mu.Unlock()

	// Cria o novo cliente fora do lock global para evitar travar outras goroutines durante o SSH Dial
	newClient, err := NewClient(server)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Se outro cliente foi criado enquanto este conectava, fecha o anterior
	if old, exists := m.clients[server.ID]; exists {
		old.ForceClose()
	}
	m.clients[server.ID] = newClient
	return newClient, nil
}

// CloseClient remove e fecha a conexão de um servidor específico
func (m *Manager) CloseClient(serverID string) {
	m.mu.Lock()
	client, exists := m.clients[serverID]
	if exists {
		delete(m.clients, serverID)
	}
	m.mu.Unlock()

	if exists && client != nil {
		client.ForceClose()
	}
}
