package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Response struct {
	Name    string    `json:"name"`
	Pokemon []Pokemon `json:"pokemon_entries"`
}

// A Pokemon Struct to map every pokemon to.
type Pokemon struct {
	EntryNo int            `json:"entry_number"`
	Species PokemonSpecies `json:"pokemon_species"`
}

// A struct to map our Pokemon's Species which includes it's name
type PokemonSpecies struct {
	Name string `json:"name"`
}

type Employee struct {
	Id   int
	Name string
	City string
}

type Mon struct {
	Name string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "123456"
	dbName := "golang"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("form/*"))
}

//var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	//MySQL
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}
	tpl.ExecuteTemplate(w, "index", res)

	// err2 := tpl.ExecuteTemplate(w, "index", res)
	// if err2 != nil {
	// 	log.Println(err2)
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// }

	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	tpl.ExecuteTemplate(w, "show", emp)
	defer db.Close()

	//var commitbranc string
}

func New(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "new", nil)
}

// func Edit(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn()
// 	nId := r.URL.Query().Get("id")
// 	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	emp := Employee{}
// 	for selDB.Next() {
// 		var id int
// 		var name, city string
// 		err = selDB.Scan(&id, &name, &city)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		emp.Id = id
// 		emp.Name = name
// 		emp.City = city
// 	}
// 	tmpl.ExecuteTemplate(w, "Edit", emp)
// 	defer db.Close()
// }

// func Insert(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn()
// 	if r.Method == "POST" {
// 		name := r.FormValue("name")
// 		city := r.FormValue("city")
// 		insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)")
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		insForm.Exec(name, city)
// 		log.Println("INSERT: Name: " + name + " | City: " + city)
// 	}
// 	defer db.Close()
// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
// }

// func Update(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn()
// 	if r.Method == "POST" {
// 		name := r.FormValue("name")
// 		city := r.FormValue("city")
// 		id := r.FormValue("uid")
// 		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		insForm.Exec(name, city, id)
// 		log.Println("UPDATE: Name: " + name + " | City: " + city)
// 	}
// 	defer db.Close()
// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
// }

// func Delete(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn()
// 	emp := r.URL.Query().Get("id")
// 	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	delForm.Exec(emp)
// 	log.Println("DELETE")
// 	defer db.Close()
// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
// }

func ShowPokemon(w http.ResponseWriter, r *http.Request) {
	//API
	response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)
	mon := PokemonSpecies{}
	res := []PokemonSpecies{}
	for i := 0; i < len(responseObject.Pokemon); i++ {
		mon.Name = responseObject.Pokemon[i].Species.Name
		res = append(res, mon)
		//fmt.Println(responseObject.Pokemon[i].Species.Name)
	}
	tpl.ExecuteTemplate(w, "showpokemon", res)
}

func execTemp() {

}

func main() {
	log.Println("Server started on: http://localhost:3000")
	// http.HandleFunc("/", Index)
	// http.HandleFunc("/show", Show)
	// http.HandleFunc("/new", New)
	// http.HandleFunc("/edit", Edit)
	// http.HandleFunc("/insert", Insert)
	// http.HandleFunc("/update", Update)
	// http.HandleFunc("/delete", Delete)
	// http.ListenAndServe(":8080", nil)

	// r := mux.NewRouter()
	// r.PathPrefix("/CRUD/").Handler(
	// 	http.StripPrefix("/CRUD/",
	// 		http.FileServer(http.Dir("/form/"))))
	// r.HandleFunc("/", Index)
	// r.HandleFunc("/show", Show)
	// r.HandleFunc("/new", New)
	// r.HandleFunc("/edit", Edit)
	// r.HandleFunc("/insert", Insert)
	// r.HandleFunc("/update", Update)
	// r.HandleFunc("/delete", Delete)
	// http.ListenAndServe(":3000", r)

	r := mux.NewRouter()

	cssHandler := http.FileServer(http.Dir("./pub/"))

	http.Handle("/pub/", http.StripPrefix("/pub/", cssHandler))
	r.HandleFunc("/", Index)
	r.HandleFunc("/new", New)
	r.HandleFunc("/show", Show)
	r.HandleFunc("/showpokemon", ShowPokemon)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

// func in(writer http.ResponseWriter, request *http.Request) {

// 	db := dbConn()
// 	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	emp := Employee{}
// 	res := []Employee{}
// 	for selDB.Next() {
// 		var id int
// 		var name, city string
// 		err = selDB.Scan(&id, &name, &city)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		emp.Id = id
// 		emp.Name = name
// 		emp.City = city
// 		res = append(res, emp)
// 	}

// 	err2 := tpl.ExecuteTemplate(writer, "index.html", nil)
// 	if err2 != nil {
// 		log.Println(err2)
// 		http.Error(writer, "Internal server error", http.StatusInternalServerError)
// 	}

// 	defer db.Close()

// }
