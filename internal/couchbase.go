package internal

import (
	"fmt"
	"github.com/couchbase/go-couchbase"
)

type cbase struct {
	bucket *couchbase.Bucket
}

func (db *cbase) Create(data Entity) error {
	return db.bucket.PutDDoc(data.Id, data)
}

// TODO: fix
func (db *cbase) Get() ([]Entity, error) {
	res, err := db.bucket.GetDDocs()
	if err != nil {
		return nil, err
	}
	for _, d := range res.Rows {
		def := d.DDoc.JSON.Views
		fmt.Printf("TODO: check what returned %t\n", def)
	}
	return nil, nil
}

func (db *cbase) GetOne(id string) (*Entity, error) {
	res := Entity{}
	err := db.bucket.GetDDoc(id, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (db *cbase) Update(id string, data Entity) error {
	return db.bucket.PutDDoc(data.Id, data)
}

func (db *cbase) Delete(id string) error {
	return db.bucket.Delete(id)
}

func NewCouchbase(uri string) (DBAdapter, error) {
	c, err := couchbase.Connect(uri)
	if err != nil {
		return nil, err
	}
	pool, err := c.GetPool("default")
	if err != nil {
		return nil, err
	}
	bucket, err := pool.GetBucket("entities")
	if err != nil {
		return nil, err
	}
	return &cbase{bucket}, nil
}
