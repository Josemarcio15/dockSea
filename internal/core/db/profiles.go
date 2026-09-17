package db

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Locale    string    `json:"locale"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (d *DB) ListProfiles() ([]Profile, error) {
	conn := d.masterConn
	if conn == nil {
		conn = d.conn
	}
	if conn == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	rows, err := conn.Query(`
		SELECT id, name, locale, is_active, created_at, updated_at
		FROM profiles
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var p Profile
		var isActiveInt int
		if err := rows.Scan(&p.ID, &p.Name, &p.Locale, &isActiveInt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.IsActive = isActiveInt == 1
		profiles = append(profiles, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (d *DB) GetActiveProfile() (*Profile, error) {
	conn := d.masterConn
	if conn == nil {
		conn = d.conn
	}
	if conn == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	var p Profile
	var isActiveInt int
	err := conn.QueryRow(`
		SELECT id, name, locale, is_active, created_at, updated_at
		FROM profiles
		WHERE is_active = 1
		LIMIT 1
	`).Scan(&p.ID, &p.Name, &p.Locale, &isActiveInt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		// Se nenhum estiver ativo, pega o primeiro da lista
		err = conn.QueryRow(`
			SELECT id, name, locale, is_active, created_at, updated_at
			FROM profiles
			ORDER BY created_at ASC
			LIMIT 1
		`).Scan(&p.ID, &p.Name, &p.Locale, &isActiveInt, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	p.IsActive = true
	return &p, nil
}

func (d *DB) SaveProfile(p Profile) error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return fmt.Errorf("o nome do perfil é obrigatório")
	}

	now := time.Now().UTC()
	isNew := p.ID == ""
	if isNew {
		p.ID = fmt.Sprintf("prof_%d", now.UnixNano())
	}
	if p.Locale == "" {
		p.Locale = "pt-BR"
	}

	conn := d.masterConn
	if conn == nil {
		conn = d.conn
	}

	// 1. Checar se já existe outro perfil com o mesmo nome (username.db duplicado)
	var existingID string
	err := conn.QueryRow(`SELECT id FROM profiles WHERE LOWER(name) = LOWER(?)`, p.Name).Scan(&existingID)
	if err == nil && existingID != p.ID {
		return fmt.Errorf("já existe um perfil com o nome '%s'", p.Name)
	}

	// Se for novo perfil, verificar se o arquivo .db já existe no disco
	if isNew && d.profilesDir != "" {
		targetFile := d.getProfileDBPath(p.Name)
		if _, err := os.Stat(targetFile); err == nil {
			return fmt.Errorf("o arquivo de banco '%s.db' já existe no diretório de perfis", SanitizeProfileName(p.Name))
		}
	}

	// Se for edição e o nome mudou, renomear o arquivo .db correspondente
	if !isNew && d.profilesDir != "" {
		var oldName string
		_ = conn.QueryRow(`SELECT name FROM profiles WHERE id = ?`, p.ID).Scan(&oldName)
		if oldName != "" && !strings.EqualFold(oldName, p.Name) {
			oldPath := d.getProfileDBPath(oldName)
			newPath := d.getProfileDBPath(p.Name)
			
			// Se o perfil que está sendo renomeado for o que está aberto no momento, fechar a conexão antes de renomear
			d.mu.Lock()
			if d.conn != nil && strings.EqualFold(d.activeProfile.Name, oldName) {
				_ = d.conn.Close()
				d.conn = nil
			}
			d.mu.Unlock()

			if _, err := os.Stat(oldPath); err == nil {
				_ = os.Rename(oldPath, newPath)
				_ = os.Remove(oldPath + "-wal")
				_ = os.Remove(oldPath + "-shm")
			}
		}
	}

	// 2. Salvar no master
	if err := d.saveProfileToMaster(p); err != nil {
		return err
	}

	// Se o perfil ativo foi renomeado, reabrir a conexão com o novo nome
	if !isNew && d.profilesDir != "" {
		d.mu.RLock()
		isActiveProf := strings.EqualFold(d.activeProfile.ID, p.ID) || p.IsActive
		d.mu.RUnlock()
		if isActiveProf {
			_ = d.switchProfileConnection(&p)
		}
	}

	// 3. Se for novo perfil, criar o arquivo .db e inicializar tabelas
	if isNew && d.profilesDir != "" {
		if err := d.createProfileDBFile(p); err != nil {
			return fmt.Errorf("falha ao criar arquivo de banco para o perfil: %w", err)
		}
	}

	return nil
}

func (d *DB) saveProfileToMaster(p Profile) error {
	conn := d.masterConn
	if conn == nil {
		conn = d.conn
	}
	now := time.Now().UTC()
	isActiveInt := 0
	if p.IsActive {
		isActiveInt = 1
	}

	_, err := conn.Exec(`
		INSERT INTO profiles (id, name, locale, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			locale = excluded.locale,
			updated_at = excluded.updated_at
	`, p.ID, p.Name, p.Locale, isActiveInt, now, now)
	return err
}

func (d *DB) DeleteProfile(id string) error {
	conn := d.masterConn
	if conn == nil {
		conn = d.conn
	}

	var count int
	err := conn.QueryRow(`SELECT COUNT(*) FROM profiles`).Scan(&count)
	if err != nil {
		return err
	}
	if count <= 1 {
		return fmt.Errorf("não é possível excluir o único perfil existente")
	}

	var prof Profile
	var isActiveInt int
	err = conn.QueryRow(`SELECT id, name, locale, is_active FROM profiles WHERE id = ?`, id).Scan(&prof.ID, &prof.Name, &prof.Locale, &isActiveInt)
	if err != nil {
		return fmt.Errorf("perfil não encontrado: %w", err)
	}
	prof.IsActive = isActiveInt == 1

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM profiles WHERE id = ?`, id); err != nil {
		return err
	}

	if prof.IsActive {
		if _, err := tx.Exec(`UPDATE profiles SET is_active = 1 WHERE id IN (SELECT id FROM profiles LIMIT 1)`); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Se o perfil ativo foi deletado, chavear conexão para o novo ativo
	if prof.IsActive && d.masterConn != nil {
		newActive, err := d.GetActiveProfile()
		if err == nil {
			_ = d.switchProfileConnection(newActive)
		}
	}

	// Remover arquivo .db e arquivos auxiliares do perfil deletado
	if d.profilesDir != "" {
		dbPath := d.getProfileDBPath(prof.Name)
		_ = os.Remove(dbPath)
		_ = os.Remove(dbPath + "-wal")
		_ = os.Remove(dbPath + "-shm")
	}

	return nil
}

func (d *DB) SetActiveProfile(id string) error {
	conn := d.masterConn
	if conn == nil {
		conn = d.conn
	}

	var targetProfile Profile
	var isActiveInt int
	err := conn.QueryRow(`SELECT id, name, locale, is_active, created_at, updated_at FROM profiles WHERE id = ?`, id).Scan(
		&targetProfile.ID, &targetProfile.Name, &targetProfile.Locale, &isActiveInt, &targetProfile.CreatedAt, &targetProfile.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("perfil não encontrado: %w", err)
	}

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE profiles SET is_active = 0`); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE profiles SET is_active = 1 WHERE id = ?`, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	targetProfile.IsActive = true

	// Chavear a conexão ativa para o banco deste perfil
	if d.profilesDir != "" {
		if err := d.switchProfileConnection(&targetProfile); err != nil {
			return fmt.Errorf("erro ao alternar banco de dados para '%s': %w", targetProfile.Name, err)
		}
	}

	return nil
}

func (d *DB) UpdateProfileLocale(id string, locale string) error {
	conn := d.masterConn
	if conn == nil {
		conn = d.conn
	}
	now := time.Now().UTC()
	_, err := conn.Exec(`
		UPDATE profiles SET locale = ?, updated_at = ? WHERE id = ?
	`, locale, now, id)
	return err
}
