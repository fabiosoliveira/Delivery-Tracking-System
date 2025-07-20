package delivery

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/delivery/internal/adapter"
	"github.com/fabiosoliveira/Delivery-Tracking-System/internal/app/delivery/internal/domain"
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

func (c *controllers) getDelivery(w http.ResponseWriter, r *http.Request) {
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

	userDao := adapter.NewUserDaoSqlite(c.db)

	drivers, err := userDao.ListDriversByCompanyId(id)
	if err != nil {
		utils.TrowError(err, w, r)
		return
	}

	listDelivery := domain.NewListDelivery(adapter.NewDeliveryRepositorySqlite(c.db), userDao)

	deliveries, err := listDelivery.Execute(id)
	if err != nil {
		utils.TrowError(err, w, r)
		return
	}

	data := struct {
		Drivers    []domain.User
		Deliveries []domain.ListDeliveryOutput
	}{
		Drivers:    drivers,
		Deliveries: deliveries,
	}

	tpl := template.Must(template.ParseFiles("template/layout.gohtml", "template/delivery.gohtml", "template/delivery-row.gohtml"))

	tpl.ExecuteTemplate(w, "layout", data)
}

func (c *controllers) getDeliveryLocalization(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /delivery/{id}/localization")

	tpl := template.Must(template.ParseFiles("template/map-localization.gohtml"))

	tpl.Execute(w, nil)
}

func (c *controllers) postDelivery(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("userCookie")
	if err != nil {
		log.Println("Error retrieving cookie:", err)
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	value := make(map[string]string)

	err = cookies.S.Decode("userCookie", cookie.Value, &value)
	if err != nil {
		log.Println("Error decoding cookie:", err)
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	companyId, err := strconv.Atoi(value["UserId"])
	if err != nil {
		log.Println("Error converting UserId to int:", err)
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recipient := r.FormValue("recipient")
	address := r.FormValue("address")
	driver_id := r.FormValue("driver_id")

	driverId, err := strconv.Atoi(driver_id)
	if err != nil {
		log.Println("Error converting DriverId to int:", err)
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
		return
	}

	log.Println("Received data - Driver ID:", driverId, "Recipient:", recipient, "Address:", address)

	createDeliveryInput := &domain.CreateDeliveryInput{
		CompanyId: uint(companyId),
		DriverId:  uint(driverId),
		Recipient: recipient,
		Address:   address,
	}

	log.Println("createDeliveryInput:", createDeliveryInput)

	createDelivery := domain.NewCreateDelivery(adapter.NewDeliveryRepositorySqlite(c.db), adapter.NewUserDaoSqlite(c.db))

	err = createDelivery.Execute(createDeliveryInput)
	if err != nil {
		log.Println("Error executing CreateDelivery:", err)
		utils.TrowError(err, w, r)
		return
	}

	http.Redirect(w, r, "/delivery", http.StatusSeeOther)
}

func (c *controllers) getDeliveryHistory(w http.ResponseWriter, r *http.Request) {
	deliveryId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
		return
	}

	getDeliveryHistory := domain.NewGetDeliveryHistory(adapter.NewDeliveryRepositorySqlite(c.db))
	history, err := getDeliveryHistory.Execute(deliveryId)
	if err != nil {
		http.Error(w, "Error fetching delivery history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
