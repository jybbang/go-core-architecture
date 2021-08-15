package gorms

import (
	"sync"

	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func getSqlServerClient(connectionString string) (*gorm.DB, *sync.RWMutex) {
	clientsInstance := getClients()

	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	_, ok := clientsInstance.clients[connectionString]
	if !ok {
		db, err := gorm.Open(sqlserver.Open(connectionString), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		core.Log.Info("sqlserver database connected")
		clientsInstance.clients[connectionString] = tx
		clientsInstance.mutexes[connectionString] = new(sync.RWMutex)
	}

	client := clientsInstance.clients[connectionString]
	mutex := clientsInstance.mutexes[connectionString]
	return client, mutex
}

func NewSqlServerAdapter(connectionString string) *adapter {
	conn, mutex := getMySqlClient(connectionString)
	sqlserver := &adapter{
		conn: conn,
		rw:   mutex,
	}

	return sqlserver
}
