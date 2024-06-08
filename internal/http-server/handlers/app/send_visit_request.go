package app

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type RequestAdder interface {
	AddVisitRequest(data RequestData) (int64, error)
}

//type RequestData struct {
//	ID                 int64
//	Status             uint8
//	RequestedAt        time.Time
//	RouteID            int
//	VisitDate          time.Time
//	GroupID            int64
//	VisitReason        int
//	VisitFormat        int
//	PhotoType          int
//	FirstName          string
//	LastName           string
//	MiddleName         string
//	Citizenship        string
//	RegistrationRegion string
//	IsMale             bool
//	passport           string
//	email              string
//	phone              string
//}

type User struct {
	FCs         string `json:"FCs"`
	DateOfBirth string `json:"date_of_birth"`
	Citizenship string `json:"sitizenship"`
	Region      string `json:"region"`
	IsMale      bool   `json:"is_male"`
	Passport    string `json:"passport"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

type RequestData struct {
	Users         []User   `json:"users"`
	VisitDate     string   `json:"visit_date"`
	RouteID       int      `json:"route_id"`
	TypeOfVisit   string   `json:"type_of_visit"`
	FormatOfVisit string   `json:"format_of_visit"`
	Reason        string   `json:"reason"`
	Photo         []string `json:"photo"`
}

type ResponseData struct {
	ID int64 `json:"request_id"`
}

func SendVisitRequest(log *slog.Logger, adder RequestAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data RequestData
		var resp ResponseData
		const op = "internal.handlers.app.sendVisitRequest"
		log = log.With("operation", op)
		err := render.DecodeJSON(r.Body, &data)
		if err != nil {
			log.Error("failed to parse request", err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp)
			return
		}
		if len(data.Users) < 1 {
			log.Error("no users found")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp)
			return
		}
		ID, err := adder.AddVisitRequest(data)
		if err != nil {
			log.Error("failed to add visit request", err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp)
			return
		}

		resp.ID = ID
		render.JSON(w, r, resp)
	}
}
