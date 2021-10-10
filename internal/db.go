package internal

type DBAdapter interface {
	Create(data Entity) error
	Get() ([]Entity, error)
	GetOne(id string) (*Entity, error)
	Delete(id string) error
	Update(id string, data Entity) error
}
