package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib" // To use pgx driver
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/v4-nikishin/hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	ctx  context.Context
	cfg  config.SQLConf
	logg *logger.Logger
	db   *sql.DB
}

func New(ctx context.Context, cfg config.SQLConf, logger *logger.Logger) (*Storage, error) {
	s := &Storage{ctx: ctx, cfg: cfg, logg: logger}
	if err := s.connect(s.cfg.DSN); err != nil {
		return nil, fmt.Errorf("cannot connect to psql: %w", err)
	}
	return s, nil
}

func (s *Storage) connect(dsn string) (err error) {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	return s.db.PingContext(s.ctx)
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		s.logg.Error(fmt.Sprintf("cannot close psql connection: %v", err))
	}
}

func (s *Storage) CreateEvent(e storage.Event) error {
	query := "insert into events (uuid, title, user, date, begin, end) values ($1, $2, $3, $4, $5, $6)"
	_, err := s.db.ExecContext(s.ctx, query, e.UUID, e.Title, e.User, e.Date, e.Begin, e.End)
	if err != nil {
		return fmt.Errorf("cannot add event %w", err)
	}
	return nil
}

func (s *Storage) GetEvent(uuid string) (storage.Event, error) {
	query := "select uuid, title from events where uuid = $1"
	row := s.db.QueryRowContext(s.ctx, query, uuid)

	e := storage.Event{}

	err := row.Scan(&e.UUID, &e.Title, &e.User, &e.Date, &e.Begin, &e.End)
	if err == sql.ErrNoRows {
		return e, fmt.Errorf("cannot find event %w", err)
	} else if err != nil {
		return e, fmt.Errorf("failed to get event %w", err)
	}
	return e, nil
}

func (s *Storage) UpdateEvent(uuid string, e storage.Event) error {
	query := "update events set title=$1, date=$2, begin=$3, end=$4 where uuid = $5"
	_, err := s.db.ExecContext(s.ctx, query, e.Title, e.Date, e.Begin, e.End, uuid)
	if err != nil {
		return fmt.Errorf("cannot update event %w", err)
	}
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	query := "delete from events where uuid = $1"
	_, err := s.db.ExecContext(s.ctx, query, id)
	if err != nil {
		return fmt.Errorf("cannot delete event %w", err)
	}
	return nil
}

func (s *Storage) Events() ([]storage.Event, error) {
	query := "select uuid, title, user, date, begin, end from events"
	rows, err := s.db.QueryContext(s.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()

	var events []storage.Event

	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.UUID,
			&e.Title,
			&e.User,
			&e.Date,
			&e.Begin,
			&e.End,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (s *Storage) EventsOnDate(date string) ([]storage.Event, error) {
	query := "select uuid, title, user, date, begin, end from events where date = $1"
	rows, err := s.db.QueryContext(s.ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()

	var events []storage.Event

	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.UUID,
			&e.Title,
			&e.User,
			&e.Date,
			&e.Begin,
			&e.End,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
