package driver

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/driver/internal/adapter"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/driver/internal/domain"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/internal/utils"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/cookies"
)

type controllers struct {
	db *sql.DB
}

func newControllers(db *sql.DB) *controllers {
	return &controllers{
		db: db,
	}
}

func (c *controllers) getDrivers(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /drivers")
	cookie, err := r.Cookie("userCookie")
	if err != nil {
		log.Println("GET /drivers: cookie not found")
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	value := make(map[string]string)

	err = cookies.S.Decode("userCookie", cookie.Value, &value)
	if err != nil {
		log.Println("GET /drivers: error decoding cookie", err)
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(value["UserId"])
	if err != nil {
		log.Println("GET /drivers: error converting UserId to int", err)
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	listDrivers := domain.NewListDrivers(adapter.NewDriverRepositorySqlite(c.db))
	drivers, err := listDrivers.Execute(id)
	if err != nil {
		utils.TrowError(err, w, r)
		return

	}
	tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/drivers.gohtml", "template/drivers-row.gohtml"))

	tpl.ExecuteTemplate(w, "layout", drivers)
}

func (c *controllers) postDrivers(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("userCookie")
	if err != nil {
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	value := make(map[string]string)

	err = cookies.S.Decode("userCookie", cookie.Value, &value)
	if err != nil {
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(value["UserId"])
	if err != nil {
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	register := &domain.RegisterInput{
		Name:      name,
		Email:     email,
		Password:  password,
		CompanyId: uint(id),
	}

	registerDriver := domain.NewRegister(adapter.NewDriverRepositorySqlite(c.db))
	err = registerDriver.Execute(register)
	if err != nil {
		utils.TrowError(err, w, r)
		return
	}

	http.Redirect(w, r, "/drivers", http.StatusSeeOther)
}
