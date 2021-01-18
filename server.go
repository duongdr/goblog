package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"text/template"
	"time"
)

type PageVariables struct {
	Date         string
	Time         string
}

type Employee struct {
	Id    int
	Name  string
	City string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "12345678"
	dbName := "goblog"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/index", Index)
	http.HandleFunc("/new", New)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/hello", helloWorld)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM employees ORDER BY id ASC ")
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
	fmt.Print(res);
	t, err := template.ParseFiles("Index.html") //parse the html file homepage.html
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, res)
	//tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request)  {
	t, err := template.ParseFiles("New.html") //parse the html file homepage.html
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	t.Execute(w, nil)
}
func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		log.Print(r.FormValue("name"))
		name := r.FormValue("name")
		city := r.FormValue("city")
		insForm, err := db.Prepare("INSERT INTO employees(name, city) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city)
		log.Println("INSERT: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/index", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employees WHERE id=?", nId)
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
	t, err := template.ParseFiles("Edit.html") //parse the html file homepage.html
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	t.Execute(w, emp)
	defer db.Close()
}

func Update(w http.ResponseWriter, r *http.Request)  {
	db :=dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE employees SET name=?, city=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city, id)
		log.Print("UPDATE: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/index", 301)
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employees WHERE id=?", nId)
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
	t, err := template.ParseFiles("Show.html") //parse the html file homepage.html
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	t.Execute(w, emp)
	defer db.Close()
}

func Delete(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM employees WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/index", 301)
}

func helloWorld(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello World")
}
func HomePage(w http.ResponseWriter, r *http.Request){

	now := time.Now() // find the time right now
	HomePageVars := PageVariables{ //store the date and time in a struct
		Date: now.Format("02-01-2006"),
		Time: now.Format("15:04:05"),
	}

	t, err := template.ParseFiles("homepage.html") //parse the html file homepage.html
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}