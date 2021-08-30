package gorms

import (
	"gorm.io/driver/sqlite"
)

func NewSqliteAdapter(settings GormSettings) *adapter {
	return &adapter{
		dialector: sqlite.Open(settings.ConnectionString),
		settings:  settings,
	}
}
