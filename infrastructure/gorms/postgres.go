package gorms

import (
	"gorm.io/driver/postgres"
)

func NewPostgresAdapter(settings GormSettings) *adapter {
	return &adapter{
		dialector: postgres.New(postgres.Config{
			DSN:                  settings.ConnectionString,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}),
		settings: settings,
	}
}
