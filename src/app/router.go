package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gtadam/ashilda-common"
)

type Router struct {
	basePath string
	database models.Database
	mux		 *mux.Router
}

func NewRouter(bp string) *Router {
	return &Router {
		basePath: bp,
		database: *models.NewDatabase(),
		mux:	  mux.NewRouter().StrictSlash(true),
	}
}

func (rt *Router) getUsers(w http.ResponseWriter, r *http.Request) {
	statement := models.NewDatabaseSelect(TABLE)
	statement.AddColumn(ID_FIELD)
	statement.AddColumn(FIRST_NAME_FIELD)
	statement.AddColumn(LAST_NAME_FIELD)
	statement.AddColumn(PASSWORD_FIELD)
	statement.AddColumn(EMAIL_FIELD)
	rows, _ := rt.database.ExecuteSelect(statement)
	users := []User{}
	for rows.Next() {
		user := User{}
		user.populate(rows)
		user.Password = ""
		users = append(users, user)
	}
	rows.Close()
	json, _ := json.Marshal(users)
	fmt.Fprintf(w, string(json))
}

func (rt *Router) getUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	statement := models.NewDatabaseSelect(TABLE)
	statement.AddColumn(ID_FIELD)
	statement.AddColumn(FIRST_NAME_FIELD)
	statement.AddColumn(LAST_NAME_FIELD)
	statement.AddColumn(PASSWORD_FIELD)
	statement.AddColumn(EMAIL_FIELD)
	statement.AddCondition(ID_FIELD, models.EQUALS, id)
	rows, _ := rt.database.ExecuteSelect(statement)
	user := User{}
	rows.Next()
	user.populate(rows)
	user.Password = ""
	rows.Close()
	if (user.User_id == 0) {
		w.WriteHeader(404)
		return
	}
	json, _ := json.Marshal(user)
	fmt.Fprintf(w, string(json))
}

func (rt *Router) putUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	user := User{}
	json.Unmarshal(body, &user)
	statement := models.NewDatabaseUpdate(TABLE)
	statement.AddStatement(FIRST_NAME_FIELD, user.FirstName)
	statement.AddStatement(LAST_NAME_FIELD, user.LastName)
	statement.AddStatement(EMAIL_FIELD, user.Email)
	statement.AddCondition(ID_FIELD, models.EQUALS, strconv.Itoa(user.User_id))
	rt.database.ExecuteUpdate(statement)
}

func (rt *Router) postUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	user := User{}
	json.Unmarshal(body, &user)
	statement := models.NewDatabaseInsert(TABLE)
	statement.AddEntry(FIRST_NAME_FIELD, user.FirstName)
	statement.AddEntry(LAST_NAME_FIELD, user.LastName)
	statement.AddEntry(PASSWORD_FIELD, user.Password)
	statement.AddEntry(EMAIL_FIELD, user.Email)
	rt.database.ExecuteInsert(statement)
}

func (rt *Router) deleteUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	user := User{}
	json.Unmarshal(body, &user)
	statement := models.NewDatabaseDelete(TABLE)
	statement.AddCondition(ID_FIELD, models.EQUALS, strconv.Itoa(user.User_id))
	rt.database.ExecuteDelete(statement)
}

func (rt *Router) populateRoutes() {
	rt.database.Connect()
	rt.mux.HandleFunc(rt.basePath+"/users", rt.getUsers).Methods("GET")
	rt.mux.HandleFunc(rt.basePath+"/user/{id:[0-9]+}", rt.getUser).Methods("GET")
	rt.mux.HandleFunc(rt.basePath+"/user", rt.putUser).Methods("PUT")
	rt.mux.HandleFunc(rt.basePath+"/user", rt.postUser).Methods("POST")
	rt.mux.HandleFunc(rt.basePath+"/user", rt.deleteUser).Methods("DELETE")
}