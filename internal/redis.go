package internal

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type rdb struct {
	client *redis.Client
}

func (db *rdb) Create(data Entity) error {
	conn := db.client.Conn(context.Background())
	conn.Set(context.Background(), data.Id, data, time.Hour*1000)
	return conn.Close()
}

func (db *rdb) Get() ([]Entity, error) {
	conn := db.client.Conn(context.Background())
	var cursor uint64
	keys, cursor, err := conn.Scan(context.Background(), cursor, "", 255).Result()
	if err != nil {
		return nil, err
	}
	res := make([]Entity, 0)
	for _, k := range keys {
		e, err := db.GetOne(k)
		if err != nil {
			return nil, err
		}
		res = append(res, *e)
	}
	return res, conn.Close()
}

func (db *rdb) GetOne(id string) (*Entity, error) {
	cmd := db.client.Get(context.Background(), id)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	e := Entity{}
	if err := cmd.Scan(e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (db *rdb) Update(id string, data Entity) error {
	conn := db.client.Conn(context.Background())
	conn.Set(context.Background(), id, data, time.Hour*1000)
	return conn.Close()
}

func (db *rdb) Delete(id string) error {
	conn := db.client.Conn(context.Background())
	conn.Del(context.Background(), id)
	return conn.Close()
}

func NewRedis(uri string) DBAdapter {
	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "",
		DB:       0,
	})
	return &rdb{client}
}
