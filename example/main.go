package main

import (
	"fmt"
	"time"

	"github.com/fwidjaya20/goloquent/config"
	"github.com/fwidjaya20/goloquent/example/migration"
	"github.com/fwidjaya20/goloquent/example/migration/seeder"
	"github.com/fwidjaya20/goloquent/example/model"
	"github.com/fwidjaya20/goloquent/pkg/goloquent"
)

func main() {
	fmt.Println(" * Goloquent * ")
	fmt.Println("===============")

	migrationSample()

	// seederSample()

	insertSample()

	// insertSample2()

	// updateSample()

	// deleteSample()

	// selectSample()
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

func seederSample() {
	goloquent.Seeder(config.GetDB(), "goloquent",
		seeder.GendreSeeder(),
	)
}

func insertSample() {
	query := goloquent.DB(config.GetDB())

	// // Insert Without Transaction
	// for i := 1; i <= 5; i++ {
	// 	genre := model.GenreModel()

	// 	genre.Name = fmt.Sprintf("Testing without Transaction %02d", i)

	// 	_, err := query.Use(genre).Insert()

	// 	if nil != err {
	// 		fmt.Println(err)
	// 	}
	// }

	// // Insert Using Transaction
	// for i := 6; i <= 10; i++ {
	// 	query.BeginTransaction()

	// 	genre := model.GenreModel()

	// 	genre.Name = fmt.Sprintf("Testing with Transaction %02d", i)

	// 	_, err := query.Use(genre).Insert()

	// 	if nil != err {
	// 		query.Rollback()
	// 		fmt.Println(err)
	// 	}

	// 	query.Commit()

	// 	query.EndTransaction()
	// }

	// // Insert Bulk Without Transaction
	// var payload []*model.Genre

	// for i := 11; i <= 15; i++ {
	// 	genre := model.GenreModel()

	// 	genre.Name = fmt.Sprintf("Testing Bulk Without Transaction %02d", i)

	// 	payload = append(payload, genre)
	// }

	// _, err := query.Use(model.GenreModel()).BulkInsert(payload)

	// if nil != err {
	// 	fmt.Println(err)
	// }

	// // Insert Bulk With Transaction
	// payload = []*model.Genre{}

	// for i := 16; i <= 20; i++ {
	// 	genre := model.GenreModel()

	// 	genre.Name = fmt.Sprintf("Testing Bulk With Transaction %02d", i)

	// 	payload = append(payload, genre)
	// }

	// query.BeginTransaction()

	// _, err = query.Use(model.GenreModel()).BulkInsert(payload)

	// if nil != err {
	// 	query.Rollback()
	// 	fmt.Println(err)
	// }

	// query.Commit()

	// query.EndTransaction()

	// Insert Raw without Transaction
	payload1 := map[string]interface{}{
		"name":       "Testing Raw without Transaction 21",
		"created_at": time.Now(),
	}

	result, err := query.RawCommand(model.GenreModel(), `insert into genres ("name", "created_at") values (:name, :created_at) returning *;`, payload1)

	if nil != err {
		fmt.Println(err)
	}

	fmt.Println(result.(*model.Genre))

	// Insert Raw with Transaction
	query.BeginTransaction()

	payload2 := map[string]interface{}{
		"name":       "Testing Raw without Transaction 22",
		"created_at": time.Now(),
	}

	result, err = query.RawCommand(model.GenreModel(), `insert into genres ("name", "created_at") values (:name, :created_at) returning *;`, payload2)

	if nil != err {
		query.Rollback()
		fmt.Println(err)
	}

	fmt.Println(result.(*model.Genre))

	query.Commit()

	query.EndTransaction()
}

func insertSample2() {
	query := goloquent.DB(config.GetDB())

	genre := model.GenreModel()

	genre.Name = "Testing Return Id"

	result, err := query.Use(genre).Insert()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(result.(*model.Genre).ID)
	fmt.Println(result.(*model.Genre).Name)

	// Insert Bulk Without Transaction
	var payload []*model.Genre

	for i := 1; i <= 5; i++ {
		genre := model.GenreModel()

		genre.Name = fmt.Sprintf("Testing Return Id %02d", i)

		payload = append(payload, genre)
	}

	result, err = query.Use(model.GenreModel()).BulkInsert(payload)

	if nil != err {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func updateSample() {
	query := goloquent.DB(config.GetDB())

	id := 3

	genre, err := query.Use(model.GenreModel()).
		Find(id)

	if nil != err {
		fmt.Println(err)
		return
	}

	v := genre.(*model.Genre)

	v.Name = fmt.Sprintf("Hello World #%d", id)

	isUpdated, err := query.Use(v).Update()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(isUpdated)
}

func deleteSample() {
	query := goloquent.DB(config.GetDB())

	genre, err := query.Use(model.GenreModel()).First()

	if nil != err {
		fmt.Println(err)
		return
	}

	v := genre.(*model.Genre)

	isDeleted, err := query.Use(v).Delete()

	fmt.Println(isDeleted)
}

func selectSample() {
	query := goloquent.DB(config.GetDB())

	m := model.GenreModel()

	getStmt(query, m)
	getInStmt(query, m)
	getExceptStmt(query, m)
	getBetweenStmt(query, m)
	getNotBetweenStmt(query, m)
	getIsNullStmt(query, m)
	getIsNotNullStmt(query, m)
	getCompareColumnStmt(query, m)
	getGroupByAndHavingStmt(query, m)
	allStmt(query, m)
	firstStmt(query, m)
	paginateStmt(query, m)
	aggregateStmt(query, m)
}

func getStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		Where("name", goloquent.ILIKE, "%bulk%").
		OrWhere("id", "=", 1).
		OrderBy("DESC", "id").
		Get()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("GET - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func getInStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		WhereIn("name", []string{"Action", "Crime", "Horror"}).
		OrWhereIn("id", []int{11, 12, 13}).
		Get()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("GET IN - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func getExceptStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		Except("name", []string{"Action", "Crime", "Horror"}).
		OrExcept("id", []int{6, 7, 8, 9, 10, 20, 21, 22}).
		Get()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("GET Except - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func getBetweenStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		WhereBetween("id", 5, 10).
		OrWhereBetween("id", 21, 25).
		Get()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("GET BETWEEN - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func getNotBetweenStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		WhereNotBetween("id", 1, 18).
		OrWhereNotBetween("id", 21, 25).
		Get()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("GET NOT BETWEEN - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func getIsNullStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		WhereNull("created_at").
		OrWhereNull("created_at").
		Get()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("GET NULL - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func getIsNotNullStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		WhereNotNull("updated_at").
		OrWhereNotNull("updated_at").
		Get()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("GET NOT NULL - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func getCompareColumnStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		WhereColumn("created_at", "!=", "updated_at").
		OrWhereColumn("created_at", "!=", "updated_at").
		Get()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("GET COMPARE COLUMN - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func getGroupByAndHavingStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		GroupBy("id", "name", "created_at", "updated_at").
		OrderBy("ASC", "id", "name").
		Get()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("GET COMPARE COLUMN - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func allStmt(query *goloquent.Query, m goloquent.IModel) {
	genres, err := query.Use(m).
		All()

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("ALL - Statement")
	for i, v := range genres.([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func firstStmt(query *goloquent.Query, m goloquent.IModel) {
	genre, err := query.Use(m).
		First()

	if nil != err {
		fmt.Println(err)
		return
	}

	v := genre.(*model.Genre)

	fmt.Println("FIRST - Statement")
	fmt.Printf("Genre #%02d\n", 1)
	fmt.Println("==========")
	fmt.Printf("ID   : %d\n", v.ID)
	fmt.Printf("Name : %s\n", v.Name)
	fmt.Println("==========")
}

func paginateStmt(query *goloquent.Query, m goloquent.IModel) {
	currentPage := 3
	limit := 10

	data, err := query.Use(m).
		Paginate(currentPage, limit)

	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Printf("PAGINATE (Page #%d, Limit %d) - Total Data : %d - Statement\n", currentPage, data["total"], limit)
	for i, v := range data["data"].([]*model.Genre) {
		fmt.Printf("Genre #%02d\n", i+1)
		fmt.Println("==========")
		fmt.Printf("ID   : %d\n", v.ID)
		fmt.Printf("Name : %s\n", v.Name)
		fmt.Println("==========")
	}
}

func aggregateStmt(query *goloquent.Query, m goloquent.IModel) {
	count := query.Use(m).Where("name", "ILIKE", "%bulk%").Count()

	fmt.Println("COUNT - Aggregate Statement")
	fmt.Printf("Total : %d\n", count)

	max := query.Use(m).Where("name", "ILIKE", "%bulk%").Max("id")

	fmt.Println("MAX - Aggregate Statement")
	fmt.Printf("Max : %d\n", int64(max))

	min := query.Use(m).Where("name", "ILIKE", "%bulk%").Min("id")

	fmt.Println("MAX - Aggregate Statement")
	fmt.Printf("Max : %d\n", int64(min))

	avg := query.Use(m).Where("name", "ILIKE", "%bulk%").Avg("id")

	fmt.Println("AVG - Aggregate Statement")
	fmt.Printf("Avg : %d\n", int64(avg))

	sum := query.Use(m).Where("name", "ILIKE", "%bulk%").Sum("id")

	fmt.Println("SUM - Aggregate Statement")
	fmt.Printf("Sum : %d\n", int64(sum))
}
