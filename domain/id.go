package domain

type Id struct {
	Value string
}

type IId interface {
	Create() (Id, error)
}
