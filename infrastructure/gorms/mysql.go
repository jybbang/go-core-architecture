package gorms

import (
	"gorm.io/driver/mysql"
)

func NewMySqlAdapter(settings GormSettings) *adapter {
	return &adapter{
		dialector: mysql.Open(settings.ConnectionString),
		settings:  settings,
	}
}
