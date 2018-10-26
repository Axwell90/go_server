package entities

import (
	"encoding/json"
	"github.com/mux"
	"go_server/database"
	"io/ioutil"
	"math"
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

func (m *Model) Get(r *http.Request) (*Model, string) {
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

func (m *Model) Create(r *http.Request) string {
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

func (m *Model) Update(r *http.Request) string {
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

func (m *Model) GetMark(r *http.Request) (interface{}, int, string) {
	sum := uint(0)
	i := 0
	from := uint(0)
	to := ^uint(0)
	var resp interface{}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		return 0, 400, "Invalid request data"
	}

	parameters := r.URL.Query()

	fromDateParam := ""
	arFromDate, _ := parameters["fromDate"]
	if len(arFromDate) > 0 {
		fromDateParam = arFromDate[0]
	}

	if fromDateParam != "" {
		fromDate, err := strconv.Atoi(fromDateParam)
		if err != nil {
			return 0, 400, "Invalid request data"
		} else {
			from = uint(fromDate)
		}
	}

	toDateParam := ""
	arToDate, _ := parameters["toDate"]
	if len(arToDate) > 0 {
		toDateParam = arToDate[0]
	}

	if toDateParam != "" {
		toDate, err := strconv.Atoi(toDateParam)
		if err != nil {
			return 0, 400, "Invalid request data"
		} else {
			to = uint(toDate)
		}
	}

	sexParam := ""
	arSex, _ := parameters["sex"]
	if len(arSex) > 0 {
		sexParam = arSex[0]
		if sexParam != "male" && sexParam != "female" {
			return 0, 400, "Invalid request data"
		}
	}

	db := database.GetStorage()

	txn := db.Txn(false)

	raw, err := txn.Get("reviews", "model", uint(id))
	if err != nil {
		return 0, 400, "Error getting model"
	}

	found := false

	for item := raw.Next(); item != nil; item = raw.Next() {
		found = true
		if item.(*Review).Created > from && item.(*Review).Created < to {
			if sexParam != "" {
				rawUser, err := txn.First("users", "id", item.(*Review).User)
				if err != nil {
					return 0, 400, "Error getting user"
				}
				if sexParam == rawUser.(*User).Sex {
					sum += item.(*Review).Mark
					i += 1
				}
			} else {
				sum += item.(*Review).Mark
				i += 1
			}
		}
	}

	if found == false {
		return 0, 404, "Model not found"
	}

	var mark float64

	if i == 0 {
		mark = 0
	} else {
		mark = float64(sum) / float64(i)
		var round float64
		pow := math.Pow(10, float64(5))
		digit := pow * mark
		round = math.Ceil(digit)
		mark = round / pow
	}

	resp = struct {
		Mark float64 `json:"mark"`
	}{Mark: mark}

	return resp, 200, ""
}
