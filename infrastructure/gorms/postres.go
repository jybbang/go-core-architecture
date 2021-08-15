package gorms

import (
	"context"
	"sync"

	"github.com/jybbang/go-core-architecture/core"
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
		db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		core.Log.Infow("postgres database connected")
		clientsInstance.clients[connectionString] = tx
		clientsInstance.mutexes[connectionString] = new(sync.RWMutex)
	}

	client := clientsInstance.clients[connectionString]
	mutex := clientsInstance.mutexes[connectionString]
	return client, mutex
}

func NewPostresAdapter(ctx context.Context, settings GormSettings) *adapter {
	conn, mutex := getMySqlClient(settings)
	postgres := &adapter{
		conn: conn,
		rw:   mutex,
	}

	return postgres
}
