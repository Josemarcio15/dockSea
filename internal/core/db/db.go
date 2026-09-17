package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	mu            sync.RWMutex
	conn          *sql.DB
	masterConn    *sql.DB
	appDir        string
	profilesDir   string
	themesDir     string
	activeProfile Profile
}

var invalidFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

// SanitizeProfileName converts a profile name into a safe file name base
func SanitizeProfileName(name string) string {
	cleaned := strings.TrimSpace(name)
	cleaned = strings.ToLower(cleaned)
	cleaned = strings.ReplaceAll(cleaned, " ", "_")
	cleaned = invalidFilenameChars.ReplaceAllString(cleaned, "")
	if cleaned == "" {
		cleaned = "profile"
	}
	return cleaned
}

func getAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		configDir, err2 := os.UserConfigDir()
		if err2 != nil {
			return ".", nil
		}
		return filepath.Join(configDir, "docksea"), nil
	}
	return filepath.Join(homeDir, "Documents", "DockSea"), nil
}

func InitDB() (*DB, error) {
	appDir, err := getAppDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine app dir: %w", err)
	}

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create docksea dir in Documents: %w", err)
	}

	profilesDir := filepath.Join(appDir, "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profiles dir: %w", err)
	}

	themesDir := filepath.Join(appDir, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create themes dir: %w", err)
	}

	// Migrar dados do antigo UserConfigDir se existir
	if configDir, err := os.UserConfigDir(); err == nil {
		oldAppDir := filepath.Join(configDir, "docksea")
		if oldAppDir != appDir {
			oldTenants := filepath.Join(oldAppDir, "tenants.db")
			newTenants := filepath.Join(appDir, "tenants.db")
			if _, err := os.Stat(oldTenants); err == nil {
				if _, err := os.Stat(newTenants); os.IsNotExist(err) {
					_ = copyFile(oldTenants, newTenants)
				}
			}
			oldProfilesDir := filepath.Join(oldAppDir, "profiles")
			if entries, err := os.ReadDir(oldProfilesDir); err == nil {
				for _, entry := range entries {
					oldP := filepath.Join(oldProfilesDir, entry.Name())
					newP := filepath.Join(profilesDir, entry.Name())
					if _, err := os.Stat(newP); os.IsNotExist(err) {
						_ = copyFile(oldP, newP)
					}
				}
			}
		}
	}

	// 1. Open master registry DB (profiles catalog)
	masterDbPath := filepath.Join(appDir, "tenants.db")
	masterConn, err := sql.Open("sqlite", masterDbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tenants database: %w", err)
	}

	if _, err := masterConn.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		masterConn.Close()
		return nil, fmt.Errorf("failed to configure master sqlite pragma: %w", err)
	}

	d := &DB{
		masterConn:  masterConn,
		appDir:      appDir,
		profilesDir: profilesDir,
		themesDir:   themesDir,
	}

	// Migrate master table
	if err := d.migrateMaster(); err != nil {
		masterConn.Close()
		return nil, fmt.Errorf("failed to migrate master database: %w", err)
	}

	// Check if we need to migrate from legacy docksea.db
	legacyDbPath := filepath.Join(appDir, "docksea.db")
	if _, err := os.Stat(legacyDbPath); err == nil {
		_ = d.migrateLegacyDB(legacyDbPath)
	}

	// Ensure default profile exists if master has no profiles
	profiles, err := d.ListProfiles()
	if err != nil || len(profiles) == 0 {
		now := time.Now().UTC()
		defaultProf := Profile{
			ID:        "default",
			Name:      "default",
			Locale:    "pt-BR",
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := d.createProfileDBFile(defaultProf); err != nil {
			return nil, fmt.Errorf("failed to create default profile db: %w", err)
		}
		if err := d.saveProfileToMaster(defaultProf); err != nil {
			return nil, fmt.Errorf("failed to save default profile to master: %w", err)
		}
	}

	// Get active profile
	activeProf, err := d.GetActiveProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to get active profile: %w", err)
	}

	// Connect active profile database
	if err := d.switchProfileConnection(activeProf); err != nil {
		return nil, fmt.Errorf("failed to connect to active profile db: %w", err)
	}

	return d, nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func (d *DB) GetThemesDir() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.themesDir
}

func (d *DB) GetAppDir() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.appDir
}

func (d *DB) GetDBPath() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.getProfileDBPath(d.activeProfile.Name)
}

func (d *DB) getProfileDBPath(name string) string {
	fileName := SanitizeProfileName(name) + ".db"
	return filepath.Join(d.profilesDir, fileName)
}

func (d *DB) switchProfileConnection(p *Profile) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	profileDbPath := d.getProfileDBPath(p.Name)
	newConn, err := sql.Open("sqlite", profileDbPath)
	if err != nil {
		return fmt.Errorf("failed to open profile sqlite: %w", err)
	}

	if _, err := newConn.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		newConn.Close()
		return fmt.Errorf("failed to configure profile sqlite pragma: %w", err)
	}

	// Ensure profile tables are migrated
	if err := migrateProfileTables(newConn); err != nil {
		newConn.Close()
		return fmt.Errorf("failed to migrate profile tables: %w", err)
	}

	// Close old connection if open
	if d.conn != nil {
		_ = d.conn.Close()
	}

	d.conn = newConn
	d.activeProfile = *p
	return nil
}

func (d *DB) createProfileDBFile(p Profile) error {
	profileDbPath := d.getProfileDBPath(p.Name)
	conn, err := sql.Open("sqlite", profileDbPath)
	if err != nil {
		return fmt.Errorf("failed to create profile db: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		return err
	}

	return migrateProfileTables(conn)
}

func (d *DB) migrateLegacyDB(legacyPath string) error {
	// If legacy docksea.db exists, rename it to default.db if default.db doesn't exist
	defaultDbPath := filepath.Join(d.profilesDir, "default.db")
	if _, err := os.Stat(defaultDbPath); os.IsNotExist(err) {
		_ = os.Rename(legacyPath, defaultDbPath)
		_ = os.Remove(legacyPath + "-wal")
		_ = os.Remove(legacyPath + "-shm")
	}
	return nil
}

// Backup realiza um checkpoint completo do WAL e salva uma cópia limpa do banco do profile ativo
func (d *DB) Backup(destPath string) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.conn == nil {
		return fmt.Errorf("conexão com o banco não inicializada")
	}

	// Forçar flush do WAL antes do backup
	if _, err := d.conn.Exec(`PRAGMA wal_checkpoint(TRUNCATE);`); err != nil {
		return fmt.Errorf("falha ao sincronizar WAL: %w", err)
	}

	_ = os.Remove(destPath)
	_, err := d.conn.Exec(`VACUUM INTO ?;`, destPath)
	if err != nil {
		return fmt.Errorf("falha ao exportar banco de dados: %w", err)
	}

	return nil
}

// Restore restaura o banco de dados do profile ativo a partir de um arquivo .db externo
func (d *DB) Restore(sourcePath string) error {
	if sourcePath == "" {
		return fmt.Errorf("caminho de origem inválido")
	}

	srcConn, err := sql.Open("sqlite", sourcePath)
	if err != nil {
		return fmt.Errorf("arquivo de backup inválido: %w", err)
	}
	var testCount int
	err = srcConn.QueryRow(`SELECT count(*) FROM sqlite_master;`).Scan(&testCount)
	srcConn.Close()
	if err != nil {
		return fmt.Errorf("o arquivo selecionado não é um banco SQLite válido: %w", err)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.conn != nil {
		_ = d.conn.Close()
	}

	currentDbPath := d.getProfileDBPath(d.activeProfile.Name)
	_ = os.Remove(currentDbPath + "-wal")
	_ = os.Remove(currentDbPath + "-shm")

	srcData, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("falha ao ler arquivo de backup: %w", err)
	}

	if err := os.WriteFile(currentDbPath, srcData, 0644); err != nil {
		return fmt.Errorf("falha ao sobrescrever banco de dados: %w", err)
	}

	conn, err := sql.Open("sqlite", currentDbPath)
	if err != nil {
		return fmt.Errorf("falha ao reabrir banco após restauração: %w", err)
	}

	if _, err := conn.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		return fmt.Errorf("falha ao reconfigurar sqlite: %w", err)
	}

	d.conn = conn
	return migrateProfileTables(d.conn)
}

// Reset apaga os dados de todas as tabelas do profile ativo e reaplica o schema
func (d *DB) Reset() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.conn == nil {
		return fmt.Errorf("conexão com o banco não inicializada")
	}

	tables := []string{
		"vps_servers",
		"app_settings",
		"image_history",
		"saved_paths",
		"stacks",
		"container_configs",
	}

	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`PRAGMA foreign_keys = OFF;`); err != nil {
		return err
	}

	for _, table := range tables {
		if _, err := tx.Exec(fmt.Sprintf("DELETE FROM %s;", table)); err != nil {
			return fmt.Errorf("falha ao limpar tabela %s: %w", table, err)
		}
	}

	if _, err := tx.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return migrateProfileTables(d.conn)
}

func (d *DB) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.conn != nil {
		_ = d.conn.Close()
	}
	if d.masterConn != nil {
		_ = d.masterConn.Close()
	}
	return nil
}
