package gorms

import (
	"gorm.io/driver/sqlserver"
)

func NewSqlServerAdapter(settings GormSettings) *adapter {
	return &adapter{
		dialector: sqlserver.Open(settings.ConnectionString),
		settings:  settings,
	}
}
