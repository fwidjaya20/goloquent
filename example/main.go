package main

import (
	"fmt"

	"github.com/fwidjaya20/goloquent/config"
	"github.com/fwidjaya20/goloquent/example/migration"
	"github.com/fwidjaya20/goloquent/pkg/goloquent"
)

func main() {
	fmt.Println(" * Goloquent * ")
	fmt.Println("===============")

	migrationSample()
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
