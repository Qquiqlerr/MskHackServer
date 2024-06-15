package portal

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type UserInfoGetter interface {
	GetUserInfo(id int) (UserInfo, error)
}

type UserInfo struct {
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	MiddleName         string    `json:"middle_name"`
	Citizenship        string    `json:"citizenship"`
	RegistrationRegion string    `json:"registration_region"`
	IsMale             bool      `json:"is_male"`
	Passport           string    `json:"passport"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	DateOfBirth        string    `json:"date_of_birth"`
	RequestedAt        time.Time `json:"requested_at"`
	VisitDate          time.Time `json:"visit_date"`
	VisitReason        string    `json:"visit_reason"`
	VisitFormat        string    `json:"visit_format"`
}

func GetUserInfo(log *slog.Logger, getter UserInfoGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "portal.handlers.GetUserInfo"
		log = log.With(
			slog.String("operation", op),
		)
		ID, err := strconv.Atoi(r.URL.Query().Get("id"))
		var info UserInfo
		info, err = getter.GetUserInfo(ID)
		if err != nil {
			log.Error("failed to get user info", err.Error())
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, info)
	}
}
