package internal

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type cockroach struct {
	conf *pgx.ConnConfig
}

func (db *cockroach) Create(data Entity) error {
	conn, err := pgx.ConnectConfig(context.Background(), db.conf)
	if err != nil {
		return err
	}
	_, err = conn.Exec(
		context.Background(),
		`INSERT INTO entities (id, name, description, enabled) VALUES (?, ?, ?, ?)`,
		data.Id,
		data.Name,
		data.Description,
		data.Enabled,
	)
	if err != nil {
		return err
	}
	return conn.Close(context.Background())
}

func (db *cockroach) Get() ([]Entity, error) {
	conn, err := pgx.ConnectConfig(context.Background(), db.conf)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(
		context.Background(),
		`SELECT id, name, description, enabled FROM entities`,
	)
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
	return res, conn.Close(context.Background())
}

func (db *cockroach) GetOne(id string) (*Entity, error) {
	conn, err := pgx.ConnectConfig(context.Background(), db.conf)
	if err != nil {
		return nil, err
	}

	var (
		name, description string
		enabled           bool
	)
	if err := conn.QueryRow(
		context.Background(),
		`SELECT name, description, enabled
			FROM entities
			WHERE id = ?`,
		id,
	).Scan(name, description, enabled); err != nil {
		return nil, err
	}
	return &Entity{id, enabled, name, description}, nil
}

func (db *cockroach) Update(id string, data Entity) error {
	conn, err := pgx.ConnectConfig(context.Background(), db.conf)
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(),
		`UPDATE entities
			SET name = ?, description = ?, enabled = ?
			WHERE id = ?`,
		data.Name,
		data.Description,
		data.Enabled,
		id,
	)
	return err
}

func (db *cockroach) Delete(id string) error {
	conn, err := pgx.ConnectConfig(context.Background(), db.conf)
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(), "DELETE FROM entities WHERE id = ?", id)
	return err
}

func NewCockroach(uri string) (DBAdapter, error) {
	conf, err := pgx.ParseConfig(uri)
	if err != nil {
		return nil, err
	}
	return &cockroach{conf}, nil
}
