package db

import (
	"database/sql"
)

func (d *DB) migrateMaster() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			locale TEXT NOT NULL DEFAULT 'pt-BR',
			is_active INTEGER DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS app_settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`,
	}

	for _, query := range queries {
		if _, err := d.masterConn.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func migrateProfileTables(conn *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS vps_servers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			connection_type TEXT NOT NULL,
			host TEXT,
			port INTEGER DEFAULT 22,
			username TEXT,
			auth_type TEXT,
			ssh_key_path TEXT,
			ssh_key_passphrase TEXT,
			ssh_password TEXT,
			sudo_password TEXT,
			docker_socket_path TEXT NOT NULL DEFAULT '/var/run/docker.sock',
			docker_path TEXT NOT NULL DEFAULT '/usr/bin/docker',
			docker_compose_path TEXT,
			is_active INTEGER DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			locale TEXT NOT NULL DEFAULT 'pt-BR',
			is_active INTEGER DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS app_settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS image_history (
			id TEXT PRIMARY KEY,
			image_name TEXT NOT NULL,
			profile_id TEXT NOT NULL DEFAULT 'default',
			pulled_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS saved_paths (
			id TEXT PRIMARY KEY,
			path TEXT NOT NULL,
			label TEXT NOT NULL,
			profile_id TEXT NOT NULL DEFAULT 'default',
			created_at DATETIME NOT NULL,
			UNIQUE(profile_id, path)
		);`,
		`CREATE TABLE IF NOT EXISTS stacks (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			project_name TEXT NOT NULL,
			source_type TEXT NOT NULL DEFAULT 'editor',
			folder_path TEXT NOT NULL DEFAULT '',
			yaml_content TEXT NOT NULL DEFAULT '',
			profile_id TEXT NOT NULL DEFAULT 'default',
			last_deployed_yaml TEXT DEFAULT '',
			last_deployed_config_yaml TEXT DEFAULT '',
			last_deployed_config_json TEXT DEFAULT '',
			last_deployed_remote_dir TEXT DEFAULT '',
			last_deployed_at DATETIME,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS container_configs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			profile_id TEXT NOT NULL,
			image TEXT NOT NULL,
			container_name TEXT NOT NULL,
			ports TEXT,
			env TEXT,
			network TEXT,
			volumes TEXT,
			restart_policy TEXT,
			project_name TEXT,
			description TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			command TEXT
		);`,
	}

	for _, query := range queries {
		if _, err := conn.Exec(query); err != nil {
			return err
		}
	}

	alterColumns := []string{
		`ALTER TABLE stacks ADD COLUMN last_deployed_yaml TEXT DEFAULT '';`,
		`ALTER TABLE stacks ADD COLUMN last_deployed_config_yaml TEXT DEFAULT '';`,
		`ALTER TABLE stacks ADD COLUMN last_deployed_config_json TEXT DEFAULT '';`,
		`ALTER TABLE stacks ADD COLUMN last_deployed_remote_dir TEXT DEFAULT '';`,
		`ALTER TABLE stacks ADD COLUMN last_deployed_at DATETIME;`,
	}
	for _, alterQuery := range alterColumns {
		_, _ = conn.Exec(alterQuery)
	}

	return nil
}
