package db

import (
	"database/sql"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptSecret(t *testing.T) {
	original := "MinhaSenhaSuperSecreta123!@#"

	encrypted, err := encryptSecret(original)
	if err != nil {
		t.Fatalf("Erro ao criptografar: %v", err)
	}

	if encrypted == original {
		t.Fatalf("O texto criptografado não deveria ser igual ao original")
	}

	decrypted, err := decryptSecret(encrypted)
	if err != nil {
		t.Fatalf("Erro ao descriptografar: %v", err)
	}

	if decrypted != original {
		t.Fatalf("Esperado %q, obteve %q", original, decrypted)
	}
}

func TestDecryptMismatchMachineGraceful(t *testing.T) {
	otherCipher, err := encryptionCipherWithKey("chave-de-outro-computador-qualquer")
	if err != nil {
		t.Fatalf("Erro ao gerar cifra com outra chave: %v", err)
	}

	nonce := make([]byte, otherCipher.NonceSize())
	ciphertext := otherCipher.Seal(nil, nonce, []byte("senha-secreta-em-outro-pc"), nil)
	encoded := append(nonce, ciphertext...)
	secretEnc := "enc:v1:" + string(encodeBase64(encoded))

	decrypted, err := decryptSecret(secretEnc)
	if err != nil {
		t.Fatalf("Não deveria retornar erro fatal, mas sim tratar graciosamente: %v", err)
	}
	if decrypted != "" {
		t.Fatalf("Esperado string vazia ao decodificar em máquina diferente, obteve %q", decrypted)
	}
}

func encodeBase64(data []byte) string {
	return base64.RawStdEncoding.EncodeToString(data)
}

func TestMultiTenantProfiles(t *testing.T) {
	tmpDir := t.TempDir()
	profilesDir := filepath.Join(tmpDir, "profiles")
	_ = os.MkdirAll(profilesDir, 0755)

	masterDbPath := filepath.Join(tmpDir, "tenants.db")
	masterConn, err := sql.Open("sqlite", masterDbPath)
	if err != nil {
		t.Fatalf("Erro ao abrir master db: %v", err)
	}
	defer masterConn.Close()

	database := &DB{
		masterConn:  masterConn,
		appDir:      tmpDir,
		profilesDir: profilesDir,
	}
	defer database.Close()

	if err := database.migrateMaster(); err != nil {
		t.Fatalf("Erro ao migrar master db: %v", err)
	}

	// 1. Criar perfil "joao"
	p1 := Profile{Name: "joao", Locale: "pt-BR", IsActive: true}
	if err := database.SaveProfile(p1); err != nil {
		t.Fatalf("Erro ao criar perfil joao: %v", err)
	}

	// Verificar se joao.db foi criado
	joaoDbPath := filepath.Join(profilesDir, "joao.db")
	if _, err := os.Stat(joaoDbPath); os.IsNotExist(err) {
		t.Fatalf("Arquivo joao.db não foi criado no disco")
	}

	// 2. Tentar criar outro perfil "joao" duplicado -> deve retornar erro
	pDupe := Profile{Name: "joao", Locale: "pt-BR"}
	if err := database.SaveProfile(pDupe); err == nil {
		t.Fatalf("Deveria ter rejeitado perfil duplicado com o mesmo nome")
	}

	// 3. Conectar ao perfil joao
	profiles, _ := database.ListProfiles()
	if len(profiles) != 1 {
		t.Fatalf("Esperava 1 perfil, obteve %d", len(profiles))
	}
	joaoProfile := profiles[0]
	if err := database.switchProfileConnection(&joaoProfile); err != nil {
		t.Fatalf("Erro ao alternar para joao: %v", err)
	}

	// Salvar VPS no perfil joao
	vps1 := VpsServer{Name: "VPS do João", ConnectionType: "local"}
	if err := database.SaveVpsServer(vps1); err != nil {
		t.Fatalf("Erro ao salvar VPS no joao: %v", err)
	}

	serversJoao, _ := database.ListVpsServers()
	if len(serversJoao) != 1 || serversJoao[0].Name != "VPS do João" {
		t.Fatalf("VPS não encontrada no banco do joao")
	}

	// 4. Criar perfil "maria"
	p2 := Profile{Name: "maria", Locale: "en-US"}
	if err := database.SaveProfile(p2); err != nil {
		t.Fatalf("Erro ao criar perfil maria: %v", err)
	}

	mariaDbPath := filepath.Join(profilesDir, "maria.db")
	if _, err := os.Stat(mariaDbPath); os.IsNotExist(err) {
		t.Fatalf("Arquivo maria.db não foi criado no disco")
	}

	// 5. Ativar perfil maria
	profiles, _ = database.ListProfiles()
	var mariaProfile Profile
	for _, p := range profiles {
		if p.Name == "maria" {
			mariaProfile = p
		}
	}
	if err := database.SetActiveProfile(mariaProfile.ID); err != nil {
		t.Fatalf("Erro ao ativar perfil maria: %v", err)
	}

	// Verificar que perfil maria NÃO tem a VPS do joão (isolamento total de banco)
	serversMaria, _ := database.ListVpsServers()
	if len(serversMaria) != 0 {
		t.Fatalf("Perfil maria não deveria ter os dados do joão! Encontrou: %d", len(serversMaria))
	}

	// Salvar VPS na maria
	vpsMaria := VpsServer{Name: "VPS da Maria", ConnectionType: "ssh", Host: "10.0.0.1"}
	if err := database.SaveVpsServer(vpsMaria); err != nil {
		t.Fatalf("Erro ao salvar VPS na maria: %v", err)
	}

	// 6. Voltar para perfil joão e checar se os dados continuam intactos
	if err := database.SetActiveProfile(joaoProfile.ID); err != nil {
		t.Fatalf("Erro ao voltar para joão: %v", err)
	}
	serversJoaoReloaded, _ := database.ListVpsServers()
	if len(serversJoaoReloaded) != 1 || serversJoaoReloaded[0].Name != "VPS do João" {
		t.Fatalf("Dados do joão foram corrompidos ao alternar!")
	}

	// 7. Renomear perfil "joao" para "joaosilva" -> deve renomear joao.db para joaosilva.db
	joaoProfile.Name = "joaosilva"
	if err := database.SaveProfile(joaoProfile); err != nil {
		t.Fatalf("Erro ao renomear perfil joao: %v", err)
	}

	newJoaoDbPath := filepath.Join(profilesDir, "joaosilva.db")
	if _, err := os.Stat(newJoaoDbPath); os.IsNotExist(err) {
		t.Fatalf("Arquivo joaosilva.db não existe após renomear")
	}
	if _, err := os.Stat(joaoDbPath); !os.IsNotExist(err) {
		t.Fatalf("Arquivo antigo joao.db ainda existe após renomear")
	}

	// Verificar se a VPS do joao ainda está acessível no novo db renomeado
	serversRenamed, err := database.ListVpsServers()
	if err != nil || len(serversRenamed) != 1 || serversRenamed[0].Name != "VPS do João" {
		t.Fatalf("Dados da VPS inacessíveis após renomear perfil: %v", err)
	}
}
