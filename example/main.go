package main

import (
	"fmt"

	"github.com/fwidjaya20/goloquent/example/migration"
)

func main() {
	fmt.Println(" * Goloquent * ")
	fmt.Println("===============")

	migrationSample()
}

func migrationSample() {
	fmt.Println("\n\nRunning Migration")

	fmt.Println("========================================")
	for _, v := range migration.Migration1.Schema {
		v.Verbose()
		fmt.Println("========================================")
	}
}
