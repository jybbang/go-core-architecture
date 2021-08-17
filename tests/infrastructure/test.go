package infrastructure

import (
	"github.com/jybbang/go-core-architecture/core"
)

type testModel struct {
	core.Entity
	Expect int `bson:"expect,omitempty"`
}
