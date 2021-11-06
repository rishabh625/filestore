package cache

import (
	"fmt"

	gredis "github.com/go-redis/redis/v8"
)

var host, port, password string
var db int

var connection *gredis.Client

//var Config *object.Configuration

func GetConnection() *gredis.Client {
	if connection == nil {
		var err error
		connection, err = newConnection()
		if err != nil {
			//fatal
		}
	}
	return connection
}

// NewConnection ... Gives connection to redis host passed
func newConnection() (*gredis.Client, error) {
	host = "localhost"
	port = "6379"
	rdb := gredis.NewClient(&gredis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	status := rdb.Ping(rdb.Context())
	err := status.Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
