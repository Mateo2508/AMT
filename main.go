package main

import (
	//"log"
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "amt"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

func main() {

	http.HandleFunc("/", Inicio)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/actualizar", Actualizar)

	fmt.Println("Servidor Corriendo.....")
	http.ListenAndServe(":8080", nil)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")
	fmt.Println(idEmployee)

	conexionEstablecida := conexionBD()
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM employees WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmployee)

	http.Redirect(w, r, "/", 301)

}

type Employee struct {
	Id    int
	User  string
	Email string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM employees ")

	if err != nil {
		panic(err.Error())
	}
	employee := Employee{}
	arregloEmployee := []Employee{}

	for registros.Next() {
		var id int
		var user, email string
		err = registros.Scan(&id, &user, &email)
		if err != nil {
			panic(err.Error())
		}
		employee.Id = id
		employee.User = user
		employee.Email = email
		arregloEmployee = append(arregloEmployee, employee)
	}
	//fmt.Println(arregloEmployee)

	//fmt.Fprintf(w, "Hola Mateo")
	plantillas.ExecuteTemplate(w, "inicio", arregloEmployee)
}

func Edit(w http.ResponseWriter, r *http.Request) {

	idEmployee := r.URL.Query().Get("id")
	fmt.Println(idEmployee)

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM employees WHERE id=?", idEmployee)
	employee := Employee{}

	for registro.Next() {
		var id int
		var user, email string
		err = registro.Scan(&id, &user, &email)
		if err != nil {
			panic(err.Error())
		}
		employee.Id = id
		employee.User = user
		employee.Email = email
	}
	fmt.Println(employee)
	plantillas.ExecuteTemplate(w, "edit", employee)

}

func Add(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "add", nil)

}
func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := r.FormValue("user")
		email := r.FormValue("email")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO employees(user,email) VALUES(?,?) ")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(user, email)

		http.Redirect(w, r, "/", 301)
	}
}
func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		user := r.FormValue("user")
		email := r.FormValue("email")

		conexionEstablecida := conexionBD()
		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE employees SET user=?,email=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(user, email, id)
		http.Redirect(w, r, "/", 301)
	}
}
