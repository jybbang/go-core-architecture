package core

type CommandRepositoryAdapter interface {
	SetModel(Entitier)
	Remove(Entitier) error
	RemoveRange([]Entitier) error
	Add(Entitier) error
	AddRange([]Entitier) error
	Update(Entitier) error
	UpdateRange([]Entitier) error
}
