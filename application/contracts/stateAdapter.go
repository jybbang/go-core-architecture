package contracts

import "github.com/jybbang/go-core-architecture/domain"

type StateAdapter interface {
	Has(string) (bool, error)
	Get(string, domain.Copyable) (bool, error)
	Set(string, interface{}) error
}
