package internal

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type postgresql struct {
	client *sqlx.DB
}

func (db *postgresql) Create(data Entity) error {
	return db.Create(data)
}

func (db *postgresql) Get() ([]Entity, error) {
	rows, err := db.client.Query("SELECT * FROM entities")
	if err != nil {
		return nil, err
	}
	res := make([]Entity, 0)
	for rows.Next() {
		e := Entity{}
		if err := rows.Scan(&e); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func (db *postgresql) GetOne(id string) (*Entity, error) {
	e := Entity{}
	if err := db.client.QueryRow("SELECT * FROM entities WHERE id = ?", id).Scan(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (db *postgresql) Update(id string, data Entity) error {
	_, err := db.client.Exec(`
		UPDATE entities
		SET name = ?, description = ?, enabled = ?
		WHERE id = ?`, data.Name, data.Description, data.Enabled, id,
	)
	return err
}

func (db *postgresql) Delete(id string) error {
	_, err := db.client.Exec("DELETE FROM entities WHERE id = ?", id)
	return err
}

func NewPostgresql(uri string) (DBAdapter, error) {
	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return nil, err
	}
	return &postgresql{db}, nil
}
