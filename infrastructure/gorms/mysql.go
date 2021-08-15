package gorms

import (
	"context"
	"sync"

	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getMySqlClient(connectionString string) (*gorm.DB, *sync.RWMutex) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

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

func NewMySqlAdapter(ctx context.Context, connectionString string) *adapter {
	conn, mutex := getMySqlClient(connectionString)
	mysql := &adapter{
		conn: conn,
		rw:   mutex,
	}

	return mysql
}
