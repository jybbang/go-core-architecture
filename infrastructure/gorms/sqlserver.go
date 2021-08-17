package gorms

import (
	"context"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func getSqlServerClient(settings GormSettings) *gorm.DB {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	connectionString := settings.ConnectionString
	_, ok := clientsInstance.clients[connectionString]
	if !ok {
		db, err := gorm.Open(sqlserver.Open(connectionString), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		clientsInstance.clients[connectionString] = tx
	}

	client := clientsInstance.clients[connectionString]
	return client
}

func NewSqlServerAdapter(ctx context.Context, settings GormSettings) *adapter {
	conn := getSqlServerClient(settings)
	sqlserver := &adapter{
		conn:     conn,
		settings: settings,
	}

	return sqlserver
}
