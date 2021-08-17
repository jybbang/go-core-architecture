package infrastructure

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_mongoQueryRepositoryService_ConnectionTimeout(t *testing.T) {
	timeout := time.Duration(1 * time.Second)
	ctx, c := context.WithTimeout(context.TODO(), timeout)
	defer c()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	time.Sleep(timeout)

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	result := r.Add(ctx, dto)

	if !errors.Is(result.E, context.DeadlineExceeded) {
		t.Errorf("Test_mongoQueryRepositoryService_ConnectionTimeout() err = %v, expect %v", result.E, context.DeadlineExceeded)
	}
}

func Test_mongoQueryRepositoryService_Find(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	dto2 := new(testModel)
	count := 1000
	for i := 0; i < count; i++ {
		go r.Find(ctx, dto.ID, dto2)
	}

	time.Sleep(1 * time.Second)

	result := r.Find(ctx, dto.ID, dto2)

	if !reflect.DeepEqual(dto2.ID, dto.ID) || !reflect.DeepEqual(dto2.Expect, dto.Expect) {
		t.Errorf("Test_mongoQueryRepositoryService_Find() result = %v, expect %v", dto2, dto)
	}

	if result.E != nil {
		t.Errorf("Test_mongoQueryRepositoryService_Find() err = %v", result.E)
	}
}

func Test_mongoQueryRepositoryService_FindnotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	dto2 := new(testModel)
	result := r.Find(ctx, dto.ID, dto2)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("Test_mongoQueryRepositoryService_FindnotFoundShouldBeError() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_mongoQueryRepositoryService_Any(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	count := 1000
	for i := 0; i < count; i++ {
		go r.Any(ctx)
	}

	time.Sleep(1 * time.Second)

	result := r.Any(ctx)

	if result.V != true {
		t.Errorf("Test_mongoQueryRepositoryService_Any() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_mongoQueryRepositoryService_Any() err = %v", result.E)
	}
}

func Test_mongoQueryRepositoryService_AnyWithFilter(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	result := r.AnyWithFilter(ctx, bson.M{"entity._id": dto.ID}, "")

	if result.V != true {
		t.Errorf("Test_mongoQueryRepositoryService_AnyWithFilter() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_mongoQueryRepositoryService_AnyWithFilter() err = %v", result.E)
	}
}

func Test_mongoQueryRepositoryService_Count(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 1000
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		r.Add(ctx, dto)
	}

	count := 1000
	for i := 0; i < count; i++ {
		go r.Count(ctx)
	}

	time.Sleep(1 * time.Second)

	result := r.Count(ctx)

	if result.V.(int64) < int64(expect) {
		t.Errorf("Test_mongoQueryRepositoryService_Count() result = %v, expect %v", result.V, expect)
	}

	if result.E != nil {
		t.Errorf("Test_mongoQueryRepositoryService_Count() err = %v", result.E)
	}
}

func Test_mongoQueryRepositoryService_CountWithFilter(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 1000
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = random

		r.Add(ctx, dto)
	}

	result := r.CountWithFilter(ctx, bson.M{"expect": random}, "")

	if result.V.(int64) != int64(expect) {
		t.Errorf("Test_mongoQueryRepositoryService_CountWithFilter() result = %v, expect %v", result.V, expect)
	}

	if result.E != nil {
		t.Errorf("Test_mongoQueryRepositoryService_CountWithFilter() err = %v", result.E)
	}
}

func Test_mongoQueryRepositoryService_List(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 1000
	cntExpect := 0
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = random

		r.Add(ctx, dto)
		cntExpect += 1
	}

	var dest = make([]*testModel, 0)
	result := r.List(ctx, &dest)

	if result.E != nil {
		t.Errorf("Test_mongoQueryRepositoryService_List() err = %v", result.E)
	}

	cnt := len(dest)

	if cnt < cntExpect {
		t.Errorf("Test_mongoQueryRepositoryService_List() cnt = %v, expect %v", cnt, cntExpect)
	}
}

func Test_mongoQueryRepositoryService_ListWithFilter(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 1000
	cntExpect := 0
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = random

		r.Add(ctx, dto)
		cntExpect += 1
	}

	var dest = make([]*testModel, 0)
	result := r.ListWithFilter(ctx, bson.M{"expect": random}, "", &dest)

	cnt := len(dest)

	if cnt != cntExpect {
		t.Errorf("Test_mongoQueryRepositoryService_ListWithFilter() cnt = %v, expect %v", cnt, cntExpect)
	}

	if result.E != nil {
		t.Errorf("Test_mongoQueryRepositoryService_ListWithFilter() err = %v", result.E)
	}
}

func Test_mongoCommandRepositoryService_Remove(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	count := 1000
	for i := 0; i < count; i++ {
		go r.Remove(ctx, dto.ID)
	}

	time.Sleep(1 * time.Second)

	result := r.Remove(ctx, dto.ID)

	if result.E != nil {
		t.Errorf("Test_mongoCommandRepositoryService_Remove() err = %v", result.E)
	}

	dto2 := new(testModel)
	result = r.Find(ctx, dto.ID, dto2)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("Test_mongoCommandRepositoryService_Remove() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_mongoCommandRepositoryService_RemoveRange(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 1000
	var ids = make([]uuid.UUID, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i
		ids = append(ids, dto.ID)

		r.Add(ctx, dto)
	}

	result := r.RemoveRange(ctx, ids)

	if result.E != nil {
		t.Errorf("Test_mongoCommandRepositoryService_RemoveRange() err = %v", result.E)
	}

	dto2 := new(testModel)
	result = r.Find(ctx, ids[0], dto2)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("Test_mongoCommandRepositoryService_RemoveRange() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_mongoCommandRepositoryService_AddRange(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 1000
	var cntExpect int64 = 0
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = random

		dtos = append(dtos, dto)
		cntExpect += 1
	}

	result := r.AddRange(ctx, dtos)

	if result.E != nil {
		t.Errorf("Test_mongoCommandRepositoryService_AddRange() err = %v", result.E)
	}

	result = r.CountWithFilter(ctx, bson.M{"expect": random}, "")

	if result.V.(int64) != cntExpect {
		t.Errorf("Test_mongoCommandRepositoryService_AddRange() cnt = %v, expect %v", result.V, cntExpect)
	}
}

func Test_mongoCommandRepositoryService_Update(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 1000

	r.Add(ctx, dto)

	dto.Expect = 1
	count := 1000
	for i := 0; i < count; i++ {
		go r.Update(ctx, dto)
	}

	time.Sleep(1 * time.Second)

	result := r.Update(ctx, dto)

	if result.E != nil {
		t.Errorf("Test_mongoCommandRepositoryService_Update() err = %v", result.E)
	}

	dto2 := new(testModel)
	result = r.Find(ctx, dto.GetID(), dto2)

	if dto2.Expect != 1 {
		t.Errorf("Test_mongoCommandRepositoryService_Update() result = %v, expect %v", dto2.Expect, 1)
	}
}

func Test_mongoCommandRepositoryService_UpdateRange(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "test_model").
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 1000
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		dtos = append(dtos, dto)
		r.Add(ctx, dto)
	}

	var cntExpect int64 = 0
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for _, dto := range dtos {
		dto.(*testModel).Expect = random
		cntExpect += 1
	}

	result := r.UpdateRange(ctx, dtos)

	if result.E != nil {
		t.Errorf("Test_mongoCommandRepositoryService_UpdateRange() err = %v", result.E)
	}

	result = r.CountWithFilter(ctx, bson.M{"expect": random}, "")

	if result.V.(int64) != cntExpect {
		t.Errorf("Test_mongoCommandRepositoryService_UpdateRange() cnt = %v, expect %v", result.V, cntExpect)
	}
}
