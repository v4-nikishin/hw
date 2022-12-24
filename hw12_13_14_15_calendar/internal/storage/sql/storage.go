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

func New(ctx context.Context, cfg config.SQLConf, logger *logger.Logger) *Storage {
	return &Storage{ctx: ctx, cfg: cfg, logg: logger}
}

func (s *Storage) Connect(dsn string) (err error) {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	return s.db.PingContext(s.ctx)
}

func (s *Storage) connect() (err error) {
	if err := s.Connect(s.cfg.DSN); err != nil {
		return fmt.Errorf("cannot connect to psql: %w", err)
	}
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) close() {
	if err := s.Close(); err != nil {
		s.logg.Error(fmt.Sprintf("cannot close psql connection: %v", err))
	}
}

func (s *Storage) CreateEvent(e storage.Event) error {
	if err := s.connect(); err != nil {
		return err
	}
	defer s.close()

	query := "insert into events (uuid, title) values ($1, $2)"
	_, err := s.db.ExecContext(s.ctx, query, e.UUID, e.Title)
	if err != nil {
		return fmt.Errorf("cannot add event %w", err)
	}
	return nil
}

func (s *Storage) GetEvent(uuid string) (storage.Event, error) {
	if err := s.connect(); err != nil {
		return storage.Event{}, err
	}
	defer s.close()

	query := "select uuid, title from events where uuid = $1"
	row := s.db.QueryRowContext(s.ctx, query, uuid)

	e := storage.Event{}

	err := row.Scan(&e.UUID, &e.Title)
	if err == sql.ErrNoRows {
		return e, fmt.Errorf("cannot find event %w", err)
	} else if err != nil {
		return e, fmt.Errorf("failed to get event %w", err)
	}
	return e, nil
}

func (s *Storage) UpdateEvent(uuid string, title string) error {
	if err := s.connect(); err != nil {
		return err
	}
	defer s.close()

	query := "update events set title = $1 where uuid = $2"
	_, err := s.db.ExecContext(s.ctx, query, title, uuid)
	if err != nil {
		return fmt.Errorf("cannot update event %w", err)
	}
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	if err := s.connect(); err != nil {
		return err
	}
	defer s.close()

	query := "delete from events where uuid = $1"
	_, err := s.db.ExecContext(s.ctx, query, id)
	if err != nil {
		return fmt.Errorf("cannot delete event %w", err)
	}
	return nil
}

func (s *Storage) Events() ([]storage.Event, error) {
	if err := s.connect(); err != nil {
		return nil, err
	}
	defer s.close()

	rows, err := s.db.QueryContext(s.ctx, `select uuid, title from events`)
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()

	var events []storage.Event

	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.Title,
			&e.UUID,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
