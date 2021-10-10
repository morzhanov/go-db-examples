package internal

import (
	"fmt"
	"github.com/vanng822/go-solr/solr"
	"net/url"
)

type slr struct {
	client *solr.SolrInterface
}

func (db *slr) Create(data Entity) error {
	m := map[string]interface{}{
		"Id":          data.Id,
		"Name":        data.Name,
		"Description": data.Description,
		"Enabled":     data.Enabled,
	}
	_, err := db.client.Add([]solr.Document{m}, 10000, nil)
	return err
}

func (db *slr) Get() ([]Entity, error) {
	q := solr.NewQuery()
	q.Q("*:*")
	s := db.client.Search(q)
	r, err := s.Result(nil)
	if err != nil {
		return nil, err
	}

	res := make([]Entity, 0)
	for _, doc := range r.Results.Docs {
		e := Entity{
			Id:          doc.Get("Id").(string),
			Name:        doc.Get("Name").(string),
			Description: doc.Get("Description").(string),
			Enabled:     doc.Get("Enabled").(bool),
		}
		res = append(res, e)
	}
	return res, nil
}

func (db *slr) GetOne(id string) (*Entity, error) {
	q := solr.NewQuery()
	q.Q(fmt.Sprintf("Id:%s", id))
	s := db.client.Search(q)
	r, err := s.Result(nil)
	if err != nil {
		return nil, err
	}
	d := r.Results.Docs[0]
	return &Entity{
		Id:          d.Get("Id").(string),
		Name:        d.Get("Name").(string),
		Description: d.Get("Description").(string),
		Enabled:     d.Get("Enabled").(bool),
	}, nil
}

func (db *slr) Update(_ string, data Entity) error {
	m := map[string]interface{}{
		"Id":          data.Id,
		"Name":        data.Name,
		"Description": data.Description,
		"Enabled":     data.Enabled,
	}
	_, err := db.client.Update(m, &url.Values{})
	return err
}

func (db *slr) Delete(id string) error {
	d, err := db.GetOne(id)
	if err != nil {
		return err
	}
	m := map[string]interface{}{
		"Id":          d.Id,
		"Name":        d.Name,
		"Description": d.Description,
		"Enabled":     d.Enabled,
	}
	_, err = db.client.Delete(m, &url.Values{})
	return err
}

func NewSolr(uri string) DBAdapter {
	si, _ := solr.NewSolrInterface(uri, "entities")
	return &slr{si}
}
