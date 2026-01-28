package migrations

import (
	"database/sql"
)

// Migration 介面定義
type Migration interface {
	Up(db *sql.DB) error
	Down(db *sql.DB) error
	Version() string
	Description() string
}

// BaseMigration 基礎結構
type BaseMigration struct {
	version     string
	description string
}

func (m *BaseMigration) Version() string {
	return m.version
}

func (m *BaseMigration) Description() string {
	return m.description
}
