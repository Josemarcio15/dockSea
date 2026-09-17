package theme

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed default_theme.json
var DefaultThemeJson []byte

// EnsureDefaultThemeFile escreve o tema predefinido "default.json" na pasta ~/Documents/DockSea/themes caso não exista
func EnsureDefaultThemeFile(themesDir string) error {
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta de temas: %w", err)
	}

	defaultFile := filepath.Join(themesDir, "default.json")
	if _, err := os.Stat(defaultFile); os.IsNotExist(err) {
		if err := os.WriteFile(defaultFile, DefaultThemeJson, 0644); err != nil {
			return fmt.Errorf("falha ao gravar default.json: %w", err)
		}
	}
	return nil
}
