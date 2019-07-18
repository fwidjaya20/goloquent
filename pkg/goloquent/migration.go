package goloquent

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Migration is a struct that is used to store information about migration
type Migration struct {
	Schema []*Schema
}

// Migrate is a function that is used to execute migration command
func Migrate(db *sqlx.DB, database string, versions ...Migration) {
	forceMigrate(db.MustBegin(), true)

	runMeta(db)

	for i, v := range versions {
		tx := db.MustBegin()
		v.Run(db, i+1)
		tx.Commit()
	}
}

// Run is a function that will run all migration tables in Migration
func (m *Migration) Run(db *sqlx.DB, batch int) {
	builder := NewBuilder()

	tx := db.MustBegin()

	for _, v := range m.Schema {
		queryString := ``

		switch v.command {
		case CMD_CREATE:
			queryString = builder.BuildCreateTable(v)
		case CMD_ALTER:
			queryString = builder.BuildAlterTable(v)
		case CMD_DROP:
			queryString = builder.BuildDropTable(v)
		default:
			panic("Invalid Migration Command!")
		}

		if "" != queryString {
			_, err := tx.Exec(queryString)

			if nil != err {
				tx.Rollback()
				panic(err)
			}

			seedMetaTable(tx, v, batch)

			fmt.Printf("Migrating table %v\n", v.name)
		}
	}

	tx.Commit()
}
