package main

import (
	"encoding/json"
	"github.com/mux"
	"go_server/database"
	"go_server/entities"
	"log"
	"net/http"
	"strconv"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(code)
	w.Write(response)
}

func main() {

	database.Init()

	TestData()

	r := mux.NewRouter()
	//r.Headers("Content-Type", "application/json")
	s := r.PathPrefix("/model").Subrouter()

	s.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		item := entities.Model{}
		model, err := item.Get(r)

		if err != "" {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		respondWithJSON(w, http.StatusOK, model)
	}).Methods("GET")

	s.HandleFunc("/{id}/mark", func(w http.ResponseWriter, r *http.Request) {

		model := entities.Model{}

		respondWithJSON(w, http.StatusOK, model)
	}).Methods("GET")

	// ADD
	s.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		model := entities.Model{}
		err := model.Create(r)

		if err != "" {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		respondWithJSON(w, http.StatusOK, "{}")
	}).Methods("POST")

	// EDIT
	s.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {

		model := entities.Model{}
		err := model.Update(r)

		if err != "" {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		respondWithJSON(w, http.StatusOK, "{}")
	}).Methods("POST")

	q := r.PathPrefix("/user").Subrouter()

	q.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		item := entities.User{}
		user, err := item.Get(r)

		if err != "" {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		respondWithJSON(w, http.StatusOK, user)
	}).Methods("GET")

	q.HandleFunc("/{id}/reviews", func(w http.ResponseWriter, r *http.Request) {

		model := entities.Model{}
		respondWithJSON(w, http.StatusOK, model)
	}).Methods("GET")

	q.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		user := entities.User{}
		err := user.Create(r)

		if err != "" {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		respondWithJSON(w, http.StatusOK, "{}")
	}).Methods("POST")

	q.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		user := entities.User{}
		err := user.Update(r)

		if err != "" {
			respondWithError(w, http.StatusBadRequest, err)
			return
		}

		respondWithJSON(w, http.StatusOK, user)
	}).Methods("POST")

	w := r.PathPrefix("/review").Subrouter()

	w.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request data")
			return
		}

		item := entities.Review{}
		review, err := item.Get(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error getting review")
			return
		}

		if review == nil {
			respondWithError(w, http.StatusBadRequest, "Review not found")
			return
		}

		respondWithJSON(w, http.StatusOK, review)
	}).Methods("GET")

	err := http.ListenAndServe(":8080", r) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServer error: ", err)
	}

}

func TestData() {
	db := database.GetStorage()
	txn := db.Txn(true)

	m := &entities.Model{
		Id:    1,
		Name:  "jetta",
		Brand: "vw",
		Year:  2015,
	}
	txn.Insert("models", m)

	u := &entities.User{
		Id:        1,
		Email:     "eka@kodix.ru",
		FirstName: "E",
		LastName:  "K",
		Sex:       "male",
		BirthDate: 123123123,
	}
	txn.Insert("users", u)

	rv := &entities.Review{
		Id:      1,
		User:    1,
		Model:   1,
		Created: 123,
		Mark:    3,
	}
	txn.Insert("reviews", rv)

	txn.Commit()
}

/*
GET /<entity>/<id>

GET /user/<id>/reviews

GET /model/<id>/mark

POST /<entity>/<id>

POST /<entity>
*/
