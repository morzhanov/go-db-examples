package internal

import (
	"context"
	"github.com/gocql/gocql"
)

type cassandra struct {
	cluster *gocql.ClusterConfig
}

func (c *cassandra) Create(data Entity) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()
	ctx := context.Background()
	return session.Query(
		`INSERT INTO entities (id, name, description, enabled) VALUES (?, ?, ?, ?)`,
		data.Id,
		data.Name,
		data.Description,
		data.Enabled).WithContext(ctx).Exec()
}

func (c *cassandra) Get() ([]Entity, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	ctx := context.Background()
	scanner := session.Query("SELECT id, name, description, enabled FROM entities").WithContext(ctx).Iter().Scanner()

	res := make([]Entity, 0)
	for scanner.Next() {
		var (
			id, name, description string
			enabled               bool
		)
		if err := scanner.Scan(&id, &name, description, enabled); err != nil {
			return nil, err
		}
		res = append(res, Entity{id, enabled, name, description})
	}
	return res, scanner.Err()
}

func (c *cassandra) GetOne(id string) (*Entity, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	ctx := context.Background()
	var (
		name, description string
		enabled           bool
	)
	if err := session.Query(
		`SELECT name, description, enabled
				FROM entities
				WHERE id = ?`, id,
	).WithContext(ctx).Consistency(gocql.One).Scan(&name, &description, &enabled); err != nil {
		return nil, err
	}
	return &Entity{id, enabled, name, description}, nil
}

func (c *cassandra) Update(id string, data Entity) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	ctx := context.Background()
	return session.Query(
		`UPDATE entities
				SET name = ?, description = ?, enabled = ?
				WHERE id = ?`, data.Name, data.Description, data.Enabled, id,
	).WithContext(ctx).Exec()
}

func (c *cassandra) Delete(id string) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()
	ctx := context.Background()
	return session.Query(`DELETE FROM entitiesWHERE id = ?`, id).WithContext(ctx).Exec()
}

func NewCassandra(uri string) DBAdapter {
	cluster := gocql.NewCluster(uri)
	cluster.Keyspace = "entities"
	return &cassandra{cluster}
}
