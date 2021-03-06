package database

import "github.com/hashicorp/go-memdb"

var Storage *memdb.MemDB

func Init() {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"models": {
				Name: "models",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
				},
			},
			"users": {
				Name: "users",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
					"sex": {
						Name:    "sex",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Sex"},
					},
				},
			},
			"reviews": {
				Name: "reviews",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
					"model": {
						Name:    "model",
						Unique:  false,
						Indexer: &memdb.UintFieldIndex{Field: "Model"},
					},
					"user": {
						Name:    "user",
						Unique:  false,
						Indexer: &memdb.UintFieldIndex{Field: "User"},
					},
					"created": {
						Name:    "created",
						Unique:  false,
						Indexer: &memdb.UintFieldIndex{Field: "Created"},
					},
				},
			},
		},
	}
	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	} else {
		Storage = db
	}
}

func GetStorage() *memdb.MemDB {
	return Storage
}
