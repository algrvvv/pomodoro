package config

type AppConfig struct {
	// WorkMinutes колво минут выделенное на работу
	WorkMinutes int `yaml:"workMinutes"`

	// ShortBreakMinutes колво минут выделенное на короткий отдых между сессиями
	ShortBreakMinutes int `yaml:"shortBreakMinutes"`

	// LongBreakMinutes колво минут выделенное на длинный отдых, который начинается
	// после достижения нужного колва пройденных сессий
	LongBreakMinutes int `yaml:"longBreakMinutes"`

	// BreakAfterSessions колво сессий после которых будет запущен длинный отдых
	BreakAfterSessions int `yaml:"breakAfterSessions"`

	// SessionsGoal цель сессий на день
	SessionsGoal int `yaml:"sessionsGoal"`
}

var Config AppConfig
