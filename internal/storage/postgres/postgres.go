package postgres

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"greenkmchSever/internal/http-server/handlers/app"
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

func (s *Storage) AddVisitRequest(data app.RequestData) (int64, error) {
	const op = "storage.postgres.addVisitRequest"
	var GroupID int64
	var TypeOfVisitID, FormatOfVisitID, ReasonID uint8
	if len(data.Users) == 1 {
		GroupID = -1
	} else {
		err := s.db.QueryRow("INSERT INTO group_permits DEFAULT VALUES RETURNING id").Scan(&GroupID)
		if err != nil {
			return 0, errors.Errorf("%s: failed to insert visit request: %s", op, err)
		}
	}
	err := s.db.QueryRow("INSERT INTO group_permits DEFAULT VALUES RETURNING id").Scan(&GroupID)
	_, err := s.db.Exec(`INSERT INTO visit_permits (status, requested_at, route_id, visit_date, group_id, visit_reason, visit_format, photo_type, first_name, last_name, middle_name, citizenship, registration_region, is_male, passport, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		1, time.Now(), data.RouteID, data.VisitDate, GroupID,
	)
}
