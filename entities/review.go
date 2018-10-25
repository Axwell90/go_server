package entities

import "go_server/database"

type Review struct {
	Id      uint `json:"id"`
	User    uint `json:"user"`
	Model   uint `json:"model"`
	Created uint `json:"created"`
	Mark    uint `json:"mark"`
}

func (r *Review) Validate() bool {
	if r.Id <= 0 {
		return false
	}

	if r.Created <= 0 {
		return false
	}
	if r.Mark < 0 || r.Mark > 5 {
		return false
	}

	return true
}

func (r *Review) Get(id int) (*Review, error) {

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("reviews", "id", uint(id))
	if err != nil || raw == nil {
		return nil, err
	}

	return raw.(*Review), nil
}

func (r *Review) Create() error {

	db := database.GetStorage()

	txn := db.Txn(true)

	if err := txn.Insert("reviews", r); err != nil {
		return err
	}

	txn.Commit()

	return nil
}
