package repositories

import "time"

type SessionRepository interface {
	// GetAll метод для получения всех сессий с ограничением по месяцам.
	// Нужно для разделения данных по месяцам.
	// Здесь startLimit - верхняя граница, начиная с какого числа получить данные
	GetAll(start, end string) (sessions []Session, err error)
	// GoalAchivedToday проверка достигнута ли цель за последние 24 часа.
	GoalAchivedToday() (achived bool, err error)
	// GetTodaySessions получить сессии за последние 24 часа.
	GetTodaySessions() (sessions []Session, err error)
	// Save метод для сохранения новой сессии.
	Save(m int, t int) (duration int, err error)
	// DeletePrevSessions удалить сессии за предыдущий месяц
	DeletePrevSessions() (err error)
}

type Session struct {
	ID          int       `json:"id"`
	Minutes     int       `json:"duration"`
	SessionType int       `json:"type"`
	CreatedAt   time.Time `json:"date"`
}
