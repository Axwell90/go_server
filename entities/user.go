package entities

import "go_server/database"

type User struct {
	Id        uint
	Email     string
	FirstName string
	LastName  string
	Sex       string
	BirthDate uint
}

func (u *User) Validate() bool {
	if u.Id <= 0 {
		return false
	}

	if u.Email == "" || len(u.Email) > 100 {
		return false
	}
	if u.FirstName == "" || len(u.FirstName) > 50 {
		return false
	}

	if u.LastName == "" || len(u.LastName) > 50 {
		return false
	}

	if u.Sex != "male" && u.Sex != "female" {
		return false
	}

	if u.BirthDate <= 0 {
		return false
	}

	return true
}

func (u *User) Get(id int) (*User, error) {

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("users", "id", uint(id))
	if err != nil || raw == nil {
		return nil, err
	}

	return raw.(*User), nil
}

func (u *User) Create() error {

	db := database.GetStorage()

	txn := db.Txn(true)

	if err := txn.Insert("users", u); err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (u *User) Update() error {

	db := database.GetStorage()

	txn := db.Txn(true)

	if err := txn.Insert("users", u); err != nil {
		return err
	}

	txn.Commit()

	return nil
}
