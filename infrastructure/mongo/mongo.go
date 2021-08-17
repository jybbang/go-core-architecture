package mongo

import (
	"context"
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	tableName  string
	model      core.Entitier
	conn       *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	settings   MongoSettings
}

type clients struct {
	clients map[string]*mongo.Client
	mutex   sync.Mutex
}

type MongoSettings struct {
	ConnectionUri       string
	DatabaseName        string
	CanCreateCollection bool
}

var clientsSync sync.Once

var clientsInstance *clients

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

func getMongoClient(ctx context.Context, settings MongoSettings) *mongo.Client {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	uri := settings.ConnectionUri
	_, ok := clientsInstance.clients[uri]
	if !ok {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}

		clientsInstance.clients[uri] = client
	}

	client := clientsInstance.clients[uri]
	return client
}

func NewMongoAdapter(ctx context.Context, settings MongoSettings) *adapter {
	client := getMongoClient(ctx, settings)
	mongo := &adapter{
		conn:     client,
		database: client.Database(settings.DatabaseName),
		settings: settings,
	}

	return mongo
}

func (a *adapter) SetModel(model core.Entitier, tableName string) {
	a.model = model
	a.tableName = tableName

	ctx := context.Background()
	names, err := a.database.ListCollectionNames(
		ctx,
		bson.D{})
	if err != nil {
		panic(err)
	}

	hasCollection := false
	for _, name := range names {
		if name == a.tableName {
			hasCollection = true
			break
		}
	}

	if !hasCollection && a.settings.CanCreateCollection {
		err := a.database.CreateCollection(ctx, a.tableName)
		if err != nil {
			panic(err)
		}
	}

	a.collection = a.database.Collection(a.tableName)
}

func (a *adapter) Find(ctx context.Context, id uuid.UUID, dest core.Entitier) (err error) {
	err = a.collection.FindOne(ctx, bson.M{"entity._id": id}).Decode(dest)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return core.ErrNotFound
		}
		return err
	}
	return nil
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

func (a *adapter) List(ctx context.Context, dest interface{}) (err error) {
	cursor, err := a.collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, dest)

	return err
}

func (a *adapter) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) (err error) {
	cursor, err := a.collection.Find(ctx, query)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, dest)

	return err
}

func (a *adapter) Remove(ctx context.Context, id uuid.UUID) error {
	_, err := a.collection.DeleteOne(ctx, bson.M{"entity._id": id})
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) RemoveRange(ctx context.Context, ids []uuid.UUID) error {
	for _, id := range ids {
		err := a.Remove(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *adapter) Add(ctx context.Context, entity core.Entitier) error {
	_, err := a.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
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
	return a.collection.FindOneAndReplace(ctx, bson.M{"entity._id": entity.GetID()}, entity).Err()
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
