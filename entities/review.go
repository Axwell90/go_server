package entities

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go_server/database"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Review struct {
	Id      uint `json:"id"`
	User    uint `json:"user"`
	Model   uint `json:"model"`
	Created uint `json:"created"`
	Mark    uint `json:"mark"`
}

func (rv *Review) Validate() bool {
	if rv.Id <= 0 {
		return false
	}

	if rv.Created <= 0 {
		return false
	}
	if rv.Mark < 0 || rv.Mark > 5 {
		return false
	}

	return true
}

func checkExistReview(id int) uint {

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("reviews", "id", uint(id))
	if err != nil || raw == nil {
		return 0
	}

	return raw.(*Review).Id
}

func (rv *Review) Get(r *http.Request) (*Review, string) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, "Invalid request data"

	}

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("reviews", "id", uint(id))

	if err != nil {
		return nil, "Error getting review"
	}

	if raw == nil {
		return nil, "Review not found"
	}

	return raw.(*Review), ""
}

func (rv *Review) Create(r *http.Request) string {
	review := &Review{}

	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(b, &review)
	}
	r.Body.Close()

	if err != nil {
		return "Invalid request data"
	}

	check := checkExistReview(int(review.Id))
	if check > 0 {
		return "Review is already exist"
	}

	if review.Validate() != true {
		return "Validate error"
	}

	db := database.GetStorage()

	txn := db.Txn(true)

	if err := txn.Insert("reviews", review); err != nil {
		return "Error add review"
	}

	txn.Commit()

	return ""
}

func (rv *Review) Update(r *http.Request) string {
	updReview := Review{}

	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(b, &updReview)
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

	raw, err := txn.First("reviews", "id", uint(id))

	if err != nil {
		return "Error getting review"
	}

	if raw == nil {
		return "Review not found"
	}

	review := &Review{}
	review = raw.(*Review)

	if updReview.User >= 0 {
		review.User = updReview.User
	}

	if updReview.Model >= 0 {
		review.Model = updReview.Model
	}

	if updReview.Created >= 0 {
		review.Created = updReview.Created
	}

	if updReview.Mark >= 0 {
		review.Mark = updReview.Mark
	}

	if review.Validate() != true {
		return "New review validate error"
	}

	txn = db.Txn(true)

	if err := txn.Insert("reviews", review); err != nil {
		return "Error update review"
	}

	txn.Commit()

	return ""
}
