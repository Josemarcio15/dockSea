package dashboard

import (
	"fmt"
	"go-walis/internal/core/connection"
	"go-walis/internal/core/db"
)

func CollectSystemUsage(server db.VpsServer) (*connection.SystemUsage, error) {
	client, err := connection.GetClient(server)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no servidor: %w", err)
	}
	defer client.Close()
	return client.FetchSystemUsage()
}
