package domain

type Copyable interface {
	CopyWith(interface{}) bool
}
