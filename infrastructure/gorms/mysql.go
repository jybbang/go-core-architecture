package gorms

import (
	"context"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getMySqlClient(settings GormSettings) *gorm.DB {
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

		clientsInstance.clients[connectionString] = tx
	}

	client := clientsInstance.clients[connectionString]
	return client
}

func NewMySqlAdapter(ctx context.Context, settings GormSettings) *adapter {
	conn := getMySqlClient(settings)
	mysql := &adapter{
		conn:     conn,
		settings: settings,
	}

	return mysql
}
