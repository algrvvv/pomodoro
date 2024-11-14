package repos

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/database"
	"github.com/algrvvv/pomodoro/internal/repositories"
)

type PostgresSessionRepo struct{}

func NewPostgresRepo() repositories.SessionRepository {
	return &PostgresSessionRepo{}
}

func (p *PostgresSessionRepo) GetAll() (sessions []repositories.Session, err error) {
	var rows *sql.Rows

	query := "SELECT id, duration, type_id from sessions"
	rows, err = database.C.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var s repositories.Session
		if err = rows.Scan(&s.ID, &s.Minutes, &s.SessionType); err != nil {
			return
		}

		sessions = append(sessions, s)
	}

	return
}

func (p *PostgresSessionRepo) GetTodaySessions() (sessions []repositories.Session, err error) {
	var rows *sql.Rows

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	query := "SELECT id, duration, type_id from sessions WHERE created_at > $1 and created_at < $2"
	rows, err = database.C.Query(query, start, end)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var s repositories.Session
		if err = rows.Scan(&s.ID, &s.Minutes, &s.SessionType); err != nil {
			return
		}

		sessions = append(sessions, s)
	}

	return
}

func (p *PostgresSessionRepo) GoalAchivedToday() (achived bool, err error) {
	var sessions []repositories.Session
	sessions, err = p.GetTodaySessions()
	if err != nil {
		return
	}

	var duration int
	for _, s := range sessions {
		duration += s.Minutes
	}

	return duration >= config.Config.Pomodoro.SessionGoalMinutes, nil
}

func (p *PostgresSessionRepo) Save(m int, t int) (duration int, err error) {
	tr, err := database.C.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err = tr.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				fmt.Println("failed to rollback transaction: ", err)
			} else {
				err = nil
			}
		}
	}()

	query := "INSERT INTO sessions(duration, type_id) VALUES ($1, $2)"
	if _, err = tr.Exec(query, m, t); err != nil {
		return
	}

	query = "SELECT sum(duration) FROM sessions"
	if err = tr.QueryRow(query).Scan(&duration); err != nil {
		return
	}

	if err = tr.Commit(); err != nil {
		return
	}

	return
}
