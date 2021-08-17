package gorms

import (
	"context"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func getSqlServerClient(settings GormSettings) (*gorm.DB, *sync.RWMutex) {
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
		clientsInstance.mutexes[connectionString] = new(sync.RWMutex)
	}

	client := clientsInstance.clients[connectionString]
	mutex := clientsInstance.mutexes[connectionString]
	return client, mutex
}

func NewSqlServerAdapter(ctx context.Context, settings GormSettings) *adapter {
	conn, mutex := getSqlServerClient(settings)
	sqlserver := &adapter{
		conn:     conn,
		rw:       mutex,
		settings: settings,
	}

	return sqlserver
}
