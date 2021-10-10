package internal

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
)

type fbase struct {
	client *db.Ref
}

func (f *fbase) Create(data Entity) error {
	c := f.client.Child(data.Id)
	return c.Set(context.Background(), data)
}

func (f *fbase) Get() ([]Entity, error) {
	q := f.client.OrderByChild("0")
	qn, err := q.GetOrdered(context.Background())
	if err != nil {
		return nil, err
	}
	res := make([]Entity, 0)
	for _, n := range qn {
		e := Entity{}
		if err := n.Unmarshal(e); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func (f *fbase) GetOne(id string) (*Entity, error) {
	c := f.client.Child(id)
	res := Entity{}
	if err := c.Get(context.Background(), res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (f *fbase) Update(id string, data Entity) error {
	c := f.client.Child(id)
	upd := make(map[string]interface{}, 0)
	upd["Name"] = data.Name
	upd["Description"] = data.Description
	upd["Enabled"] = data.Enabled
	return c.Update(context.Background(), upd)
}

func (f *fbase) Delete(id string) error {
	c := f.client.Child(id)
	return c.Delete(context.Background())
}

func NewFirebase(dbUri string) (DBAdapter, error) {
	app, err := firebase.NewApp(context.Background(), &firebase.Config{DatabaseURL: dbUri})
	if err != nil {
		return nil, err
	}
	d, err := app.Database(context.Background())
	if err != nil {
		return nil, err
	}
	r := d.NewRef("entities")
	return &fbase{r}, nil
}
