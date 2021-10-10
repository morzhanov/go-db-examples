package internal

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

type n4j struct {
	uri string
}

func (db *n4j) Create(data Entity) error {
	driver, err := neo4j.NewDriver(db.uri, neo4j.BasicAuth("", "", ""))
	if err != nil {
		return err
	}
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run(
			"CREATE (e:$id) SET e.id = $id, e.name = $name, e.description = $description, e.enabled = $enabled",
			map[string]interface{}{"id": data.Id, "name": data.Name, "description": data.Description, "enabled": data.Enabled})
		if err != nil {
			return nil, err
		}
		return nil, res.Err()
	})
	if err != nil {
		return err
	}
	if err := session.Close(); err != nil {
		return err
	}
	return driver.Close()
}

func (db *n4j) Get() ([]Entity, error) {
	driver, err := neo4j.NewDriver(db.uri, neo4j.BasicAuth("", "", ""))
	if err != nil {
		return nil, err
	}
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	results, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (n) RETURN n",
			map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		results := make([]interface{}, 0)
		for result.Next() {
			results = append(results, result.Record().Values[0])
		}
		return results, result.Err()
	})
	if err != nil {
		return nil, err
	}

	arr := results.([]interface{})
	res := make([]Entity, 0)
	for _, r := range arr {
		ent := r.(map[string]interface{})
		el := Entity{
			Id:          ent["Id"].(string),
			Name:        ent["Name"].(string),
			Description: ent["Description"].(string),
			Enabled:     ent["Enabled"].(bool),
		}
		res = append(res, el)
	}

	if err := session.Close(); err != nil {
		return nil, err
	}
	return res, driver.Close()
}

func (db *n4j) GetOne(id string) (*Entity, error) {
	driver, err := neo4j.NewDriver(db.uri, neo4j.BasicAuth("", "", ""))
	if err != nil {
		return nil, err
	}
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (n {id: $id}) RETURN n",
			map[string]interface{}{"id": id})
		if err != nil {
			return nil, err
		}
		return result.Record().Values[0], result.Err()
	})
	if err != nil {
		return nil, err
	}

	ent := result.(map[string]interface{})
	el := Entity{
		Id:          ent["Id"].(string),
		Name:        ent["Name"].(string),
		Description: ent["Description"].(string),
		Enabled:     ent["Enabled"].(bool),
	}

	if err := session.Close(); err != nil {
		return nil, err
	}
	return &el, driver.Close()
}

func (db *n4j) Update(id string, data Entity) error {
	driver, err := neo4j.NewDriver(db.uri, neo4j.BasicAuth("", "", ""))
	if err != nil {
		return err
	}
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run(
			"MATCH (n {id: $id}) SET n.name = $name, n.description = $description, n.enabled = $enabled",
			map[string]interface{}{"id": id, "name": data.Name, "description": data.Description, "enabled": data.Enabled})
		if err != nil {
			return nil, err
		}
		return nil, res.Err()
	})
	if err != nil {
		return err
	}
	if err := session.Close(); err != nil {
		return err
	}
	return driver.Close()
}

func (db *n4j) Delete(id string) error {
	driver, err := neo4j.NewDriver(db.uri, neo4j.BasicAuth("", "", ""))
	if err != nil {
		return err
	}
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run(
			"MATCH (n {id: $id}) DELETE n",
			map[string]interface{}{"id": id})
		if err != nil {
			return nil, err
		}
		return nil, res.Err()
	})
	if err != nil {
		return err
	}
	if err := session.Close(); err != nil {
		return err
	}
	return driver.Close()
}

func NewNeo4j(uri string) DBAdapter {
	return &n4j{uri}
}
