package gorms

import (
	"github.com/jybbang/go-core-architecture/core"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func getSqlServerClient(connectionString string) *gorm.DB {
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
	}

	client := clientsInstance.clients[connectionString]
	return client
}

func NewSqlServerAdapter(connectionString string) *adapter {
	sqlserver := &adapter{
		conn: getSqlServerClient(connectionString),
	}

	return sqlserver
}
