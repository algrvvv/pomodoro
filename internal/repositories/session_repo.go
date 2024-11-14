package repositories

type SessionRepository interface {
	GetAll() (sessions []Session, err error)
	GoalAchivedToday() (achived bool, err error)
	GetTodaySessions() (sessions []Session, err error)
	Save(m int, t int) (duration int, err error)
}

type Session struct {
	ID          int
	Minutes     int
	SessionType int
}
