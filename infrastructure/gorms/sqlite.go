package gorms

import (
	"sync"

	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getSqliteClient(connectionString string) (*gorm.DB, *sync.RWMutex) {
	clientsInstance := getClients()

	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	_, ok := clientsInstance.clients[connectionString]
	if !ok {
		db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		core.Log.Info("sqlite database connected")
		clientsInstance.clients[connectionString] = tx
		clientsInstance.mutexes[connectionString] = new(sync.RWMutex)
	}

	client := clientsInstance.clients[connectionString]
	mutex := clientsInstance.mutexes[connectionString]
	return client, mutex
}

func NewSqliteAdapter(connectionString string) *adapter {
	conn, mutex := getMySqlClient(connectionString)
	sqlite := &adapter{
		conn: conn,
		rw:   mutex,
	}

	return sqlite
}
