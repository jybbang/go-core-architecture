package gorms

import (
	"context"

	"gorm.io/driver/sqlite"
)

func NewSqliteAdapter(ctx context.Context, settings GormSettings) *adapter {
	connectionString := settings.ConnectionString

	sqlite := &adapter{
		dialector: sqlite.Open(connectionString),
		settings:  settings,
	}

	sqlite.open(ctx)
	return sqlite
}
