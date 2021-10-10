package main

import (
	"fmt"
	"github.com/morzhanov/go-db-examples/internal"
	"log"
)

var dbs map[string]internal.DBAdapter
var srv internal.Service

func init() {
	config, err := internal.NewConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	dbs = make(map[string]internal.DBAdapter, 12)
	dbs["cassandra"] = internal.NewCassandra(config.CassandraUri)
	dbs["clickhouse"], err = internal.NewClickhouse(config.ClickhouseUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbs["cockroach"], err = internal.NewCockroach(config.CockroachUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbs["couchbase"], err = internal.NewCouchbase(config.CouchbaseUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbs["couchdb"], err = internal.NewCouchdb(config.CouchdbUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbs["firebase"], err = internal.NewFirebase(config.FirebaseUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbs["leveldb"], err = internal.NewLeveldb(config.LeveldbUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbs["mongodb"], err = internal.NewMongodb(config.MongodbUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbs["neo4j"] = internal.NewNeo4j(config.Neo4jUri)
	dbs["postgresql"], err = internal.NewPostgresql(config.PostgresqlUri)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbs["redis"] = internal.NewRedis(config.RedisUri)
	dbs["solr"] = internal.NewSolr(config.SolrUri)

	srv = internal.NewService()
}

func test(adapter internal.DBAdapter, dbName string, enabled bool, name string, desc string) {
	fmt.Printf("Testing %s...\n", dbName)
	if err := srv.Test(adapter, enabled, name, desc); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%s successfully tested\n", dbName)
}

func main() {
	for name, db := range dbs {
		test(
			db,
			name,
			true,
			fmt.Sprintf("%s-name", name),
			fmt.Sprintf("%s-description", name),
		)
	}
}
