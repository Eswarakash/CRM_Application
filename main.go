package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Customer struct {
	CustomerId              int
	CustomerName            string
	CustomerNumber          int
	CustomerComplaintDesc   string
	CustomerComplaintStatus string
	ResolutionMessage       string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "golangapp"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Customer ORDER BY customerId DESC")
	if err != nil {
		panic(err.Error())
	}
	cust := Customer{}
	res := []Customer{}
	for selDB.Next() {
		var id, number int
		var name, desc, status, resolution string
		err = selDB.Scan(&id, &name, &number, &desc, &status, &resolution)
		if err != nil {
			panic(err.Error())
		}
		cust.CustomerId = id
		cust.CustomerName = name
		cust.CustomerNumber = number
		cust.CustomerComplaintDesc = desc
		cust.CustomerComplaintStatus = status
		cust.ResolutionMessage = resolution
		res = append(res, cust)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("select * from Customer where customerId=?", nId)
	if err != nil {
		panic(err.Error())
	}
	cust := Customer{}
	for selDB.Next() {
		var id, number int
		var name, desc string
		status := "OPEN"
		resolution := "NO-COMMENT"
		err = selDB.Scan(&id, &name, &number, &desc, &status, &resolution)
		if err != nil {
			panic(err.Error())
		}
		cust.CustomerId = id
		cust.CustomerName = name
		cust.CustomerNumber = number
		cust.CustomerComplaintDesc = desc
		cust.CustomerComplaintStatus = status
		cust.ResolutionMessage = resolution
	}
	tmpl.ExecuteTemplate(w, "Show", cust)
	defer db.Close()
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	seldb, err := db.Query("select * from Customer where customerId=?", nId)
	if err != nil {
		panic(err.Error())
	}
	Cust := Customer{}
	for seldb.Next() {
		var id, number int
		var name, desc, status, resolution string
		err = seldb.Scan(&id, &name, &number, &desc, &status, &resolution)
		if err != nil {
			panic(err.Error())
		}
		Cust.CustomerId = id
		Cust.CustomerName = name
		Cust.CustomerNumber = number
		Cust.CustomerComplaintDesc = desc
		Cust.CustomerComplaintStatus = status
		Cust.ResolutionMessage = resolution
		fmt.Println(Cust.CustomerId)
		fmt.Println(Cust)
	}
	tmpl.ExecuteTemplate(w, "Edit", Cust)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("CustomerName")
		number := r.FormValue("CustomerNumber")
		description := r.FormValue("CustomerComplaintDesc")
		resolution := "NO-COMMENT"
		status := "open"
		insForm, err := db.Prepare("INSERT INTO Customer(customerName, customerNumber, customerComplaintDesc, customerComplaintStatus, resolutionMessage) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, number, description, status, resolution)
		log.Println("INSERT: Name: " + name + " | Number: " + number + " | Description: " + description + " | Status: " + status + " | Resolution: " + resolution)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("customerName")
		number := r.FormValue("customerNumber")
		description := r.FormValue("customerComplaintDesc")
		status := r.FormValue("customerComplaintStatus")
		resolution := r.FormValue("resolutionMessage")
		id := r.FormValue("customerId")
		insForm, err := db.Prepare("UPDATE Customer SET customerName=?, customerNumber=?, customerComplaintDesc=?, customerComplaintStatus=?, resolutionMessage= ?  WHERE customerId=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, number, description, status, resolution, id)
		log.Println("UPDATE: Name: " + name + " | Number: " + number + " | Description: " + description + " | Status: " + status + " | Resolution: " + resolution)
	}

	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func TicketSearch(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "TicketSearch", nil)
}

func Search1(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	cust := Customer{}
	c := []Customer{}
	name := r.FormValue("customerName")
	selDB, err := db.Query("Select *from Customer where customerName= (?) ", name)
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		var id, number int
		var name, desc, status, resolution string
		err = selDB.Scan(&id, &name, &number, &desc, &status, &resolution)
		if err != nil {
			panic(err.Error())
		}
		cust.CustomerId = id
		cust.CustomerName = name
		cust.CustomerNumber = number
		cust.CustomerComplaintDesc = desc
		cust.CustomerComplaintStatus = status
		cust.ResolutionMessage = resolution
		fmt.Println(cust)
		c = append(c, cust)
	}
	tmpl.ExecuteTemplate(w, "Index", c)
}

func Search2(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	cust := Customer{}
	c := []Customer{}
	name := r.FormValue("customerComplaintStatus")
	selDB, err := db.Query("Select *from Customer where customerComplaintStatus= (?) ", name)
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		var id, number int
		var name, desc, status, resolution string
		err = selDB.Scan(&id, &name, &number, &desc, &status, &resolution)
		if err != nil {
			panic(err.Error())
		}
		cust.CustomerId = id
		cust.CustomerName = name
		cust.CustomerNumber = number
		cust.CustomerComplaintDesc = desc
		cust.CustomerComplaintStatus = status
		cust.ResolutionMessage = resolution
		fmt.Println(cust)
		c = append(c, cust)
	}
	tmpl.ExecuteTemplate(w, "Index", c)
}

func Search3(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	cust := Customer{}
	c := []Customer{}
	name := r.FormValue("customerId")
	selDB, err := db.Query("Select *from Customer where customerId= (?) ", name)
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		var id, number int
		var name, desc, status, resolution string
		err = selDB.Scan(&id, &name, &number, &desc, &status, &resolution)
		if err != nil {
			panic(err.Error())
		}
		cust.CustomerId = id
		cust.CustomerName = name
		cust.CustomerNumber = number
		cust.CustomerComplaintDesc = desc
		cust.CustomerComplaintStatus = status
		cust.ResolutionMessage = resolution
		fmt.Println(cust)
		c = append(c, cust)
	}
	tmpl.ExecuteTemplate(w, "Index", c)
}

func Name(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Name", nil)
}

func Status(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Status", nil)
}

func Id(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Id", nil)
}

func ExitApplication(w http.ResponseWriter, r *http.Request) {
	//db := dbConn()
	fmt.Println("THANK YOU USING GeNiE")
	os.Exit(0)
}

func main() {
	log.Println("Server started on: http://localhost:31000")
	http.HandleFunc("/", Index)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/ticketSearch", TicketSearch)
	http.HandleFunc("/name", Name)
	http.HandleFunc("/status", Status)
	http.HandleFunc("/id", Id)
	http.HandleFunc("/search1", Search1)
	http.HandleFunc("/search2", Search2)
	http.HandleFunc("/search3", Search3)
	http.HandleFunc("/exitApplication", ExitApplication)
	http.ListenAndServe(":31000", nil)
}
