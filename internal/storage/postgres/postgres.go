package postgres

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"greenkmchSever/internal/http-server/handlers/app"
	"strconv"
	"strings"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(url string) (*Storage, error) {
	const op = "storage.postgres.new"
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, errors.Errorf("%s: failed to connect to postgres: %s", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Errorf("%s: failed to ping postgres: %s", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) GetAllRoutes() ([]app.Routes, error) {
	const op = "storage.postgres.getAllRoutes"
	var routes []app.Routes
	rows, err := s.db.Query("SELECT id, name, length_hours, zone_id, description FROM routes")
	if err != nil {
		return nil, errors.Errorf("%s: failed to query all routes: %s", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var route app.Routes
		err = rows.Scan(&route.ID, &route.Name, &route.Duration, &route.Zone_id, &route.Desc)
		if err != nil {
			return nil, errors.Errorf("%s: failed to scan row: %s", op, err)
		}
		routes = append(routes, route)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Errorf("%s: failed to iterate rows: %s", op, err)
	}
	return routes, nil
}

func (s *Storage) AddVisitRequest(data app.RequestData) error {
	const op = "storage.postgres.addVisitRequest"
	var GroupID int64
	var FormatOfVisitID, ReasonID uint8

	err := s.db.QueryRow("SELECT id FROM visit_reasons WHERE name=$1", data.Reason).Scan(&ReasonID)
	if err != nil {
		return errors.Errorf("%s: failed to query visit request: %s", op, err)
	}
	err = s.db.QueryRow("SELECT id FROM visit_format WHERE name=$1", data.FormatOfVisit).Scan(&FormatOfVisitID)
	if err != nil {
		return errors.Errorf("%s: failed to query visit request: %s", op, err)
	}

	if len(data.Users) == 1 {
		GroupID = -1
		_, err = s.db.Exec(`INSERT INTO visit_permits (status, requested_at, route_id, visit_date, group_id, visit_reason, visit_format, first_name, last_name, middle_name, citizenship, registration_region, is_male, passport, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
			1, time.Now(), data.RouteID, data.VisitDate, GroupID, ReasonID, FormatOfVisitID, data.Users[0].FirstName, data.Users[0].LastName, data.Users[0].MiddleName, data.Users[0].Citizenship, data.Users[0].Region, data.Users[0].IsMale, data.Users[0].Passport, data.Users[0].Email, data.Users[0].Phone,
		)

	} else {
		err := s.db.QueryRow("INSERT INTO group_permits DEFAULT VALUES RETURNING id").Scan(&GroupID)
		if err != nil {
			return errors.Errorf("%s: failed to insert visit request: %s", op, err)
		}
		for _, user := range data.Users {
			var ID int64
			err = s.db.QueryRow(`INSERT INTO visit_permits (status, requested_at, route_id, visit_date, group_id, visit_reason, visit_format, first_name, last_name, middle_name, citizenship, registration_region, is_male, passport, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id`,
				1, time.Now(), data.RouteID, data.VisitDate, GroupID, ReasonID, FormatOfVisitID, user.FirstName, user.LastName, user.MiddleName, user.Citizenship, user.Region, user.IsMale, user.Passport, user.Email, user.Phone,
			).Scan(&ID)
			if err != nil {
				return errors.Errorf("%s: failed to insert visit request: %s", op, err)
			}
			for _, photoType := range data.Photo {
				_, err = s.db.Exec(`INSERT INTO visit_permits_photo_types (visit_permit_id, photo_type_id)
                VALUES ($1, (SELECT id FROM photo_types WHERE name = $2))`, ID, photoType)
				if err != nil {
					return errors.Errorf("%s: failed to insert visit request: %s", op, err)
				}
			}
		}
	}
	return nil
}

func (s *Storage) GetStatus(fn, sn, mn string) (string, error) {
	const op = "storage.postgres.getStatus"
	var Status string
	err := s.db.QueryRow("SELECT ps.name FROM visit_permits vp JOIN permit_statuses ps ON vp.status = ps.id WHERE vp.first_name = $1 AND vp.last_name = $2 AND vp.middle_name = $3;", fn, sn, mn).Scan(&Status)
	if err != nil {
		return "", errors.Errorf("%s: failed to query visit status: %s", op, err)
	}
	return Status, nil
}

func (s *Storage) SaveReport(data app.Report) (int64, error) {
	const op = "storage.postgres.saveReport"
	var ID int64
	err := s.db.QueryRow("INSERT INTO reports(sent_at, reported_at, location, type, photo, comment, email, phone, statusid) VALUES ($1, $2, $3, (SELECT id FROM type_of_reports WHERE name=$4), $5, $6, $7, $8, 1) RETURNING id", time.Now(), data.Time, data.Location, data.TypeOfReport, data.Photo, data.Comment, data.Email, data.Phone).Scan(&ID)
	if err != nil {
		return 0, errors.Errorf("%s: failed to save report: %s", op, err)
	}
	return ID, nil
}

func (s *Storage) GetAllReports() ([]app.ReportCutted, error) {
	const op = "storage.postgres.getAllReports"
	var Reports []app.ReportCutted
	rows, err := s.db.Query("SELECT statusid, type, location, rs.name, rt.name from reports as r JOIN reports_statuses as rs ON r.statusID = rs.id JOIN type_of_reports as rt ON r.type = rt.id")
	defer rows.Close()
	if err != nil {
		return Reports, errors.Errorf("%s - %s", op, err.Error())
	}
	for rows.Next() {
		var report app.ReportCutted
		var location string
		err = rows.Scan(&report.StatusID, &report.TypeID, &location, &report.StatusName, &report.TypeOfReport)
		if err != nil {
			return Reports, errors.Errorf("%s: failed to scan report: %s", op, err)
		}
		coords := strings.Split(location, " ")
		for _, coord := range coords {
			coordf, err := strconv.ParseFloat(coord, 64)
			if err != nil {
				return Reports, errors.Errorf("failed to parse float")
			}
			report.Coords = append(report.Coords, coordf)
		}
		Reports = append(Reports, report)
	}
	return Reports, nil
}
