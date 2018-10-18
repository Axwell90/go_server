package main

import (
	"fmt"
	"go_server/database"
	"go_server/entities"
)

//"github.com/hashicorp/go-memdb"

type User struct {
	id         int32
	email      string
	first_name string
	last_name  string
	sex        string
	birth_date string
}

func main() {

	database.Init()

	m := &entities.Model{
		Id:    1,
		Name:  "jetta",
		Brand: "vw",
		Year:  2015,
	}

	//mm := Model{2,"tiguan", "vw", 2016}

	//// Create a write transaction
	//txn := db.Txn(true)

	if m.Validate() != true {
		// TODO: validate error
	}
	if err := m.Create(m); err != nil {
		panic(err)
	}

	m.Brand = "nissan"
	m.Update(m)

	emptyModel := entities.Model{}
	model, err := emptyModel.Get(1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello TYT %v!", model)

	//if err := txn.Insert("models", m); err != nil {
	//	panic(err)
	//}

	//if mm.Validate() != true {
	//	// TODO: validate error
	//}
	//if err := txn.Insert("models", mm); err != nil {
	//	panic(err)
	//}
	//
	//// Commit the transaction
	//txn.Commit()

	{
		/*txn := db.Txn(false)
		defer txn.Abort()

		//raw, err := txn.First("models", "id", "joe@aol.com")
		//if err != nil {
		//	panic(err)
		//}

		raw, err := txn.Get("models", "id", uint(1))
		if err != nil {
			panic(err)
		}

		for item := raw.Next(); item != nil; item = raw.Next() {
			fmt.Printf("Hello %v!", item)
		}*/

		//fmt.Printf("Hello %s!", raw.(*Model))

		//mmm := Model{3,"jetta", "vw", 2015}
	}

	// Insert a new person
	//p := &Person{"joe@aol.com", "Joe", 30}
	//if err := txn.Insert("person", p); err != nil {
	//	panic(err)
	//}

	// Commit the transaction
	//txn.Commit()
	//
	//// Create read-only transaction
	//txn = db.Txn(false)
	//defer txn.Abort()
	//
	//raw, err := txn.Get("models", "id")
	//if err != nil {
	//	panic(err)
	//}
	//
	//for item := raw.Next(); item != nil; item = raw.Next() {
	//	fmt.Printf("Hello %v!", item)
	//}

	// Lookup by email
	//raw, err := txn.First("person", "id", "joe@aol.com")
	//if err != nil {
	//	panic(err)
	//}

	//fmt.Printf("Hello %s!", raw.(*Person))

	//raw, err := txn.Get("person", "id", "joe@aol.com")
	//if err != nil {
	//	panic(err)
	//}
	//
	//for item := raw.Next(); item != nil; item = raw.Next() {
	//	fmt.Printf("Hello %s!", item)
	//}

}
