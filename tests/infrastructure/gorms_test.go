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
	"github.com/jybbang/go-core-architecture/infrastructure/gorms"
)

func Test_gormsQueryRepositoryService_ConnectionTimeout(t *testing.T) {
	timeout := time.Duration(1 * time.Second)
	ctx, c := context.WithTimeout(context.TODO(), timeout)
	defer c()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
		Create()

	time.Sleep(timeout)

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	err := r.Add(ctx, dto)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Test_gormsQueryRepositoryService_ConnectionTimeout() err = %v, expect %v", err, context.DeadlineExceeded)
	}
}

func Test_gormsQueryRepositoryService_Find(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
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

	time.Sleep(1 * time.Second)

	err := r.Find(ctx, dto.ID, dto2)

	if !reflect.DeepEqual(dto2.Expect, dto.Expect) {
		t.Errorf("Test_gormsQueryRepositoryService_Find() result = %v, expect %v", dto2, dto)
	}

	if err != nil {
		t.Errorf("Test_gormsQueryRepositoryService_Find() err = %v", err)
	}
}

func Test_gormsQueryRepositoryService_FindnotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
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

func Test_gormsQueryRepositoryService_Any(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	count := 100000
	for i := 0; i < count; i++ {
		go r.Any(ctx)
	}

	ok, err := r.Any(ctx)

	if ok != true {
		t.Errorf("Test_gormsQueryRepositoryService_Any() ok = %v, expect %v", ok, true)
	}

	if err != nil {
		t.Errorf("Test_gormsQueryRepositoryService_Any() err = %v", err)
	}
}

func Test_gormsQueryRepositoryService_AnyWithFilter(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	ok, err := r.AnyWithFilter(ctx, "id = ?", dto.ID)

	if ok != true {
		t.Errorf("Test_gormsQueryRepositoryService_AnyWithFilter() ok = %v, expect %v", ok, true)
	}

	if err != nil {
		t.Errorf("Test_gormsQueryRepositoryService_AnyWithFilter() err = %v", err)
	}
}

func Test_gormsQueryRepositoryService_Count(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
		Create()

	expect := 100
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		r.Add(ctx, dto)
	}

	count := 100000
	for i := 0; i < count; i++ {
		go r.Count(ctx)
	}

	result, err := r.Count(ctx)

	if result < int64(expect) {
		t.Errorf("Test_gormsQueryRepositoryService_Count() result = %v, expect %v", result, expect)
	}

	if err != nil {
		t.Errorf("Test_gormsQueryRepositoryService_Count() err = %v", err)
	}
}

func Test_gormsQueryRepositoryService_CountWithFilter(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
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

	result, err := r.CountWithFilter(ctx, "expect = ?", random)

	if result != int64(expect) {
		t.Errorf("Test_gormsQueryRepositoryService_CountWithFilter() result = %v, expect %v", result, expect)
	}

	if err != nil {
		t.Errorf("Test_gormsQueryRepositoryService_CountWithFilter() err = %v", err)
	}
}

func Test_gormsQueryRepositoryService_List(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
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
		t.Errorf("Test_gormsQueryRepositoryService_List() err = %v", err)
	}

	cnt := len(dest)

	if cnt < cntExpect {
		t.Errorf("Test_gormsQueryRepositoryService_List() cnt = %v, expect %v", cnt, cntExpect)
	}
}

func Test_gormsQueryRepositoryService_ListWithFilter(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
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
	err := r.ListWithFilter(ctx, "expect = ?", random, &dest)

	cnt := len(dest)

	if cnt != cntExpect {
		t.Errorf("Test_gormsQueryRepositoryService_ListWithFilter() cnt = %v, expect %v", cnt, cntExpect)
	}

	if err != nil {
		t.Errorf("Test_gormsQueryRepositoryService_ListWithFilter() err = %v", err)
	}
}

func Test_gormsCommandRepositoryService_Remove(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 123

	r.Add(ctx, dto)

	count := 100000
	for i := 0; i < count; i++ {
		go r.Remove(ctx, dto)
	}

	err := r.Remove(ctx, dto)

	if err != nil {
		t.Errorf("Test_gormsCommandRepositoryService_Remove() err = %v", err)
	}

	dto2 := new(testModel)
	err = r.Find(ctx, dto.ID, dto2)

	if !errors.Is(err, core.ErrNotFound) {
		t.Errorf("Test_gormsCommandRepositoryService_Remove() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func Test_gormsCommandRepositoryService_RemoveRange(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
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
		t.Errorf("Test_gormsCommandRepositoryService_RemoveRange() err = %v", err)
	}

	dto2 := new(testModel)
	err = r.Find(ctx, dtos[0].GetID(), dto2)

	if !errors.Is(err, core.ErrNotFound) {
		t.Errorf("Test_gormsCommandRepositoryService_RemoveRange() err = %v, expect %v", err, core.ErrNotFound)
	}
}

func Test_gormsCommandRepositoryService_AddRange(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
		Create()

	expect := 100
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
		t.Errorf("Test_gormsCommandRepositoryService_AddRange() err = %v", err)
	}

	cnt, _ := r.CountWithFilter(ctx, "expect = ?", random)

	if cnt != int64(cntExpect) {
		t.Errorf("Test_gormsCommandRepositoryService_AddRange() cnt = %v, expect %v", cnt, cntExpect)
	}
}

func Test_gormsCommandRepositoryService_Update(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
		Create()

	dto := new(testModel)
	dto.ID = uuid.New()
	dto.Expect = 100

	r.Add(ctx, dto)

	dto.Expect = 1
	count := 100000
	for i := 0; i < count; i++ {
		go r.Update(ctx, dto)
	}

	time.Sleep(1 * time.Second)

	err := r.Update(ctx, dto)

	if err != nil {
		t.Errorf("Test_gormsCommandRepositoryService_Update() err = %v", err)
	}

	dto2 := new(testModel)
	r.Find(ctx, dto.GetID(), dto2)

	result := dto2.Expect
	if result != 1 {
		t.Errorf("Test_gormsCommandRepositoryService_Update() result = %v, expect %v", result, 1)
	}
}

func Test_gormsCommandRepositoryService_UpdateRange(t *testing.T) {
	ctx := context.Background()

	gorms := gorms.NewPostresAdapter(ctx, gorms.GormSettings{
		ConnectionString: "postgres://postgres:admin@localhost:5432/test",
		CanCreateTable:   true,
	})
	r := core.NewRepositoryServiceBuilder(new(testModel), "T_TESTMODEL").
		CommandRepositoryAdapter(gorms).
		QueryRepositoryAdapter(gorms).
		Create()

	expect := 100
	var dtos = make([]core.Entitier, 0)
	for i := 0; i < expect; i++ {
		dto := new(testModel)
		dto.ID = uuid.New()
		dto.Expect = i

		r.Add(ctx, dto)
		dtos = append(dtos, dto)
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
		t.Errorf("Test_gormsCommandRepositoryService_UpdateRange() err = %v", err)
	}

	cnt, _ := r.CountWithFilter(ctx, "expect = ?", random)

	if cnt != int64(cntExpect) {
		t.Errorf("Test_gormsCommandRepositoryService_UpdateRange() cnt = %v, expect %v", cnt, cntExpect)
	}
}
