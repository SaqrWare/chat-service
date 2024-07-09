package data

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
)

// CassandraSession Cassandra connection
var CassandraSession *gocql.Session

func InitiateCassandraSession() *gocql.Session {
	hosts := os.Getenv("CASSANDRA_HOSTS")
	if hosts == "" {
		hosts = "127.0.0.1" // Default host
	}
	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	if keyspace == "" {
		keyspace = "chat" // Default keyspace
	}

	cluster := gocql.NewCluster(hosts)
	cluster.Keyspace = keyspace
	var err error
	CassandraSession, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	return CassandraSession
}

// RedisClient Redis Connection
var RedisClient *redis.Client
var RCTX = context.Background()

func InitiateRedisClient() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // Default address
	}
	password := os.Getenv("REDIS_PASSWORD")

	dbStr := os.Getenv("REDIS_DB")
	db := 0 // Default DB
	if dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			log.Fatalf("Invalid REDIS_DB value: %v", err)
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
