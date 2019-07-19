package main

import (
	"fmt"

	"github.com/fwidjaya20/goloquent/config"
	"github.com/fwidjaya20/goloquent/example/migration"
	"github.com/fwidjaya20/goloquent/example/model"
	"github.com/fwidjaya20/goloquent/pkg/goloquent"
)

func main() {
	fmt.Println(" * Goloquent * ")
	fmt.Println("===============")

	migrationSample()

	insertSample()
}

func migrationSample() {
	// fmt.Println("========================================")
	// for _, v := range migration.Migration1.Schema {
	// 	v.Verbose()
	// 	fmt.Println("========================================")
	// }

	goloquent.Migrate(config.GetDB(), "goloquent",
		migration.Migration1,
		migration.Migration2,
		migration.Migration3,
	)
}

func insertSample() {
	query := goloquent.DB(config.GetDB())

	// Insert Without Transaction
	for i := 1; i <= 5; i++ {
		genre := model.GenreModel()

		genre.Name = fmt.Sprintf("Testing without Transaction %02d", i)

		_, err := query.Use(genre).Insert()

		if nil != err {
			fmt.Println(err)
		}
	}

	// Insert Using Transaction
	for i := 6; i <= 10; i++ {
		query.BeginTransaction()

		genre := model.GenreModel()

		genre.Name = fmt.Sprintf("Testing with Transaction %02d", i)

		_, err := query.Use(genre).Insert()

		if nil != err {
			query.Rollback()
			fmt.Println(err)
		}

		query.Commit()

		query.EndTransaction()
	}

	// Insert Bulk Without Transaction
	var payload []*model.Genre

	for i := 11; i <= 15; i++ {
		genre := model.GenreModel()

		genre.Name = fmt.Sprintf("Testing Bulk Without Transaction %02d", i)

		payload = append(payload, genre)
	}

	_, err := query.Use(model.GenreModel()).BulkInsert(payload)

	if nil != err {
		fmt.Println(err)
	}

	// Insert Bulk With Transaction
	payload = []*model.Genre{}

	for i := 16; i <= 20; i++ {
		genre := model.GenreModel()

		genre.Name = fmt.Sprintf("Testing Bulk With Transaction %02d", i)

		payload = append(payload, genre)
	}

	query.BeginTransaction()

	_, err = query.Use(model.GenreModel()).BulkInsert(payload)

	if nil != err {
		query.Rollback()
		fmt.Println(err)
	}

	query.Commit()

	query.EndTransaction()

	// Insert Raw without Transaction
	_, err = query.RawCommand(`insert into genres ("name") values ($1);`, "Testing Raw without Transaction 21")

	if nil != err {
		fmt.Println(err)
	}

	// Insert Raw with Transaction
	query.BeginTransaction()

	_, err = query.RawCommand(`insert into genres ("name") values ($1);`, "Testing Raw with Transaction 22")

	if nil != err {
		query.Rollback()
		fmt.Println(err)
	}

	query.Commit()

	query.EndTransaction()
}
