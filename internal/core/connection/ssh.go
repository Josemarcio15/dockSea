package connection

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-walis/internal/core/db"

	"golang.org/x/crypto/ssh"
)

var (
	knownHostsMu sync.Mutex
)

// getKnownHostsPath retorna o caminho do arquivo de fingerprints conhecidos do DockSea
func getKnownHostsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	dir := filepath.Join(home, ".docksea")
	_ = os.MkdirAll(dir, 0700)
	return filepath.Join(dir, "known_hosts")
}

// tofuHostKeyCallback implementa TOFU (Trust On First Use) seguro para validação de chaves SSH
func tofuHostKeyCallback(targetAddr string) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		knownHostsPath := getKnownHostsPath()
		if knownHostsPath == "" {
			// Fallback se não conseguir obter diretório do usuário
			return nil
		}

		knownHostsMu.Lock()
		defer knownHostsMu.Unlock()

		keyFingerprint := hex.EncodeToString(sha256.New().Sum(key.Marshal()))
		keyType := key.Type()
		entry := fmt.Sprintf("%s %s %s\n", targetAddr, keyType, keyFingerprint)

		data, err := os.ReadFile(knownHostsPath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("falha ao ler known_hosts: %w", err)
		}

		if len(data) > 0 {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				parts := strings.Fields(strings.TrimSpace(line))
				if len(parts) >= 3 && parts[0] == targetAddr {
					if parts[1] == keyType && parts[2] == keyFingerprint {
						// Chave idêntica e confiável
						return nil
					}
					// Chave mudou! Alerta de possível Man-In-The-Middle
					return fmt.Errorf("ATENÇÃO: A chave do host SSH para '%s' mudou! Possível ataque Man-in-the-Middle ou servidor reinstalado. Fingerprint gravado: %s, Fingerprint recebido: %s", targetAddr, parts[2], keyFingerprint)
				}
			}
		}

		// Primeiro uso (TOFU): registra a chave conhecida
		f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			_, _ = f.WriteString(entry)
			_ = f.Close()
		}
		return nil
	}
}

// createSshClient constrói a conexão SSH usando a regra inteligente de autenticação
func createSshClient(server db.VpsServer) (*ssh.Client, error) {
	host := strings.TrimSpace(server.Host)
	if host == "" {
		return nil, fmt.Errorf("o host ou IP da VPS é obrigatório")
	}

	port := server.Port
	if port == 0 {
		port = 22
	}
	// Utiliza net.JoinHostPort para suportar IPv4, IPv6 (ex: [2001:db8::1]:22) e domínios com segurança (SSH-006)
	targetAddr := net.JoinHostPort(host, strconv.Itoa(port))

	var authMethods []ssh.AuthMethod

	// 1. Chave privada (se informada)
	keyPath := strings.TrimSpace(server.SshKeyPath)
	if keyPath != "" {
		if strings.HasPrefix(keyPath, "~/") {
			home, _ := os.UserHomeDir()
			keyPath = home + keyPath[1:]
		}

		keyBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("não foi possível ler o arquivo da chave SSH '%s': %w", keyPath, err)
		}

		var signer ssh.Signer
		passphrase := strings.TrimSpace(server.SshKeyPassphrase)
		if passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(keyBytes)
		}

		if err != nil {
			return nil, fmt.Errorf("chave privada SSH inválida ou senha incorreta: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// 2. Senha SSH (se informada)
	password := strings.TrimSpace(server.SshPassword)
	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	// Se nada foi informado
	if len(authMethods) == 0 {
		return nil, fmt.Errorf("nenhum método de autenticação SSH fornecido (preencha a chave privada ou senha)")
	}

	username := strings.TrimSpace(server.Username)
	if username == "" {
		username = "root"
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: tofuHostKeyCallback(targetAddr),
		Timeout:         8 * time.Second,
	}

	return ssh.Dial("tcp", targetAddr, config)
}

func isLocal(server db.VpsServer) bool {
	if server.ConnectionType == "local" {
		return true
	}
	host := strings.TrimSpace(strings.ToLower(server.Host))
	return host == "" || host == "localhost" || host == "127.0.0.1"
}

func escapeShell(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
