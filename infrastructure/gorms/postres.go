package gorms

import (
	"context"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getPostresClient(settings GormSettings) (*gorm.DB, *sync.RWMutex) {
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
		clientsInstance.mutexes[connectionString] = new(sync.RWMutex)
	}

	client := clientsInstance.clients[connectionString]
	mutex := clientsInstance.mutexes[connectionString]
	return client, mutex
}

func NewPostresAdapter(ctx context.Context, settings GormSettings) *adapter {
	conn, mutex := getPostresClient(settings)
	postgres := &adapter{
		conn:     conn,
		rw:       mutex,
		settings: settings,
	}

	return postgres
}
