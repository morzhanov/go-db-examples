package internal

import (
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"log"
)

type chouse struct {
	client *sqlx.DB
}

func (c *chouse) Create(data Entity) error {
	_, err := c.client.MustExec(
		"INSERT INTO entities (id, name, description, enabled) VALUES (?, ?, ?, ?)",
		data.Id,
		data.Name,
		data.Description,
		data.Enabled,
	).RowsAffected()
	return err
}

func (c *chouse) Get() ([]Entity, error) {
	rows, err := c.client.Query("SELECT id, name, description, enabled FROM entities")
	if err != nil {
		return nil, err
	}

	res := make([]Entity, 0)
	for rows.Next() {
		var (
			id, name, description string
			enabled               bool
		)
		if err := rows.Scan(&id, &name, &description, &enabled); err != nil {
			return nil, err
		}
		res = append(res, Entity{id, enabled, name, description})
	}
	return res, rows.Err()
}

func (c *chouse) GetOne(id string) (*Entity, error) {
	var (
		name, description string
		enabled           bool
	)
	if err := c.client.QueryRow(`
		SELECT name, description, enabled
		FROM entities
		WHERE id = ?`, id).Scan(name, description, enabled); err != nil {
		return nil, err
	}
	return &Entity{id, enabled, name, description}, nil
}

func (c *chouse) Update(id string, data Entity) error {
	_, err := c.client.Exec(`
		UPDATE entities
		SET name = ?, description = ?, enabled = ?
		WHERE id = ?`, data.Name, data.Description, data.Enabled, id,
	)
	return err
}

func (c *chouse) Delete(id string) error {
	_, err := c.client.Exec("DELETE FROM entities WHERE id = ?", id)
	return err
}

func NewClickhouse(uri string) (DBAdapter, error) {
	connect, err := sqlx.Open("clickhouse", uri)
	if err != nil {
		log.Fatal(err)
	}
	return &chouse{connect}, nil
}
