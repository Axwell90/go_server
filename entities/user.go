package entities

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go_server/database"
	"io/ioutil"
	"net/http"
	"strconv"
)

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

func checkExistUser(id int) uint {

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("users", "id", uint(id))
	if err != nil || raw == nil {
		return 0
	}

	return raw.(*User).Id
}

func (u *User) Get(r *http.Request) (*User, string) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		return nil, "Invalid request data"
	}

	db := database.GetStorage()

	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("users", "id", uint(id))

	if err != nil {
		return nil, "Error getting user"
	}

	if raw == nil {
		return nil, "User not found"
	}

	return raw.(*User), ""
}

func (u *User) Create(r *http.Request) string {
	user := &User{}

	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(b, &user)
	}
	r.Body.Close()

	if err != nil {
		return "Invalid request data"
	}

	check := checkExistUser(int(user.Id))
	if check > 0 {
		return "User is already exist"
	}

	if user.Validate() != true {
		return "Validate error"
	}

	db := database.GetStorage()

	txn := db.Txn(true)

	if err := txn.Insert("users", user); err != nil {
		return "Error add user"
	}

	txn.Commit()

	return ""
}

func (u *User) Update(r *http.Request) string {
	updUser := User{}

	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(b, &updUser)
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

	raw, err := txn.First("users", "id", uint(id))

	if err != nil {
		return "Error getting user"
	}

	if raw == nil {
		return "User not found"
	}

	user := &User{}
	user = raw.(*User)

	if updUser.Email != "" {
		user.Email = updUser.Email
	}

	if updUser.FirstName != "" {
		user.FirstName = updUser.FirstName
	}

	if updUser.LastName != "" {
		user.LastName = updUser.LastName
	}

	if updUser.Sex != "" {
		user.Sex = updUser.Sex
	}

	if user.Validate() != true {
		return "New user validate error"
	}

	txn = db.Txn(true)

	if err := txn.Insert("users", user); err != nil {
		return "Error update user"
	}

	txn.Commit()

	return ""
}

func (u *User) GetReviews(r *http.Request) (interface{}, int, string) {
	arReviews := make([]*Review, 0)
	var resp interface{}
	from := uint(0)
	to := ^uint(0)

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

	db := database.GetStorage()

	txn := db.Txn(false)

	raw, err := txn.Get("reviews", "user", uint(id))
	if err != nil {
		return 0, 400, "Error getting reviews"
	}

	found := false

	for item := raw.Next(); item != nil; item = raw.Next() {
		found = true
		if item.(*Review).Created > from && item.(*Review).Created < to {
			arReviews = append(arReviews, item.(*Review))
		}
	}

	if found == false {
		rawUser, err := txn.First("users", "id", uint(id))
		if err != nil {
			return 0, 400, "Error getting user"
		}
		if rawUser == nil {
			return 0, 404, "User not found"
		}
	}

	resp = struct {
		Reviews []*Review `json:"reviews"`
	}{Reviews: arReviews}

	return resp, 200, ""
}
