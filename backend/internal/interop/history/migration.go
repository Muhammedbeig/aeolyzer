package history

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"
)

//go:embed migrations/001_history.sql
var historyMigration string

//go:embed migrations/002_knowledge.sql
var knowledgeMigration string

func Migrate(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return errors.New("database is nil")
	}
	for _, migration := range []string{historyMigration, knowledgeMigration} {
		for _, statement := range strings.Split(migration, ";") {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}
			if _, err := db.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("apply data migration: %w", err)
			}
		}
	}
	return nil
}
