package infrastructure

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/leveldb"
)

func Test_leveldbStateService_Has(t *testing.T) {
	ctx := context.Background()

	leveldb := leveldb.NewLevelDbAdapter(ctx, leveldb.LevelDbSettings{
		Path: "_test.db",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(leveldb).
		Create()

	count := 10000
	key := "qwe"
	expect := testModel{
		Expect: 123,
	}

	s.Set(ctx, key, &expect)

	for i := 0; i < count; i++ {
		go s.Has(ctx, key)
	}

	time.Sleep(1 * time.Second)

	result := s.Has(ctx, key)

	if result.V != true {
		t.Errorf("Test_leveldbStateService_Has() ok = %v, expect %v", result.V, true)
	}

	if result.E != nil {
		t.Errorf("Test_leveldbStateService_Has() err = %v", result.E)
	}
}

func Test_leveldbStateService_Get(t *testing.T) {
	ctx := context.Background()

	leveldb := leveldb.NewLevelDbAdapter(ctx, leveldb.LevelDbSettings{
		Path: "_test.db",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(leveldb).
		Create()

	count := 10000
	key := "qwe"
	expect := &testModel{
		Expect: 123,
	}

	s.Set(ctx, key, expect)

	dest := &testModel{}
	for i := 0; i < count; i++ {
		go s.Get(ctx, key, dest)
	}

	time.Sleep(1 * time.Second)

	result := s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("Test_leveldbStateService_Get() dest = %v, expect %v", dest, expect)
	}

	if result.E != nil {
		t.Errorf("Test_leveldbStateService_Get() err = %v", result.E)
	}
}

func Test_leveldbStateService_GetNotFoundShouldBeError(t *testing.T) {
	ctx := context.Background()

	leveldb := leveldb.NewLevelDbAdapter(ctx, leveldb.LevelDbSettings{
		Path: "_test.db",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(leveldb).
		Create()

	dest := &testModel{}
	result := s.Get(ctx, "zxc", dest)

	if !errors.Is(result.E, core.ErrNotFound) {
		t.Errorf("TestStateService_GetNotFoundShouldBeError() err = %v, expect %v", result.E, core.ErrNotFound)
	}
}

func Test_leveldbStateService_Set(t *testing.T) {
	ctx := context.Background()

	leveldb := leveldb.NewLevelDbAdapter(ctx, leveldb.LevelDbSettings{
		Path: "_test.db",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(leveldb).
		Create()

	count := 10000
	key := "qwe"
	expect := &testModel{
		Expect: 123,
	}
	for i := 0; i < count; i++ {
		go s.Set(ctx, key, expect)
	}

	time.Sleep(1 * time.Second)

	result := s.Set(ctx, key, expect)

	if result.E != nil {
		t.Errorf("Test_leveldbStateService_Set() err = %v", result.E)
	}

	dest := &testModel{}
	s.Get(ctx, key, dest)

	if !reflect.DeepEqual(dest, expect) {
		t.Errorf("Test_leveldbStateService_Set() dest = %v, expect %v", dest, expect)
	}
}

func Test_leveldbStateService_BatchSet(t *testing.T) {
	ctx := context.Background()

	leveldb := leveldb.NewLevelDbAdapter(ctx, leveldb.LevelDbSettings{
		Path: "_test.db",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(leveldb).
		Create()

	expect := 10000

	kvs := make([]core.Kvs, 0)
	for i := 0; i < expect; i++ {
		kvs = append(
			kvs,
			core.Kvs{
				K: strconv.Itoa(i),
				V: &testModel{
					Expect: i,
				}})
	}

	result := s.BatchSet(ctx, kvs)

	if result.E != nil {
		t.Errorf("Test_leveldbStateService_BatchSet() err = %v", result.E)
	}

	dest := &testModel{}
	s.Get(ctx, strconv.Itoa(expect), dest)

	if dest.Expect == expect {
		t.Errorf("Test_leveldbStateService_BatchSet() dest.Expect = %d, expect %d", dest.Expect, expect)
	}
}

func Test_leveldbStateService_Delete(t *testing.T) {
	ctx := context.Background()

	leveldb := leveldb.NewLevelDbAdapter(ctx, leveldb.LevelDbSettings{
		Path: "_test.db",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(leveldb).
		Create()

	count := 10000
	key := "qwe"
	expect := testModel{
		Expect: 123,
	}

	s.Set(ctx, key, &expect)
	for i := 0; i < count; i++ {
		go s.Delete(ctx, key)
	}

	time.Sleep(1 * time.Second)

	s.Delete(ctx, key)

	result := s.Has(ctx, key)

	if result.E != nil {
		t.Errorf("Test_leveldbStateService_Delete() err = %v", result.E)
	}

	if result.V != false {
		t.Errorf("Test_leveldbStateService_Delete() ok = %v, expect %v", result.V, false)
	}
}

func Test_leveldbStateService_DeleteNotFoundShouldBeNoError(t *testing.T) {
	ctx := context.Background()

	leveldb := leveldb.NewLevelDbAdapter(ctx, leveldb.LevelDbSettings{
		Path: "_test.db",
	})
	s := core.NewStateServiceBuilder().
		StateAdapter(leveldb).
		Create()

	key := "qwe"

	result := s.Delete(ctx, key)

	if result.E != nil {
		t.Errorf("Test_leveldbStateService_DeleteNotFoundShouldBeNoError() err = %v", result.E)
	}
}
