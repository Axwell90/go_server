package entities

import (
	"encoding/json"
	"github.com/mux"
	"go_server/database"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Model struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Brand string `json:"brand"`
	Year  uint   `json:"year"`
}

func (m *Model) Validate() bool {
	if m.Id <= 0 {
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

func (m *Model) Get_(r *http.Request) (*Model, string) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		return nil, "Invalid request data"
	}

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("models", "id", uint(id))

	if err != nil {
		return nil, "Error getting model"
	}

	if raw == nil {
		return nil, "Model not found"
	}

	return raw.(*Model), ""
}

//func (m *Model) Get(id int) (*Model, error) {
//
//	db := database.GetStorage()
//
//	txn := db.Txn(false)
//	defer txn.Abort()
//
//	raw, err := txn.First("models", "id", uint(id))
//	if err != nil || raw == nil {
//		return nil, err
//	}
//
//	return raw.(*Model), nil
//}

func checkExistModel(id int) uint {

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("models", "id", uint(id))
	if err != nil || raw == nil {
		return 0
	}

	return raw.(*Model).Id
}

func (m *Model) Create_(r *http.Request) string {
	model := &Model{}

	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(b, &model)
	}
	r.Body.Close()

	if err != nil {
		return "Invalid request data"
	}

	check := checkExistModel(int(model.Id))
	if check > 0 {
		return "Model is already exist"
	}

	if model.Validate() != true {
		return "Validate error"
	}

	db := database.GetStorage()

	txn := db.Txn(true)

	if err := txn.Insert("models", model); err != nil {
		return "Error add model"
	}

	txn.Commit()

	return ""
}

func (m *Model) Update_(r *http.Request) string {
	updModel := Model{}

	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(b, &updModel)
	}
	r.Body.Close()

	if err != nil {
		return "Invalid request data"
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		return "Invalid request data"
	}

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("models", "id", uint(id))

	if err != nil {
		return "Error getting model"
	}

	if raw == nil {
		return "Model not found"
	}

	model := &Model{}
	model = raw.(*Model)

	if updModel.Name != "" {
		model.Name = updModel.Name
	}

	if updModel.Brand != "" {
		model.Brand = updModel.Brand
	}

	if updModel.Year >= 0 {
		model.Year = updModel.Year
	}

	if model.Validate() != true {
		return "New model validate error"
	}

	txn = db.Txn(true)

	if err := txn.Insert("models", model); err != nil {
		return "Error update model"
	}

	txn.Commit()

	return ""
}

func (m *Model) Create() error {

	db := database.GetStorage()

	txn := db.Txn(true)

	if err := txn.Insert("models", m); err != nil {
		return err
	}

	txn.Commit()

	return nil
}

//func (m *Model) Update() error {
//
//	db := database.GetStorage()
//
//	txn := db.Txn(true)
//
//	if err := txn.Insert("models", m); err != nil {
//		return err
//	}
//
//	txn.Commit()
//
//	return nil
//}
