package config

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var once sync.Once

// GetDB .
func GetDB() *sqlx.DB {
	once.Do(func() {
		var err error
		conn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			DbConf["DB_HOST"],
			DbConf["DB_PORT"],
			DbConf["DB_USER"],
			DbConf["DB_PASS"],
			DbConf["DB_NAME"],
		)
		db, err = sqlx.Connect("postgres", conn)
		if nil != err {
			panic(err)
		}
	})
	return db
}
