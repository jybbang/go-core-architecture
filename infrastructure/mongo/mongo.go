package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
)

type adapter struct {
	tableName string
	model     core.Entitier
	client    *clientProxy
	settings  MongoSettings
}

type clientProxy struct {
	conn        *mongo.Client
	database    *mongo.Database
	collection  *mongo.Collection
	isConnected bool
}

type clients struct {
	clients map[string]*clientProxy
	sync.Mutex
}

type MongoSettings struct {
	ConnectionUri       string
	DatabaseName        string
	CanCreateCollection bool
}

var clientsInstance *clients

func init() {
	clientsInstance = &clients{
		clients: make(map[string]*clientProxy),
	}
}

func (a *adapter) migration() {
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

func NewMongoAdapter(settings MongoSettings) *adapter {
	return &adapter{
		settings: settings,
	}
}

func (a *adapter) IsConnected() bool {
	return a.client.isConnected
}

func (a *adapter) Connect(ctx context.Context) error {
	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	uri := a.settings.ConnectionUri

	if strings.TrimSpace(uri) == "" {
		return fmt.Errorf("uri is required")
	}

	cli, ok := clientsInstance.clients[uri]

	if !ok || !cli.isConnected {
		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

		if err != nil {
			return err
		}

		// Check context cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		clientsInstance.clients[uri] = &clientProxy{
			conn:        mongoClient,
			database:    mongoClient.Database(a.settings.DatabaseName),
			isConnected: true,
		}
	}

	a.client = clientsInstance.clients[uri]

	if a.tableName != "" {
		a.migration()
	}

	return nil
}

func (a *adapter) Disconnect() {
	clientsInstance.Lock()
	defer clientsInstance.Unlock()

	a.client.conn.Disconnect(context.Background())

	a.client.isConnected = false
}

func (a *adapter) SetModel(model core.Entitier, tableName string) {
	a.model = model
	a.tableName = tableName
}

func (a *adapter) Find(ctx context.Context, id uuid.UUID, dest core.Entitier) error {
	err := a.client.collection.FindOne(ctx, bson.M{"entity._id": id}).Decode(dest)

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
	count, err = a.client.collection.CountDocuments(ctx, bson.M{})

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a *adapter) CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error) {
	count, err = a.client.collection.CountDocuments(ctx, query)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a *adapter) List(ctx context.Context, dest interface{}) error {
	cursor, err := a.client.collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	if err != nil {
		return err
	}

	err = cursor.All(ctx, dest)

	return err
}

func (a *adapter) ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) error {
	cursor, err := a.client.collection.Find(ctx, query)
	defer cursor.Close(ctx)

	if err != nil {
		return err
	}

	err = cursor.All(ctx, dest)

	return err
}

func (a *adapter) Remove(ctx context.Context, id uuid.UUID) error {
	_, err := a.client.collection.DeleteOne(ctx, bson.M{"entity._id": id})

	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) RemoveRange(ctx context.Context, ids []uuid.UUID) error {
	for _, id := range ids {
		_, err := a.client.collection.DeleteOne(ctx, bson.M{"entity._id": id})

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *adapter) Add(ctx context.Context, entity core.Entitier) error {
	_, err := a.client.collection.InsertOne(ctx, entity)

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

	_, err := a.client.collection.InsertMany(ctx, vals)

	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) Update(ctx context.Context, entity core.Entitier) error {
	return a.client.collection.FindOneAndReplace(ctx, bson.M{"entity._id": entity.GetID()}, entity).Err()
}

func (a *adapter) UpdateRange(ctx context.Context, entities []core.Entitier) error {
	for _, entity := range entities {
		err := a.client.collection.FindOneAndReplace(ctx, bson.M{"entity._id": entity.GetID()}, entity).Err()

		if err != nil {
			return err
		}
	}

	return nil
}
