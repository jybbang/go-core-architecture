package mongos

import (
	"context"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	model      core.Entitier
	conn       *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	rw         *sync.RWMutex
}

type clients struct {
	clients map[string]*mongo.Client
	mutexes map[string]*sync.RWMutex
	sync.Mutex
}

var clientsSync sync.Once

var clientsInstance *clients

func getClients() *clients {
	if clientsInstance == nil {
		clientsSync.Do(
			func() {
				clientsInstance = &clients{
					clients: make(map[string]*mongo.Client),
					mutexes: make(map[string]*sync.RWMutex),
				}
			})
	}
	return clientsInstance
}

func getMongoClient(ctx context.Context, connectionUri string) (*mongo.Client, *sync.RWMutex) {
	clientsInstance := getClients()

	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	_, ok := clientsInstance.clients[connectionUri]
	if !ok {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUri))
		if err != nil {
			panic(err)
		}
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}

		core.Log.Info("mongo database connected")
		clientsInstance.clients[connectionUri] = client
		clientsInstance.mutexes[connectionUri] = new(sync.RWMutex)
	}

	client := clientsInstance.clients[connectionUri]
	mutex := clientsInstance.mutexes[connectionUri]
	return client, mutex
}

func NewMongoAdapter(ctx context.Context, connectionUri string, dbName string) *adapter {
	client, mutex := getMongoClient(ctx, connectionUri)
	mongo := &adapter{
		conn:     client,
		database: client.Database(dbName),
		rw:       mutex,
	}

	return mongo
}

func (a *adapter) SetModel(model core.Entitier) {
	valueOf := reflect.ValueOf(model)
	key := valueOf.Type().Name()

	a.model = model
	a.collection = a.database.Collection(key)
}

func (a *adapter) Find(ctx context.Context, dest core.Entitier, id uuid.UUID) (ok bool, err error) {
	err = a.collection.FindOne(ctx, bson.M{"_id": id}).Decode(dest)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *adapter) Any(ctx context.Context) (ok bool, err error) {
	count, err := a.Count(ctx)
	return count > 0, err
}

func (a *adapter) AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error) {
	count, err := a.CountWithFilter(ctx, query, args)
	return count > 0, err
}

func (a *adapter) Count(ctx context.Context) (count int64, err error) {
	count, err = a.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	count, err = a.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a *adapter) List(ctx context.Context, dest []core.Entitier) error {
	cursor, err := a.collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	return cursor.All(ctx, dest)
}

func (a *adapter) ListWithFilter(ctx context.Context, dest []core.Entitier, query interface{}, args interface{}) error {
	cursor, err := a.collection.Find(ctx, query)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	return cursor.All(ctx, dest)
}

func (a *adapter) Remove(ctx context.Context, entity core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	_, err := a.collection.DeleteOne(ctx, bson.M{"_id": entity.GetID()})
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) RemoveRange(ctx context.Context, entities []core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	vals := make([]bson.M, len(entities))
	for i, entity := range entities {
		vals[i] = bson.M{"_id": entity.GetID()}
	}
	_, err := a.collection.DeleteMany(ctx, vals)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) Add(ctx context.Context, entity core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	_, err := a.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	vals := make([]interface{}, len(entities))
	for i, entity := range entities {
		vals[i] = entity
	}
	_, err := a.collection.InsertMany(ctx, vals)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) Update(ctx context.Context, entity core.Entitier) error {
	a.rw.Lock()
	defer a.rw.Unlock()

	_, err := a.collection.UpdateOne(ctx, bson.M{"_id": entity.GetID()}, entity)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	for _, entity := range entities {
		err := a.Update(ctx, entity)
		if err != nil {
			return err
		}
	}

	return nil
}
