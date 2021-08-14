package gorms

import (
	"github.com/jybbang/go-core-architecture/infrastructure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getPostresClient(connectionString string) *gorm.DB {
	clientsInstance := getClients()

	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	_, ok := clientsInstance.clients[connectionString]
	if !ok {
		db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		if err != nil {
			panic("failed to connect database")
		}

		infrastructure.Log.Info("postgres database connected")
		clientsInstance.clients[connectionString] = tx
	}

	client := clientsInstance.clients[connectionString]
	return client
}

func NewPostresAdapter(connectionString string) *adapter {
	sqlite := &adapter{
		conn: getPostresClient(connectionString),
	}

	return sqlite
}
