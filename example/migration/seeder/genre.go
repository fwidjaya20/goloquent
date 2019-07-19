package seeder

import (
	"github.com/fwidjaya20/goloquent/example/model"
	"github.com/fwidjaya20/goloquent/pkg/goloquent"
)

// GenreSeed .
type GenreSeed struct {
	goloquent.Seed
}

// GendreSeeder .
func GendreSeeder() *GenreSeed {
	return &GenreSeed{
		Seed: goloquent.Seed{
			Seeds: []map[string]interface{}{
				{"name": "Action"},
				{"name": "Crime"},
				{"name": "Horror"},
				{"name": "Thriller"},
				{"name": "Commedy"},
			},
		},
	}
}

// Model .
func (s *GenreSeed) Model() goloquent.IModel {
	return model.GenreModel()
}

// BulkModel .
func (s *GenreSeed) BulkModel() []goloquent.IModel {
	var bulks []goloquent.IModel

	for _, v := range s.Seeds {
		bulk := model.GenreModel()

		bulk.Name = v["name"].(string)

		bulks = append(bulks, bulk)
	}

	return bulks
}
