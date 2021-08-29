package gorms

import (
	"context"

	"gorm.io/driver/postgres"
)

func NewPostgresAdapter(ctx context.Context, settings GormSettings) *adapter {
	connectionString := settings.ConnectionString

	postgres := &adapter{
		dialector: postgres.New(postgres.Config{
			DSN:                  connectionString,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}),
		settings: settings,
	}

	postgres.open(ctx)
	return postgres
}
