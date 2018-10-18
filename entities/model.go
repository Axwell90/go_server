package entities

import "go_server/database"

type Model struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Brand string `json:"brand"`
	Year  uint   `json:"year"`
}

func (m *Model) Validate() bool {
	if m.Id == 0 {
		return false
	}

	if m.Name == "" || len(m.Name) > 50 {
		return false
	}

	if m.Brand == "" || len(m.Brand) > 50 {
		return false
	}

	if m.Year == 0 {
		return false
	}

	return true
}

func (m *Model) Get(id uint) (*Model, error) {

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("models", "id", id)
	if err != nil {
		return nil, err
	}

	return raw.(*Model), nil
}

func (m *Model) Create(data *Model) error {

	db := database.GetStorage()

	txn := db.Txn(true)

	if err := txn.Insert("models", m); err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (m *Model) Update(data *Model) error {

	db := database.GetStorage()

	txn := db.Txn(true)

	if err := txn.Insert("models", m); err != nil {
		return err
	}

	txn.Commit()

	return nil
}
