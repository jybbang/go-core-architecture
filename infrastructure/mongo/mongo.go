package mongo

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	tableName string
	model     core.Entitier
	client    *client
	settings  MongoSettings
	mutex     sync.Mutex
}

type client struct {
	conn       *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	isOpened   bool
}

type clients struct {
	clients map[string]*client
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
					clients: make(map[string]*client),
				}
			})
	}
	return clientsInstance
}

func NewMongoAdapter(ctx context.Context, settings MongoSettings) *adapter {
	mongoService := &adapter{
		settings: settings,
	}
	mongoService.open(ctx)

	return mongoService
}

func (a *adapter) open(ctx context.Context) {
	clientsInstance := getClients()

	clientsInstance.mutex.Lock()
	defer clientsInstance.mutex.Unlock()

	uri := a.settings.ConnectionUri

	if strings.TrimSpace(uri) == "" {
		panic("uri is required")
	}

	cli, ok := clientsInstance.clients[uri]
	if !ok || !cli.isOpened {
		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			panic(err)
		}

		clientsInstance.clients[uri] = &client{
			conn:     mongoClient,
			database: mongoClient.Database(a.settings.DatabaseName),
			isOpened: true,
		}
	}

	client := clientsInstance.clients[uri]

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.client = client
	if a.model != nil {
		a.SetModel(a.model, a.tableName)
	}
}

func (a *adapter) OnCircuitOpen() {
	a.client.isOpened = false
}

func (a *adapter) Open() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.open(ctx)
}

func (a *adapter) Close() {
	a.client.conn.Disconnect(context.Background())
}

func (a *adapter) SetModel(model core.Entitier, tableName string) {
	a.model = model
	a.tableName = tableName

	ctx := context.Background()
	names, err := a.client.database.ListCollectionNames(
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
		err := a.client.database.CreateCollection(ctx, a.tableName)
		if err != nil {
			panic(err)
		}
	}

	a.client.collection = a.client.database.Collection(a.tableName)
}

func (a *adapter) Find(ctx context.Context, id uuid.UUID, dest core.Entitier) (err error) {
	if !a.client.isOpened {
		a.Open()
	}

	err = a.client.collection.FindOne(ctx, bson.M{"entity._id": id}).Decode(dest)
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
	if !a.client.isOpened {
		a.Open()
	}

	count, err = a.client.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	if !a.client.isOpened {
		a.Open()
	}

	count, err = a.client.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a *adapter) List(ctx context.Context, dest interface{}) (err error) {
	if !a.client.isOpened {
		a.Open()
	}

	cursor, err := a.client.collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, dest)

	return err
}

func (a *adapter) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) (err error) {
	if !a.client.isOpened {
		a.Open()
	}

	cursor, err := a.client.collection.Find(ctx, query)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, dest)

	return err
}

func (a *adapter) Remove(ctx context.Context, id uuid.UUID) error {
	if !a.client.isOpened {
		a.Open()
	}

	_, err := a.client.collection.DeleteOne(ctx, bson.M{"entity._id": id})
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) RemoveRange(ctx context.Context, ids []uuid.UUID) error {
	if !a.client.isOpened {
		a.Open()
	}

	for _, id := range ids {
		_, err := a.client.collection.DeleteOne(ctx, bson.M{"entity._id": id})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *adapter) Add(ctx context.Context, entity core.Entitier) error {
	if !a.client.isOpened {
		a.Open()
	}

	_, err := a.client.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) AddRange(ctx context.Context, entities []core.Entitier) error {
	if !a.client.isOpened {
		a.Open()
	}

	vals := make([]interface{}, len(entities))
	for i, entity := range entities {
		vals[i] = entity
	}
	_, err := a.client.collection.InsertMany(ctx, vals)
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) Update(ctx context.Context, entity core.Entitier) error {
	if !a.client.isOpened {
		a.Open()
	}

	return a.client.collection.FindOneAndReplace(ctx, bson.M{"entity._id": entity.GetID()}, entity).Err()
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	if !a.client.isOpened {
		a.Open()
	}

	for _, entity := range entities {
		err := a.client.collection.FindOneAndReplace(ctx, bson.M{"entity._id": entity.GetID()}, entity).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
