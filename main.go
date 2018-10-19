package main

import (
	"encoding/json"
	"fmt"
	"github.com/mux"
	"go_server/database"
	"go_server/entities"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	id         int32
	email      string
	first_name string
	last_name  string
	sex        string
	birth_date string
}

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

	m := &entities.Model{
		Id:    1,
		Name:  "jetta",
		Brand: "vw",
		Year:  2015,
	}
	if err := m.Create(); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Headers("Content-Type", "application/json")
	s := r.PathPrefix("/model").Subrouter()

	s.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request data")
			return
		}

		item := entities.Model{}
		model, err := item.Get(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error getting model")
			return
		}

		if model == nil {
			respondWithError(w, http.StatusBadRequest, "Model not found")
			return
		}

		respondWithJSON(w, http.StatusOK, model)
	}).Methods("GET")

	s.HandleFunc("/{id}/mark", func(w http.ResponseWriter, r *http.Request) {
		data := []byte(`{"Id":1,"Year":2010}`)

		updModel := entities.Model{}

		b, err := ioutil.ReadAll(r.Body)
		b = data
		if err == nil {
			err = json.Unmarshal(b, &updModel)
		}
		r.Body.Close()

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request data")
			return
		}
		id := updModel.Id

		item := entities.Model{}
		model, err := item.Get(int(id))
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error getting model")
			return
		}

		if model == nil {
			respondWithError(w, http.StatusBadRequest, "Model not found")
			return
		}

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
			respondWithError(w, http.StatusBadRequest, "New model validate error")
			return
		}

		err = model.Update()
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error update model")
			return
		}

		respondWithJSON(w, http.StatusOK, model)

		/*
			model := entities.Model{}

			b, err := ioutil.ReadAll(r.Body)
			b = data
			if err == nil {
				err = json.Unmarshal(b, &model)
			}
			r.Body.Close()

			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid request data")
				return
			}

			if model.Validate() != true {
				respondWithError(w, http.StatusBadRequest, "Validate error")
				return
			}

			err = model.Create()
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Error add model")
				return
			}

			respondWithJSON(w, http.StatusOK, "{}")
		*/

	}).Methods("GET")

	// ADD
	s.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {

		data := []byte(`{"Id":1, "Name":"tiguan","Brand":"vw","Year":2016}`)

		model := entities.Model{}

		b, err := ioutil.ReadAll(r.Body)
		b = data
		if err == nil {
			err = json.Unmarshal(b, &model)
		}
		r.Body.Close()

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request data")
			return
		}

		item := entities.Model{}
		check, err := item.Get(int(model.Id))
		if check != nil {
			respondWithError(w, http.StatusBadRequest, "Model is already exist")
			return
		}

		if model.Validate() != true {
			respondWithError(w, http.StatusBadRequest, "Validate error")
			return
		}

		err = model.Create()
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error add model")
			return
		}

		respondWithJSON(w, http.StatusOK, check)

	}).Methods("GET")

	// EDIT
	s.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		data := []byte(`{"Id":1,"Year":2010}`)

		updModel := entities.Model{}

		b, err := ioutil.ReadAll(r.Body)
		b = data
		if err == nil {
			err = json.Unmarshal(b, &updModel)
		}
		r.Body.Close()

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request data")
			return
		}
		id := updModel.Id

		item := entities.Model{}
		model, err := item.Get(int(id))
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error getting model")
			return
		}

		if model == nil {
			respondWithError(w, http.StatusBadRequest, "Model not found")
			return
		}

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
			respondWithError(w, http.StatusBadRequest, "New model validate error")
			return
		}

		err = model.Update()
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error update model")
			return
		}

		respondWithJSON(w, http.StatusOK, model)
	}).Methods("POST")

	err := http.ListenAndServe(":8080", r) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServer error: ", err)
	}

	//.Set("Content-Type", "application/json")

	//mm := Model{2,"tiguan", "vw", 2016}

	//// Create a write transaction
	//txn := db.Txn(true)

	m.Brand = "nissan"
	m.Update()

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

/*
GET /<entity>/<id>

GET /user/<id>/reviews

GET /model/<id>/mark

POST /<entity>/<id>

POST /<entity>
*/
