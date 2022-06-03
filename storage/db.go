package storage

import (
	"database/sql"

	_ "github.com/lib/pq" // postgres driver
	"github.com/okobsamoht/tomato/config"
	"github.com/okobsamoht/tomato/test"
	"gopkg.in/mgo.v2"
)

// OpenMongoDB 打开 MongoDB
func OpenMongoDB() *mgo.Database {
	// 此处仅用于测试
	if config.TConfig.DatabaseURI == "" {
		config.TConfig.DatabaseURI = test.MongoDBTestURL
	}

	session, err := mgo.Dial(config.TConfig.DatabaseURI)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("")
}

// OpenPostgreSQL 打开 PostgreSQL
func OpenPostgreSQL() *sql.DB {
	db, err := sql.Open("postgres", config.TConfig.DatabaseURI)
	if err != nil {
		panic(err)
	}
	return db
}
