package mongos

import (
	"context"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/domain"
	"go.uber.org/zap"
)

type adapter struct {
	model      domain.Entity
	conn       *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

type clients struct {
	clients map[string]*mongo.Client
	sync.Mutex
}

var log *zap.SugaredLogger

var clientsSync sync.Once

var clientsInstance *clients

var ctx context.Context

func init() {
	logger, _ := zap.NewProduction()
	log = logger.Sugar()
}

func getClients() *clients {
	if clientsInstance == nil {
		clientsSync.Do(
			func() {
				clientsInstance = &clients{
					clients: make(map[string]*mongo.Client),
				}
			})
	}
	return clientsInstance
}

func getMongoClient(connectionUri string) *mongo.Client {
	clientsInstance := getClients()

	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	_, ok := clientsInstance.clients[connectionUri]
	if !ok {
		ctx = context.Background()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUri))
		if err != nil {
			log.Fatal(err)
		}

		log.Info("mongo database connected")
		clientsInstance.clients[connectionUri] = client
	}

	client := clientsInstance.clients[connectionUri]
	return client
}

func NewMongoAdapter(connectionUri string, dbName string) *adapter {
	client := getMongoClient(connectionUri)
	mongo := &adapter{
		conn:     client,
		database: client.Database(dbName),
	}

	return mongo
}

func (a *adapter) SetModel(model domain.Entity) {
	valueOf := reflect.ValueOf(model)
	key := valueOf.Type().Name()

	a.model = model
	a.collection = a.database.Collection(key)
}

func (a *adapter) Find(model *domain.Entity, dto domain.Copyable, id uuid.UUID) error {
	err := a.collection.FindOne(ctx, bson.M{"_id": id}).Decode(dto)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) Any() (bool, error) {
	count, err := a.Count()
	return count > 0, err
}

func (a *adapter) AnyWithFilter(query interface{}, args interface{}) (bool, error) {
	count, err := a.CountWithFilter(query, args)
	return count > 0, err
}

func (a *adapter) Count() (int64, error) {
	count, err := a.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a *adapter) CountWithFilter(query interface{}, args interface{}) (int64, error) {
	count, err := a.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a *adapter) List(dtos []domain.Copyable) error {
	cursor, err := a.collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	return cursor.All(ctx, dtos)
}

func (a *adapter) ListWithFilter(dtos []domain.Copyable, query interface{}, args interface{}) error {
	cursor, err := a.collection.Find(ctx, query)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	return cursor.All(ctx, dtos)
}

func (a *adapter) Remove(entity *domain.Entity) error {
	_, err := a.collection.DeleteOne(ctx, bson.M{"_id": entity.ID})
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) RemoveRange(entities []*domain.Entity) error {
	vals := make([]bson.M, len(entities))
	for i, entity := range entities {
		vals[i] = bson.M{"_id": entity.ID}
	}
	_, err := a.collection.DeleteMany(ctx, vals)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) Add(entity *domain.Entity) error {
	_, err := a.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) AddRange(entities []*domain.Entity) error {
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

func (a *adapter) Update(entity *domain.Entity) error {
	_, err := a.collection.UpdateOne(ctx, bson.M{"_id": entity.ID}, entity)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) UpdateRange(entities []*domain.Entity) error {
	for _, entity := range entities {
		err := a.Update(entity)
		if err != nil {
			return err
		}
	}

	return nil
}
