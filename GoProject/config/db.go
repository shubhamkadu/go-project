package config

// import (
// 	"database/sql"
// 	"log"

// 	_ "github.com/go-sql-driver/mysql"
// )

// var DB *sql.DB

// func ConnectDB() {
// 	var err error
// 	DB, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/go_project")
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}
// 	if err = DB.Ping(); err != nil {
// 		log.Fatal("Database unreachable:", err)
// 	}
// 	log.Println("Database connected")
// }

// package database
// package database

import (
	"database/sql"
	"fmt"
	"go-gqlgen/constants"
	"os"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

var MySQLDB *sql.DB

// connect to database
func ConnectMySQLDB() (*sql.DB, error) {
	dbDriver := "mysql"
	dbUser := os.Getenv("MYSQL_DB_USERNAME")
	dbPass := os.Getenv("MYSQL_DB_PASSWORD")
	dbName := os.Getenv("MYSQL_DB_NAME")
	dbHost := os.Getenv("MYSQL_DB_HOST")
	dbPort := os.Getenv("MYSQL_DB_PORT")
	// Create the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Open a connection to the database
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(constants.DatabaseConnected)

	return db, nil
}

// connects to Redis
func ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})
	return rdb
}