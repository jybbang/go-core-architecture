package gorms

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getPostgresClient(settings GormSettings) *gorm.DB {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	connectionString := settings.ConnectionString
	_, ok := clientsInstance.clients[connectionString]
	if !ok {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  connectionString,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		clientsInstance.clients[connectionString] = tx
	}

	client := clientsInstance.clients[connectionString]
	return client
}

func NewPostgresAdapter(ctx context.Context, settings GormSettings) *adapter {
	conn := getPostgresClient(settings)
	postgres := &adapter{
		conn:     conn,
		settings: settings,
	}

	return postgres
}
