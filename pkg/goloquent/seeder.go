package goloquent

import (
	"github.com/jmoiron/sqlx"
)

// SeederInterface is an interface used for creating seeder
type SeederInterface interface {
	Model() IModel
	BulkModel() []IModel
}

// Seed is a struct that is used to store information about migration
type Seed struct {
	Seeds []map[string]interface{}
}

// Seeder is a function that is used to execute seeder command
func Seeder(db *sqlx.DB, table string, seeds ...SeederInterface) {
	query := DB(db)

	for _, v := range seeds {
		query.BeginTransaction()

		_, err := query.Use(v.Model()).BulkInsert(v.BulkModel())

		if nil != err {
			query.Rollback()
			panic(err)
		}

		query.Commit()

		query.EndTransaction()
	}
}
