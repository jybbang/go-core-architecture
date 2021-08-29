package gorms

import (
	"context"

	"gorm.io/driver/sqlserver"
)

func NewSqlServerAdapter(ctx context.Context, settings GormSettings) *adapter {
	connectionString := settings.ConnectionString

	sqlserver := &adapter{
		dialector: sqlserver.Open(connectionString),
		settings:  settings,
	}

	err := sqlserver.setClient(ctx)
	if err != nil {
		panic(err)
	}

	return sqlserver
}
