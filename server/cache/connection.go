package cache

import (
	"fmt"
	"os"

	gredis "github.com/go-redis/redis/v8"
)

var host, port, password string
var db int

var connection *gredis.Client

//var Config *object.Configuration
//  GetConnection ...  Gives connection to redis at a time only one connection to redis is maintained
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
	host = os.Getenv("REDIS_HOST")
	port = os.Getenv("REDIS_PORT")
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
