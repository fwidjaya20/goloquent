package goloquent

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func forceMigrate(tx *sqlx.Tx, forced bool) {
	if !forced {
		return
	}

	queryString := `DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON ALL TABLES IN SCHEMA public TO public;`

	_, err := tx.Exec(queryString)

	if nil != err {
		_ = tx.Rollback()
		panic(err)
	}

	tx.Commit()
}

// RunMeta is a function that will create metadata for migrations
func runMeta(db *sqlx.DB) {
	meta := Migration{
		Schema: []*Schema{
			Create("migrations", func(table *Schema) {
				table.Text("command")
				table.Text("migrate")
				table.Text("batch")
			}),
		},
	}

	if !isMetaExists(db) {
		meta.Run(db, 1)
	}
}

func isMetaExists(db *sqlx.DB) bool {
	return isTableExist(db, "migrations")
}

func isTableExist(db *sqlx.DB, table string) bool {
	query := fmt.Sprintf("SELECT * FROM %s", table)

	_, err := db.Exec(query)

	if nil != err {
		return false
	}

	return true
}

func seedMetaTable(tx *sqlx.Tx, blueprint *Schema, batch int) {
	var query string

	query = "INSERT INTO migrations VALUES (:command, :migrate, :batch)"

	_, err := tx.NamedExec(query, map[string]interface{}{
		"command": blueprint.command,
		"migrate": blueprint.name,
		"batch":   batch,
	})

	if nil != err {
		tx.Rollback()
		panic(err)
	}
}
