package model

import "github.com/fwidjaya20/goloquent/pkg/goloquent"

// Genre .
type Genre struct {
	goloquent.Model `json:"-"`
	ID              int64  `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`
}

// GenreModel .
func GenreModel() *Genre {
	return &Genre{
		Model: goloquent.AutoIncrementModel("genres", "id", true, false),
	}
}
