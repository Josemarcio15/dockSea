package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ThemeService struct {
	themesDir string
}

func NewThemeService(themesDir string) *ThemeService {
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		fmt.Printf("Aviso: Falha ao criar diretório de temas '%s': %v\n", themesDir, err)
	}
	return &ThemeService{themesDir: themesDir}
}

// GetThemesDir retorna o caminho absoluto da pasta themes
func (s *ThemeService) GetThemesDir() string {
	return s.themesDir
}

// ListThemes lista todos os temas .json disponíveis na pasta ~/Documents/DockSea/themes/
func (s *ThemeService) ListThemes() ([]string, error) {
	if err := os.MkdirAll(s.themesDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao abrir pasta de temas: %w", err)
	}

	entries, err := os.ReadDir(s.themesDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler temas: %w", err)
	}

	var themes []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			themes = append(themes, strings.TrimSuffix(entry.Name(), ".json"))
		}
	}
	return themes, nil
}

// LoadTheme carrega o conteúdo de um arquivo de tema JSON
func (s *ThemeService) LoadTheme(name string) (string, error) {
	cleanName := strings.TrimSpace(name)
	if !strings.HasSuffix(strings.ToLower(cleanName), ".json") {
		cleanName += ".json"
	}

	filePath := filepath.Join(s.themesDir, cleanName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("falha ao ler arquivo de tema '%s': %w", cleanName, err)
	}

	return string(data), nil
}

// SaveTheme grava ou atualiza um arquivo de tema JSON na pasta themes
func (s *ThemeService) SaveTheme(name string, jsonContent string) error {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return fmt.Errorf("nome do tema não pode ser vazio")
	}
	if !strings.HasSuffix(strings.ToLower(cleanName), ".json") {
		cleanName += ".json"
	}

	if err := os.MkdirAll(s.themesDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta de temas: %w", err)
	}

	filePath := filepath.Join(s.themesDir, cleanName)
	if err := os.WriteFile(filePath, []byte(jsonContent), 0644); err != nil {
		return fmt.Errorf("falha ao salvar tema '%s': %w", cleanName, err)
	}

	return nil
}

// DeleteTheme remove um arquivo de tema
func (s *ThemeService) DeleteTheme(name string) error {
	cleanName := strings.TrimSpace(name)
	if !strings.HasSuffix(strings.ToLower(cleanName), ".json") {
		cleanName += ".json"
	}

	filePath := filepath.Join(s.themesDir, cleanName)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("falha ao excluir tema: %w", err)
	}
	return nil
}
