package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var plantillas = template.Must(template.ParseGlob("templates/*"))

type Dev struct {
	Id       int
	Name     string
	Lastname string
	Email    string
}

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/reset-table", resetTable)
	fmt.Println("Server runnig")
	http.ListenAndServe(":8080", nil)
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conecction_setup := connectionDB()
	data, err := conecction_setup.Query("SELECT * from Devps")
	if err != nil {
		panic(err.Error())
	}
	dev := Dev{}

	devs := []Dev{}

	for data.Next() {
		var id int
		var name, lastname, email string
		err := data.Scan(&id, &name, &lastname, &email)
		if err != nil {
			panic(err.Error())
		}

		dev.Id = id
		dev.Name = name
		dev.Lastname = lastname
		dev.Email = email

		devs = append(devs, dev)
	}
	// fmt.Println(devs)

	plantillas.ExecuteTemplate(w, "index", devs)
}

func Create(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "create", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		email := r.FormValue("email")
		conecction_setup := connectionDB()
		insert, err := conecction_setup.Prepare("INSERT INTO Devps (firstname, lastname, email) VALUES (?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}

		insert.Exec(firstname, lastname, email)

		http.Redirect(w, r, "/", 301)

	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	// fmt.Println(id)
	conecction_setup := connectionDB()
	delete, err := conecction_setup.Prepare("DELETE FROM Devps WHERE id =  (?)")
	if err != nil {
		panic(err.Error())
	}

	delete.Exec(id)

	http.Redirect(w, r, "/", 301)
}
func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	conecction_setup := connectionDB()
	data, err := conecction_setup.Query("SELECT * FROM Devps WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	dev := Dev{}

	for data.Next() {
		var id int
		var name, lastname, email string
		err := data.Scan(&id, &name, &lastname, &email)
		if err != nil {
			panic(err.Error())
		}

		dev.Id = id
		dev.Name = name
		dev.Lastname = lastname
		dev.Email = email

	}
	fmt.Println(dev)
	plantillas.ExecuteTemplate(w, "edit", dev)

}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		email := r.FormValue("email")
		fmt.Println(id)
		conecction_setup := connectionDB()
		update, err := conecction_setup.Prepare("UPDATE Devps SET firstname = ?,lastname = ?,email = ? WHERE id = ?")
		if err != nil {
			panic(err.Error())
		}

		update.Exec(firstname, lastname, email, id)

		http.Redirect(w, r, "/", 301)

	}

}

func resetTable(w http.ResponseWriter, r *http.Request) {

	conecction_setup := connectionDB()
	delete, err := conecction_setup.Prepare("DELETE FROM Devps")

	if err != nil {
		panic(err.Error())
	}
	delete.Exec()
	http.Redirect(w, r, "/", 301)

}

func connectionDB() (connection *sql.DB) {
	Driver := "mysql"
	User := "sanetrox"
	Password := "CAMIlodev1994"
	Name := "mysql"

	connection, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Name)

	if err != nil {
		panic(err.Error())
	}
	return connection

}
