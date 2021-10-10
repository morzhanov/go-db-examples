package internal

import (
	"github.com/zemirco/couchdb"
	"net/url"
)

type cdbEntity struct {
	couchdb.Document
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type cdb struct {
	client couchdb.DatabaseService
}

func (db *cdb) Create(data Entity) error {
	ent := cdbEntity{Id: data.Id, Name: data.Name, Description: data.Description, Enabled: data.Enabled}
	_, err := db.client.Post(&ent)
	return err
}

func (db *cdb) Get() ([]Entity, error) {
	v, err := db.client.AllDocs(&couchdb.QueryParameters{})
	if err != nil {
		return nil, err
	}
	res := make([]Entity, 0)
	for _, r := range v.Rows {
		ent := r.Value.(cdbEntity)
		el := Entity{Id: ent.Id, Name: ent.Name, Description: ent.Description, Enabled: ent.Enabled}
		res = append(res, el)
	}
	return res, nil
}

func (db *cdb) getEntity(id string) (*cdbEntity, error) {
	ent := cdbEntity{}
	if err := db.client.Get(&ent, id); err != nil {
		return nil, err
	}
	return &ent, nil
}

func (db *cdb) GetOne(id string) (*Entity, error) {
	ent, err := db.getEntity(id)
	if err != nil {
		return nil, err
	}
	return &Entity{Id: ent.Id, Name: ent.Name, Description: ent.Description, Enabled: ent.Enabled}, nil
}

func (db *cdb) Update(_ string, data Entity) error {
	ent := cdbEntity{Id: data.Id, Name: data.Name, Description: data.Description, Enabled: data.Enabled}
	_, err := db.client.Put(&ent)
	return err
}

func (db *cdb) Delete(id string) error {
	ent, err := db.getEntity(id)
	if err != nil {
		return err
	}
	_, err = db.client.Delete(ent)
	return err
}

func NewCouchdb(uri string) (DBAdapter, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	client, err := couchdb.NewClient(u)
	if err != nil {
		return nil, err
	}
	_, err = client.Create("entities")
	if err != nil {
		return nil, err
	}
	db := client.Use("entities")
	return &cdb{db}, nil
}
