package internal

import "github.com/spf13/viper"

type Config struct {
	CassandraUri  string `mapstructure:"CASSANDRA_URI"`
	ClickhouseUri string `mapstructure:"CLICKHOUSE_URI"`
	CockroachUri  string `mapstructure:"COCKROACH_URI"`
	CouchbaseUri  string `mapstructure:"COUCHBASE_URI"`
	CouchdbUri    string `mapstructure:"COUCHDB_URI"`
	FirebaseUri   string `mapstructure:"FIREBASE_URI"`
	LeveldbUri    string `mapstructure:"LEVELDB_URI"`
	MongodbUri    string `mapstructure:"MONGODB_URI"`
	Neo4jUri      string `mapstructure:"NEO4J_URI"`
	PostgresqlUri string `mapstructure:"POSTGRESQL_URI"`
	RedisUri      string `mapstructure:"REDIS_URI"`
	SolrUri       string `mapstructure:"SOLR_URI"`
}

func NewConfig() (config *Config, err error) {
	viper.AddConfigPath("./configs/")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
