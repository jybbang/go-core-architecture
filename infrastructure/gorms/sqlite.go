package gorms

import (
	"github.com/jybbang/go-core-architecture/infrastructure"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getSqliteClient(connectionString string) *gorm.DB {
	clientsInstance := getClients()

	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	_, ok := clientsInstance.clients[connectionString]
	if !ok {
		db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		if err != nil {
			panic("failed to connect database")
		}

		infrastructure.Log.Info("sqlite database connected")
		clientsInstance.clients[connectionString] = tx
	}

	client := clientsInstance.clients[connectionString]
	return client
}

func NewSqliteAdapter(connectionString string) *adapter {
	sqlite := &adapter{
		conn: getSqliteClient(connectionString),
	}

	return sqlite
}
