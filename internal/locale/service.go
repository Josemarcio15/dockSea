package locale

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed default_locales/*.json
var defaultLocalesFS embed.FS

type LocaleService struct {
	localesDir string
}

func NewLocaleService(localesDir string) *LocaleService {
	s := &LocaleService{localesDir: localesDir}
	if err := s.EnsureDefaultLocales(); err != nil {
		fmt.Printf("Aviso: Falha ao inicializar pasta de locales '%s': %v\n", localesDir, err)
	}
	return s
}

// mergeMaps mescla recursivamente novas chaves de src para dst sem apagar o que já existe em dst
func mergeMaps(dst, src map[string]any) bool {
	changed := false
	for k, srcVal := range src {
		if dstVal, exists := dst[k]; exists {
			if srcMap, ok := srcVal.(map[string]any); ok {
				if dstMap, ok := dstVal.(map[string]any); ok {
					if mergeMaps(dstMap, srcMap) {
						changed = true
					}
				}
			}
		} else {
			dst[k] = srcVal
			changed = true
		}
	}
	return changed
}

// EnsureDefaultLocales cria a pasta locales e provisiona ou atualiza os arquivos com novas chaves do binário
func (s *LocaleService) EnsureDefaultLocales() error {
	if err := os.MkdirAll(s.localesDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta de locales: %w", err)
	}

	entries, err := defaultLocalesFS.ReadDir("default_locales")
	if err != nil {
		return fmt.Errorf("falha ao ler locales embutidos: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		targetPath := filepath.Join(s.localesDir, entry.Name())
		embeddedData, err := defaultLocalesFS.ReadFile("default_locales/" + entry.Name())
		if err != nil {
			continue
		}

		// Se não existe, escreve o arquivo padrão do binário
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			_ = os.WriteFile(targetPath, embeddedData, 0644)
		} else {
			// Se já existe, mescla chaves novas que ainda não existam no arquivo do usuário
			existingData, err := os.ReadFile(targetPath)
			if err == nil {
				var existingMap, embeddedMap map[string]any
				if json.Unmarshal(existingData, &existingMap) == nil && json.Unmarshal(embeddedData, &embeddedMap) == nil {
					if mergeMaps(existingMap, embeddedMap) {
						if mergedJSON, err := json.MarshalIndent(existingMap, "", "  "); err == nil {
							_ = os.WriteFile(targetPath, mergedJSON, 0644)
						}
					}
				}
			}
		}
	}
	return nil
}

// GetLocalesDir retorna o caminho da pasta ~/Documents/DockSea/locales
func (s *LocaleService) GetLocalesDir() string {
	return s.localesDir
}

// ListLocales lista todas as linguagens disponíveis no disco (ex: ["pt-BR", "en-US", "es-ES"])
func (s *LocaleService) ListLocales() ([]string, error) {
	if err := os.MkdirAll(s.localesDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao abrir pasta de locales: %w", err)
	}

	entries, err := os.ReadDir(s.localesDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler pasta de locales: %w", err)
	}

	var list []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			list = append(list, strings.TrimSuffix(entry.Name(), ".json"))
		}
	}
	return list, nil
}

// LoadAllLocales carrega todos os arquivos JSON de locale em um mapa de chave-valor { "pt-BR": {...}, "en-US": {...} }
func (s *LocaleService) LoadAllLocales() (map[string]string, error) {
	if err := os.MkdirAll(s.localesDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao abrir pasta de locales: %w", err)
	}

	entries, err := os.ReadDir(s.localesDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler locales: %w", err)
	}

	result := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			localeKey := strings.TrimSuffix(entry.Name(), ".json")
			filePath := filepath.Join(s.localesDir, entry.Name())
			data, err := os.ReadFile(filePath)
			if err == nil {
				result[localeKey] = string(data)
			}
		}
	}
	return result, nil
}

// LoadLocale carrega o conteúdo JSON de um locale específico
func (s *LocaleService) LoadLocale(name string) (string, error) {
	cleanName := strings.TrimSpace(name)
	if !strings.HasSuffix(strings.ToLower(cleanName), ".json") {
		cleanName += ".json"
	}

	filePath := filepath.Join(s.localesDir, cleanName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("falha ao ler arquivo de tradução '%s': %w", cleanName, err)
	}

	return string(data), nil
}

// SaveLocale salva ou atualiza um arquivo de locale no disco
func (s *LocaleService) SaveLocale(name string, jsonContent string) error {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return fmt.Errorf("nome do locale não pode ser vazio")
	}

	if !strings.HasSuffix(strings.ToLower(cleanName), ".json") {
		cleanName += ".json"
	}

	if err := os.MkdirAll(s.localesDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta de locales: %w", err)
	}

	filePath := filepath.Join(s.localesDir, cleanName)
	if err := os.WriteFile(filePath, []byte(jsonContent), 0644); err != nil {
		return fmt.Errorf("falha ao salvar locale '%s': %w", cleanName, err)
	}

	return nil
}
