package gorms

import (
	"context"
	"sync"

	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getMySqlClient(settings GormSettings) (*gorm.DB, *sync.RWMutex) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	connectionString := settings.ConnectionString
	_, ok := clientsInstance.clients[connectionString]
	if !ok {
		db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		core.Log.Infow("mySql database connected")
		clientsInstance.clients[connectionString] = tx
		clientsInstance.mutexes[connectionString] = new(sync.RWMutex)
	}

	client := clientsInstance.clients[connectionString]
	mutex := clientsInstance.mutexes[connectionString]
	return client, mutex
}

func NewMySqlAdapter(ctx context.Context, settings GormSettings) *adapter {
	conn, mutex := getMySqlClient(settings)
	mysql := &adapter{
		conn:     conn,
		rw:       mutex,
		settings: settings,
	}

	return mysql
}
