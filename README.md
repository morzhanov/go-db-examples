# Go Database Examples

Golang Database examples contains base code for working with some popular databases using Go libraries.

<img src="https://i.ibb.co/gv27Xz4/IMG-0121.jpg" alt="arch"/>

## Project structure

- `cmd/main` - main go file
- `deploy` - contains docker-compose file to run databases locally
- `internal/service` - contains main business logic
- `internal/config` - configuration setup using viper and .env file
- `internal/db` - database adapter interface
- `internal/*` - other go files contain database adapter implementations

## DB types:

- Relational (PostgreSQL, Cockroach)
- Document (MongoDB, Couchbase)
- Key-Value (Redis)
- Graph (Neo4J)
- Columnar (Cassandra)
- Time-series (Cassandra, Solr)
- Realtime (Firebase)
