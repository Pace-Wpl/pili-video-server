package dbops

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
)

var (
	dbConn *sql.DB
	Pool   redis.Pool
	err    error
)

func init() {
	//mysql db init
	fmt.Println("Entering conn.go init function...")
	dbConn, err = sql.Open("mysql", "pace:123@/sample")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("dbConn open +%v\n", dbConn)

	//redis pool init
	Pool = redis.Pool{
		MaxIdle:     16,
		MaxActive:   32,
		IdleTimeout: 120,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}
