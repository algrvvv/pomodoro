package repositories

import "time"

type SessionRepository interface {
	GetAll() (sessions []Session, err error)
	GoalAchivedToday() (achived bool, err error)
	GetTodaySessions() (sessions []Session, err error)
	Save(m int, t int) (duration int, err error)
	DeletePrevSessions() (err error)
}

type Session struct {
	ID          int       `json:"id"`
	Minutes     int       `json:"duration"`
	SessionType int       `json:"type"`
	CreatedAt   time.Time `json:"date"`
}
