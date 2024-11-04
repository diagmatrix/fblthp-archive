package db

import (
	"os"
	"time"
)

type Migration struct {
	ID        uint
	Name      string
	CreatedAt time.Time
}

func (m *Migration) GetOperation(operation MigrationOperation) (string, error) {
	fileName := MIGRATION_DIRECTORY + "/" + MIGRATION_IDENTIFIER + m.Name + "/" + string(operation) + ".sql"
	file, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
