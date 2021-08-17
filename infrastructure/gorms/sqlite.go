package gorms

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getSqliteClient(settings GormSettings) *gorm.DB {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	connectionString := settings.ConnectionString
	_, ok := clientsInstance.clients[connectionString]
	if !ok {
		db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		clientsInstance.clients[connectionString] = tx
	}

	client := clientsInstance.clients[connectionString]
	return client
}

func NewSqliteAdapter(ctx context.Context, settings GormSettings) *adapter {
	conn := getSqliteClient(settings)
	sqlite := &adapter{
		conn:     conn,
		settings: settings,
	}

	return sqlite
}
