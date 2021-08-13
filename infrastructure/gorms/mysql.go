package gorms

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getMySqlClient(connectionString string) *gorm.DB {
	clientsInstance := getClients()

	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	_, ok := clientsInstance.clients[connectionString]
	if !ok {
		db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
		tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})

		if err != nil {
			panic("failed to connect database")
		}

		log.Info("mySql database connected")
		clientsInstance.clients[connectionString] = tx
	}

	client := clientsInstance.clients[connectionString]
	return client
}

func NewMySqlAdapter(connectionString string) *adapter {
	sqlite := &adapter{
		conn: getMySqlClient(connectionString),
	}

	return sqlite
}
