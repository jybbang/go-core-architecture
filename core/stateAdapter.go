package core

type StateAdapter interface {
	Has(string) (bool, error)
	Get(string, Entitier) (bool, error)
	Set(string, interface{}) error
}
