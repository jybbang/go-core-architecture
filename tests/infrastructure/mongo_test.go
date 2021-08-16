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

func Test_queryRepositoryService_ConnectionTimeout(t *testing.T) {
	timeout := time.Duration(1 * time.Second)
	ctx, c := context.WithTimeout(context.TODO(), timeout)
	defer c()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	time.Sleep(timeout)

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	err := r.Add(ctx, dto)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Test_queryRepositoryService_ConnectionTimeout() err = %v, expect %v", err, context.DeadlineExceeded)
	}
}

func Test_queryRepositoryService_Find(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	dto2 := new(testModel)
	count := 100
	for i := 0; i < count; i++ {
		go r.Find(ctx, dto.ID, dto2)
	}

	err := r.Find(ctx, dto.ID, dto2)

	if !reflect.DeepEqual(dto2.Expect, dto.Expect) {
		t.Errorf("Test_queryRepositoryService_Find() result = %v, expect %v", dto2, dto)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_Find() err = %v", err)
	}
}

func Test_queryRepositoryService_FindnotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	dto2 := new(testModel)
	err := r.Find(ctx, dto.ID, dto2)

	if !errors.Is(err, core.ErrNotFound) {
		t.Errorf("TestStateService_GetNotFoundShouldBeError() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func Test_queryRepositoryService_Any(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	count := 100
	for i := 0; i < count; i++ {
		go r.Any(ctx)
	}

	ok, err := r.Any(ctx)

	if ok != true {
		t.Errorf("Test_queryRepositoryService_Any() ok = %v, expect %v", ok, true)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_Any() err = %v", err)
	}
}

func Test_queryRepositoryService_AnyWithFilter(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	ok, err := r.AnyWithFilter(ctx, bson.M{"entity._id": dto.ID}, "")

	if ok != true {
		t.Errorf("Test_queryRepositoryService_AnyWithFilter() ok = %v, expect %v", ok, true)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_AnyWithFilter() err = %v", err)
	}
}

func Test_queryRepositoryService_Count(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 100
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		r.Add(ctx, dto)
	}

	count := 100
	for i := 0; i < count; i++ {
		go r.Count(ctx)
	}

	result, err := r.Count(ctx)

	if result < int64(expect) {
		t.Errorf("Test_queryRepositoryService_Count() result = %v, expect %v", result, expect)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_Count() err = %v", err)
	}
}

func Test_queryRepositoryService_CountWithFilter(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 100
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = random

		r.Add(ctx, dto)
	}

	result, err := r.CountWithFilter(ctx, bson.M{"expect": random}, "")

	if result != int64(expect) {
		t.Errorf("Test_queryRepositoryService_CountWithFilter() result = %v, expect %v", result, expect)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_CountWithFilter() err = %v", err)
	}
}

func Test_queryRepositoryService_List(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 100
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
	err := r.List(ctx, &dest)

	if err != nil {
		t.Errorf("Test_queryRepositoryService_List() err = %v", err)
	}

	cnt := len(dest)

	if cnt < cntExpect {
		t.Errorf("Test_queryRepositoryService_List() cnt = %v, expect %v", cnt, cntExpect)
	}
}

func Test_queryRepositoryService_ListWithFilter(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 100
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
	err := r.ListWithFilter(ctx, bson.M{"expect": random}, "", &dest)

	cnt := len(dest)

	if cnt != cntExpect {
		t.Errorf("Test_queryRepositoryService_ListWithFilter() cnt = %v, expect %v", cnt, cntExpect)
	}

	if err != nil {
		t.Errorf("Test_queryRepositoryService_ListWithFilter() err = %v", err)
	}
}

func Test_commandRepositoryService_Remove(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	count := 100
	for i := 0; i < count; i++ {
		go r.Remove(ctx, dto)
	}

	err := r.Remove(ctx, dto)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_Remove() err = %v", err)
	}

	dto2 := new(testModel)
	err = r.Find(ctx, dto.ID, dto2)

	if !errors.Is(err, core.ErrNotFound) {
		t.Errorf("Test_commandRepositoryService_Remove() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func Test_commandRepositoryService_RemoveRange(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 100
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i
		dtos = append(dtos, dto)

		r.Add(ctx, dto)
	}

	err := r.RemoveRange(ctx, dtos)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_RemoveRange() err = %v", err)
	}

	dto2 := new(testModel)
	err = r.Find(ctx, dtos[0].GetID(), dto2)

	if !errors.Is(err, core.ErrNotFound) {
		t.Errorf("Test_commandRepositoryService_RemoveRange() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func Test_commandRepositoryService_AddRange(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 10000
	cntExpect := 0
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

	err := r.AddRange(ctx, dtos)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_AddRange() err = %v", err)
	}

	var dest = make([]*testModel, 0)
	r.ListWithFilter(ctx, bson.M{"expect": random}, "", &dest)

	cnt := len(dest)

	if cnt != cntExpect {
		t.Errorf("Test_commandRepositoryService_AddRange() cnt = %v, expect %v", cnt, cntExpect)
	}
}

func Test_commandRepositoryService_Update(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 100

	r.Add(ctx, dto)

	dto.Expect = 1
	count := 100
	for i := 0; i < count; i++ {
		go r.Update(ctx, dto)
	}

	time.Sleep(1 * time.Second)

	err := r.Update(ctx, dto)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_Update() err = %v", err)
	}

	dto2 := new(testModel)
	r.Find(ctx, dto.GetID(), dto2)

	result := dto2.Expect
	if result != 1 {
		t.Errorf("Test_commandRepositoryService_Update() result = %v, expect %v", result, 1)
	}
}

func Test_commandRepositoryService_UpdateRange(t *testing.T) {
	ctx := context.Background()

	mongo := mongo.NewMongoAdapter(ctx, mongo.MongoSettings{
		ConnectionUri:       "mongodb://admin:admin@localhost:27017",
		DatabaseName:        "testdb",
		CanCreateCollection: true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel)).
		CommandRepositoryAdapter(mongo).
		QueryRepositoryAdapter(mongo).
		Create()

	expect := 100
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		dtos = append(dtos, dto)
		r.Add(ctx, dto)
	}

	cntExpect := 0
	rand.Seed(time.Now().UnixNano())
	random := rand.Int()
	for _, dto := range dtos {
		dto.(*testModel).Expect = random
		cntExpect += 1
	}

	err := r.UpdateRange(ctx, dtos)

	if err != nil {
		t.Errorf("Test_commandRepositoryService_UpdateRange() err = %v", err)
	}

	var dest = make([]*testModel, 0)
	r.ListWithFilter(ctx, bson.M{"expect": random}, "", &dest)

	cnt := len(dest)

	if cnt != cntExpect {
		t.Errorf("Test_commandRepositoryService_UpdateRange() cnt = %v, expect %v", cnt, cntExpect)
	}
}
