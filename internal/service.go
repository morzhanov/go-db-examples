package internal

import (
	"fmt"
	"github.com/google/uuid"
)

type Service interface {
	Test(db DBAdapter, enabled bool, name string, description string) error
}

type service struct {
}

type Entity struct {
	Id          string `client:"id"`
	Enabled     bool   `client:"enabled"`
	Name        string `client:"name"`
	Description string `client:"description"`
}

func (s *service) Test(db DBAdapter, enabled bool, name string, description string) error {
	// CREATE
	uid, err := uuid.NewUUID()
	id := uid.String()
	if err != nil {
		return err
	}
	e := Entity{
		Id:          id,
		Enabled:     enabled,
		Name:        name,
		Description: description,
	}

	if err := db.Create(e); err != nil {
		return err
	}
	fmt.Printf("Created entity with id %s and name %s\n", id, name)

	// GET
	res, err := db.Get()
	if err != nil {
		return err
	}
	fmt.Printf("Fetched all entities...\n")
	for i, e := range res {
		fmt.Printf("[%d]: %s - %s\n", i, e.Id, e.Name)
	}

	// UPDATE
	e.Name = fmt.Sprintf("%s-updated", e.Name)
	e.Description = fmt.Sprintf("%s-updated", e.Description)
	e.Enabled = false
	if err := db.Update(id, e); err != nil {
		return err
	}
	fmt.Printf("Updated entity with id %s name = %s, description = %s\n", id, e.Name, e.Description)

	// GET ONE
	res2, err := db.GetOne(id)
	fmt.Printf(
		"Fetched entity with id %s: name: %s, description: %s, enabled: %t\n",
		id,
		res2.Name,
		res2.Description,
		res2.Enabled,
	)

	// DELETE
	if err := db.Delete(id); err != nil {
		return err
	}
	fmt.Printf("Deleted entity with id %s\n", id)
	return nil
}

func NewService() Service {
	return &service{}
}
