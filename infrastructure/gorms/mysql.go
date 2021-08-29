package gorms

import (
	"context"

	"gorm.io/driver/mysql"
)

func NewMySqlAdapter(ctx context.Context, settings GormSettings) *adapter {
	connectionString := settings.ConnectionString

	mysql := &adapter{
		dialector: mysql.Open(connectionString),
		settings:  settings,
	}

	mysql.open(ctx)
	return mysql
}
