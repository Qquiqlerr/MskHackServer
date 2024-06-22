package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"greenkmchSever/internal/http-server/handlers/app"
	"greenkmchSever/internal/http-server/handlers/portal"
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

// AddVisitRequest adds a visit request to the storage.
//
// It takes a data parameter of type app.RequestData, which contains the details of the visit request.
// The function scans the visit_reasons and visit_format tables to get the ReasonID and FormatOfVisitID.
// If there is only one user in the data.Users slice, the function inserts the visit request into the visit_permits table.
// Otherwise, it inserts the visit request into the group_permits table and then inserts each user's visit request into the visit_permits table.
// For each user, it also inserts the photo_types into the visit_permits_photo_types table.
//
// Returns an error if any database query fails.
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
		err = s.db.QueryRow(`INSERT INTO visit_permits (status, requested_at, route_id, visit_date, group_id, visit_reason, visit_format, first_name, last_name, middle_name, citizenship, registration_region, is_male, passport, email, phone, date_of_birth) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,$17) returning id`,
			1, time.Now(), data.RouteID, data.VisitDate, GroupID, ReasonID, FormatOfVisitID, data.Users[0].FirstName, data.Users[0].LastName, data.Users[0].MiddleName, data.Users[0].Citizenship, data.Users[0].Region, data.Users[0].IsMale, data.Users[0].Passport, data.Users[0].Email, data.Users[0].Phone, data.Users[0].DateOfBirth,
		).Scan(&GroupID)
		if err != nil {
			return errors.Errorf("%s: failed to insert visit request: %s", op, err)
		}
		fmt.Println(data.Photo)
		for _, photoType := range data.Photo {
			_, err = s.db.Exec(`INSERT INTO visit_permits_photo_types (visit_permit_id, photo_type_id)
                VALUES ($1, (SELECT id FROM photo_types WHERE name = $2))`, GroupID, photoType)
			if err != nil {
				return errors.Errorf("%s: failed to insert visit request: %s", op, err)
			}
		}

	} else {
		err := s.db.QueryRow("INSERT INTO group_permits DEFAULT VALUES RETURNING id").Scan(&GroupID)
		if err != nil {
			return errors.Errorf("%s: failed to insert visit request: %s", op, err)
		}
		for _, user := range data.Users {
			var ID int64
			err = s.db.QueryRow(`INSERT INTO visit_permits (status, requested_at, route_id, visit_date, group_id, visit_reason, visit_format, first_name, last_name, middle_name, citizenship, registration_region, is_male, passport, email, phone,date_of_birth) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,$17) RETURNING id`,
				1, time.Now(), data.RouteID, data.VisitDate, GroupID, ReasonID, FormatOfVisitID, user.FirstName, user.LastName, user.MiddleName, user.Citizenship, user.Region, user.IsMale, user.Passport, user.Email, user.Phone, user.DateOfBirth,
			).Scan(&ID)
			if err != nil {
				return errors.Errorf("%s: failed to insert visit request: %s", op, err)
			}
			fmt.Println(data.Photo[0])
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

func (s *Storage) GetProblems() (portal.Problems, error) {
	const op = "storage.postgres.getAllProblems"
	var problems portal.Problems
	rows, err := s.db.Query(`SELECT r.id, rt.name, rs.name, comment from reports as r join reports_statuses as rs on r.statusID = rs.id  join type_of_reports as rt on r.type = rt.id where r.statusID = 1`)
	if err != nil {
		return problems, errors.Errorf("%s - %s", op, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var problem portal.Problem
		err = rows.Scan(&problem.Id, &problem.Type, &problem.Status, &problem.Comment)
		if err != nil {
			return problems, errors.Errorf("%s: failed to scan problem: %s", op, err)
		}
		problems.Problems = append(problems.Problems, problem)
	}
	return problems, nil
}

func (s *Storage) UpdateProblem(id int64, newStatus string) error {
	const op = "storage.postgres.updateProblem"
	_, err := s.db.Exec("UPDATE reports SET statusID = (SELECT id FROM reports_statuses WHERE name = $1) WHERE id = $2", newStatus, id)
	if err != nil {
		return errors.Errorf("%s: failed to update problem: %s", op, err)
	}
	return nil
}

func (s *Storage) GetAllOopts() ([]portal.Oopt, error) {
	const op = "storage.postgres.getAllOopts"
	var oopts []portal.Oopt
	rows, err := s.db.Query("SELECT name, id, stress FROM zones")
	if err != nil {
		return oopts, errors.Errorf("%s - %s", op, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var oopt portal.Oopt
		err = rows.Scan(&oopt.Name, &oopt.Id, &oopt.Stress)
		if err != nil {
			return oopts, errors.Errorf("%s: failed to scan oopt: %s", op, err)
		}
		oopts = append(oopts, oopt)
	}
	return oopts, nil
}

func (s *Storage) GetAllRoutesFromZone(zoneID int) ([]portal.Route, error) {
	const op = "storage.postgres.getAllRoutes"
	var routes []portal.Route
	rows, err := s.db.Query("SELECT id, name, stress FROM routes WHERE zone_id = $1", zoneID)
	if err != nil {
		return routes, errors.Errorf("%s - %s", op, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var route portal.Route
		err = rows.Scan(&route.ID, &route.Name, &route.Stress)
		if err != nil {
			return routes, errors.Errorf("%s: failed to scan route: %s", op, err)
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func (s *Storage) GetRouteLines(id int) ([]float64, error) {
	const op = "storage.postgres.getRouteLines"
	var lines []float64
	s.db.QueryRow("SELECT lines FROM routes WHERE id = $1", id).Scan(
		pq.Array(&lines),
	)
	return lines, nil
}

func (s *Storage) SendRouteStress(id int, stress float64) error {
	const op = "storage.postgres.sendRouteStress"
	_, err := s.db.Exec("UPDATE routes SET stress = $1 WHERE id = $2", stress, id)
	if err != nil {
		return errors.Errorf("%s - %s", op, err.Error())
	}
	return nil
}

func (s *Storage) SendOOPTStress(stress float64, ID int) error {
	const op = "storage.postgres.sendOOPTStress"
	_, err := s.db.Exec("UPDATE zones SET stress = $1 where id = $2", stress, ID)
	if err != nil {
		return errors.Errorf("%s - %s", op, err.Error())
	}
	return nil
}

func (s *Storage) GetAllZonesForDate() ([]portal.ZonesForDate, error) {
	const op = "storage.postgres.getAllZonesForDate"
	var zones []portal.ZonesForDate
	rows, err := s.db.Query("SELECT id, visit_date, route_id FROM visit_permits where status = 1")
	if err != nil {
		return zones, errors.Errorf("%s - %s", op, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id          int64
			visitDate   time.Time
			routeID     int
			count       int
			stress      float64
			routeString string
		)
		err = rows.Scan(&id, &visitDate, &routeID)
		if err != nil {
			return zones, errors.Errorf("%s: failed to scan zone: %s", op, err)
		}

		err = s.db.QueryRow("SELECT COUNT(*) from visit_permits WHERE status = 3 AND visit_date = $1 AND route_id = $2", visitDate, routeID).Scan(&count)
		if err != nil {
			return zones, errors.Errorf("%s: failed to scan count: %s", op, err)
		}
		err = s.db.QueryRow("SELECT stress from routes WHERE id = $1", routeID).Scan(&stress)
		if err != nil {
			return zones, errors.Errorf("%s: failed to scan stress: %s", op, err)
		}
		err = s.db.QueryRow("SELECT name from routes WHERE id = $1", routeID).Scan(&routeString)
		currStress := float64(count) / stress * 100.0
		stressIfSubmit := float64(count+1) / stress * 100.0
		zones = append(zones, portal.ZonesForDate{
			ID:             id,
			Route:          routeString,
			Stress:         currStress,
			StressIfSubmit: stressIfSubmit,
		})
	}
	return zones, nil
}
func (s *Storage) GetUserInfo(id int) (portal.UserInfo, error) {
	const op = "storage.postgres.getUserInfo"
	var info portal.UserInfo
	err := s.db.QueryRow("SELECT requested_at, visit_date, (SELECT name FROM visit_reasons WHERE id = visit_reason), (select name from visit_format WHERE id = visit_format), first_name, last_name, middle_name, citizenship, registration_region, is_male, passport, email, phone, date_of_birth FROM visit_permits WHERE id = $1", id).Scan(
		&info.RequestedAt, &info.VisitDate, &info.VisitReason, &info.VisitFormat, &info.FirstName,
		&info.LastName, &info.MiddleName, &info.Citizenship, &info.RegistrationRegion, &info.IsMale, &info.Passport,
		&info.Email, &info.Phone, &info.DateOfBirth)
	if err != nil {
		return info, errors.Errorf("%s - %s", op, err.Error())
	}
	return info, nil
}

func (s *Storage) UpdateRequestStatus(id int, status int) error {
	const op = "storage.postgres.updateRequestStatus"
	_, err := s.db.Exec("UPDATE visit_permits SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return errors.Errorf("%s - %s", op, err.Error())
	}
	return nil
}
