package internal

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type ldb struct {
	client *leveldb.DB
}

func (l *ldb) Create(data Entity) error {
	b, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return l.client.Put([]byte(data.Id), b, &opt.WriteOptions{})
}

func (l *ldb) Get() ([]Entity, error) {
	itr := l.client.NewIterator(
		&util.Range{Start: []byte("0"), Limit: []byte("255")},
		&opt.ReadOptions{},
	)
	res := make([]Entity, 0)
	for itr.Next() {
		e := Entity{}
		b := itr.Value()
		if err := json.Unmarshal(b, &e); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func (l *ldb) GetOne(id string) (*Entity, error) {
	b, err := l.client.Get([]byte(id), &opt.ReadOptions{})
	if err != nil {
		return nil, err
	}
	e := Entity{}
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (l *ldb) Update(id string, data Entity) error {
	b, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return l.client.Put([]byte(id), b, &opt.WriteOptions{})
}

func (l *ldb) Delete(id string) error {
	return l.client.Delete([]byte(id), &opt.WriteOptions{})
}

func NewLeveldb(uri string) (DBAdapter, error) {
	db, err := leveldb.OpenFile(uri, &opt.Options{})
	if err != nil {
		return nil, err
	}
	return &ldb{db}, nil
}
